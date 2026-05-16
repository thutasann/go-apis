package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func FilePath_Sample_One() {
	path := filepath.Join("dir1", "dir2/../dir3", "text.txt")
	// fmt.Println(path)
	fmt.Println(filepath.IsAbs(path))
	fmt.Println(filepath.IsAbs("/dir/file"))
	fmt.Println(filepath.Ext(path))
}

func OS_Stat_Sample() {
	fileInfo, err := os.Stat(FILE_PATH)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fileInfo.Name())
	fmt.Println(fileInfo.Size())
	fmt.Println(fileInfo.Mode())
	fmt.Println(fileInfo.IsDir())
}
