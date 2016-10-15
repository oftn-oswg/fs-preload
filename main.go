package main

import (
	"flag"
	"log"
	"os"
	"syscall"

	"os/signal"
	"path/filepath"

	"golang.org/x/exp/mmap"
)

func mapper(files chan string) {
	for path := range files {
		_, err := mmap.Open(path)
		log.Printf("%s\n", path)
		if err != nil {
			log.Printf("Error mapping file: %s\n", err)
		}
	}
}

func main() {
	files := make(chan string)
	defer close(files)

	flag.Parse()

	rootDirectory := flag.Arg(0)

	go mapper(files)

	err := filepath.Walk(rootDirectory, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			files <- path
		}
		return nil
	})
	if err != nil {
		log.Printf("Error traversing directory: %s\n", err)
	}

	log.Println("Files mapped to memory.")

	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<- exit
}
