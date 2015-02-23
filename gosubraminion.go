package main

import (
	"flag"
	"fmt"
	//"strings"
	"hash/crc32"
	"io/ioutil"
	"os"
	"path/filepath"
)

// cli args
var target_dir string
var file_types string
var delete_files bool
var prompt_user bool
var skip_dot bool

func getHash(filename string) (uint32, error) {
	//fmt.Println(filename)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	h := crc32.NewIEEE()
	h.Write(data)
	return h.Sum32(), nil
}

func deleteFiles(hmap map[uint32][]string, key uint32) {
	fmt.Println("deleting files:", hmap[key])
	fmt.Println("confirm y/n?")
	var input string
	iflag := false
	for iflag == false {
		fmt.Scanln(&input)
		switch input {
		case "y":
			iflag = true
		case "n":
			iflag = true
		default:
			fmt.Println("invalid choice")
		}
	}
	if input == "n" {
		return
	}
	for _, v := range hmap[key] {
		if os.Remove(v) != nil {
			fmt.Println("error while deleting")
			return
		}
		delete(hmap, key)
	}
	fmt.Println("deleted!")
	return
}

func processDuplicates(hmap map[uint32][]string) {
	for k, v := range hmap {
		if len(v) > 1 {
			fmt.Println("found duplicates")
			c := 0
			for i, fi := range v {
				fmt.Println(i, fi)
				c++
			}
			var choice int
			cflag := false
			for cflag == false {
				fmt.Println("enter index of file that you want to keep or", c, "to delete all")
				fmt.Scanln(&choice)
				if choice > c {
					fmt.Println("invalid choice")
				} else {
					cflag = true
				}
			}
			if choice < len(hmap[k]) {
				hmap[k] = append(v[:choice], v[choice+1:]...)
			}
			deleteFiles(hmap, k)
		}
		//fmt.Println(k,v)
	}
}

func main() {
	flag.StringVar(&target_dir, "t", ".", "target directory")
	flag.StringVar(&file_types, "f", "all", "comma separated file types")
	flag.BoolVar(&delete_files, "d", false, "dry run if this is false")
	flag.BoolVar(&prompt_user, "p", true, "prompt user before deleting files")
	flag.BoolVar(&skip_dot, "s", true, "skip dot files")
	flag.Parse()
	// get file ext types
	// f_types := strings.Split(file_types, ",")
	// clean target dir
	t_dir := filepath.Clean(target_dir)
	// check if dir exists
	if _, err := os.Stat(t_dir); os.IsNotExist(err) {
		fmt.Println("incorrect path: ", target_dir)
		return
	}
	// map to store hash values like map[hash_val] = []string{path1,path2}
	uniques_map := make(map[uint32][]string)
	// traverse dir and populate map
	filepath.Walk(t_dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			h, err := getHash(path)
			if err != nil {
				return err
			}
			dup_files, ok := uniques_map[h]
			if ok {
				uniques_map[h] = append(dup_files, path)
			} else {
				uniques_map[h] = []string{path}
			}
		}
		return nil
	})
	//
	processDuplicates(uniques_map)
	//fmt.Println(uniques_map)
}
