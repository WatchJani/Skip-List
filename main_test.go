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
