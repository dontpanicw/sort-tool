package sorter

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dontpanicw/sort-tool/pkg/flags"
)

// Comparator сравнивает строки согласно настройкам
type Comparator struct {
	opts *flags.Options
}

// NewComparator создает новый компаратор
func NewComparator(opts *flags.Options) (*Comparator, error) {
	return &Comparator{
		opts: opts,
	}, nil
}

// Compare сравнивает две строки
func (c *Comparator) Compare(a, b string) (bool, error) {
	// Извлечение колонки для сравнения
	keyA := c.extractKey(a)
	keyB := c.extractKey(b)

	// В зависимости от типа сортировки
	switch {
	case c.opts.NumericSort:
		return c.compareNumeric(keyA, keyB)
	case c.opts.MonthSort:
		return c.compareMonth(keyA, keyB)
	case c.opts.HumanNumeric:
		return c.compareHumanNumeric(keyA, keyB)
	default:
		// Лексикографическое сравнение
		return keyA < keyB, nil
	}
}

// extractKey извлекает ключ для сравнения из строки
func (c *Comparator) extractKey(line string) string {
	if c.opts.KeyColumn == 0 {
		return line
	}

	fields := strings.Split(line, c.opts.FieldSeparator)
	columnIndex := c.opts.KeyColumn - 1 // -k 1 означает первую колонку

	if columnIndex < 0 || columnIndex >= len(fields) {
		return ""
	}

	return strings.TrimSpace(fields[columnIndex])
}

// compareNumeric сравнивает строки как числа
func (c *Comparator) compareNumeric(a, b string) (bool, error) {
	// Удаляем пробелы
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)

	// Пытаемся преобразовать в числа
	numA, errA := strconv.ParseFloat(a, 64)
	numB, errB := strconv.ParseFloat(b, 64)

	// Если обе строки - числа
	if errA == nil && errB == nil {
		return numA < numB, nil
	}

	// Если только одна строка - число, числа идут первыми
	if errA == nil && errB != nil {
		return true, nil
	}
	if errA != nil && errB == nil {
		return false, nil
	}

	// Если ни одна не число, сравниваем как строки
	return a < b, nil
}

// compareMonth сравнивает строки как названия месяцев
func (c *Comparator) compareMonth(a, b string) (bool, error) {
	monthA, errA := parseMonth(a)
	monthB, errB := parseMonth(b)

	// Если оба - месяцы
	if errA == nil && errB == nil {
		return monthA < monthB, nil
	}

	// Если только один месяц
	if errA == nil && errB != nil {
		return true, nil
	}
	if errA != nil && errB == nil {
		return false, nil
	}

	// Если ни один не месяц, сравниваем как строки
	return a < b, nil
}

// parseMonth парсит название месяца
func parseMonth(s string) (time.Month, error) {
	s = strings.TrimSpace(strings.ToLower(s))

	months := map[string]time.Month{
		"jan": time.January, "january": time.January,
		"feb": time.February, "february": time.February,
		"mar": time.March, "march": time.March,
		"apr": time.April, "april": time.April,
		"may": time.May,
		"jun": time.June, "june": time.June,
		"jul": time.July, "july": time.July,
		"aug": time.August, "august": time.August,
		"sep": time.September, "september": time.September,
		"oct": time.October, "october": time.October,
		"nov": time.November, "november": time.November,
		"dec": time.December, "december": time.December,
	}

	if month, ok := months[s]; ok {
		return month, nil
	}

	return 0, fmt.Errorf("неизвестный месяц: %s", s)
}

// compareHumanNumeric сравнивает строки с суффиксами (K, M, G и т.д.)
func (c *Comparator) compareHumanNumeric(a, b string) (bool, error) {
	valueA, errA := parseHumanNumber(a)
	valueB, errB := parseHumanNumber(b)

	// Если оба - числа с суффиксами
	if errA == nil && errB == nil {
		return valueA < valueB, nil
	}

	// Если только одно число
	if errA == nil && errB != nil {
		return true, nil
	}
	if errA != nil && errB == nil {
		return false, nil
	}

	// Если ни одно не число, сравниваем как строки
	return a < b, nil
}

// parseHumanNumber парсит числа с суффиксами
func parseHumanNumber(s string) (float64, error) {
	s = strings.TrimSpace(strings.ToLower(s))

	// Регулярное выражение для поиска чисел с суффиксами
	re := regexp.MustCompile(`^([0-9]+(?:\.[0-9]+)?)\s*([kmgtpe])?[i]?b?$`)
	matches := re.FindStringSubmatch(s)

	if matches == nil {
		return 0, fmt.Errorf("не является числом с суффиксом: %s", s)
	}

	// Преобразуем число
	num, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, err
	}

	// Умножаем на множитель в зависимости от суффикса
	if len(matches) > 2 && matches[2] != "" {
		multipliers := map[string]float64{
			"k": 1024,
			"m": 1024 * 1024,
			"g": 1024 * 1024 * 1024,
			"t": 1024 * 1024 * 1024 * 1024,
			"p": 1024 * 1024 * 1024 * 1024 * 1024,
			"e": 1024 * 1024 * 1024 * 1024 * 1024 * 1024,
		}
		if mult, ok := multipliers[matches[2]]; ok {
			num *= mult
		}
	}

	return num, nil
}
