package main

import "testing"

// 150ns
func BenchmarkInsert(b *testing.B) {
	b.StopTimer()

	sl := New(32, b.N*2, 0.25)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		sl.Insert(i, i)
	}
}

// 30ns search
func BenchmarkSearch(b *testing.B) {
	b.StopTimer()

	sl := New(32, 12315*2, 0.25)

	for index := range 12315 {
		sl.Insert(index, index)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		sl.Search(152)
	}
}
