package sorter

import (
	"container/heap"
	"fmt"
	"sort"
	"strings"

	"github.com/dontpanicw/sort-tool/pkg/flags"
)

// Sorter отвечает за сортировку строк
type Sorter struct {
	opts *flags.Options
	comp *Comparator
}

// NewSorter создает новый сортировщик
func NewSorter(opts *flags.Options) (*Sorter, error) {
	comp, err := NewComparator(opts)
	if err != nil {
		return nil, err
	}

	return &Sorter{
		opts: opts,
		comp: comp,
	}, nil
}

// Sort сортирует строки согласно настройкам
func (s *Sorter) Sort(lines []string) ([]string, error) {
	// Обработка флага -b (игнорирование хвостовых пробелов)
	if s.opts.IgnoreBlanks {
		for i := range lines {
			lines[i] = strings.TrimRight(lines[i], " \t")
		}
	}

	// Создание копии для сортировки
	sortedLines := make([]string, len(lines))
	copy(sortedLines, lines)

	// Сортировка
	sort.SliceStable(sortedLines, func(i, j int) bool {
		less, err := s.comp.Compare(sortedLines[i], sortedLines[j])
		if err != nil {
			// В случае ошибки сравниваем как строки
			return sortedLines[i] < sortedLines[j]
		}
		return less
	})

	// Обратный порядок (флаг -r)
	if s.opts.ReverseSort {
		reverse(sortedLines)
	}

	// Уникальные строки (флаг -u)
	if s.opts.Unique {
		sortedLines = unique(sortedLines)
	}

	return sortedLines, nil
}

// IsSorted проверяет, отсортированы ли строки
func (s *Sorter) IsSorted(lines []string) (bool, error) {
	for i := 1; i < len(lines); i++ {
		less, err := s.comp.Compare(lines[i-1], lines[i])
		if err != nil {
			return false, err
		}
		if !less {
			return false, nil
		}
	}
	return true, nil
}

// reverse переворачивает слайс
func reverse(lines []string) {
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}
}

// unique возвращает только уникальные строки
func unique(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}

	result := make([]string, 0, len(lines))
	result = append(result, lines[0])

	for i := 1; i < len(lines); i++ {
		if lines[i] != lines[i-1] {
			result = append(result, lines[i])
		}
	}

	return result
}

// ExternalMergeSort для больших файлов (опционально)
type externalMergeSort struct {
	chunkSize int
	opts      *flags.Options
}

// NewExternalMergeSort создает сортировщик для больших файлов
func NewExternalMergeSort(opts *flags.Options, chunkSize int) *externalMergeSort {
	return &externalMergeSort{
		chunkSize: chunkSize,
		opts:      opts,
	}
}

// Sort сортирует большие файлы с использованием внешней сортировки
func (ems *externalMergeSort) Sort(lines []string) ([]string, error) {
	if len(lines) <= ems.chunkSize {
		// Если данные помещаются в память, сортируем обычным способом
		sorter, err := NewSorter(ems.opts)
		if err != nil {
			return nil, err
		}
		return sorter.Sort(lines)
	}

	// Для очень больших файлов нужно реализовать внешнюю сортировку
	// с разбиением на чанки и их последующим слиянием
	// Здесь упрощенная реализация для демонстрации

	// Просто сортируем обычным способом (для демо)
	sorter, err := NewSorter(ems.opts)
	if err != nil {
		return nil, err
	}
	return sorter.Sort(lines)
}
