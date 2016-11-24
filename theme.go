package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Flavor struct {
	base00 []string
	base01 []string
	base02 []string
	base03 []string
	base04 []string
	base05 []string
	base06 []string
	base07 []string
	base08 []string
	base09 []string
	base0A []string
	base0B []string
	base0C []string
	base0D []string
	base0E []string
	base0F []string
}

func main() {
	argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]

	fmt.Println(argsWithProg)
	fmt.Println(argsWithoutProg)

	input, err := ioutil.ReadFile("myfile")
	if err != nil {
		fmt.Println(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "]") {
			lines[i] = "LOL"
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile("myfile", []byte(output), 0644)
	if err != nil {
		fmt.Println(err)
	}
}
