/*
Package webmercator does the translation to and from WebMercator and WGS84
Gotten from: http://wiki.openstreetmap.org/wiki/Mercator#C.23
*/
package webmercator

import (
	"fmt"
	"math"

	"github.com/rikardNL/tegola/maths"
)

const (
	RMajor = 6378137.0
	RMinor = 6356752.3142
	Ratio  = RMinor / RMajor
)

var Eccent float64
var Com float64

func init() {
	Eccent = math.Sqrt(1.0 - (Ratio * Ratio))
	Com = 0.5 * Eccent
}

func con(phi float64) float64 {
	v := Eccent * math.Sin(phi)
	return math.Pow(((1.0 - v) / (1.0 + v)), Com)
}

//LonToX converts from a Longitude to a X coordinate in WebMercator.
func LonToX(lon float64) float64 {
	return RMajor * maths.DegToRad(lon)
}

// LatToY converts from Latitude to a Y coordinate in WebMercator.
func LatToY(lat float64) float64 {
	lat = math.Min(89.5, math.Max(lat, -89.5))
	phi := maths.DegToRad(lat)
	ts := math.Tan(0.5*((math.Pi*0.5)-phi)) / con(phi)
	return 0 - RMajor*math.Log(ts)
}

// XToLon converts from X coordinate in WebMercator to Lontitude in WGS84
func XToLon(x float64) float64 {
	return maths.RadToDeg(x) / RMajor
}

// YToLat converts from Y coordinate in WebMercator to Latitude in WGS84
func YToLat(y float64) float64 {
	ts := math.Exp(-y / RMajor)
	phi := maths.PiDiv2 - 2*math.Atan(ts)
	dphi := 1.0
	i := 0
	for (math.Abs(dphi) > 0.000000001) && (i < 15) {
		dphi = maths.PiDiv2 - 2*math.Atan(ts*con(phi)) - phi
		phi += dphi
		i++
	}
	return maths.RadToDeg(phi)
}

// ToLonLat given a set of coordinates (x,y) it will convert them to Lon/Lat coordinates. If more then x,y is given (i.e. z, and m) they will be returned untransformed.
func ToLonLat(c ...float64) ([]float64, error) {
	if len(c) < 2 {
		return c, fmt.Errorf("Coords should have at least 2 coords")
	}
	crds := []float64{XToLon(c[0]), YToLat(c[1])}
	crds = append(crds, c[2:]...)
	return crds, nil
}

// ToXY given a set of coordinates (lon,lat) it will convert them to X,Y coordinates. If more then lon/lat is given (i.e. z, and m) they will be returned untransformed.
func ToXY(c ...float64) ([]float64, error) {
	if len(c) < 2 {
		return c, fmt.Errorf("Coords should have at least 2 coords")
	}
	crds := []float64{LonToX(c[0]), LatToY(c[1])}
	crds = append(crds, c[2:]...)
	return crds, nil
}
