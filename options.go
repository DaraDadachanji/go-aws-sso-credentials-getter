package main

import (
	"flag"
	"fmt"
)

func DoOptions() (halt bool) {
	version := flag.Bool("version", false, "check program version")

	flag.Parse()
	if *version {
		halt = true
		fmt.Println(VERSION)
	}
	return halt
}
