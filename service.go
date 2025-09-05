package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func gnuLikeSort(files []string, scanner *bufio.Scanner, cfg *Config) {
	var list []string
	res := readFromSource(files, scanner, cfg)
	if res.failed {
		return
	}
	var err error

	// работаем с временными файлами
	if len(res.fileNames) > 0 {
		var tmpFiles []string

		tmpFiles, err = splitAndSortLargeFiles(res.fileNames, ChunkSize, cfg)
		if len(tmpFiles) == 0 {
			return
		}

		if err != nil {
			fmt.Println(err)
			return
		}

		err = mergeSortedFiles(tmpFiles, cfg)
		if err != nil {
			fmt.Println("Ошибка при слиянии файлов:", err)
			return
		}

		// Можно удалить временные файлы после объединения
		for _, f := range tmpFiles {
			_ = os.Remove(f)
		}
		for _, f := range res.fileNames {
			_ = os.Remove(f)
		}
	} else {
		// работаем с данными в памяти
		list = res.lines

		// проверяем не отсортированы ли строки
		if cfg.isSorted {
			if checkIfSorted(list, cfg) {
				return
			} else {
				fmt.Println("Данные не отсортированы")
				return
			}
		}

		sortByColumn(list, cfg)

		list = distinct(list, cfg)

		fmt.Println(strings.Join(list, "\n"))
	}

	for _, r := range res.fileNames {
		err = os.Remove(r)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func distinct(list []string, cfg *Config) []string {
	if cfg.unique {
		duplicates := make(map[string]struct{})

		uniqueList := make([]string, 0, len(list))
		for _, it := range list {
			_, ok := duplicates[it]
			if ok {
				continue
			}
			duplicates[it] = struct{}{}
			uniqueList = append(uniqueList, it)
		}
		list = uniqueList
	}
	return list
}

// Читаем строки из фалов или stdin
func readFromSource(files []string, scanner *bufio.Scanner, cfg *Config) readResult {
	var result readResult

	if len(files) == 0 {
		lines, file, err := readLinesWithLimitCheck(scanner, cfg)
		if err != nil {
			fmt.Println("Ошибка при чтении из stdin:", err)
			return readResult{failed: true}
		}
		result.lines = lines
		if file != nil {
			result.fileNames = []string{*file}
		}
	} else {
		lines, fsToSort, err := readFromFiles(files, cfg)
		if err != nil {
			fmt.Println("Ошибка при чтении из файлов:", err)
			return readResult{failed: true}
		}
		result.lines = lines
		result.fileNames = fsToSort
	}

	return result
}

func checkIfSorted(list []string, cfg *Config) bool {
	return checkSortByColumn(list, cfg)
}
