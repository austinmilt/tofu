package array

import (
	"slices"
)

// https://stackoverflow.com/a/57213476
func RemoveAt[T interface{}](s []T, index int) ([]T, T) {
	val := s[index]
	return slices.Delete(s, index, index+1), val
}

// https://stackoverflow.com/a/58069984/3314063
func Shift[T interface{}](arr []T, n int) []T {
	n = n % len(arr)
	return append(arr[n:len(arr):len(arr)], arr[:n]...)
}
