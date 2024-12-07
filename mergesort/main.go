package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	algorithm := flag.String("algorithm", "goroutines", "Sorting algorithm to use")
	listLength := flag.Int("listLength", 10000000, "Length of the list")
	//60000000
	maxDepth := flag.Int("maxDepth", 4, "Max tree depth")
	runs := flag.Int("runs", 1, "Number of runs")
	warmUpRuns := flag.Int("warmUpRuns", 0, "Number of warm-up runs")
	flag.Parse()

	fmt.Printf("Go - Algorithm: %s, List length: %d, Max Depth: %d, Runs: %d, War Up Runs: %d\n",
		*algorithm, *listLength, *maxDepth, *runs, *warmUpRuns)

	// Datei lesen
	list := importData(fmt.Sprintf("List%d.txt", *listLength))
	//fmt.Println("Unsortierte Zahlen:", list)

	fmt.Println("File imported.")

	for i := 0; i < *warmUpRuns; i++ {
		copyList := make([]int, len(list))
		copy(copyList, list)
		runAlgorithm(*algorithm, copyList, *maxDepth)
	}

	fmt.Println("warum up runs finished")

	start := time.Now()

	for i := 0; i < *runs; i++ {
		copyList := make([]int, len(list))
		copy(copyList, list)
		runAlgorithm(*algorithm, copyList, *maxDepth)
	}

	elapsed := time.Since(start)
	fmt.Printf("Go %s, Time: %s\n", *algorithm, elapsed)

}

func runAlgorithm(algorithm string, list []int, maxDepth int) {
	switch algorithm {
	case "single":
		RunMergeSortSingle(list)
	case "goroutines":
		RunMergeSortGoroutines(list, maxDepth)
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
	return numbers
}
