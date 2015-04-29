package main

import (
	"log"
	"os"

	"github.com/jecolon/fn"
)

func main() {
	// Open source dir
	dir, err := os.Open(".")
	check(err)
	defer dir.Close()

	// Read filenames
	names, err := dir.Readdirnames(0)
	check(err)

	// Fix filenames for shell
	format := "%-20s\t%-20s\t%-20s\n"
	log.Printf(format, "Original", "Shell", "URL")
	for _, n := range names {
		log.Printf(format, n, fn.FixForShell(n), fn.FixForURL(n))
	}
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
