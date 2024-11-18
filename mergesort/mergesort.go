package main

func MergeSort(array []int) []int {
	if len(array) <= 1 {
		return array
	}

	mid := len(array) / 2
	left := MergeSort(array[:mid])
	right := MergeSort(array[mid:])

	return merge(left, right)
}

func MergeSortGoroutine(array []int, chunkSize int) <-chan []int {
	c := make(chan []int)
	go func() {
		if len(array) <= 1 {
			c <- array
			return
		}

		mid := len(array) / 2
		var left []int
		var right []int

		if mid >= chunkSize {
			leftChan := MergeSortGoroutine(array[:mid], chunkSize)
			rightChan := MergeSortGoroutine(array[mid:], chunkSize)

			for i := 0; i < 2; i++ {
				select {
				case msg1 := <-leftChan:
					left = msg1
				case msg2 := <-rightChan:
					right = msg2
				}
			}
		} else {
			left = MergeSort(array[:mid])
			right = MergeSort(array[mid:])
		}

		c <- merge(left, right)
	}()
	return c
}

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	result = append(result, left[i:]...)
	result = append(result, right[j:]...)

	return result
}
