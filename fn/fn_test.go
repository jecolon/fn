package main

import "testing"

func TestFn(t *testing.T) {
	conf := &config{
		dir:   "TEST",
		mv:    false,
		out:   "TEST",
		input: []string{"    ", " _fíle$-.txt "},
	}
	process(conf)
	w := []string{"FN_NO_NAME_0", "fíle.txt"}
	for i, n := range conf.output {
		if n != w[i] {
			t.Fatalf("wanted %v but got %v", w, conf.output)
		}
	}
}

func BenchmarkFn(b *testing.B) {
	conf := &config{
		dir:   "TEST",
		mv:    false,
		out:   "TEST",
		input: []string{"    ", " _fíle$-.txt "},
	}
	for i := 0; i < b.N; i++ {
		process(conf)
	}
}
