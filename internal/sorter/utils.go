package sorter

import (
	"container/heap"
	"os"
	"path/filepath"
	"strings"
)

// minHeap для слияния отсортированных чанков
type minHeap struct {
	data     []string
	comparer func(a, b string) bool
}

func (h minHeap) Len() int           { return len(h.data) }
func (h minHeap) Less(i, j int) bool { return h.comparer(h.data[i], h.data[j]) }
func (h minHeap) Swap(i, j int)      { h.data[i], h.data[j] = h.data[j], h.data[i] }

func (h *minHeap) Push(x interface{}) {
	h.data = append(h.data, x.(string))
}

func (h *minHeap) Pop() interface{} {
	old := h.data
	n := len(old)
	x := old[n-1]
	h.data = old[0 : n-1]
	return x
}

// mergeSortedChunks сливает отсортированные чанки
func mergeSortedChunks(chunks [][]string, comparer func(a, b string) bool) []string {
	if len(chunks) == 0 {
		return nil
	}
	if len(chunks) == 1 {
		return chunks[0]
	}

	// Используем кучу для слияния
	h := &minHeap{
		comparer: comparer,
	}
	heap.Init(h)

	// Индексы текущих элементов в каждом чанке
	indices := make([]int, len(chunks))

	// Добавляем первые элементы каждого чанка в кучу
	for i, chunk := range chunks {
		if len(chunk) > 0 {
			heap.Push(h, chunk[0])
		}
	}

	var result []string

	// Процесс слияния
	for h.Len() > 0 {
		// Извлекаем минимальный элемент
		minElem := heap.Pop(h).(string)
		result = append(result, minElem)

		// Находим, из какого чанка был этот элемент
		for i, chunk := range chunks {
			if indices[i] < len(chunk) && chunk[indices[i]] == minElem {
				indices[i]++
				// Добавляем следующий элемент из этого чанка
				if indices[i] < len(chunk) {
					heap.Push(h, chunk[indices[i]])
				}
				break
			}
		}
	}

	return result
}

// splitIntoChunks разбивает большой файл на чанки
func splitIntoChunks(filePath string, chunkSize int) ([]string, error) {
	// Читаем файл
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var chunks []string

	for i := 0; i < len(lines); i += chunkSize {
		end := i + chunkSize
		if end > len(lines) {
			end = len(lines)
		}

		// Создаем временный файл для чанка
		tempFile, err := os.CreateTemp("", "chunk_*.txt")
		if err != nil {
			return nil, err
		}
		defer tempFile.Close()

		// Записываем чанк во временный файл
		chunk := lines[i:end]
		tempFile.WriteString(strings.Join(chunk, "\n"))
		chunks = append(chunks, tempFile.Name())
	}

	return chunks, nil
}

// cleanupTempFiles удаляет временные файлы
func cleanupTempFiles(filePaths []string) {
	for _, filePath := range filePaths {
		os.Remove(filePath)
	}
}
