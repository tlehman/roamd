package main

import (
	"fmt"
	"os"
)

const (
	OK = 0
	ERR_NOARGS = 1  // error when user doesn't pass in note through args
	ERR_NOAUTH = 2  // error when user doesn't set auth env vars for roam research
)

func usage() {
	fmt.Printf("Usage:\n\troamd 'example note to be appended as a roam block'\n")
	os.Exit(ERR_NOARGS)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	note := os.Args[1]
	fmt.Printf("Appending '%s' to the queue\n", note)
}


