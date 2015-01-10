package main

import (
	"flag"
	"fmt"
	"github.com/ftrvxmtrx/tga"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	// "runtime"
	"strings"
	"sync"
	"sync/atomic"
)

var wg sync.WaitGroup
var filecount int64
var dircount int64
var errorcount int64

//path values from flags
var path string
var outdir string

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic: ", r)
		}
	}()
	doflags()

	walkdir("/Users/Josh")
	wg.Wait()
	fmt.Printf("Found %d file %d directories and encountered %d errors", filecount, dircount, errorcount)
}

func walkdir(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Error at: ", dir, "\n>>", err.Error())
		atomic.AddInt64(&errorcount, 1)
		return
	}

	for _, v := range files {
		if v.IsDir() {
			newdir := filepath.Join(dir, v.Name())
			atomic.AddInt64(&dircount, 1)
			walkdir(newdir)
		} else {
			wg.Add(1)
			atomic.AddInt64(&filecount, 1)
			go func(f os.FileInfo, dir string) {
				defer wg.Done()
				name := filepath.Join(dir, f.Name())
				fmt.Println(name)
				//call a convert method here passing in the filename "name"
			}(v, dir)
		}
	}
}

func doflags() {
	flag.StringVar(&path, "path", "", "Folder with targa files")
	flag.StringVar(&outdir, "out", "", "Output folder where png files will be saved")
	inplace := flag.Bool("i", false, "If set, will ignore output directory and convert files inplace")
	flag.Parse()

	//check flags
	if f, err := os.Stat(path); os.IsNotExist(err) {
		panic(err)
	} else if err != nil {
		panic(err)
	} else if !f.IsDir() {
		panic(f.Name() + " is not a valid directory!")
	}

}

//converttga should always be called in a new go function
func converttga(fname string) {
	defer wg.Done()
	file, err := os.Open(fname)
	if err != nil {
		fmt.Println("Could not open file: ", fname, "\n>>", err.Error())
		return
	}
	img, err := tga.Decode(file)
	if err != nil {
		fmt.Println("Error decoding file: ", fname, "\n>>", err.Error())
		return
	}
}

func getfinaldir(fname string) string {
	return "" + "hi"
}
