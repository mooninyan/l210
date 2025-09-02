package main

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func checkSortByColumn(list []string, config *Config) bool {
	listOfList := make([][]string, len(list))
	for i, str := range list {
		listOfList[i] = strings.Fields(str)
	}

	return sliceIsSorted(listOfList, config)
}

func sliceIsSorted(listOfList [][]string, config *Config) bool {
	comp := CompImpl{config.asInt, config.sortByMonth, config.humanSize}
	sorted := sort.SliceIsSorted(listOfList, lessByColumn(listOfList, config.column, comp, !config.reverse))
	return sorted
}

type Comparator interface {
	Less(o1, o2 any) bool
}

type CompImpl struct {
	asInt     bool
	byMonth   bool
	humanSize bool
}

func (ci CompImpl) Less(o1, o2 any) bool {
	if ci.humanSize {
		return compareHumanSize(o1.(string), o2.(string))
	}
	if ci.byMonth {
		return lessMonth(o1.(string), o2.(string))
	}
	if ci.asInt {
		return mustAtoi(o1.(string)) < mustAtoi(o2.(string))
	}
	return o1.(string) < o2.(string)
}

// Вспомогательная функция сравнения по размерам с суффиксами
func compareHumanSize(a, b string) bool {
	aSize := parseSize(a)
	bSize := parseSize(b)
	if aSize < bSize {
		return true
	} else if aSize > bSize {
		return false
	}
	return false
}

// Парсинг человекочитаемых размеров
func parseSize(s string) float64 {
	s = strings.TrimSpace(s)
	re := regexp.MustCompile(`(?i)^([\d.]+)\s*([KMGTPE]?B)?$`)
	matches := re.FindStringSubmatch(s)
	if len(matches) < 2 {
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0
		}
		return val
	}
	val, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0
	}
	switch strings.ToUpper(matches[2]) {
	case "KB":
		return val * 1024
	case "MB":
		return val * 1024 * 1024
	case "GB":
		return val * 1024 * 1024 * 1024
	case "TB":
		return val * 1024 * 1024 * 1024 * 1024
	case "PB":
		return val * 1024 * 1024 * 1024 * 1024 * 1024
	default:
		return val
	}
}

// сравнение месяцев
func lessMonth(a, b string) bool {
	months := map[string]int{
		"Jan": 1, "Feb": 2, "Mar": 3, "Apr": 4,
		"May": 5, "Jun": 6, "Jul": 7, "Aug": 8,
		"Sep": 9, "Oct": 10, "Nov": 11, "Dec": 12,
	}
	getMonthNum := func(s string) int {
		s = strings.TrimSpace(s)
		if val, ok := months[s]; ok {
			return val
		}
		return 0
	}
	aMonth := getMonthNum(a)
	bMonth := getMonthNum(b)
	if aMonth < bMonth {
		return true
	} else if aMonth > bMonth {
		return false
	}
	return false
}

// сравнение по колонке
func lessByColumn(listOfList [][]string,
	colNumber int,
	comparator Comparator,
	ascending bool) func(i, j int) bool {
	return func(i, j int) bool {
		if (len(listOfList[i]) <= colNumber) || (len(listOfList[j]) <= colNumber) {
			return false
		}
		less := comparator.Less(listOfList[i][colNumber], listOfList[j][colNumber])
		if ascending {
			return less
		}
		return !less
	}
}

// сортировка по колонке
func sortByColumn(list []string, config *Config) {
	listOfList := make([][]string, len(list))
	for i, str := range list {
		listOfList[i] = strings.Fields(str)
	}
	comp := CompImpl{config.asInt, config.sortByMonth, config.humanSize}

	sort.Slice(listOfList, lessByColumn(listOfList, config.column, comp, !config.reverse))

	for i, it := range listOfList {
		list[i] = strings.Join(it, " ")
	}
}
