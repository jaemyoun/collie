package config

func ToggleDetails() bool {
	setting.details = !setting.details
	return GetDetails()
}

func GetDetails() bool {
	return setting.details
}
