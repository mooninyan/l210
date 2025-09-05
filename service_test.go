package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestSortWithSmallChunks(t *testing.T) {
	ChunkSizeOld := ChunkSize
	MaxMemoryBytesOld := MaxMemoryBytes

	ChunkSize = 4
	MaxMemoryBytes = 10
	// Массив с буквами английского алфавита
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	expected := strings.Join(strings.Split(string(letters), ""), "\n") + "\n"
	// Инициализация генератора случайных чисел
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Перемешивание слайса
	rand.Shuffle(len(letters), func(i, j int) {
		letters[i], letters[j] = letters[j], letters[i]
	})

	randomized := strings.Join(strings.Split(string(letters), ""), "\n")

	scanner := bufio.NewScanner(strings.NewReader(randomized))

	actual := readStdout(func() {
		gnuLikeSort([]string{}, scanner, &Config{})
	})

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("По первому столбцу: ожидаемый: \n%#v, \nполученный: \n%#v",
			expected,
			actual)
	}

	ChunkSize = ChunkSizeOld
	MaxMemoryBytes = MaxMemoryBytesOld
}

func TestSortWithSmallChunksReverse(t *testing.T) {
	ChunkSizeOld := ChunkSize
	MaxMemoryBytesOld := MaxMemoryBytes

	ChunkSize = 4
	MaxMemoryBytes = 10
	// Массив с буквами английского алфавита
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	expected := reverseString(string(letters))
	expected = strings.Join(strings.Split(expected, ""), "\n") + "\n"
	// Инициализация генератора случайных чисел
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Перемешивание слайса
	rand.Shuffle(len(letters), func(i, j int) {
		letters[i], letters[j] = letters[j], letters[i]
	})

	randomized := strings.Join(strings.Split(string(letters), ""), "\n")

	scanner := bufio.NewScanner(strings.NewReader(randomized))

	actual := readStdout(func() {
		gnuLikeSort([]string{}, scanner, &Config{reverse: true})
	})

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("По первому столбцу: ожидаемый: \n%#v, \nполученный: \n%#v",
			expected,
			actual)
	}

	ChunkSize = ChunkSizeOld
	MaxMemoryBytes = MaxMemoryBytesOld
}

func reverseString(s string) string {
	runes := []rune(s) // преобразуем строку в срез руны
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i] // меняем местами
	}
	return string(runes) // возвращаем строку из рун
}

func readStdout(f func()) string {
	originalStdout := os.Stdout

	var buf bytes.Buffer

	r, w, _ := os.Pipe()

	os.Stdout = w

	defer func() { os.Stdout = originalStdout }()

	done := make(chan struct{})
	go func() {
		var bufBytes bytes.Buffer
		_, err := bufBytes.ReadFrom(r)
		if err != nil {
			fmt.Println(err)
		}
		buf.Write(bufBytes.Bytes())
		close(done)
	}()
	f()
	err := w.Close()
	if err != nil {
		fmt.Println(err)
	}

	<-done
	return buf.String()
}
