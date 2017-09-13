package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("one...")
	time.Sleep(400 * time.Millisecond)

	fmt.Println("two...")
	time.Sleep(400 * time.Millisecond)

	fmt.Println("three!")
}
