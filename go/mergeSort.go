package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var threshold int

func mergeSort(arr []int, wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	if len(arr) <= 1 {
		return
	}

	mid := len(arr) / 2
	var leftWg, rightWg sync.WaitGroup
	var left, right = arr[:mid], arr[mid:]

	// Use goroutines only for large enough slices
	if len(left) > threshold {
		leftWg.Add(1)
		go mergeSort(left, &leftWg)
	} else {
		mergeSort(left, nil)
	}
	if len(right) > threshold {
		rightWg.Add(1)
		go mergeSort(right, &rightWg)
	} else {
		mergeSort(right, nil)
	}

	leftWg.Wait()
	rightWg.Wait()
	merge(arr, left, right)
}

func merge(arr, left, right []int) {
	i, j, k := 0, 0, 0
	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
			arr[k] = left[i]
			i++
		} else {
			arr[k] = right[j]
			j++
		}
		k++
	}
	for i < len(left) {
		arr[k] = left[i]
		i++
		k++
	}
	for j < len(right) {
		arr[k] = right[j]
		j++
		k++
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	size := 1000000
	arr := make([]int, size)
	for i := range arr {
		arr[i] = rand.Intn(10000000)
	}

	thresholds := make([]int, 10)
	for i := 0; i < 10; i++ {
		thresholds[i] = 1024 << i
	}
	runs := 5

	for _, t := range thresholds {
		total := time.Duration(0)
		for i := 0; i < runs; i++ {
			a := make([]int, size)
			copy(a, arr)

			start := time.Now()
			var wg sync.WaitGroup
			threshold = t
			wg.Add(1)
			mergeSort(a, &wg)
			wg.Wait()
			elapsed := time.Since(start)
			total += elapsed
		}
		avg := total / time.Duration(runs)
		fmt.Printf("Threshold: %6d | Avg time: %v\n", t, avg)
	}
}
