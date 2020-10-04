package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	//"path/filepath"
)

func printDirs(out io.Writer, root os.FileInfo, curPath string, isLast bool, printerSimple string, printerLast string) {
	var files []os.FileInfo
	var err error
	if root.IsDir() {
		files, err = ioutil.ReadDir(curPath)
	}
	if err != nil {
		panic(err)
	}
	if isLast {
		var size int64
		size = -1
		sizePrint := ""
		if !root.IsDir() {
			file, err := os.Open(curPath)
			if err != nil {
				panic(err)
			}
			stat, err := file.Stat()
			if err != nil {
				panic(err)
			}
			size = stat.Size()
		}
		if size == 0 {
			sizePrint = "(empty)"
		} else {
			sizePrint = "(" + strconv.FormatInt(size, 10) + "b)"
		}
		if size == -1 {
			fmt.Fprintln(out, printerLast+root.Name())
		} else {
			fmt.Fprintln(out, printerLast+root.Name(), sizePrint)
		}
		temp := strings.Split(printerSimple, "")
		temp[len(temp)-4] = "\t" + temp[len(temp)-4]
		printerSimple = strings.Join(temp, "")

		temp = strings.Split(printerLast, "")
		temp[len(temp)-4] = "\t" + temp[len(temp)-4]
		printerLast = strings.Join(temp, "")
	} else {
		var size int64
		size = -1
		sizePrint := ""
		if !root.IsDir() {
			file, err := os.Open(curPath)
			if err != nil {
				panic(err)
			}
			stat, err := file.Stat()
			if err != nil {
				panic(err)
			}
			size = stat.Size()
		}
		if size == 0 {
			sizePrint = "(empty)"
		} else {
			sizePrint = "(" + strconv.FormatInt(size, 10) + "b)"
		}
		if size == -1 {

			fmt.Fprintln(out, printerSimple+root.Name())
		} else {
			fmt.Fprintln(out, printerSimple+root.Name(), sizePrint)
		}
		printerSimple = "│\t" + printerSimple
		printerLast = "│\t" + printerLast
	}
	for i := 0; i < len(files); i++ {
		var size int64
		size = -1
		sizePrint := ""
		if !files[i].IsDir() {
			file, err := os.Open(path.Join(curPath, files[i].Name()))
			if err != nil {
				panic(err)
			}
			stat, err := file.Stat()
			if err != nil {
				panic(err)
			}
			size = stat.Size()
		}
		if size == 0 {
			sizePrint = "(empty)"
		} else {
			sizePrint = "(" + strconv.FormatInt(size, 10) + "b)"
		}
		if files[i].IsDir() {
			if i == len(files)-1 {
				printDirs(out, files[i], path.Join(curPath, files[i].Name()), true, printerSimple, printerLast)
			} else {
				printDirs(out, files[i], path.Join(curPath, files[i].Name()), false, printerSimple, printerLast)
			}
		} else if i == len(files)-1 {
			if size == -1 {
				fmt.Fprintln(out, printerLast+files[i].Name(), sizePrint)
			} else {
				fmt.Fprintln(out, printerLast+files[i].Name(), sizePrint)
			}
		} else {
			if size == -1 {
				fmt.Fprintln(out, printerSimple+files[i].Name(), sizePrint)
			} else {
				fmt.Fprintln(out, printerSimple+files[i].Name(), sizePrint)
			}
		}
	}
}

func printDirsNoFiles(out io.Writer, root os.FileInfo, curPath string, isLast bool, printerSimple string, printerLast string) {
	var files []os.FileInfo
	var err error
	if root.IsDir() {
		files, err = ioutil.ReadDir(curPath)
		files = removFiles(files)
		if isLast {
			fmt.Fprintln(out, printerLast+root.Name())
			temp := strings.Split(printerSimple, "")
			temp[len(temp)-4] = "\t" + temp[len(temp)-4]
			printerSimple = strings.Join(temp, "")

			temp = strings.Split(printerLast, "")
			temp[len(temp)-4] = "\t" + temp[len(temp)-4]
			printerLast = strings.Join(temp, "")
		} else {
			fmt.Fprintln(out, printerSimple+root.Name())
			printerSimple = "│\t" + printerSimple
			printerLast = "│\t" + printerLast
		}
	}
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(files); i++ {
		if files[i].IsDir() {
			if i == len(files)-1 {
				printDirsNoFiles(out, files[i], path.Join(curPath, files[i].Name()), true, printerSimple, printerLast)
			} else {
				printDirsNoFiles(out, files[i], path.Join(curPath, files[i].Name()), false, printerSimple, printerLast)
			}
		}
	}
}

func removFiles(files []os.FileInfo) []os.FileInfo {
	f := true
	for f {
		f = false
		var index int
		for i := range files {
			if !files[i].IsDir() {
				index = i
				f = true
				break
			}
		}
		if f {
			for i := index; i < len(files)-1; i++ {
				files[i] = files[i+1]
			}
			files = files[:len(files)-1]
		}
	}
	return files
}

func dirTree(out io.Writer, rootPath string, condition bool) error {
	printerSimple := "├───"
	printerLast := "└───"
	files, err := ioutil.ReadDir(rootPath)
	if err != nil {
		panic(err)
	}
	if condition {
		for i := range files {
			if i == len(files)-1 {
				printDirs(out, files[i], path.Join(rootPath, files[i].Name()), true, printerSimple, printerLast)
			} else {
				printDirs(out, files[i], path.Join(rootPath, files[i].Name()), false, printerSimple, printerLast)
			}
		}
	} else {
		files = removFiles(files)
		for i := range files {
			if i == len(files)-1 {
				printDirsNoFiles(out, files[i], path.Join(rootPath, files[i].Name()), true, printerSimple, printerLast)
			} else {
				printDirsNoFiles(out, files[i], path.Join(rootPath, files[i].Name()), false, printerSimple, printerLast)
			}
		}
	}
	return nil
}
func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
