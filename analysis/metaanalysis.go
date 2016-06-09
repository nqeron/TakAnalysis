package analysis

import (
  "log"
  "time"
  "math"
  "os"
  "fmt"
  "github.com/nelhage/taktician/ptn"
  "github.com/nelhage/taktician/ai"
  "github.com/nelhage/taktician/tak"
  "github.com/montanaflynn/stats"
)

var(
  //depth = 6 //flag.Int("depth", 6, "minimax depth")
  seed int64 = 0 //flag.Int64("seed",0,"random seed")
  //debug = 0 //flag.Int("debug", 1, "debug level")
  //timeLimit = time.Minute //flag.Duration("limit", time.Minute, "analysis time limit")

  sort = true
  table = true
)

type metaFlag int

type Config struct{
  Sensitivity int
  Depth int
  TimeLimit time.Duration
  Debug bool
}

const (
  NoFlag = 0
  Minimum = 1
  Maximum = 2
  Mean = 3
  Q1 = 4
  Q2 = 5
  Q3 = 6
  /*OneStd = 4
  TwoStd = 5
  ThreeStd = 6*/
  //NumFlags = 6
)

//meta-analysis on a ptn file
//assume ptn is parsed?
func Meta(parsed *ptn.PTN,file *os.File, cfg Config){
  p,e := parsed.InitialPosition()
  if e != nil{
    log.Fatalln("initial: ", e)
  }
  w,b := MakeAI(p,cfg.Depth), MakeAI(p,cfg.Depth)

  var values stats.Float64Data = make(stats.Float64Data,0,len(parsed.Ops))
  var moves  []tak.Move = make([]tak.Move, 0, len(parsed.Ops))
  //index := -1
  if cfg.Debug{log.Println("Analyzing ...");}
  for _, o := range parsed.Ops {
    m, ok := o.(*ptn.Move)
    if !ok { continue
    }
    //index++;

    if cfg.Debug && (len(moves)%5 == 0){ log.Println("...  ", len(moves)); }
    moves = append(moves,m.Move)
    //moves[index] = m.Move
    //var pmoves []tak.Move
    var val int64
    switch {
    case p.ToMove() == tak.White:
      //log here?
        _, val, _ = w.Analyze(p,cfg.TimeLimit)
    case p.ToMove() == tak.Black :
      //log?
        _, val, _ = b.Analyze(p,cfg.TimeLimit)
    }
    values = append(values,float64(val))
    //values[index] = float64(val)

    var e error
    p, e = p.Move(&m.Move)
    if e != nil {
      log.Fatalf("illegal move %s: %v", ptn.FormatMove(&m.Move),e)
    }
  }

  //extract flagging for clarity
  var flags []metaFlag = make([]metaFlag, len(values))
  flagVals(&values, flags)
  if cfg.Debug { log.Println("flags: ",flags); }

  //write flags to output file

  // print tags
  for _, tag := range parsed.Tags{
    file.WriteString(fmt.Sprintf("[%s \"%s\"]\n",tag.Name,tag.Value))
  }
  file.WriteString(fmt.Sprintf("\n"))
  //for each line, print modified
  for i := 0; i + 1 < len(moves); i = i+2  {
    wMoveM := moves[i]
    wFlag := flags[i]
    bMoveM := moves[i+1]
    bFlag := flags[i+1]

    wMoveS := ptn.FormatMove(&wMoveM)
    wMoveS = flagMove(wMoveS,wFlag,cfg.Sensitivity)
    bMoveS := ptn.FormatMove(&bMoveM)
    bMoveS = flagMove(bMoveS,bFlag,cfg.Sensitivity)

    var moveNum int = (i + 2) /2
    if cfg.Debug{ log.Printf("%d. %s %s: vals: %f, %f",moveNum,wMoveS,bMoveS,values[i],values[i+1]) ; }
    file.WriteString(fmt.Sprintf("%d. %s %s\n",moveNum,wMoveS,bMoveS))
  }
}

func flagMove(moveS string,flag metaFlag,sens int) string{
  //log.Printf("flag: ", flag, "sens: ", sens)
  if (sens <0) && ( (flag >= Q1) || (flag <= -Q1) ) {
    return moveS
  } else if (sens >0) && (sens <=3) && ( (flag>=Q1) || (flag <= -Q1) )  {
    if (flag < 0 ){
      flag = flag + metaFlag(sens)
    }else{
      flag = flag - metaFlag(sens) //reduce level of flag
    }
    if flag < Q1 && flag > -Q1 { //if level goes under, then don't modify move
      return moveS
    }
  }
  if flag == NoFlag{
    return moveS
  // } else if flag == Minimum {
  //   return "(" + moveS
  // } else if flag == Maximum {
  //   return ")" + moveS
  } else if flag == Q1{ // StandardDeviation of some sort
    return moveS + "!"
  } else if flag == -Q1{
    return moveS + "?"
  }else if (flag == Q2){
    return  moveS +"!!"
  }else if (flag == -Q2)  {
    return  moveS +"??"
  }else if flag == Q3{
    return moveS + "!*"
  } else if flag == -Q3{
    return moveS + "?*"
  }
  return moveS
}

func flagVals(values *stats.Float64Data,flags []metaFlag)  {
  //values := (*vs)[:len(*vs)-1] //chop off last value?
  //chop of last two positions for min/max
  // min / max / iqr of values
  //mean,_ := stats.Mean(values)
  //std,_ := stats.StandardDeviation(values)

  var dVals stats.Float64Data = make(stats.Float64Data,0,len(*values)-1) // collect the differences between successive values
  for i:=1; i < len(*values); i++{
    dVals = append(dVals,(*values)[i] + (*values)[i-1]) //add because adjacent vals are opposite signs
  }

  min,_ := stats.Min(dVals)
  max,_ := stats.Max(dVals)

  //quartiles,_ := stats.Quartile(*values)
  iqr,_ := stats.InterQuartileRange(dVals)

  //log.Printf("min: %f, max: %f, half iqr: %f  \n",min,max,iqr/2)

  for j, dVal := range dVals {
    i := j+1 //assign flag to move at the end of the d
    //var dVal float64 = 1
    //if i >= 1 {
        //dVal = val + (*values)[i-1] //get last play value
        //dVal = val + (*values)[i-1] //last val of given player
    //}
    dAdjust := dVal / (iqr/2) //use half the iqr - that is the 25% between median and either Q1 or Q3
    var dir int = int(dAdjust / math.Abs(dAdjust)) //-1 or +1 depending on direction

    if dVal == min{
      flags[i] = Minimum
    } else if dVal == max {
      flags[i] = Maximum
    } else if math.Abs(dAdjust) >= 1.5{ //val < quartiles.Q1{
      flags[i] = metaFlag(Q3 * dir)
    } else if math.Abs(dAdjust) >= 1{
      flags[i] = metaFlag(Q2 * dir)
    } else if math.Abs(dAdjust) >= 0.5 {
      flags[i] = metaFlag(Q1 * dir)
    } else{
      flags[i] = NoFlag
    }
  }
}

func MakeAI(p *tak.Position, depth int) *ai.MinimaxAI{
  return ai.NewMinimax(ai.MinimaxConfig{
    Size: p.Size(),
    Depth: depth,
    Seed: seed,
    Debug: 0,

    NoSort: !sort,
    NoTable: !table,
  })
}
