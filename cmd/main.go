package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"

	tiler "github.com/mfbonfigli/gocesiumtiler/v2"
	"github.com/mfbonfigli/gocesiumtiler/v2/internal/utils"
	"github.com/urfave/cli/v2"
)

// this global variable controls the tiler that will be used. Useful to inject mocks during tests.
var tilerProvider func() (tiler.Tiler, error) = func() (tiler.Tiler, error) {
	return tiler.NewGoCesiumTiler()
}

var version = "2.0.0-alpha"

const logo = `
                           _                 _   _ _
  __ _  ___   ___ ___  ___(_)_   _ _ __ ___ | |_(_) | ___ _ __ 
 / _  |/ _ \ / __/ _ \/ __| | | | | '_   _ \| __| | |/ _ \ '__|
| (_| | (_) | (_|  __/\__ \ | |_| | | | | | | |_| | |  __/ |   
 \__, |\___/ \___\___||___/_|\__,_|_| |_| |_|\__|_|_|\___|_|   
  __| | A Cesium Point Cloud tile generator written in golang
 |___/  Copyright YYYY - Massimo Federico Bonfigli    
`

func main() {
	printBanner()
	getCli(defaultCliOptions()).Run(os.Args)
}

func getCli(c *cliOpts) *cli.App {
	return &cli.App{
		Name:    "gocesiumtiler",
		Usage:   "transforms LAS files into Cesium.JS 3D Tiles",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:  "file",
				Usage: "convert a LAS file into 3D tiles",
				Flags: getFileFlags(c),
				Action: func(cCtx *cli.Context) error {
					fileCommand(c, cCtx.Args().First())
					return nil
				},
			},
			{
				Name:  "folder",
				Usage: "convert all LAS files in a folder file into 3D tiles",
				Flags: getFolderFlags(c),
				Action: func(cCtx *cli.Context) error {
					folderCommand(c, cCtx.Args().First())
					return nil
				},
			},
		},
		EnableBashCompletion: true,
	}
}

func getFileFlags(c *cliOpts) []cli.Flag {
	return getFlags(c)
}

func getFolderFlags(c *cliOpts) []cli.Flag {
	stdFlags := getFlags(c)
	joinFlag := &cli.BoolFlag{
		Name:        "join",
		Aliases:     []string{"j"},
		Value:       c.join,
		Usage:       "merge the input LAS files in the folder into a single cloud. The LAS files must have the same properties (CRS etc)",
		Destination: &c.join,
	}
	return append(stdFlags, joinFlag)
}

func getFlags(c *cliOpts) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "out",
			Aliases:     []string{"o"},
			Value:       c.output,
			Usage:       "full path of the output folder where to save the resulting Cesium tilesets",
			Destination: &c.output,
		},
		&cli.IntFlag{
			Name:        "epsg",
			Aliases:     []string{"e"},
			Value:       c.epsg,
			Usage:       "EPSG code of the input coordinate system",
			Destination: &c.epsg,
		},
		&cli.Float64Flag{
			Name:        "resolution",
			Aliases:     []string{"r"},
			Value:       c.resolution,
			Usage:       "minimum resolution of the 3d tiles, in meters. approximately represets the maximum sampling distance between any two points at the lowest level of detail",
			Destination: &c.resolution,
		},
		&cli.Float64Flag{
			Name:        "z-offset",
			Aliases:     []string{"z"},
			Value:       c.zOffset,
			Usage:       "z offset to apply to the point, in meters. only use it if the input elevation is referred to the WGS84 ellipsoid or geoid",
			Destination: &c.zOffset,
		},
		&cli.IntFlag{
			Name:        "depth",
			Aliases:     []string{"d"},
			Value:       c.maxDepth,
			Usage:       "maximum depth of the output tree.",
			Destination: &c.maxDepth,
		},
		&cli.IntFlag{
			Name:        "min-points-per-tile",
			Aliases:     []string{"m"},
			Value:       c.minPoints,
			Usage:       "minimum number of points to enforce in each 3D tile",
			Destination: &c.minPoints,
		},
		&cli.BoolFlag{
			Name:        "geoid",
			Aliases:     []string{"g"},
			Value:       c.geoid,
			Usage:       "set to interpret input points elevation as relative to the Earth geoid",
			Destination: &c.geoid,
		},
		&cli.BoolFlag{
			Name:        "8-bit",
			Value:       c.eightBit,
			Usage:       "set to interpret the input points color as part of a 8bit color space",
			Destination: &c.eightBit,
		},
	}
}

type cliOpts struct {
	output     string
	epsg       int
	maxDepth   int
	minPoints  int
	resolution float64
	zOffset    float64
	geoid      bool
	eightBit   bool
	join       bool
}

func defaultCliOptions() *cliOpts {
	return &cliOpts{
		epsg:       -1,
		maxDepth:   10,
		minPoints:  5000,
		resolution: 20,
		zOffset:    0,
		geoid:      false,
		eightBit:   false,
		join:       false,
	}
}

func (c *cliOpts) validate() {
	if c.output == "" {
		log.Fatal("output flag must be set")
	}
	if c.epsg <= 0 {
		log.Fatal("epsg code is invalid")
	}
	if c.maxDepth <= 1 || c.maxDepth > 20 {
		log.Fatal("depth should be between 1 and 20")
	}
	if c.minPoints < 1 {
		log.Fatal("min-points-per-tile should be at least 1")
	}
	if c.resolution < 0.5 || c.resolution > 1000 {
		log.Fatal("resolution should be between 1 and 1000 meters")
	}
}

func (c *cliOpts) print() {
	fmt.Printf(`*** Execution settings:
- EPSG Code: %d,
- Max Depth: %d,
- Resolution: %f meters,
- Min Points per tile: %d
- Z-Offset: %f meters,
- Geoid elevation: %v,
- 8Bit Color: %v
- Join Clouds: %v

`, c.epsg, c.maxDepth, c.resolution, c.minPoints, c.zOffset, c.geoid, c.eightBit, c.join)
}

func (c *cliOpts) getTilerOptions() *tiler.TilerOptions {
	c.validate()
	return tiler.NewTilerOptions(
		tiler.WithEightBitColors(c.eightBit),
		tiler.WithGeoidElevation(c.geoid),
		tiler.WithElevationOffset(c.zOffset),
		tiler.WithGridSize(c.resolution),
		tiler.WithMaxDepth(c.maxDepth),
		tiler.WithMinPointsPerTile(c.minPoints),
		tiler.WithCallback(eventListener),
	)
}

func fileCommand(opts *cliOpts, filepath string) {
	t, err := tilerProvider()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("*** Mode: File, process LAS file at %s\n", filepath)
	opts.print()
	tilerOpts := opts.getTilerOptions()
	runnable := func(ctx context.Context) error {
		return t.ProcessFiles([]string{filepath}, opts.output, opts.epsg, tilerOpts, ctx)
	}
	launch(runnable)
}

func folderCommand(opts *cliOpts, folderpath string) {
	t, err := tilerProvider()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("*** Mode: Folder, process all files in %s\n", folderpath)
	opts.print()
	tilerOpts := opts.getTilerOptions()
	runnable := func(ctx context.Context) error {
		if opts.join {
			files, err := utils.FindLasFilesInFolder(folderpath)
			if err != nil {
				return err
			}
			return t.ProcessFiles(files, opts.output, opts.epsg, tilerOpts, ctx)
		}
		return t.ProcessFolder(folderpath, opts.output, opts.epsg, tilerOpts, ctx)
	}
	launch(runnable)
}

func launch(function func(ctx context.Context) error) {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := function(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()
	wg.Wait()
}

func eventListener(e tiler.TilerEvent, filename string, elapsed int64, msg string) {
	fmt.Printf("[%s] [%s] %s\n", time.Now().UTC().Format("2006-01-02 15:04:05.000"), filename, msg)
}

func printBanner() {
	fmt.Println(strings.ReplaceAll(logo, "YYYY", strconv.Itoa(time.Now().Year())))
}
