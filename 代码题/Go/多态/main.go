package main

import "fmt"

type Usb interface {
	Study()
}

type StuA struct {
	sub string
}

func (a StuA) Study() {
	fmt.Printf("a同学在学习%s\n", a.sub)
}

type StuB struct {
	sub string
}

func (b StuB) Study() {
	fmt.Printf("b同学在学习%s\n", b.sub)
}

type Working struct {
}

func (w Working) StartStudy(u Usb) {
	u.Study()
}
func main() {
	var w Working
	var u Usb
	var a StuA = StuA{"数学"}
	u = a
	w.StartStudy(u)
	var b StuB = StuB{"语文"}
	u = b
	w.StartStudy(u)

}
