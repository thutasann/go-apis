package main

import (
	"fmt"
	"path/filepath"
)

func FilePath_Sample_One() {
	path := filepath.Join("dir1", "dir2/../dir3", "text.txt")
	fmt.Println(path)
}
