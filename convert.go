package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

var wg sync.WaitGroup
var filecount int64
var dircount int64
var errorcount int64

func main() {
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
