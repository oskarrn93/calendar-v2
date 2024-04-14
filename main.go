package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")

	config := InitializeConfig()

	fmt.Println("config", config)
}
