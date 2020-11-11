package main

import (
	"bufio"
	"fmt"
	"golang.org/x/tools/benchmark/parse"
	"os"
)

func main() {
	file, err := os.Open("leftist.txt")
	if err != nil {
		fmt.Println("Error")
		return
	}

	io := bufio.NewReader(file)

	set, err := parse.ParseSet(io)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, s := range set {
		fmt.Println(s)
	}
}
