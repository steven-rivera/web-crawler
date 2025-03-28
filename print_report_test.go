package main

import (
	"reflect"
	"testing"
)

func TestSortPages(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		expected []page
	}{
		{
			name: "order count descending",
			input: map[string]int{
				"url1": 5,
				"url2": 1,
				"url3": 3,
				"url4": 10,
				"url5": 7,
			},
			expected: []page{
				{url: "url4", count: 10},
				{url: "url5", count: 7},
				{url: "url1", count: 5},
				{url: "url3", count: 3},
				{url: "url2", count: 1},
			},
		},
		{
			name: "alphabetize",
			input: map[string]int{
				"d": 1,
				"a": 1,
				"e": 1,
				"b": 1,
				"c": 1,
			},
			expected: []page{
				{url: "a", count: 1},
				{url: "b", count: 1},
				{url: "c", count: 1},
				{url: "d", count: 1},
				{url: "e", count: 1},
			},
		},
		{
			name: "order count then alphabetize",
			input: map[string]int{
				"d": 2,
				"a": 1,
				"e": 3,
				"b": 1,
				"c": 2,
			},
			expected: []page{
				{url: "e", count: 3},
				{url: "c", count: 2},
				{url: "d", count: 2},
				{url: "a", count: 1},
				{url: "b", count: 1},
			},
		},
		{
			name:     "empty map",
			input:    map[string]int{},
			expected: []page{},
		},
		{
			name:     "nil map",
			input:    nil,
			expected: []page{},
		},
		{
			name: "one key",
			input: map[string]int{
				"url1": 1,
			},
			expected: []page{
				{url: "url1", count: 1},
			},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := sortPages(tc.input)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected url: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}