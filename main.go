package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

// Config - входные параметры
type Config struct {
	column      int
	asInt       bool
	reverse     bool
	unique      bool
	sortByMonth bool
	trim        bool
	isSorted    bool
	humanSize   bool
}

type readResult struct {
	lines     []string
	fileNames []string
	failed    bool
}

func main() {
	// Определяем флаги
	config, files, err := initConfig(os.Args[1:])
	if err != nil {
		fmt.Println("Ошибка парсинга флагов:", err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	gnuLikeSort(files, scanner, config)
}

func initConfig(args []string) (*Config, []string, error) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	conf := &Config{}
	fs.IntVar(&conf.column, "k", 1, "сортировать по столбцу (колонке) №N (разделитель — табуляция по умолчанию)")
	fs.BoolVar(&conf.asInt, "n", false, "сортировать по числовому значению (строки интерпретируются как числа)")
	fs.BoolVar(&conf.reverse, "r", false, "сортировать в обратном порядке (reverse)")
	fs.BoolVar(&conf.unique, "u", false, "не выводить повторяющиеся строки (только уникальные)")
	fs.BoolVar(&conf.sortByMonth, "M", false, "сортировать по названию месяца (Jan, Feb, ... Dec), т.е. распознавать специфический формат дат")
	fs.BoolVar(&conf.trim, "b", false, "игнорировать хвостовые пробелы (trailing blanks)")
	fs.BoolVar(&conf.isSorted, "c", false, "проверить, отсортированы ли данные; если нет, вывести сообщение об этом")
	fs.BoolVar(&conf.humanSize, "h", false, "сортировать по числовому значению с учётом суффиксов "+
		"(например, К = килобайт, М = мегабайт — человекочитаемые размеры)")

	err := fs.Parse(args)
	conf.column -= 1
	if err != nil {
		return nil, nil, err
	}
	files := fs.Args()
	return conf, files, nil
}
