package las

import (
	"fmt"

	"github.com/mfbonfigli/gocesiumtiler/v2/internal/geom"
)

type MockLasReader struct {
	Cur  int
	Pts  []geom.Point64
	Srid int
}

// NumberOfPoints returns the number of points stored in the LAS file
func (m *MockLasReader) NumberOfPoints() int {
	return len(m.Pts)
}
func (m *MockLasReader) GetSrid() int {
	return m.Srid
}
func (m *MockLasReader) GetNext() (geom.Point64, error) {
	if m.Cur < len(m.Pts) {
		m.Cur++
		return m.Pts[m.Cur-1], nil
	}
	return geom.Point64{}, fmt.Errorf("point not available")
}
