package main

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// Для Windows
		cmd = exec.Command("cmd", "/C", "go build -o build/l210.exe")
	} else {
		// Для Unix-подобных ОС (Linux, macOS)
		cmd = exec.Command("sh", "-c", "go build -o build/l210")
	}

	// Перенаправляем вывод
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Если сборка не удалась
		os.Exit(1)
	}

	// Запуск остальных тестов
	code := m.Run()
	os.Exit(code)
}

func TestParseFlags(t *testing.T) {
	args := []string{"-k=2", "-n", "-r", "-u"}
	config, files, err := initConfig(args)
	if err != nil {
		t.Fatalf("parseFlags error: %v", err)
	}

	if config.column != 1 {
		t.Errorf("expected column=2, got %d", config.column)
	}
	if !config.asInt {
		t.Errorf("expected asInt=true")
	}
	if !config.reverse {
		t.Errorf("expected reverse=true")
	}
	if !config.unique {
		t.Errorf("expected unique=true")
	}
	if len(files) != 0 {
		t.Errorf("expected no files, got %v", files)
	}
}

func TestSortProgramFromBinary(t *testing.T) {
	binaryPath := "./build/l210"

	t.Run("Basic gnuLikeSort by first column", func(t *testing.T) {
		input := "banana 2\napple 1\ncherry 3\n"
		cmd := exec.Command(binaryPath)
		cmd.Stdin = strings.NewReader(input)

		outputBytes, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Error running program: %v", err)
		}

		output := string(outputBytes)
		expected := "apple 1\nbanana 2\ncherry 3\n"
		if output != expected {
			t.Errorf("Expected:\n%sGot:\n%s", expected, output)
		}
	})

	t.Run("Sort numerically by second column", func(t *testing.T) {
		input := "item1 20\nitem2 3\nitem3 100\n"
		cmd := exec.Command(binaryPath, "-k", "2", "-n")
		cmd.Stdin = strings.NewReader(input)

		outputBytes, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Error running program: %v", err)
		}

		output := string(outputBytes)
		expected := "item2 3\n" +
			"item1 20\n" +
			"item3 100\n"
		if output != expected {
			t.Errorf("Expected:\n%sGot:\n%s", expected, output)
		}
	})

	t.Run("Check if input is sorted (-c)", func(t *testing.T) {
		input := "a\nb\nc\n"
		cmd := exec.Command(binaryPath, "-c")
		cmd.Stdin = strings.NewReader(input)

		outputBytes, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Error running program: %v", err)
		}

		output := string(outputBytes)
		if strings.TrimSpace(output) != "" {
			t.Errorf("Expected empty output, got: %s", output)
		}

		// Проверка несортированных данных
		input2 := "b\na\n"
		cmd2 := exec.Command(binaryPath, "-c")
		cmd2.Stdin = strings.NewReader(input2)

		outputBytes, _ = cmd2.CombinedOutput()
		output = string(outputBytes)
		expected := "Данные не отсортированы\n"
		if output != expected {
			t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
		}
	})
}
