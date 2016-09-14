package main

//command line direct for meta-analysis call
//run stuff built
import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nelhage/taktician/ptn"
	"github.com/nqeron/Takanalysis/analysis"
)

var (
	out         = flag.String("out", "", "ptn file to export results to")
	depth       = flag.Int("depth", 6, "minimax depth")
	debug       = flag.Bool("debug", false, "debug logs for analysis")
	timeLimit   = flag.Duration("limit", time.Minute, "analysis time limit")
	sensitivity = flag.Int("sensitivity", 1, "level to highlight")
	recurs      = flag.Bool("r", false, "recurs over directory")
	verbose     = flag.Bool("v", false, "verbose")
	noteOnly    = flag.Bool("noteOnly", false, "only generate basic notation")
)

func main() {
	flag.Parse()

	if *recurs {
		dir := flag.Arg(0)
		if *out == "" {
			*out = dir[:len(dir)] + "-analysis/"
		}
		_, err := os.Stat(*out)
		if os.IsNotExist(err) { // if directory doesn't exist
			err = os.Mkdir(*out, os.ModePerm)
			if err != nil {
				log.Fatal("create dir: ", err)
			}

		}

		fmt.Println("Outputting to: ", *out)

		e := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			fmt.Println("Analyzing: ", path)
			if strings.HasSuffix(path, ".ptn") {
				outfile := filepath.Join(*out, filepath.Base(path))
				if *verbose {
					fmt.Println("---> ", outfile)
				}
				err := analyzeFile(path, outfile)
				return err
			}
			return err
		})

		if e != nil {
			log.Fatal("directory walk: ", e)
		}
	} else { // single file
		fileName := flag.Arg(0)
		if *out == "" {
			*out = strings.Replace(fileName, ".ptn", "-analysis.ptn", -1)
		}
		analyzeFile(fileName, *out)
	}

}

func analyzeFile(fileName string, out string) error {

	f, e := os.Open(fileName)
	defer f.Close()
	if e != nil {
		log.Fatal("open: ", e)
		return e
	}

	outfile, e := os.Create(out)
	defer outfile.Close()
	if e != nil {
		log.Fatal("open out file: ", e)
		return e
	}

	parsed, e := ptn.ParsePTN(f)
	if e != nil {
		log.Fatal("parse: ", e)
		return e
	}
	analysis.Meta(parsed, outfile, analysis.Config{
		Depth:          *depth,
		Sensitivity:    *sensitivity,
		TimeLimit:      *timeLimit,
		Debug:          *debug,
		Verbose:        *verbose,
		AnnotationOnly: *noteOnly,
	})

	return nil
}
