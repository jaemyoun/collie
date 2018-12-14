package config

import (
	"log"
	"regexp"
)

func AddFilter(expr, desc string) {
	regExpr, err := regexp.CompilePOSIX(expr)
	if err != nil {
		log.Fatalln("cannot compile the regular expresion to add")
	}
	setting.addedFilters = append(setting.addedFilters, FilterInfo{Filter: regExpr, Desc: desc})
}

func GetFilters() []FilterInfo {
	return setting.addedFilters
}

func SetRegExprOfFilter(regExpr, newRegExpr string) bool {
	for idx, e := range setting.addedFilters {
		if e.Filter.String() == regExpr {
			var err error
			setting.addedFilters[idx].Filter, err = regexp.CompilePOSIX(newRegExpr)
			if err != nil {
				log.Fatalln("cannot compile the regular expresion to add")
			}
			return true
		}
	}
	return false
}

func SetDescOfFilter(regExpr, newDesc string) bool {
	for idx, e := range setting.addedFilters {
		if e.Filter.String() == regExpr {
			setting.addedFilters[idx].Desc = newDesc
			return true
		}
	}
	return false
}

func DeleteFilter(regExpr string) bool {
	for idx, e := range setting.addedFilters {
		if e.Filter.String() == regExpr {
			setting.addedFilters = append(setting.addedFilters[:idx], setting.addedFilters[idx+1:]...)
			return true
		}
	}
	return false
}
