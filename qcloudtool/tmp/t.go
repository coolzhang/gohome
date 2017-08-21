package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	content := "temporary file's content"
	fmt.Println(content)
	dir, err := ioutil.TempDir(".", "example")
	fmt.Println(dir)
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(dir) // clean up

	tmpfn := filepath.Join(dir, "tmpfile")
	fmt.Println(tmpfn)
	if err := ioutil.WriteFile(tmpfn, []byte(content), 0666); err != nil {
		log.Fatal(err)
	}

}
