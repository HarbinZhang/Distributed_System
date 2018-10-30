package main

import (
	"fmt"
)

type help struct {
	ch chan []string
}

func main() {
	// w := help{make(chan []string)}
	w := help{make(chan []string, 1)}
	temp := []string{"hi"}
	w.ch <- temp

	go run(w)

	temp = []string{"hi"}
	w.ch <- temp

	// go con(w)
	
	select {
	case msg := <-w.ch:
		fmt.Println(msg)q
	default:
		fmt.Println("Wrong")
	}
}

func con() {

}

func run(w help) {
	wh := <-w.ch
	wh = append(wh, "hello")
	fmt.Println(wh)

	// w.ch <- wh
}
