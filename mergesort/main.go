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

	// Retrieve the algorithm, list length, max depth, runs and warm up runs from the command line arguments
	algorithm := flag.String("algorithm", "goroutines", "Sorting algorithm to use")
	listLength := flag.Int("listLength", 10000000, "Length of the list")
	maxDepth := flag.Int("maxDepth", 4, "Max tree depth")
	runs := flag.Int("runs", 1, "Number of runs")
	warmUpRuns := flag.Int("warmUpRuns", 0, "Number of warm-up runs")
	flag.Parse()

	fmt.Printf("Go - Algorithm: %s, List length: %d, Max Depth: %d, Runs: %d, War Up Runs: %d\n",
		*algorithm, *listLength, *maxDepth, *runs, *warmUpRuns)

	// Import data
	list := importData(fmt.Sprintf("List%d.txt", *listLength))

	fmt.Println("File imported.")

	// Warm up runs
	for i := 0; i < *warmUpRuns; i++ {
		copyList := make([]int, len(list))
		copy(copyList, list)
		runAlgorithm(*algorithm, copyList, *maxDepth)
	}

	fmt.Println("warum up runs finished")

	// Runs the algorithm
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
}

func importData(fileName string) []int {
	// Read the content of the file
	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Fehler beim Lesen der Datei:", err)
		return nil
	}

	// Convert the content to a slice of strings
	numberStrings := strings.Split(strings.TrimSpace(string(content)), ",")

	// Slice to store the numbers
	numbers := make([]int, 0, len(numberStrings))

	// convert the strings to numbers
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
