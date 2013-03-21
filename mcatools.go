package main

import (
	"runtime"

	"github.com/pa001024/GoCraft/nbt"
	"io/ioutil"
	"os"

	// "compress/gzip"
	// "compress/zlib"
	// "encoding/binary"
	// "encoding/hex"
	// "bufio"
	// "bytes"
	// "errors"
	// "fmt"
	// "io"
	"log"
)

var (
	clearTile   bool
	clearEntity bool
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	seq := make(chan int, runtime.NumCPU())
	exit := make(chan int)
	a := os.Args[1:]
	log.Println("Processing..")
	go func() {
		for _, v := range a {
			println(v)
			switch v {
			case "-t":
				clearTile = true
			case "-e":
				clearEntity = true
			case "*":
				fs, _ := ioutil.ReadDir(".")
				for _, v := range fs {
					if !v.IsDir() {
						seq <- 1
						err := Convert(v.Name())
						if err != nil {
							log.Fatal(err)
						}
						<-seq
					}
				}
			default:
				seq <- 1
				err := Convert(v)
				if err != nil {
					log.Fatal(err)
				}
				<-seq
			}
		}
		exit <- 0
		log.Println("Fin")
	}()
	<-exit
}
func Convert(file string) error {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return err
	}
	o, err := nbt.ReadRegion(f)
	if err != nil {
		return err
	}
	for _, root := range o.Chunks {
		for _, v := range root.Data {
			if v == nil {
				return err
			}
			ls := v.(*nbt.Compound).List("TileEntities")
			ls.Type = nbt.TagCompound
			ls.Data = []*nbt.Compound{}
		}
	}
	nfile := file + ".bak"
	f.Close()
	os.Remove(nfile)
	err = os.Rename(file, nfile)
	if err != nil {
		return err
	}
	nf, err := os.Create(file)
	defer nf.Close()
	if err != nil {
		return err
	}
	o.WriteRegion(nf)
	return err
}
