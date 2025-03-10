package tiler

import "runtime"

type TilerEvent int

const (
	EventReadLasHeaderStarted TilerEvent = iota
	EventReadLasHeaderCompleted
	EventReadLasHeaderError
	EventPointLoadingStarted
	EventPointLoadingCompleted
	EventPointLoadingError
	EventBuildStarted
	EventBuildCompleted
	EventBuildError
	EventExportStarted
	EventExportCompleted
	EventExportError
)

type TilerOptions struct {
	gridSize         float64
	maxDepth         int
	elevationOffset  float64
	eightBitColors   bool
	geoidElevation   bool
	numWorkers       int
	minPointsPerTile int
	callback         TilerCallback
}

type tilerOptionsFn func(*TilerOptions)

type TilerCallback func(event TilerEvent, inputDesc string, elapsed int64, msg string)

// NewDefaultTilerOptions returns sensible defaults for tiling options
func NewDefaultTilerOptions() *TilerOptions {
	return &TilerOptions{
		gridSize:         20,
		maxDepth:         10,
		elevationOffset:  0,
		numWorkers:       runtime.NumCPU(),
		minPointsPerTile: 5000,
		eightBitColors:   false,
		geoidElevation:   false,
		callback:         nil,
	}
}

// NewTilerOptions returns default tiler options modified using the
// provided manipulating functions
func NewTilerOptions(optFn ...tilerOptionsFn) *TilerOptions {
	opts := NewDefaultTilerOptions()
	for _, fn := range optFn {
		fn(opts)
	}
	return opts
}

// WithGridSize sets the max grid size, i.e. the approximate max allowed spacing between
// any two points at the coarser level of detail. Expressed in meters.
func WithGridSize(size float64) tilerOptionsFn {
	return func(opt *TilerOptions) {
		opt.gridSize = size
	}
}

// WithMaxDepth sets the max depth, i.e. the maximum number of levels the tree can reach.
func WithMaxDepth(maxDepth int) tilerOptionsFn {
	return func(opt *TilerOptions) {
		opt.maxDepth = maxDepth
	}
}

// WithElevationOffset sets the Z offset to force on points, in meters. Only use this
// if the input coordinates are expressed as elevation above the geoid or ellipsoid.
func WithElevationOffset(offset float64) tilerOptionsFn {
	return func(opt *TilerOptions) {
		opt.elevationOffset = offset
	}
}

// WithWorkerNumber sets the number of workers to use to read the las files or to
// run the export jobs
func WithWorkerNumber(numWorkers int) tilerOptionsFn {
	return func(opt *TilerOptions) {
		opt.numWorkers = numWorkers
	}
}

// WithMinPointsPerTile returns the minimum number of points a tile must store to exist.
// Used to avoid almost empty tiles that could be consolidated with their parent.
func WithMinPointsPerTile(minPointsPerTile int) tilerOptionsFn {
	return func(opt *TilerOptions) {
		opt.minPointsPerTile = minPointsPerTile
	}
}

// WithCallback sets a function that should be invoked as the tiler job runs
func WithCallback(callback TilerCallback) tilerOptionsFn {
	return func(opt *TilerOptions) {
		opt.callback = callback
	}
}

// WithEightBitColors true forces the tiler to interpret the color info on the file as eight bit colors
func WithEightBitColors(eightBit bool) tilerOptionsFn {
	return func(opt *TilerOptions) {
		opt.eightBitColors = eightBit
	}
}

// WithGeoidElevation true tells the tiler to interpret the Z elevation as elevation over the geoid
func WithGeoidElevation(geoid bool) tilerOptionsFn {
	return func(opt *TilerOptions) {
		opt.geoidElevation = geoid
	}
}
