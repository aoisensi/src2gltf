package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/aoisensi/src2gltf/smd"
)

func main() {
	for _, fn := range flag.Args() {
		f, err := os.Open(fn)
		if err != nil {
			log.Println(err)
			continue
		}
		defer f.Close()
		switch filepath.Ext(fn) {
		case ".smd":
			smd, err := smd.Decode(f)
			if err != nil {
				panic(err)
			}
			out := fn[:len(fn)-4] + ".gltf"
			if err := smd2gltf(smd, out); err != nil {
				panic(err)
			}
		default:
			log.Printf("the filetype is not supported.")
		}
	}
}
