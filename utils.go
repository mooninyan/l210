package main

import (
	"fmt"
	"strconv"
)

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("не удалось преобразовать '%s' в число: %v", s, err))
	}
	return i
}
