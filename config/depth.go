package config

func SetDepth(d int) int {
	setting.depth = d
	return GetDepth()
}

func GetDepth() int {
	return setting.depth
}

func RemoveDepth() {
	setting.depth = -1
}
