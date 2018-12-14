package config

import "strings"

func SetCurrentLocation(location string) {
	location = strings.TrimSpace(location)
	if len(location) >= 1 && location[len(location)-1] != '/' {
		location += "/"
	}
	setting.currentLocation = location
}

func GetCurrentLocation() string {
	return setting.currentLocation
}
