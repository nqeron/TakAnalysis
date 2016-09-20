package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/nelhage/taktician/ptn"
)

const (
	fileSuffix = "-taktics.json"
)

var (
	out    = flag.String("out", "", "ptn file to export results to")
	recurs = flag.Bool("r", false, "recurse over directory")
)

func main() {
	flag.Parse()

	if *recurs {
		dir := flag.Arg(0)
		if *out == "" {
			*out = dir[:len(dir)] + fileSuffix //single suffix / json for a directory sweep?
		}

		_, err := os.Stat(*out)
		if os.IsNotExist(err) {
			err = os.Mkdir(*out, os.ModePerm)
			if err != nil {
				log.Fatal("create out file", err)
			}
		}
		//file path walk of some sort?
	} else {
		fileName := flag.Arg(0)
		if *out == "" {
			*out = strings.Replace(fileName, ".ptn", fileSuffix, -1)
		}

		genTakticsForPTN(fileName, *out)
	}
}

func genTakticsForPTN(fileName string, out string) error {
	f, e := os.Open(fileName)
	defer f.Close()

	if e != nil {
		log.Fatal("open file: ", e)
		return e
	}

	outfile, e := os.Create(out)
	defer outfile.Close()
	if e != nil {
		log.Fatal("Outfile creation: ", e)
	}

	parsed, e := ptn.ParsePTN(f)
	if e != nil {
		log.Fatal("parse: ", e)
		return e
	}

	//TODO
	log.Println(parsed.Render())

	//pos := parsed.InitialPosition()
	//parsed.PositionAtMove(move, color)

	return nil
}
