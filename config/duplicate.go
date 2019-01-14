package config

func ToggleDuplication() bool {
	setting.duplication = !setting.duplication
	return GetDetails()
}

func GetDuplication() bool {
	return setting.duplication
}
