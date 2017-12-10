package main

import (
	"bytes"
	"compress/zlib"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func visit(searchDir string, fileList map[string]string) {
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			contents := decode(path)
			if strings.HasPrefix(contents, "tree") || strings.HasPrefix(contents, "commit") {
				contents = "N/A\n"
			}
			if strings.HasSuffix(contents, "\n") {
				fileList[path] = contents
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	return
}

func decode(filename string) string {
	buf := new(bytes.Buffer)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	reader, err := zlib.NewReader(file)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer reader.Close()
	buf.ReadFrom(reader)
	return buf.String()
}

func main() {
	flag.Parse()
	searchDir := flag.Arg(0)
	fmt.Println("searchdir:", searchDir)
	fileList := make(map[string]string)
	visit(searchDir, fileList)
	for file, content := range fileList {
		fmt.Println(file, "content:", content)
	}
}
