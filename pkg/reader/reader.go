package reader

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// ReadInput читает строки из файлов или STDIN
func ReadInput(filepaths []string) ([]string, error) {
	var lines []string

	if len(filepaths) == 0 {
		// Чтение из STDIN
		return readFromReader(os.Stdin)
	}

	// Чтение из файлов
	for _, filepath := range filepaths {
		file, err := os.Open(filepath)
		if err != nil {
			return nil, fmt.Errorf("ошибка открытия файла %s: %w", filepath, err)
		}
		defer file.Close()

		fileLines, err := readFromReader(file)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения файла %s: %w", filepath, err)
		}

		lines = append(lines, fileLines...)
	}

	return lines, nil
}

// readFromReader читает строки из io.Reader
func readFromReader(reader io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(reader)

	// Установка размера буфера для больших файлов
	const maxCapacity = 1024 * 1024 // 1MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка сканирования: %w", err)
	}

	return lines, nil
}
