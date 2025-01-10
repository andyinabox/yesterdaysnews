//go:build objectstoretest
// +build objectstoretest

package objectstoreservice

import (
	"testing"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func makeTestErrorStream(t *testing.T) chan<- domain.Error {
	errs := make(chan domain.Error)
	t.Cleanup(func() {
		close(errs)
	})
	go func() {
		err := <-errs
		t.Error(err)
	}()
	return errs
}

func makeTestStringStream(count int, gen func(int) string) <-chan string {
	stream := make(chan string)
	go func() {
		defer close(stream)
		for i := 0; i < count; i++ {
			stream <- gen(i)
		}
	}()
	return stream
}

func makeTest2StringStream(count int, gen1 func(int) string, gen2 func(int, string) string) <-chan [2]string {
	stream := make(chan [2]string)
	go func() {
		defer close(stream)
		for i := 0; i < count; i++ {
			s1 := gen1(i)
			s2 := gen2(i, s1)
			stream <- [2]string{s1, s2}
		}
	}()
	return stream
}

func makeTestStringTo2StringStream(inStream <-chan string, gen func(string) string) <-chan [2]string {
	stream := make(chan [2]string)
	go func() {
		defer close(stream)
		for s1 := range inStream {
			s2 := gen(s1)
			stream <- [2]string{s1, s2}
		}
	}()
	return stream
}
