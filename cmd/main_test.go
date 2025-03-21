package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	tiler "github.com/mfbonfigli/gocesiumtiler/v2"
	"github.com/mfbonfigli/gocesiumtiler/v2/internal/utils"
)

func TestDefaultTiler(t *testing.T) {
	tl, err := tilerProvider()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	switch tl.(type) {
	case *tiler.GoCesiumTiler:
	default:
		t.Errorf("unexpected tiler type returned")
	}
}

func TestMainProcessFile(t *testing.T) {
	mockTiler := &tiler.MockTiler{}
	tilerProvider = func() (tiler.Tiler, error) {
		return mockTiler, nil
	}
	os.Args = []string{"gocesiumtiler", "file",
		"-out", ".\\abc",
		"-epsg", "4979",
		"-resolution", "11.1",
		"-z-offset", "-1",
		"-depth", "13",
		"-min-points-per-tile", "1200",
		"-geoid", "-8-bit",
		"myfile.las"}
	main()
	if mockTiler.ProcessFilesCalled != true {
		t.Error("expected processFiles called but was not")
	}
	if actual := mockTiler.InputFiles; !reflect.DeepEqual(actual, []string{"myfile.las"}) {
		t.Errorf("expected tiler to be called with %v but got %v", []string{"myfile.las"}, actual)
	}
	if actual := mockTiler.EpsgCode; actual != 4979 {
		t.Errorf("expected tiler to be called with epsg %v but got epsg %v", 4979, actual)
	}
	if actual := mockTiler.OutputFolder; actual != ".\\abc" {
		t.Errorf("expected tiler to be called with output folder %v but got %v", ".\\abc", actual)
	}
	if actual := mockTiler.EightBit; actual != true {
		t.Errorf("expected tiler to be called with EightBit %v but got %v", true, actual)
	}
	if actual := mockTiler.GeoidElev; actual != true {
		t.Errorf("expected tiler to be called with GeoidElev %v but got %v", true, actual)
	}
	if actual := mockTiler.GridSize; actual != 11.1 {
		t.Errorf("expected tiler to be called with GridSize %v but got %v", 11.1, actual)
	}
	if actual := mockTiler.PtsPerTile; actual != 1200 {
		t.Errorf("expected tiler to be called with PtsPerTile %v but got %v", 1200, actual)
	}
	if actual := mockTiler.Depth; actual != 13 {
		t.Errorf("expected tiler to be called with Depth %v but got %v", 13, actual)
	}
	if actual := mockTiler.ElevOffset; actual != -1 {
		t.Errorf("expected tiler to be called with ElevOffset %v but got %v", -1, actual)
	}
}

func TestMainProcessFolder(t *testing.T) {
	mockTiler := &tiler.MockTiler{}
	tilerProvider = func() (tiler.Tiler, error) {
		return mockTiler, nil
	}
	os.Args = []string{"gocesiumtiler", "folder",
		"-out", ".\\abc",
		"-epsg", "4979",
		"-resolution", "11.1",
		"-z-offset", "-1",
		"-depth", "13",
		"-min-points-per-tile", "1200",
		"-geoid", "-8-bit",
		"myfolder"}
	main()
	if mockTiler.ProcessFolderCalled != true {
		t.Error("expected processFolder called but was not")
	}
	if actual := mockTiler.InputFolder; !reflect.DeepEqual(actual, "myfolder") {
		t.Errorf("expected tiler to be called with %v but got %v", "myfolder", actual)
	}
	if actual := mockTiler.EpsgCode; actual != 4979 {
		t.Errorf("expected tiler to be called with epsg %v but got epsg %v", 4979, actual)
	}
	if actual := mockTiler.OutputFolder; actual != ".\\abc" {
		t.Errorf("expected tiler to be called with output folder %v but got %v", ".\\abc", actual)
	}
	if actual := mockTiler.EightBit; actual != true {
		t.Errorf("expected tiler to be called with EightBit %v but got %v", true, actual)
	}
	if actual := mockTiler.GeoidElev; actual != true {
		t.Errorf("expected tiler to be called with GeoidElev %v but got %v", true, actual)
	}
	if actual := mockTiler.GridSize; actual != 11.1 {
		t.Errorf("expected tiler to be called with GridSize %v but got %v", 11.1, actual)
	}
	if actual := mockTiler.PtsPerTile; actual != 1200 {
		t.Errorf("expected tiler to be called with PtsPerTile %v but got %v", 1200, actual)
	}
	if actual := mockTiler.Depth; actual != 13 {
		t.Errorf("expected tiler to be called with Depth %v but got %v", 13, actual)
	}
	if actual := mockTiler.ElevOffset; actual != -1 {
		t.Errorf("expected tiler to be called with ElevOffset %v but got %v", -1, actual)
	}
}

func TestMainProcessFolderJoin(t *testing.T) {
	tmp, err := os.MkdirTemp(os.TempDir(), "tst")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(tmp)
	})

	utils.TouchFile(filepath.Join(tmp, "test0.las"))
	utils.TouchFile(filepath.Join(tmp, "test0.xyz"))
	utils.TouchFile(filepath.Join(tmp, "test1.LAS"))
	utils.TouchFile(filepath.Join(tmp, "test2.LAS"))

	mockTiler := &tiler.MockTiler{}
	tilerProvider = func() (tiler.Tiler, error) {
		return mockTiler, nil
	}
	os.Args = []string{"gocesiumtiler", "folder",
		"-out", ".\\abc",
		"-epsg", "4979",
		"-resolution", "11.1",
		"-z-offset", "-1",
		"-depth", "13",
		"-min-points-per-tile", "1200",
		"-geoid", "-8-bit",
		"-join",
		tmp}
	main()
	if mockTiler.ProcessFolderCalled != false {
		t.Error("expected processFolder to not be called but it was")
	}
	if mockTiler.ProcessFilesCalled != true {
		t.Error("expected processFiles called but was not")
	}
	expected := []string{
		filepath.Join(tmp, "test0.las"),
		filepath.Join(tmp, "test1.LAS"),
		filepath.Join(tmp, "test2.LAS"),
	}
	if actual := mockTiler.InputFiles; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected tiler to be called with %v but got %v", expected, actual)
	}
	if actual := mockTiler.EpsgCode; actual != 4979 {
		t.Errorf("expected tiler to be called with epsg %v but got epsg %v", 4979, actual)
	}
	if actual := mockTiler.OutputFolder; actual != ".\\abc" {
		t.Errorf("expected tiler to be called with output folder %v but got %v", ".\\abc", actual)
	}
	if actual := mockTiler.EightBit; actual != true {
		t.Errorf("expected tiler to be called with EightBit %v but got %v", true, actual)
	}
	if actual := mockTiler.GeoidElev; actual != true {
		t.Errorf("expected tiler to be called with GeoidElev %v but got %v", true, actual)
	}
	if actual := mockTiler.GridSize; actual != 11.1 {
		t.Errorf("expected tiler to be called with GridSize %v but got %v", 11.1, actual)
	}
	if actual := mockTiler.PtsPerTile; actual != 1200 {
		t.Errorf("expected tiler to be called with PtsPerTile %v but got %v", 1200, actual)
	}
	if actual := mockTiler.Depth; actual != 13 {
		t.Errorf("expected tiler to be called with Depth %v but got %v", 13, actual)
	}
	if actual := mockTiler.ElevOffset; actual != -1 {
		t.Errorf("expected tiler to be called with ElevOffset %v but got %v", -1, actual)
	}
}
