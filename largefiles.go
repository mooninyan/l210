package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

func splitAndSortLargeFiles(inputFiles []string, chunkSize int, cfg *Config) ([]string, error) {
	var files []string
	for _, fileName := range inputFiles {
		f, err := splitAndSortLargeFile(fileName, chunkSize, cfg)
		if err != nil {
			return nil, err
		}
		files = append(files, f...)
	}
	return files, nil
}

// Разделение файла на части и сортировка
func splitAndSortLargeFile(inputFile string, chunkSize int, cfg *Config) ([]string, error) {
	inFile, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer func(inFile *os.File) {
		err = inFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(inFile)

	scanner := bufio.NewScanner(inFile)
	var chunk []string
	var tempFiles []string
	var outFile *os.File
	partIdx := 0
	err = nil // объявляем переменную err один раз

	for {
		chunk = chunk[:0]
		for len(chunk) < chunkSize && scanner.Scan() {
			line := scanner.Text()

			chunk = append(chunk, line)
		}

		if err = scanner.Err(); err != nil {
			return nil, err
		}

		if len(chunk) == 0 {
			break
		}

		if cfg.isSorted {
			if checkIfSorted(chunk, cfg) {
				return nil, nil
			} else {
				fmt.Println("Данные не отсортированы")
				return nil, nil
			}
		}
		// сортируем текущий кусок
		sortByColumn(chunk, cfg)

		// записываем в временный файл
		tempFileName := fmt.Sprintf("chunk_%d.txt", partIdx)
		outFile, err = os.Create(tempFileName)
		if err != nil {
			return nil, err
		}

		for _, line := range chunk {
			_, err = outFile.WriteString(line + "\n")
			if err != nil {
				err = outFile.Close()
				if err != nil {
					fmt.Println(err)
					return nil, err
				}
				return nil, err
			}
		}
		err = outFile.Close()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		tempFiles = append(tempFiles, tempFileName)
		partIdx++
	}

	return tempFiles, nil
}

// K-way merge отсортированных файлов
func mergeSortedFiles(files []string, cfg *Config) error {
	var readers []*bufio.Scanner
	var fileHandles []*os.File

	for _, fname := range files {
		f, err := os.Open(fname)
		if err != nil {
			return err
		}
		fileHandles = append(fileHandles, f)
		scanners := bufio.NewScanner(f)
		if scanners.Scan() {
			readers = append(readers, scanners)
		} else {
			// файл пустой, закрываем
			err = f.Close()
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	// Создаем приоритетную очередь
	impl := CompImpl{asInt: cfg.asInt, byMonth: cfg.sortByMonth, humanSize: cfg.humanSize, reverse: cfg.reverse}
	pq := &ItemHeap{Comparator: impl}
	heap.Init(pq)

	// Инициализация очереди начальными элементами
	for i, scanner := range readers {
		if scanner != nil && scanner.Text() != "" {
			heap.Push(pq, &Item{
				value:   scanner.Text(),
				fileIdx: i,
			})
		}
	}

	isFirst := true
	lastElem := ""
	for pq.Len() > 0 {
		minItem := heap.Pop(pq).(*Item)

		if cfg.unique {
			if !isFirst && lastElem == minItem.value {
			} else {
				fmt.Println(minItem.value)
				lastElem = minItem.value
				isFirst = false
			}
		} else {
			fmt.Println(minItem.value)
		}
		fIdx := minItem.fileIdx
		if readers[fIdx].Scan() {
			heap.Push(pq, &Item{
				value:   readers[fIdx].Text(),
				fileIdx: fIdx,
			})
		}
	}

	// закрываем файлы
	for _, f := range fileHandles {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

// Item структура для heap
type Item struct {
	value   string
	fileIdx int
}

type ItemHeap struct {
	items []*Item
	Comparator
}

func (h *ItemHeap) Len() int           { return len(h.items) }
func (h *ItemHeap) Less(i, j int) bool { return h.Comparator.Less(h.items[i].value, h.items[j].value) }
func (h *ItemHeap) Swap(i, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }
func (h *ItemHeap) Push(x interface{}) {
	(*h).items = append((*h).items, x.(*Item))
}
func (h *ItemHeap) Pop() interface{} {
	old := *h
	n := len(old.items)
	item := old.items[n-1]
	h.items = old.items[:n-1]
	return item
}
