package sort

import (
	"fmt"
	"os"

	"github.com/dontpanicw/sort-tool/internal/sorter"
	"github.com/dontpanicw/sort-tool/pkg/flags"
	"github.com/dontpanicw/sort-tool/pkg/reader"
)

func main() {
	// Парсинг флагов
	opts, filepaths, err := flags.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка парсинга флагов: %v\n", err)
		os.Exit(1)
	}

	// Чтение данных
	lines, err := reader.ReadInput(filepaths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения: %v\n", err)
		os.Exit(1)
	}

	// Создание сортировщика
	sorter, err := sorter.NewSorter(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка создания сортировщика: %v\n", err)
		os.Exit(1)
	}

	// Проверка сортировки (флаг -c)
	if opts.CheckSorted {
		isSorted, err := sorter.IsSorted(lines)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка проверки сортировки: %v\n", err)
			os.Exit(1)
		}
		if isSorted {
			os.Exit(0)
		} else {
			fmt.Fprintln(os.Stderr, "Данные не отсортированы")
			os.Exit(1)
		}
	}

	// Сортировка
	sortedLines, err := sorter.Sort(lines)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка сортировки: %v\n", err)
		os.Exit(1)
	}

	// Вывод результата
	for _, line := range sortedLines {
		fmt.Println(line)
	}
}
