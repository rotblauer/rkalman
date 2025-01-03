package kalman

import (
	"math"

	"github.com/regnull/kalman/geo"
)

const (
	_LAT             = 0
	_LNG             = 1
	_ALTITUDE        = 2
	_VLAT            = 3
	_VLNG            = 4
	minSpeedAccuracy = 0.1 // Meters per second.
	incline          = 5   // Degrees, used to estimate altitude random step.
)

var sqrtOf2 = math.Sqrt(2)                              // Pre-computed to speed up computations.
var inclineFactor = math.Sin(incline * math.Pi / 180.0) // Pre-computed to speed up computations.

// GeoFilter is a Kalman filter that deals with geographic coordinates and altitude.
type GeoFilter struct {
	filter *Filter
}

// GeoProcessNoise is used to initialize the process noise.
type GeoProcessNoise struct {
	// Base latitude to use for computing distances.
	BaseLat float64
	// DistancePerSecond is the expected random walk distance per second.
	DistancePerSecond float64
	// SpeedPerSecond is the expected speed per second change.
	SpeedPerSecond float64
}

// GeoObserved represents a single observation, in geographical coordinates and altitude.
type GeoObserved struct {
	Lat, Lng, Altitude float64 // Geographical coordinates (in degrees) and latitude.
	Speed              float64 // Speed, in meters per second.
	SpeedAccuracy      float64 // Speed accuracy, in meters per second.
	Direction          float64 // Travel direction, in degrees from North, 0 to 360 range.
	DirectionAccuracy  float64 // Direction accuracy, in degrees.
	HorizontalAccuracy float64 // Horizontal accuracy, in meters.
	VerticalAccuracy   float64 // Vertical accuracy, in meters.
}

// GeoEstimated contains estimated location, obtained by processing several observed locations.
type GeoEstimated struct {
	Lat, Lng, Altitude float64
	Speed              float64
	Direction          float64
	HorizontalAccuracy float64
}

// NewGeoFilter creates and returns a new GeoFilter.
func NewGeoFilter(d *GeoProcessNoise) (*GeoFilter, error) {
	metersPerDegreeLat := geo.FastMetersPerDegreeLat(d.BaseLat)
	metersPerDegreeLng := geo.FastMetersPerDegreeLng(d.BaseLat)

	dx := d.DistancePerSecond / sqrtOf2 / metersPerDegreeLat
	dy := d.DistancePerSecond / sqrtOf2 / metersPerDegreeLng
	dz := d.DistancePerSecond * inclineFactor
	dsvx := d.SpeedPerSecond / sqrtOf2 / metersPerDegreeLat
	dsvy := d.SpeedPerSecond / sqrtOf2 / metersPerDegreeLng
	dsvz := d.SpeedPerSecond * inclineFactor
	f, err := NewFilter(&ProcessNoise{
		ST:  1.0,
		SX:  dx,
		SY:  dy,
		SZ:  dz,
		SVX: dsvx,
		SVY: dsvy,
		SVZ: dsvz})
	if err != nil {
		return nil, err
	}
	return &GeoFilter{filter: f}, nil
}

func (g *GeoFilter) Observe(td float64, ob *GeoObserved) error {
	metersPerDegreeLat := geo.FastMetersPerDegreeLat(ob.Lat)
	metersPerDegreeLng := geo.FastMetersPerDegreeLng(ob.Lat)
	directionRad := ob.Direction * math.Pi / 180.0
	directionRadAccuracy := ob.DirectionAccuracy * math.Pi / 180.0
	speedLat := ob.Speed * math.Cos(directionRad) / metersPerDegreeLat
	speedLng := ob.Speed * math.Sin(directionRad) / metersPerDegreeLng
	ob1 := &Observed{
		X:   ob.Lat,
		Y:   ob.Lng,
		Z:   ob.Altitude,
		VX:  speedLat,
		VY:  speedLng,
		VZ:  0.0, // There is no way to estimate vertical speed.
		XA:  ob.HorizontalAccuracy / metersPerDegreeLat,
		YA:  ob.HorizontalAccuracy / metersPerDegreeLng,
		ZA:  ob.VerticalAccuracy,
		VXA: speedLatAccuracy(ob.Speed, ob.SpeedAccuracy, directionRad, directionRadAccuracy, metersPerDegreeLat),
		VYA: speedLngAccuracy(ob.Speed, ob.SpeedAccuracy, directionRad, directionRadAccuracy, metersPerDegreeLng),
		VZA: minSpeedAccuracy,
	}
	return g.filter.Observe(td, ob1)
}

// Estimate returns the best location estimate.
func (g *GeoFilter) Estimate() *GeoEstimated {
	if g.filter.state == nil {
		return nil
	}
	lat := g.filter.state.AtVec(_LAT)
	/*
		panic: runtime error: index out of range [120] with length 91

		goroutine 6 [running]:
		github.com/regnull/kalman/geo.FastMetersPerDegreeLat(...)
		        /home/ia/dev/regnull/kalman/geo/fast_geo_factors.go:200
		github.com/regnull/kalman.(*GeoFilter).Estimate(0xc000070120)
		        /home/ia/dev/regnull/kalman/geo_filter.go:111 +0x3f8
		main.(*RKalmanFilterT).EstimateFromObservation(0xc000010270, 0xc0003d44b0)
		        /home/ia/dev/rotblauer/catvector/main.go:501 +0x62d
		main.readStreamRKalmanFilter.func1()
		        /home/ia/dev/rotblauer/catvector/main.go:587 +0x1a9
		created by main.readStreamRKalmanFilter in goroutine 1
		        /home/ia/dev/rotblauer/catvector/main.go:558 +0x21f
		---
		panic: runtime error: index out of range [-40]

		goroutine 1255 [running]:
		github.com/regnull/kalman/geo.FastMetersPerDegreeLat(...)
			/home/ia/go/pkg/mod/github.com/regnull/kalman@v0.0.0-20200908141424-10753ec93999/geo/fast_geo_factors.go:200
		github.com/regnull/kalman.(*GeoFilter).Estimate(0xc01ab9cbb0)
			/home/ia/go/pkg/mod/github.com/regnull/kalman@v0.0.0-20200908141424-10753ec93999/geo_filter.go:111 +0x3ad
		github.com/rotblauer/catd/geo/act.(*ProbableCat).Add(0xc033ab0ad0, {{0xe8bf40, 0x155f7f0}, {0xc0a184f1a0, 0x7}, {0x0, 0x0, 0x0}, {0x156dd18, 0xc0a184f200}, ...})
			/home/ia/dev/rotblauer/catd/geo/act/act.go:183 +0x173
		github.com/rotblauer/catd/api.(*Cat).ImprovedActTracks.func1()
			/home/ia/dev/rotblauer/catd/api/act.go:43 +0x209
		created by github.com/rotblauer/catd/api.(*Cat).ImprovedActTracks in goroutine 1250
			/home/ia/dev/rotblauer/catd/api/act.go:31 +0x254
		---
		panic: runtime error: index out of range [464] with length 91

		goroutine 440 [running]:
		github.com/regnull/kalman/geo.FastMetersPerDegreeLat(...)
				/home/ia/dev/regnull/kalman/geo/fast_geo_factors.go:200
		github.com/regnull/kalman.(*GeoFilter).Estimate(0xc02f836190)
				/home/ia/dev/regnull/kalman/geo_filter.go:128 +0x3f8
		github.com/rotblauer/catd/geo/act.(*ProbableCat).Add(0xc0de909860, {{0xe8cf60, 0x1560ae8}, {0xc08e3422d0, 0x7}, {0x0, 0x0, 0x0}, {0x156f018, 0xc08e342330}, ...})
				/home/ia/dev/rotblauer/catd/geo/act/act.go:191 +0x25d
		github.com/rotblauer/catd/api.(*Cat).ImprovedActTracks.func1()
				/home/ia/dev/rotblauer/catd/api/act.go:43 +0x209
		created by github.com/rotblauer/catd/api.(*Cat).ImprovedActTracks in goroutine 435
				/home/ia/dev/rotblauer/catd/api/act.go:31 +0x254
	*/
	if lat < -90 || lat > 90 {
		return nil
	}

	metersPerDegreeLat := geo.FastMetersPerDegreeLat(lat)
	metersPerDegreeLng := geo.FastMetersPerDegreeLng(lat)
	speedLatMeters := g.filter.state.AtVec(_VLAT) * metersPerDegreeLat
	speedLngMeters := g.filter.state.AtVec(_VLNG) * metersPerDegreeLng
	speed := math.Sqrt(speedLatMeters*speedLatMeters + speedLngMeters*speedLngMeters)
	haLatSquared := g.filter.cov.At(_LAT, _LAT) * metersPerDegreeLat * metersPerDegreeLat
	haLngSquared := g.filter.cov.At(_LNG, _LNG) * metersPerDegreeLng * metersPerDegreeLng
	ha := math.Max(math.Sqrt(haLatSquared), math.Sqrt(haLngSquared))

	return &GeoEstimated{
		Lat:                g.filter.state.AtVec(_LAT),
		Lng:                g.filter.state.AtVec(_LNG),
		Altitude:           g.filter.state.AtVec(_ALTITUDE),
		Speed:              speed,
		HorizontalAccuracy: ha,
	}
}

func speedLatAccuracy(speed float64, speedAccuracy float64, directionRad float64, directionRadAccuracy float64, metersPerDegreeLat float64) float64 {
	ds := math.Cos(directionRad) / metersPerDegreeLat * speedAccuracy
	dr := -speed * math.Sin(directionRad) / metersPerDegreeLat * directionRadAccuracy
	return math.Max(math.Sqrt(ds*ds+dr*dr), minSpeedAccuracy/metersPerDegreeLat)
}

func speedLngAccuracy(speed float64, speedAccuracy float64, directionRad float64, directionRadAccuracy float64, metersPerDegreeLng float64) float64 {
	ds := math.Sin(directionRad) / metersPerDegreeLng * speedAccuracy
	dr := speed * math.Cos(directionRad) / metersPerDegreeLng * directionRadAccuracy
	return math.Max(math.Sqrt(ds*ds+dr*dr), minSpeedAccuracy/metersPerDegreeLng)
}
