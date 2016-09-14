package analysis

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/montanaflynn/stats"
	"github.com/nelhage/taktician/ai"
	"github.com/nelhage/taktician/ptn"
	"github.com/nelhage/taktician/tak"
)

type metaFlag struct {
	Name       string
	Value      string
	Annotation string
	Level      int
}

//Config -- configuration for the analysis:
// Sensitivity - sensitivity to Annotation
// Depth - level to search for moves and Tinues
// TimeLimit - max # time to wait
// Debug - print debug info
// Verbose - show progress / additional text
// AnnotationOnly - show annotations only (no comments)
type Config struct {
	Sensitivity    int
	Depth          int
	TimeLimit      time.Duration
	Debug          bool
	Verbose        bool
	AnnotationOnly bool
}

const (
	noFlag = 0
)

//Cfg - configuration passed to analysis program
var (
	//depth = 6 //flag.Int("depth", 6, "minimax depth")
	seed int64 //flag.Int64("seed",0,"random seed")
	//debug = 0 //flag.Int("debug", 1, "debug level")
	//timeLimit = time.Minute //flag.Duration("limit", time.Minute, "analysis time limit")

	sort  = true
	table = true
	Cfg   = Config{Sensitivity: 1, Depth: 8, TimeLimit: time.Minute, Debug: false, Verbose: false, AnnotationOnly: false}
)

//Meta - meta-analysis on a ptn file
//assume ptn is parsed?
func Meta(parsed *ptn.PTN, file *os.File, customCfg Config) {
	Cfg = customCfg
	p, e := parsed.InitialPosition()
	if e != nil {
		log.Fatalln("initial: ", e)
	}
	w, b := MakeAI(p, Cfg.Depth), MakeAI(p, Cfg.Depth)

	var values = make(stats.Float64Data, 0, len(parsed.Ops))
	var moves = make([]tak.Move, 0, len(parsed.Ops))
	var movePos = make([]*tak.Position, 0, len(parsed.Ops))
	//index := -1
	if Cfg.Debug {
		log.Println("Analyzing ...")
	}
	for i, o := range parsed.Ops {
		m, ok := o.(*ptn.Move)
		if !ok {
			continue
		}
		//index++;

		if (Cfg.Debug || Cfg.Verbose) && (i%2 == 0) {
			log.Println("...  ", i, "/", len(parsed.Ops))
		}
		moves = append(moves, m.Move)
		//moves[index] = m.Move
		//var pmoves []tak.Move

		if !Cfg.AnnotationOnly {
			var val int64
			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(Cfg.TimeLimit))
			defer cancel()
			switch {
			case p.ToMove() == tak.White:
				//log here?
				_, val, _ = w.Analyze(ctx, p)
			case p.ToMove() == tak.Black:
				//log?
				_, val, _ = b.Analyze(ctx, p)
			}
			values = append(values, float64(val))
			//values[index] = float64(val)
		}
		var e error
		p, e = p.Move(&m.Move)
		if e != nil {
			log.Fatalf("illegal move %s: %v", ptn.FormatMove(&m.Move), e)
		}
		movePos = append(movePos, p)
	}

	//extract flagging for clarity
	var flags = make([][]metaFlag, len(moves))
	if !Cfg.AnnotationOnly {
		flagVals(&values, flags)
		flagMoves(&moves, flags)
	}
	flagPositions(&movePos, flags, w)
	if Cfg.Debug {
		log.Println("flags: ", flags)
		log.Println("len vals: ", len(values), "moves: ", len(moves))
	}

	//write flags to output file

	// print tags
	for _, tag := range parsed.Tags {
		file.WriteString(fmt.Sprintf("[%s \"%s\"]\n", tag.Name, tag.Value))
	}
	file.WriteString(fmt.Sprintf("\n"))
	//for each line, print modified
	for i := 0; i < len(moves); i = i + 2 {
		wMoveM := moves[i]
		wFlag := flags[i]

		wMoveS := ptn.FormatMove(&wMoveM)
		wMoveS = writeMoveFlags(wMoveS, wFlag)

		bMoveS := ""
		var bVal float64
		if i+1 < len(moves) {
			bMoveM := moves[i+1]
			bFlag := flags[i+1]
			bVal = values[i+1]
			bMoveS = ptn.FormatMove(&bMoveM)
			bMoveS = writeMoveFlags(bMoveS, bFlag)
		}

		var moveNum = (i + 2) / 2
		if Cfg.Debug {
			log.Printf("%d. %s %s: vals: %f, %f", moveNum, wMoveS, bMoveS, values[i], bVal)
		}
		file.WriteString(fmt.Sprintf("%d. %s %s\n", moveNum, wMoveS, bMoveS))
	}
}

//writes a move's flags to the PTN
func writeMoveFlags(moveS string, flags []metaFlag) string {
	annot := ""
	comm := ""
	for _, flag := range flags {
		if Cfg.Sensitivity < 0 {
			continue
		} else if Cfg.Sensitivity < flag.Level {
			continue
		}
		if flag.Annotation != "" {
			annot = annot + flag.Annotation
		}
		if flag.Value == "" { //only print value if exists
			comm = comm + fmt.Sprintf("%s, ", flag.Name)
		} else {
			comm = comm + fmt.Sprintf("%s = %s, ", flag.Name, flag.Value)
		}
	}

	if comm == "" || Cfg.AnnotationOnly {
		return moveS + annot
	}
	return moveS + annot + " {" + comm + "}"
}

//flags based on position of board at given move
func flagPositions(movePos *[]*tak.Position, flags [][]metaFlag, ai *ai.MinimaxAI) {

	if Cfg.Verbose {
		log.Println("Flagging positions...")
	}
	prevTin := 0
	prevYield := false
	for i := range *movePos {

		if Cfg.Verbose && i%5 == 0 {
			log.Println(i, "/", len(*movePos))
		}
		pos := (*movePos)[i]

		move, curTin, depth := HasTinue(pos, ai)
		fMove := ""

		if move != nil {
			fMove = fmt.Sprintf("Depth: %d, Move: ", depth)
			fMove += ptn.FormatMove(move)
		}
		if (prevTin != 0 || prevYield) && curTin == 0 && i != len(flags)-1 {
			flags[i] = append(flags[i], metaFlag{Name: "wasTinue", Annotation: "", Value: fMove, Level: 1})
			prevYield = false
		} else if prevTin != 0 && (curTin == -prevTin) { //same tinue (opposite sign)
			if prevYield {
				flags[i] = append(flags[i], metaFlag{Name: "newTinue", Annotation: "''", Value: fMove, Level: 1})
			} else {
				flags[i] = append(flags[i], metaFlag{Name: "stepTinue", Annotation: "", Value: fMove, Level: 1})
			}
			prevYield = false
		} else if prevTin != 0 && (curTin == prevTin) { //different tinue?
			flags[i] = append(flags[i], metaFlag{Name: "yieldsTinueA", Annotation: "??", Value: fMove, Level: 1})
			prevYield = true
		} else if (prevTin == 0 || prevYield) && (curTin < 0) {
			flags[i] = append(flags[i], metaFlag{Name: "newTinue", Annotation: "''", Value: fMove, Level: 1})
			prevYield = false
		} else if (prevTin == 0 || prevYield) && (curTin > 0) {
			flags[i] = append(flags[i], metaFlag{Name: "yieldsTinueB", Annotation: "??", Value: fMove, Level: 1})
			prevYield = true //set to 0 to force a newTinue or hadTinue next
		} else if Cfg.Debug {
			flags[i] = append(flags[i], metaFlag{Name: "Debug", Value: fmt.Sprintf("i: %d, curTin: %v", i, curTin), Level: 1})
		}

		if move, ok := IsTak(pos); ok {
			flags[i] = append(flags[i], metaFlag{Name: "isTak", Annotation: "'", Value: ptn.FormatMove(move), Level: 1})
		}

		prevTin = curTin
	}

}

//flags moves based on moves
func flagMoves(moves *[]tak.Move, flags [][]metaFlag) {
	//does nothing right now
}

// flags moves based on calculated values
func flagVals(values *stats.Float64Data, flags [][]metaFlag) {
	//values := (*vs)[:len(*vs)-1] //chop off last value?
	//chop of last two positions for min/max
	// min / max / iqr of values
	//mean,_ := stats.Mean(values)
	//std,_ := stats.StandardDeviation(values)

	var dVals = make(stats.Float64Data, 0, len(*values)-1) // collect the differences between successive values
	for i := 1; i < len(*values); i++ {
		dVals = append(dVals, (*values)[i]+(*values)[i-1]) //add because adjacent vals are opposite signs
	}

	min, _ := stats.Min(dVals)
	max, _ := stats.Max(dVals)

	//quartiles,_ := stats.Quartile(*values)
	iqr, _ := stats.InterQuartileRange(dVals)

	//log.Printf("min: %f, max: %f, half iqr: %f  \n",min,max,iqr/2)

	for j, dVal := range dVals {
		i := j + 1 //assign flag to move at the end of the d
		//var dVal float64 = 1
		//if i >= 1 {
		//dVal = val + (*values)[i-1] //get last play value
		//dVal = val + (*values)[i-1] //last val of given player
		//}
		dAdjust := dVal / (iqr / 2) //use half the iqr - that is the 25% between median and either Q1 or Q3
		//var dir int = int(dAdjust / math.Abs(dAdjust)) //-1 or +1 depending on direction

		if dVal == min {
			flags[i] = append(flags[i], metaFlag{Name: "dV Minimum", Level: 3}) //Minimum
		} else if dVal == max {
			flags[i] = append(flags[i], metaFlag{Name: "dV Maximum", Level: 3}) //Maximum
		} else if math.Abs(dAdjust) >= 1.5 { //val < quartiles.Q1{
			if dAdjust > 0 {
				flags[i] = append(flags[i], metaFlag{Name: "dV+ Q3", Annotation: "!!!", Level: 3}) //metaFlag(Q3 * dir)
			} else {
				flags[i] = append(flags[i], metaFlag{Name: "dV- Q3", Annotation: "???", Level: 3}) //metaFlag(Q3 * dir)
			}
		} else if math.Abs(dAdjust) >= 1 {
			if dAdjust > 0 {
				flags[i] = append(flags[i], metaFlag{Name: "dV+ Q2", Annotation: "!!", Level: 2}) //metaFlag(Q2 * dir)
			} else {
				flags[i] = append(flags[i], metaFlag{Name: "dV- Q2", Annotation: "??", Level: 2}) //metaFlag(Q2 * dir)
			}
		} else if math.Abs(dAdjust) >= 0.5 {
			if dAdjust > 0 {
				flags[i] = append(flags[i], metaFlag{Name: "dV+ Q1", Annotation: "!", Level: 1}) //metaFlag(Q1 * dir)
			} else {
				flags[i] = append(flags[i], metaFlag{Name: "dV- Q1", Annotation: "?", Level: 1}) //metaFlag(Q1 * dir)
			}

		} else {
			//flags[i] = NoFlag
		}
	}
}

//MakeAI -- creates an ai w/ given depth
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
