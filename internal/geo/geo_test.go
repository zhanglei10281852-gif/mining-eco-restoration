package geo_test

import (
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/geo"
	"math"
	"testing"
)

func TestPointBounds(t *testing.T) {
	p := geo.Point{Latitude: 39.9, Longitude: 116.4}
	if !p.Valid() {
		t.Fatal()
	}
	if (geo.Point{Latitude: 91}).Valid() {
		t.Fatal()
	}
	b := geo.Bounds{SouthWest: geo.Point{39, 116}, NorthEast: geo.Point{40, 117}}
	if !b.Contains(p) || b.Contains(geo.Point{41, 116}) {
		t.Fatal()
	}
	if !b.Valid() {
		t.Fatal()
	}
}
func TestDistanceAndClamp(t *testing.T) {
	d := geo.Distance(geo.Point{0, 0}, geo.Point{0, 1})
	if d < 100000 || d > 120000 {
		t.Fatal(d)
	}
	p := geo.Clamp(geo.Point{100, -200})
	if p.Latitude != 90 || p.Longitude != -180 {
		t.Fatal(p)
	}
}
func TestPolygonGeometry(t *testing.T) {
	pts := []geo.Point{{0, 0}, {0, 1}, {1, 1}, {1, 0}}
	if e := geo.ValidatePolygon(pts); e != nil {
		t.Fatal(e)
	}
	if geo.Area(pts) != 1 {
		t.Fatal(geo.Area(pts))
	}
	c := geo.Centroid(pts)
	if math.Abs(c.Latitude-.5) > .001 || math.Abs(c.Longitude-.5) > .001 {
		t.Fatal(c)
	}
	if geo.ValidatePolygon(pts[:2]) == nil {
		t.Fatal()
	}
}
