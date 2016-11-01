package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/context"

	tai "github.com/nelhage/taktician/ai"
	"github.com/nelhage/taktician/ptn"
	"github.com/nelhage/taktician/tak"
	"github.com/nqeron/TakAnalysis/analysis"
)

const (
	fileSuffix = "-taktics.json"
)

var (
	recurs  = flag.Bool("r", false, "recurse over directory")
	out     = flag.String("out", "", "File to export JSON / CSV to")
	debug   = flag.Bool("debug", false, "debug")
	verbose = flag.Bool("v", false, "verbose handling")
)

//TakTPS : structure that maps a given tps to a given set of moves that are tak?
type TakTPS struct {
	Tps   string
	Moves []string
}

//type TakTPS map[string][]string

//var colTPS []TakTPS

func main() {
	flag.Parse()
	var colTPS = make([]TakTPS, 0)
	if *recurs {
		dir := flag.Arg(0)
		if *out == "" {
			*out = dir[:len(dir)] + fileSuffix
		}

		var errors []error

		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".ptn") || strings.HasSuffix(path, ".tps") {
				if *verbose {
					fmt.Println(">> ", path)
				}
				tps, moves, e := getTaksFromPTNTPS(path)

				if e != nil {
					errors = append(errors, fmt.Errorf("Error at path: %s: %s", path, e))
					return nil //fmt.Errorf("Error at path: %s: %s", path, e)
				}
				colTPS = append(colTPS, TakTPS{Tps: tps, Moves: moves})
				//takTPS[tps] = moves
			}

			return nil
		})
		if len(errors) >= 0 {
			//log.Fatalln("Error walk: ", e)
			if *verbose {
				fmt.Println("Errors ")
				for _, err := range errors {
					fmt.Println("\t ", err)
				}
			} else {
				fmt.Println("Errors in some files, use -v option to list all errors")
			}
		}

	} else {
		fileName := flag.Arg(0)
		if *out == "" {
			*out = strings.Replace(fileName, ".ptn", fileSuffix, -1)
		}
		tps, moves, e := getTaksFromPTNTPS(fileName)
		if e != nil {
			log.Fatalln("Error get Taks: ", e)
		}
		colTPS = append(colTPS, TakTPS{Tps: tps, Moves: moves})
		//takTPS[tps] = moves
	}
	//output takTPS to JSON in file
	outfile, e := os.Create(*out)
	defer outfile.Close()

	if e != nil {
		log.Fatal("Open out file: ", e)
	}
	outJSON := json.NewEncoder(outfile)
	outJSON.SetEscapeHTML(false)
	outJSON.Encode(colTPS)

}

func getTaksFromPTNTPS(fileName string) (tpn string, outmoves []string, err error) {
	f, e := os.Open(fileName)
	defer f.Close()

	if e != nil {
		log.Fatal("open: ", e)
		return "", nil, e
	}

	var pos *tak.Position
	var tps string
	if strings.HasSuffix(fileName, ".ptn") { //PTN parser
		parsed, e := ptn.ParsePTN(f)
		if e != nil {
			log.Fatal("parse: ", e)
			return "", nil, e
		}

		tps = parsed.FindTag("TPS")

		if tps == "" { //if there is no tps, then skip this one!
			return "", nil, fmt.Errorf("No TPS tag in PTN!")
		}
		//otherwise use tps in tag
		pos, err = ptn.ParseTPS(tps)
	} else if strings.HasSuffix(fileName, ".tps") {
		buf, err := ioutil.ReadAll(f)
		tps = string(buf)
		pos, err = ptn.ParseTPS(tps)
		if err != nil {
			log.Fatal("Error parsing TPS: ", err)
			return "", nil, fmt.Errorf("Error parsing TPS: %s", tps)
		}
	} else {
		//log.Fatal("Not PTN or TPS")
		return "", nil, fmt.Errorf("Not PTN or TPS: %s", fileName)
	}

	ai := analysis.MakeAI(pos, 1)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))
	defer cancel()
	moves, v, _ := ai.AnalyzeAll(ctx, pos)

	if *debug {
		fmt.Println("possible move(s): ", moves)
	}
	var outMoves []string
	if (v >= tai.WinThreshold || v <= -tai.WinThreshold) && moves != nil && len(moves) > 0 {
		for _, mSet := range moves {

			if len(mSet) > 1 || len(mSet) < 1 {
				continue
			}

			p, _ := pos.Move(&mSet[0])
			if p.WinDetails().Reason != tak.RoadWin { //only include road / tak wins
				continue
			}

			sMove := ptn.FormatMove(&mSet[0])
			outMoves = append(outMoves, sMove)
		}
	}
	return tps, outMoves, nil
}
