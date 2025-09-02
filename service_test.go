package main

import (
	"bufio"
	"bytes"
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
	rand.Seed(time.Now().UnixNano())

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

func readStdout(f func()) string {
	originalStdout := os.Stdout

	var buf bytes.Buffer

	r, w, _ := os.Pipe()

	os.Stdout = w

	defer func() { os.Stdout = originalStdout }()

	done := make(chan struct{})
	go func() {
		var bufBytes bytes.Buffer
		bufBytes.ReadFrom(r)
		buf.Write(bufBytes.Bytes())
		close(done)
	}()
	f()
	w.Close()

	<-done
	return buf.String()
}
