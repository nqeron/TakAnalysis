package main

import (
	"flag"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/nelhage/taktician/ai"
	"github.com/nelhage/taktician/ptn"
	"github.com/nelhage/taktician/tak"
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

		fPos := parsed.PositionAtMove(0, tak.NoColor)

		over, win := fPos.GameOver()
		if !over || (win != tak.RoadWin) {
			return nil // if the game not over, or not in a road win, skip this game
		}

		// determine penultimate move
		mNum := fPos.MoveNumber()
		lMoveNum := mNum / 2
		lMoveColor := (mNum % 2 * tak.White) + ((mNum + 1) % 2 * tak.Black)
		lMovePos := parsed.PositionAtMove(lMoveNum, lMoveColor)

		//check to see if win is inevitable
		pai := MakeAI(lMovePos, 1)
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))
		defer cancel()
		moves, val, _ := pai.AnalyzeAll(ctx, lMovePos)
		if math.Abs(float64(val)) <= ai.WinThreshold {
			return nil //if not in the win threshhold, then don't worry about it
		}

		var tinueMoves = make([]tak.Move, 0)
		for _, m := range moves {
			if len(moves) > 1 || len(moves) <= 0 {
				continue
			}
			endPos := lMovePos.Move(m[0]) //do I need to worry about tracking lMovePos to make sure it doesn't change?
			if endPos.GameOver() && (endPos.WinDetails().Reason == tak.RoadWin) {
				tinueMoves = append(tinueMoves, m[0])
			}
		}

		if len(tinueMoves) > 0 {
			//output tinueMoves
		}
		return nil
		//pIt.
	}

	//TODO
	//log.Println(parsed.Render())

	//pos := parsed.InitialPosition()
	//parsed.PositionAtMove(move, color)

	return nil
}

func MakeAI(p *tak.Position, depth int) *ai.MinimaxAI {
	return ai.NewMinimax(ai.MinimaxConfig{
		Size:  p.Size(),
		Depth: depth,
		Seed:  seed,
		Debug: 0,

		NoSort:     !sort,
		NoTable:    !table,
		NoNullMove: false,
	})
}
