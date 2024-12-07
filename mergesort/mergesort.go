package main

func RunMergeSortSingle(array []int) {
	tempArray := make([]int, len(array))
	mergeSortSingle(array, tempArray, 0, len(array)-1)
}

func mergeSortSingle(array, tempArray []int, left, right int) {
	if left >= right {
		return
	}

	mid := (left + right) / 2
	mergeSortSingle(array, tempArray, left, mid)
	mergeSortSingle(array, tempArray, mid+1, right)
	merge(array, tempArray, left, mid, right)
}

func RunMergeSortGoroutines(array []int, maxDepth int) {
	tempArray := make([]int, len(array))
	done := make(chan struct{})
	go mergeSortGoroutines(array, tempArray, 0, len(array)-1, 0, maxDepth, done)
	<-done
}

func mergeSortGoroutines(array, tempArray []int, left, right, currentDepth, maxDepth int, done chan struct{}) {
	defer close(done)
	if left >= right {
		return
	}

	mid := (left + right) / 2
	doneLeft := make(chan struct{})
	doneRight := make(chan struct{})

	if currentDepth < maxDepth {
		go mergeSortGoroutines(array, tempArray, left, mid, currentDepth+1, maxDepth, doneLeft)
		go mergeSortGoroutines(array, tempArray, mid+1, right, currentDepth+1, maxDepth, doneRight)

		<-doneLeft
		<-doneRight
	} else {
		mergeSortSingle(array, tempArray, left, mid)
		mergeSortSingle(array, tempArray, mid+1, right)
	}

	merge(array, tempArray, left, mid, right)
}

func merge(array, tempArray []int, left, mid, right int) {
	for i := left; i <= right; i++ {
		tempArray[i] = array[i]
	}

	i, j, k := left, mid+1, left

	for i <= mid && j <= right {
		if tempArray[i] <= tempArray[j] {
			array[k] = tempArray[i]
			i++
		} else {
			array[k] = tempArray[j]
			j++
		}
		k++
	}

	for i <= mid {
		array[k] = tempArray[i]
		i++
		k++
	}
}
