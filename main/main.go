package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"performancetest/mergesort"
)

func main() {

	algorithm := flag.String("algorithm", "goroutine", "Sorting algorithm to use")
	listLength := flag.Int("listLength", 10, "Length of the list")
	//60000000
	chunkNumber := flag.Int("chunkNumber", 16, "Number of chunks")
	runs := flag.Int("runs", 1, "Number of runs")
	flag.Parse()

	chunkSize := *listLength / *chunkNumber

	fmt.Printf("Go - Algorithm: %s, List length: %d, Chunk number: %d, Runs: %d\n",
		*algorithm, *listLength, *chunkNumber, *runs)

	// Datei lesen
	list := importData(fmt.Sprintf("List%d.txt", *listLength))
	//fmt.Println("Unsortierte Zahlen:", list)

	for i := 0; i < *runs; i++ {
		copyList := make([]int, len(list))
		copy(copyList, list)
		runAlgorithm(*algorithm, copyList, chunkSize)
	}
}

func runAlgorithm(algorithm string, list []int, chunkSize int) {
	switch algorithm {
	case "single":
		list = mergesort.MergeSort(list)
	case "goroutine":
		sortChan := mergesort.MergeSortGoroutine(list, chunkSize)
		list = <-sortChan
	default:
		fmt.Println("Unknown algorithm")
	}

	//fmt.Println("Sortierte Zahlen:", list)
}

func importData(fileName string) []int {
	// Datei lesen
	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Fehler beim Lesen der Datei:", err)
		return nil
	}

	// Inhalt in einen String umwandeln und am Komma aufteilen
	numberStrings := strings.Split(strings.TrimSpace(string(content)), ",")

	// Slice für die Zahlen erstellen
	numbers := make([]int, 0, len(numberStrings))

	// Strings in Integers umwandeln und zum Slice hinzufügen
	for _, numStr := range numberStrings {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			fmt.Printf("Fehler beim Konvertieren von '%s': %v\n", numStr, err)
			continue
		}
		numbers = append(numbers, num)
	}
	fmt.Println("File imported.")
	return numbers
}
