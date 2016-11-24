package main

import "os"
import "fmt"

func main() {
	argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]


	fmt.Println(argsWithProg)
	fmt.Println(argsWithoutProg)
}
