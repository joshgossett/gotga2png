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

	walkdir(path)
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
			go func(f os.FileInfo, dir string) {
				defer wg.Done()
				name := filepath.Join(dir, f.Name())
				//call a convert method here passing in the filename "name"
				if strings.ToLower(filepath.Ext(name)) == ".tga" {
					atomic.AddInt64(&filecount, 1)
					fmt.Println(name)
					wg.Add(1)
					go converttga(name)
				}
			}(v, dir)
		}
	}
}

func doflags() {
	flag.StringVar(&path, "path", "", "Folder with targa files")
	flag.StringVar(&outdir, "out", "", "Output folder where png files will be saved")
	flag.Parse()

	//if path is not set then use the current working directory
	if path == "" {
		path = filepath.Dir(os.Args[0])
	}

	fmt.Println("Walking path: ", path)

	//check flags
	if f, err := os.Stat(path); os.IsNotExist(err) {
		panic(err)
	} else if err != nil {
		panic(err)
	} else if !f.IsDir() {
		panic(f.Name() + " is not a valid directory!")
	}

	// if outdir == "" {
	outdir = path
	// }

	// err := os.Mkdir(outdir, os.ModeDir)
	// if err != nil {
	// 	panic(outdir)
	// }
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

	newname := getfinaldir(fname)
	pngfile, err := os.Create(newname)
	if err != nil {
		fmt.Println("Error, could not create file: ", newname)
		return
	}
	err = png.Encode(pngfile, img)
	if err != nil {
		fmt.Println("Could not Encode file to png: ", fname)
	}
}

func getfinaldir(fname string) string {
	partialp := fname[len(path):len(fname)-len(filepath.Ext(fname))] + ".png"
	newpath := filepath.Join(outdir, partialp)
	// fmt.Println(newpath)
	return newpath
}
