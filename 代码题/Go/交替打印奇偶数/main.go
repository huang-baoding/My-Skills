package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 1; i < 101; i++ {

			//奇数
			if i%2 == 1 {
				fmt.Println("线程1打印:", i)
			}
			ch <- struct{}{}
		}
	}()
	go func() {
		defer wg.Done()
		for i := 1; i < 101; i++ {

			//偶数
			if i%2 == 0 {
				fmt.Println("线程2打印:", i)
			}
			<-ch
		}
	}()
	wg.Wait()
}
