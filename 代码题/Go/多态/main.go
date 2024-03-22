// 多态的体现在于，把接口作为参数传入一个方法，根据接口指向的对象的不同而表现出不同的行为
package main

import (
	"fmt"
)

type Anymal interface {
	Voice()
}

type Dog struct{}

func (dog *Dog) Voice() { fmt.Println("汪汪汪~") }

type Cat struct{}

func (cat *Cat) Voice() { fmt.Println("喵喵喵~") }

func MakeSound(a Anymal) { a.Voice() }

func main() {
	var inter Anymal
	d := &Dog{}
	inter = d
	MakeSound(inter)
	c := &Cat{}
	inter = c
	MakeSound(inter)
}
