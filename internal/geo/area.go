package geo

import "math"

func Area(points []Point) float64 {
	if len(points) < 3 {
		return 0
	}
	sum := 0.
	for i := range points {
		j := (i + 1) % len(points)
		sum += points[i].Longitude*points[j].Latitude - points[j].Longitude*points[i].Latitude
	}
	return math.Abs(sum) / 2
}
func Centroid(points []Point) Point {
	if len(points) == 0 {
		return Point{}
	}
	var lat, lon float64
	for _, p := range points {
		lat += p.Latitude
		lon += p.Longitude
	}
	n := float64(len(points))
	return Point{lat / n, lon / n}
}
func Clamp(p Point) Point {
	if p.Latitude > 90 {
		p.Latitude = 90
	}
	if p.Latitude < -90 {
		p.Latitude = -90
	}
	if p.Longitude > 180 {
		p.Longitude = 180
	}
	if p.Longitude < -180 {
		p.Longitude = -180
	}
	return p
}
