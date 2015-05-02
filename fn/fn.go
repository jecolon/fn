package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"sync"

	"github.com/jecolon/fn"
	flag "github.com/ogier/pflag"
)

type request struct {
	conf  *config
	index int
}

type config struct {
	dir    string
	mv     bool
	out    string
	report bool
	input  []string
	output []string
}

const multiplier = 8

var dir = flag.StringP("dir", "d", ".", "Directory to process. Default is current directory.")
var mv = flag.BoolP("move", "m", false, "Move (rename) instead of copy files.")
var out = flag.StringP("out", "o", ".", "Output directory (relative to --dir) to save copies.")
var report = flag.BoolP("report", "r", false, "Just report the filename changes that would occur.")
var numCPU int
var MaxOutstanding int
var wg sync.WaitGroup

func main() {
	// Setup
	numCPU = runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	MaxOutstanding = numCPU * multiplier
	flag.Parse()

	// Run
	process(&config{
		dir:    *dir,
		mv:     *mv,
		out:    *out,
		report: *report,
	})

	// Wait for gophers
	wg.Wait()
}

func process(conf *config) {
	// If not testing, get input filenames from filesystem
	if conf.dir != "TEST" {
		// Prepare path
		conf.dir = path.Clean(conf.dir)
		// Open source dir
		d, err := os.Open(conf.dir)
		check(err)
		defer d.Close()

		// Read filenames
		conf.input, err = d.Readdirnames(0)
		check(err)
	}

	// Fix filenames for shell
	conf.output = make([]string, len(conf.input))
	for i := range conf.input {
		conf.output[i] = fn.FixForShell(conf.input[i])
	}

	// Testing ends here... for now
	if conf.dir == "TEST" {
		return
	}

	// If just reporting
	if conf.report {
		fmt.Printf("\nfn fix filenames report for %q directory:\n", conf.dir)
		fmt.Printf("Move (rename):\t%t\n", conf.mv)
		fmt.Printf("Output dir:\t%q\n\n", path.Join(conf.dir, conf.out))
		for i := range conf.input {
			fmt.Printf("%q -> %q\n", conf.input[i], conf.output[i])
		}
		fmt.Println()
		return
	}

	// Copy or move files
	err := os.Chdir(conf.dir)
	check(err)

	// Copy uses --out directory (which could be same as --dir)
	if !conf.mv {
		err = os.Mkdir(conf.out, 0750)
		// If --out exists, carry on, stop on other errors
		if err != nil && !os.IsExist(err) {
			check(err)
		}
	}

	requests := make(chan *request)
	for i := 0; i < MaxOutstanding; i++ {
		wg.Add(1)
		go handle(requests)
	}
	for i := range conf.input {
		requests <- &request{conf, i}
	}
	close(requests)
}

func handle(queue <-chan *request) {
	for r := range queue {
		processFile(r)
	}
	wg.Done()
}

func processFile(r *request) {
	// TODO: Refactor these out later
	i := r.index
	conf := r.conf
	// Skip non-fixed names
	if conf.input[i] == conf.output[i] {
		log.Printf("%q looks OK, skipping", conf.input[i])
		return
	}
	// Skip directories
	si, err := os.Stat(conf.input[i])
	check(err)
	if si.IsDir() {
		log.Printf("%q is a directory, skipping", conf.input[i])
		return
	}
	// Open files
	src, err := os.Open(conf.input[i])
	check(err)
	defer src.Close()

	var dst *os.File
	// Output dir
	if conf.mv {
		// Move (rename) stays in source directory
		dst, err = os.Create(conf.output[i])
		check(err)
	} else {
		// Copy to --out directory
		dst, err = os.Create(path.Join(conf.out, conf.output[i]))
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
		err = os.Remove(conf.input[i])
		check(err)
	}
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
