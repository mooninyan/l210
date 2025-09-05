package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// readFromFiles читает строки из файлов
func readFromFiles(files []string, cfg *Config) (lines []string, f []string, err error) {
	var totalBytes int64 = 0

	for _, fileName := range files {
		var fileLines []string
		fileLines, err = readFromFile(fileName, cfg)
		if err != nil {
			return nil, f, err
		}

		// Подсчет размера новых строк
		var sizeOfNewLines int64 = 0
		for _, line := range fileLines {
			sizeOfNewLines += int64(len(line))
		}

		// Проверка, не превышает ли добавление новых строк лимит
		if totalBytes+sizeOfNewLines > MaxMemoryBytes {
			return nil, files, nil
		}

		lines = append(lines, fileLines...)
		totalBytes += sizeOfNewLines
	}
	return
}

// readFromFile читает строки из файла
func readFromFile(fileName string, cfg *Config) (lines []string, err error) {
	// Открываем файл для чтения
	var file *os.File

	file, err = os.Open(fileName)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	// Создаём сканер для построчного чтения
	scanner := bufio.NewScanner(file)

	// Читаем строки по одной
	for scanner.Scan() {
		line := scanner.Text()

		if cfg.trim {
			line = strings.Trim(line, "\t")
		}
		lines = append(lines, line)
		fmt.Println("Прочитана строка:", line)
	}

	// Проверка ошибок при чтении
	if err = scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}
	return
}

func readLinesWithLimitCheck(scanner *bufio.Scanner, cfg *Config) ([]string, *string, error) {
	var totalBytes int64 = 0
	tempFile, err := os.CreateTemp("", "lines")
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(tempFile)

	list := make([]string, 0)

	exceeded := false // флаг, что лимит уже превышен

	for scanner.Scan() {
		line := scanner.Text()
		if cfg.trim {
			line = strings.Trim(line, "\t")
		}
		lineBytes := int64(len(line))
		if !exceeded && totalBytes+lineBytes > MaxMemoryBytes {
			// Лимит превышен, переносим все сохраненные строки в файл
			for _, s := range list {
				_, err = writer.WriteString(s + "\n")
				if err != nil {
					panic(err)
				}
			}
			list = nil
			exceeded = true
		}

		if exceeded {
			// Уже превышен лимит, пишем прямо в файл
			_, err = writer.WriteString(line + "\n")
			if err != nil {
				panic(err)
			}
		} else {
			// В памяти
			list = append(list, line)
			totalBytes += lineBytes
		}
	}

	err = writer.Flush()
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	err = tempFile.Close()
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	tempFileAddr := tempFile.Name()
	if err = scanner.Err(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Ошибка при чтении:", err)
		return nil, &tempFileAddr, err
	}
	if len(list) > 0 {
		return list, nil, nil
	}
	return nil, &tempFileAddr, nil
}
