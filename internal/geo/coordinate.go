package geo

import (
	"errors"
	"math"
)

type Point struct{ Latitude, Longitude float64 }
type Bounds struct{ SouthWest, NorthEast Point }

func (p Point) Valid() bool {
	return p.Latitude >= -90 && p.Latitude <= 90 && p.Longitude >= -180 && p.Longitude <= 180 && !math.IsNaN(p.Latitude) && !math.IsNaN(p.Longitude)
}
func (b Bounds) Valid() bool {
	return b.SouthWest.Valid() && b.NorthEast.Valid() && b.SouthWest.Latitude <= b.NorthEast.Latitude && b.SouthWest.Longitude <= b.NorthEast.Longitude
}
func (b Bounds) Contains(p Point) bool {
	return b.Valid() && p.Valid() && p.Latitude >= b.SouthWest.Latitude && p.Latitude <= b.NorthEast.Latitude && p.Longitude >= b.SouthWest.Longitude && p.Longitude <= b.NorthEast.Longitude
}
func Distance(a, b Point) float64 {
	const r = 6371000.
	p1, p2 := a.Latitude*math.Pi/180, b.Latitude*math.Pi/180
	dp := (b.Latitude - a.Latitude) * math.Pi / 180
	dl := (b.Longitude - a.Longitude) * math.Pi / 180
	h := math.Sin(dp/2)*math.Sin(dp/2) + math.Cos(p1)*math.Cos(p2)*math.Sin(dl/2)*math.Sin(dl/2)
	return 2 * r * math.Asin(math.Sqrt(h))
}
func ValidatePolygon(points []Point) error {
	if len(points) < 3 {
		return errors.New("polygon needs three points")
	}
	for _, p := range points {
		if !p.Valid() {
			return errors.New("invalid point")
		}
	}
	return nil
}
