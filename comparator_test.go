package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestSortByColumn(t *testing.T) {
	// Исходный список
	list := []string{
		"apple red",
		"banana yellow",
		"kiwi green",
		"strawberry red",
		"orange orange",
		"pear",
	}

	// Тестируем сортировку по первому столбцу (индекс 0)
	expectedByFirstColumn := []string{
		"apple red",
		"banana yellow",
		"kiwi green",
		"orange orange",
		"pear",
		"strawberry red",
	}
	sortByColumn(list, &Config{})
	if !reflect.DeepEqual(list, expectedByFirstColumn) {
		t.Errorf("По первому столбцу: ожидаемый %v, полученный %v",
			strings.Join(expectedByFirstColumn, ", "),
			strings.Join(list, ", "))
	}

	// Восстановим список для следующего теста
	list = []string{
		"banana yellow",
		"apple red",
		"kiwi green",
		"strawberry red",
		"orange orange",
		"pear",
	}

	// Тестируем сортировку по второму столбцу (индекс 1)
	expectedBySecondColumn := []string{
		"kiwi green",
		"orange orange",
		"apple red",
		"strawberry red",
		"banana yellow",
		"pear",
	}
	sortByColumn(list, &Config{column: 1})
	if !reflect.DeepEqual(list, expectedBySecondColumn) {
		t.Errorf("\nПо второму столбцу: \nожидаемый %v, \nполученный %v",
			strings.Join(expectedBySecondColumn, ", "),
			strings.Join(list, ", "))
	}

	// Восстановим список для следующего теста
	list = []string{
		"apple red",
		"banana yellow",
		"kiwi green",
		"strawberry red",
		"orange orange",
		"pear",
	}

	// Тестируем сортировку по первому столбцу (индекс 0) с флагом reverse
	expectedByFirstColumnWithReverse := []string{
		"strawberry red",
		"pear",
		"orange orange",
		"kiwi green",
		"banana yellow",
		"apple red",
	}
	sortByColumn(list, &Config{reverse: true})
	if !reflect.DeepEqual(list, expectedByFirstColumnWithReverse) {
		t.Errorf("По первому столбцу: ожидаемый %v, полученный %v",
			strings.Join(expectedByFirstColumnWithReverse, ", "),
			strings.Join(list, ", "))
	}

	// Восстановим список для следующего теста
	list = []string{
		"22",
		"41",
		"3",
		"7",
		"53",
		"10",
	}

	// Тестируем сортировку по первому столбцу (индекс 0) по int значению
	expectedByInt := []string{
		"3",
		"7",
		"10",
		"22",
		"41",
		"53",
	}
	sortByColumn(list, &Config{asInt: true})
	if !reflect.DeepEqual(list, expectedByInt) {
		t.Errorf("По первому столбцу: ожидаемый %v, полученный %v",
			strings.Join(expectedByFirstColumnWithReverse, ", "),
			strings.Join(list, ", "))
	}

	// Восстановим список для следующего теста
	list = []string{
		"apple red",
		"banana yellow",
		"kiwi green",
		"strawberry red",
		"orange orange",
		"pear",
	}

}

func TestCheckSort(t *testing.T) {
	// Тестируем проверку на сортировку
	list := []string{
		"a",
		"b",
	}
	sorted := checkSortByColumn(list, &Config{})
	if true != sorted {
		t.Errorf("ожидалось: %v. получено: %v",
			true,
			sorted)
	}

	// Тестируем проверку на сортировку
	list = []string{
		"b",
		"a",
	}
	sorted = checkSortByColumn(list, &Config{})
	if true == sorted {
		t.Errorf("ожидалось: %v. получено: %v",
			true,
			sorted)
	}
}
