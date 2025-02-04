package streams

import (
	"context"
	"sync"
	"time"
)

// func StringPipe(ctx context.Context, dest chan<- string, strs ...string) <-chan struct{} {
// 	return StringPipeThrottled(ctx, 0, dest, strs...)
// }

// func StringPipeThrottled(ctx context.Context, d time.Duration, dest chan<- string, strs ...string) <-chan struct{} {
// 	done := make(chan struct{})

// 	go func() {
// 		defer func() {
// 			done <- struct{}{}
// 			close(done)
// 		}()
// 		for _, s := range strs {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			default:
// 				dest <- s
// 				time.Sleep(d)
// 			}
// 		}
// 	}()
// 	return done
// }

func StringPipeThrottled(ctx context.Context, d time.Duration, inStream <-chan string) <-chan string {
	stream := make(chan string)
	go func() {
		defer close(stream)
		for s := range inStream {
			select {
			case <-ctx.Done():
				return
			default:
				stream <- s
				time.Sleep(d)
			}
		}
	}()
	return stream
}

func StringStream(ctx context.Context, strs ...string) <-chan string {
	return StringStreamThrottled(ctx, 0, strs...)
}

func StringStreamThrottled(ctx context.Context, d time.Duration, strs ...string) <-chan string {
	stream := make(chan string)

	go func() {
		defer close(stream)
		for _, s := range strs {
			select {
			case <-ctx.Done():
				return
			default:
				stream <- s
				time.Sleep(d)
			}
		}
	}()

	return stream
}

func StringSlice(ctx context.Context, stream <-chan string) (slice []string) {
	var mu sync.Mutex
	slice = []string{}

	for s := range stream {
		select {
		case <-ctx.Done():
			return
		default:
			mu.Lock()
			slice = append(slice, s)
			mu.Unlock()
		}
	}

	return slice
}

func StringTransformStream(ctx context.Context, inStream <-chan string, fn func(string) string) <-chan string {
	outStream := make(chan string)

	go func() {
		defer close(outStream)
		for s := range inStream {
			select {
			case <-ctx.Done():
				return
			default:
				outStream <- fn(s)
			}
		}
	}()

	return outStream
}

func StringFilterStream(ctx context.Context, inStream <-chan string, fn func(string) bool) <-chan string {
	outStream := make(chan string)

	go func() {
		defer close(outStream)
		for s := range inStream {
			select {
			case <-ctx.Done():
				return
			default:
				if fn(s) {
					outStream <- s
				}
			}
		}
	}()

	return outStream
}

func StringStreamTo2StringSliceStream(ctx context.Context, inStream <-chan string, fn func(string) [2]string) <-chan [2]string {
	outStream := make(chan [2]string)

	go func() {
		defer close(outStream)
		for s := range inStream {
			select {
			case <-ctx.Done():
				return
			default:
				outStream <- fn(s)
			}
		}
	}()

	return outStream
}

func MergeStringStreams(ctx context.Context, streams ...<-chan string) <-chan string {
	merged := make(chan string)

	var wg sync.WaitGroup

	wg.Add(len(streams))
	for _, stream := range streams {
		go func() {
			defer wg.Done()
			for s := range stream {
				select {
				case <-ctx.Done():
					return
				default:
					merged <- s
				}
			}
		}()
	}

	// close merge stream
	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

func DuplicateStringStream(ctx context.Context, stream <-chan string) (<-chan string, <-chan string) {
	stream1 := make(chan string)
	stream2 := make(chan string)

	cleanup := func() {
		close(stream1)
		close(stream2)
	}

	go func() {
		defer cleanup()
		for s := range stream {
			select {
			case <-ctx.Done():
				return
			default:
				stream1 <- s
				stream2 <- s
			}
		}
	}()

	return stream1, stream2
}

func BifurcatedStringStream(ctx context.Context, stream <-chan string, fn func(s string) bool) (<-chan string, <-chan string) {
	trueStream := make(chan string)
	falseStream := make(chan string)

	var wg sync.WaitGroup

	cleanup := func() {
		wg.Wait()
		close(trueStream)
		close(falseStream)
	}

	check := func(s string) {
		defer wg.Done()
		if fn(s) {
			trueStream <- s
		} else {
			falseStream <- s
		}
	}

	go func() {
		defer cleanup()
		for s := range stream {
			select {
			case <-ctx.Done():
				return
			default:
				wg.Add(1)
				go check(s)
			}
		}
	}()

	return trueStream, falseStream
}
