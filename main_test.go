package main

import "testing"

//170ns multithread ryzen 5 3500x
//136ns multithread ryzen 5 5600x

// 156ns //single thread
func BenchmarkInsert(b *testing.B) {
	b.StopTimer()

	sl := New(32, b.N*2, 0.25)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		sl.Insert(i, i)
	}
}

//180ns for multi core

// 30ns search
//20ns ryzen 5 5600x

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

// 4 130 000ns
func BenchmarkLoopMemTable(b *testing.B) {
	b.StopTimer()

	sl := New(32, b.N*2, 0.25)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for range 40_000 {
			sl.Insert(i, i)
		}
		sl.Clear()
	}
}
