package flags

import (
	"flag"
	"fmt"
	"os"
)

// Options содержит все параметры сортировки
type Options struct {
	KeyColumn      int    // -k N
	NumericSort    bool   // -n
	ReverseSort    bool   // -r
	Unique         bool   // -u
	MonthSort      bool   // -M
	IgnoreBlanks   bool   // -b
	CheckSorted    bool   // -c
	HumanNumeric   bool   // -h
	FieldSeparator string // -t (дополнительно)
}

// ParseFlags парсит аргументы командной строки
func ParseFlags() (*Options, []string, error) {
	var opts Options

	// Определение флагов
	flag.IntVar(&opts.KeyColumn, "k", 0, "сортировать по колонке N (0 - вся строка)")
	flag.BoolVar(&opts.NumericSort, "n", false, "сортировать по числовому значению")
	flag.BoolVar(&opts.ReverseSort, "r", false, "сортировать в обратном порядке")
	flag.BoolVar(&opts.Unique, "u", false, "выводить только уникальные строки")
	flag.BoolVar(&opts.MonthSort, "M", false, "сортировать по названию месяца")
	flag.BoolVar(&opts.IgnoreBlanks, "b", false, "игнорировать хвостовые пробелы")
	flag.BoolVar(&opts.CheckSorted, "c", false, "проверить, отсортированы ли данные")
	flag.BoolVar(&opts.HumanNumeric, "h", false, "сортировать с учётом суффиксов (K,M,G)")
	separator := flag.String("t", "\t", "разделитель полей")

	// Кастомное использование
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Использование: %s [ОПЦИИ]... [ФАЙЛ]...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Сортирует строки из всех ФАЙЛОВ или стандартного ввода.\n\n")
		fmt.Fprintf(os.Stderr, "Обязательные опции:\n")
		fmt.Fprintf(os.Stderr, "  -k N    сортировать по колонке N (разделитель - табуляция)\n")
		fmt.Fprintf(os.Stderr, "  -n      сортировать по числовому значению\n")
		fmt.Fprintf(os.Stderr, "  -r      сортировать в обратном порядке\n")
		fmt.Fprintf(os.Stderr, "  -u      выводить только уникальные строки\n")
		fmt.Fprintf(os.Stderr, "\nДополнительные опции:\n")
		fmt.Fprintf(os.Stderr, "  -M      сортировать по названию месяца\n")
		fmt.Fprintf(os.Stderr, "  -b      игнорировать хвостовые пробелы\n")
		fmt.Fprintf(os.Stderr, "  -c      проверить, отсортированы ли данные\n")
		fmt.Fprintf(os.Stderr, "  -h      сортировать с учётом суффиксов (K,M,G,T,P,E)\n")
		fmt.Fprintf(os.Stderr, "  -t SEP  использовать SEP как разделитель полей\n")
		fmt.Fprintf(os.Stderr, "\nПримеры:\n")
		fmt.Fprintf(os.Stderr, "  %s -k 2 file.txt      Сортировать по второму столбцу\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -nr file.txt       Числовая сортировка в обратном порядке\n", os.Args[0])
	}

	flag.Parse()

	// Установка разделителя
	opts.FieldSeparator = *separator
	if len(opts.FieldSeparator) > 1 {
		return nil, nil, fmt.Errorf("разделитель должен быть одним символом")
	}

	// Валидация комбинаций флагов
	if opts.NumericSort && opts.MonthSort {
		return nil, nil, fmt.Errorf("флаги -n и -M не могут использоваться вместе")
	}
	if opts.NumericSort && opts.HumanNumeric {
		return nil, nil, fmt.Errorf("флаги -n и -h не могут использоваться вместе")
	}
	if opts.MonthSort && opts.HumanNumeric {
		return nil, nil, fmt.Errorf("флаги -M и -h не могут использоваться вместе")
	}
	if opts.KeyColumn < 0 {
		return nil, nil, fmt.Errorf("номер колонки должен быть положительным числом")
	}

	return &opts, flag.Args(), nil
}
