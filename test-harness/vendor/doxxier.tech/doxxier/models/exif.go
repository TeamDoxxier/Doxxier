package models

import "time"

type GpsInfo struct {
	Latitude     float64
	Longitude    float64
	Date         time.Time
	Time         uint32
	Altitude     float32
	AltitudeRef  bool
	LatitudeRef  bool
	LongitudeRef bool
}
