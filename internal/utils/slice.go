package utils

import "math"

type Page[T any] struct {
	Page             int `json:"page"`
	Size             int `json:"size"`
	NumberOfElements int `json:"numberOfElements"`
	TotalElements    int `json:"totalElements"`
	TotalPages       int `json:"totalPages"`
	Content          []T `json:"content"`
}

func PaginateSlice[T any](slice []T, page int, size int) []T {
	return Paginate(slice, page, size).Content
}

func Paginate[T any](slice []T, page int, size int) Page[T] {
	p := page - 1
	length := len(slice)
	start := p * size

	if start > length {
		start = length
	}
	end := start + size
	if end > length {
		end = length
	}
	part := slice[start:end]
	return Page[T]{
		Page:             page,
		Size:             size,
		NumberOfElements: len(part),
		TotalElements:    length,
		TotalPages:       int(math.Ceil(float64(length) / float64(size))),
		Content:          part,
	}
}

func MapSlice[T any, U any](input []T, mapper func(T) U) []U {
	output := make([]U, len(input))
	for i, item := range input {
		output[i] = mapper(item)
	}
	return output
}

func FilterSlice[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}
