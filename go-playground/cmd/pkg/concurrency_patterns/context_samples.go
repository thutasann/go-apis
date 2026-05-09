/*
1. context is immutable

2. context can't be altered or mutated

3. context can't be cancelled by its children
*/
package concurrencypatterns

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func func1(ctx context.Context, parentWg *sync.WaitGroup, stream <-chan interface{}) {
	defer parentWg.Done()
	var wg sync.WaitGroup

	doWork := func(ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-stream:
				if !ok {
					fmt.Println("channel closed")
					return
				}
				fmt.Println(d)
			}
		}
	}

	newCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go doWork(newCtx)
	}

	wg.Wait()
}

func genericFunc(ctx context.Context, wg *sync.WaitGroup, stream <-chan interface{}) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-stream:
			if !ok {
				fmt.Println("channel closed")
				return
			}
			fmt.Println(d)
		}
	}
}

func ContextSampleOne() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	generator := func(dataItem string, stream chan any) {
		for {
			select {
			case <-ctx.Done():
				return
			case stream <- dataItem:
			}
		}
	}

	infiniteApples := make(chan interface{})
	go generator("apple", infiniteApples)

	infiniteOranges := make(chan interface{})
	go generator("orange", infiniteOranges)

	infinitePeaches := make(chan interface{})
	go generator("peach", infinitePeaches)

	wg.Add(1)
	go func1(ctx, &wg, infiniteApples)

	func2 := genericFunc
	func3 := genericFunc

	wg.Add(1)
	go func2(ctx, &wg, infiniteOranges)

	wg.Add(1)
	go func3(ctx, &wg, infinitePeaches)

	wg.Wait()
}
