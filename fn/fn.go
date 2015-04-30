package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jecolon/fn"
	flag "github.com/ogier/pflag"
)

type config struct {
	dir    string
	mv     bool
	out    string
	report bool
	input  []string
	output []string
}

var dir = flag.StringP("dir", "d", ".", "Directory to process. Default is current directory.")
var mv = flag.BoolP("move", "m", false, "Move (rename) instead of copy files.")
var out = flag.StringP("out", "o", ".", "Output directory for copies. Default is --dir flag's value.")
var report = flag.BoolP("report", "r", false, "Just report the filename changes that would occur.")

func main() {
	flag.Parse()
	process(&config{
		dir:    *dir,
		mv:     *mv,
		out:    *out,
		report: *report,
	})
}

func process(conf *config) {
	// Input slice
	var names []string
	// For tests
	if conf.dir == "TEST" {
		names = conf.input
	} else {
		// Open source dir
		d, err := os.Open(conf.dir)
		check(err)
		defer d.Close()

		// Read filenames
		names, err = d.Readdirnames(0)
		check(err)
	}

	// Fix filenames for shell
	conf.output = make([]string, len(names))
	for i, n := range names {
		conf.output[i] = fn.FixForShell(n)
	}

	// Testing ends here... for now
	if conf.dir == "TEST" {
		return
	}

	// If just reporting
	if conf.report {
		fmt.Printf("\nfn fix filenames report for %q directory:\n", conf.dir)
		fmt.Printf("Move (rename):\t%t\n", conf.mv)
		fmt.Printf("Output dir:\t%q\n\n", conf.dir+"/"+conf.out)
		for i, n := range names {
			fmt.Printf("%q -> %q\n", n, conf.output[i])
		}
		fmt.Println()
		return

	}

	// Copy or move files
	err := os.Chdir(conf.dir)
	check(err)
	for i, n := range names {
		// Skip non-fixed names
		if n == conf.output[i] {
			continue
		}
		// Open files
		src, err := os.Open(n)
		check(err)
		defer src.Close()

		var dst *os.File
		// Output dir
		if conf.mv {
			dst, err = os.Create(conf.output[i])
			check(err)
		} else {
			err = os.MkdirAll(conf.out, 0750)
			if err != nil && !os.IsExist(err) {
				check(err)
			}
			dst, err = os.Create(conf.out + "/" + conf.output[i])
			check(err)
		}

		defer dst.Close()
		// Copy
		_, err = io.Copy(dst, src)
		check(err)
		// Close now; maybe many files to process
		src.Close()
		dst.Close()
		// Move is really a copy + delete source
		if conf.mv {
			err = os.Remove(n)
			check(err)
		}
	}
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
