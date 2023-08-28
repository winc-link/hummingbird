package dtos

import "fmt"

type Gps struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func (gps Gps) ToLocation() string {
	return fmt.Sprintf("%f,%f", gps.Latitude, gps.Longitude)
}
