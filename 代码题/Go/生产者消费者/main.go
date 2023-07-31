package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Product struct {
	id     int
	pro_id int
}

var ch chan Product = make(chan Product, 1)
var stop bool

func provide(wgp *sync.WaitGroup, i int) {
	defer wgp.Done()
	for !stop {
		pro_id := rand.Intn(100000)
		m := Product{i, rand.Intn(100000)}
		ch <- m
		fmt.Printf("生产者 %d 生产了产品： %d\n", i, pro_id)
		time.Sleep(1 * time.Microsecond)
	}
}

func consume(wgc *sync.WaitGroup, i int) {
	defer wgc.Done()
	for m := range ch {
		fmt.Printf("消费者 %d 消费了生产者 %d 生产的产品： %d\n", i, m.id, m.pro_id)
		time.Sleep(1 * time.Microsecond)
	}
}

func main() {
	var wgp, wgc sync.WaitGroup
	wgp.Add(5)
	wgc.Add(5)

	for i := 0; i < 5; i++ {
		go provide(&wgp, i)
		go consume(&wgc, i)
	}

	go func() {
		time.Sleep(3 * time.Second)
		stop = true
	}()
	wgp.Wait()
	close(ch)
	wgc.Wait()
}
