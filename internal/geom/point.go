package geom

import (
	"fmt"
)

// Point64 contains data of a Point Cloud Point, namely X,Y,Z coords,
// R,G,B color components, Intensity and Classification. Coordinates are expressed
// as double precision float64 numbers.
type Point64 struct {
	X              float64
	Y              float64
	Z              float64
	R              uint8
	G              uint8
	B              uint8
	Intensity      uint8
	Classification uint8
}

// ToPointFromBaseline returns a Point from this Point64 with coordinates expressed as
// offset from a baseline
func (p Point64) ToPointFromBaseline(baseline Point64) Point32 {
	return NewPoint32(
		float32(p.X-baseline.X),
		float32(p.Y-baseline.Y),
		float32(p.Z-baseline.Z),
		p.R,
		p.G,
		p.B,
		p.Intensity,
		p.Classification,
	)
}

// Point32 Contains data of a Point32 Cloud Point32, namely X,Y,Z coords,
// R,G,B color components, Intensity and Classification. X,Y,Z coordinates
// are expressed as float32 single precision numbers
type Point32 struct {
	X              float32
	Y              float32
	Z              float32
	R              uint8
	G              uint8
	B              uint8
	Intensity      uint8
	Classification uint8
}

// Builds a new Point from the given coordinates, colors, intensity and classification values
func NewPoint32(X, Y, Z float32, R, G, B, Intensity, Classification uint8) Point32 {
	return Point32{
		X:              X,
		Y:              Y,
		Z:              Z,
		R:              R,
		G:              G,
		B:              B,
		Intensity:      Intensity,
		Classification: Classification,
	}
}

// Point32List models a list of Point32. Points are immutable and returned by value.
type Point32List interface {
	Len() int
	Next() (Point32, error)
	Reset()
}

// LinkedPoint wraps a Point32 to create a Linked List
type LinkedPoint struct {
	Next *LinkedPoint
	Pt   Point32
}

// LinkedPointStream is a wrapper helper that allows a LinkedPoint to implement the PointList interface
type LinkedPointStream struct {
	len     int
	current *LinkedPoint
	start   *LinkedPoint
}

// NewLinkedPointStream initializes a linked stream from the given root.
// the Count is not cross-verified, it must be coherent with the actual point count in the linked list.
func NewLinkedPointStream(root *LinkedPoint, len int) *LinkedPointStream {
	return &LinkedPointStream{
		len:     len,
		current: root,
		start:   root,
	}
}

func (l *LinkedPointStream) Next() (Point32, error) {
	if l.current == nil {
		return Point32{}, fmt.Errorf("no more points")
	}
	pt := l.current.Pt
	l.current = l.current.Next
	return pt, nil
}

func (l *LinkedPointStream) Len() int {
	return l.len
}

func (l *LinkedPointStream) Reset() {
	l.current = l.start
}
