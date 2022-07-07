package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	lines := ReadStdIn()
	if len(lines) != 4 {
		log.Fatal("expected 4 lines")
	}
	profile := lines[0]
	code := profile[:len(profile)-1]
	fmt.Print("updated: ", code)
}

func ReadStdIn() []string {
	reader := bufio.NewReader(os.Stdin)
	var lines []string
	for {
		line, readErr := reader.ReadString('\n')
		lines = append(lines, line)
		if readErr == io.EOF {
			return lines
		}
	}
}
