package main

import (
  "fmt"
  "os"
  "flag"
  "log"
  "strings"
  "time"
  //"regexp"
  "golang.org/x/net/context"
  "github.com/nqeron/TakAnalysis/analysis"
  "github.com/nelhage/taktician/ptn"
  "github.com/nelhage/taktician/cli"
  "github.com/nelhage/taktician/tak"
)

const(
  VERSION = "0.1.1"
)

//regular expressions for commands
var(
  //goRE = regexp.MustCompile("^go ([1-9]+)")
  //t = 8
  level = flag.Int("level", 8, "minimax level")
)


func main() {
  flag.Parse()
  fileName := flag.Arg(0)

  f, err := os.Open(fileName)
  defer f.Close()

  useFile := false
  if err == nil{
    useFile = true
  }

  fmt.Printf("Welcome to the Tak Analysis console, version: %s \n", VERSION)
  if useFile{
    fmt.Printf("Using file: %s \n",fileName)
  }else{
    fmt.Println("Using blank PTN \n")
  }

  //Only handles files for now
  if err != nil {
    log.Fatal(err)
  }

  // first parse the PTN
  parsed, e := ptn.ParsePTN(f)
  if e!= nil {
    log.Fatal("parse: ",e)
  }

  doAnalysis := YN("Analyze file now?")
  //if !doAnalysis

  if doAnalysis{
    fmt.Println("Analyzing file ...")
    out := strings.Replace(fileName,".ptn","-analysis.ptn",-1)

    //var outfile *os.File
    //if outfile,err := os.Open(out); err != nil{



    outfile, _ :=  os.Create(out)


    analysis.Meta(parsed,outfile,analysis.Config{
        Depth: *level,
        Sensitivity: 2,
        TimeLimit: time.Minute,
        Debug: false,
    })

    outfile.Close()
    fmt.Printf("Analysis Complete - output to file: %s \n",out)

    anFile, err := os.Open(out)
    if err != nil{
      log.Fatal("analysis file: ",err)
    }
    //} else{
      //fmt.Println("Using existant file")
    //}
    parsed, err = ptn.ParsePTN(anFile)
    if err != nil{
      log.Fatal("analysis to ptn: ",err)
    }
  }

  var parsedSize int
  if len(parsed.Ops) % 3 == 0 {
    parsedSize = len(parsed.Ops)/ 3.0
  } else{
    parsedSize = len(parsed.Ops)/ 3.0 + 1
  }
  //options
  var (
    showAll = false //alt - show current
    displayBoard = true
    //depth = 6
    moveNum = 1
    moveColor = tak.White
    moves = make([]*tak.Position,0,5)
    isExplore = false
  )

  //pos, err := parsed.InitialPosition()
  pos,err := parsed.PositionAtMove(moveNum, moveColor)
  if err != nil {
    log.Fatal("position: ", err)
  }
  //var in string
  for {
    if showAll{
      //hl := false
      parsed.Render()
    } else{
      low := (moveNum-1)*3
      up := moveNum*3-1
      //if moveColor == tak.White && up >
      if (moveColor == tak.Black) { up = up +2}
      //fmt.Printf("moveNum: %d, len(ops): %d, parsedSize: %s, moveColor: %s",moveNum,len(parsed.Ops),parsedSize, moveColor)

      fmt.Printf("  (%d / %d)",moveNum, parsedSize)
      for i := low; i <= up ; i++{
        if i >= len(parsed.Ops){
          if moveColor == tak.Black {
            fmt.Printf(" [__]")
            break // don't print more
          }
          fmt.Printf(" __")
        } else {
          op := parsed.Ops[i]
          switch o := op.(type) {
            case *ptn.MoveNumber:
              fmt.Printf("\n%d.", o.Number)
            case *ptn.Move:
              if (i%3 == 2 && moveColor == tak.Black) ||
                 (i%3 == 1 && moveColor == tak.White){
                fmt.Printf(" [%s%s]", ptn.FormatMove(&o.Move), o.Modifiers)
              }else{
                fmt.Printf(" %s%s", ptn.FormatMove(&o.Move), o.Modifiers)
              }
            default:
          }
        }
      }


    }

    if displayBoard{
      cli.RenderBoard(os.Stdout, pos)
    }

    //handle commands
    for {
      var cmd string
      var num int
      var mod string
      fmt.Scanf("%s",&cmd)
      switch cmd {
        case "go":
          fmt.Scanf("%d %s",&num,&mod)
          //fmt.Printf("Called Go :",groups[1])
          moveNum = num
          moveColor = tak.White
          if(mod == "b" || mod == "black"){
            moveColor = tak.Black
          }
        case "next", "n", "d":
          if(moveColor == tak.Black){
            moveNum++
            moveColor = tak.White
          } else{
            moveColor = tak.Black
          }
        case "prev","p","a":
          if moveColor == tak.White{
            moveNum--
            moveColor= tak.Black
          } else{
            moveColor = tak.White
          }
        case "w":
          moveNum++
        case "s":
          moveNum--
        case "ai":
          ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))
	        defer cancel()
          ai := analysis.MakeAI(pos,*level)
          moves,val,_ := ai.Analyze(ctx,pos)
          fmt.Println("current value: ",val)
          fmt.Printf("anticipated moves: ")
          for _, m := range moves{
            fmt.Printf("%s, ",ptn.FormatMove(&m))
          }
          fmt.Printf("\n")
        case "print","render":
          fmt.Println(parsed.Render())
        case "undo","u","r":
          if len(moves) <=0{
            fmt.Println("Nothing to undo!")
          }else{
            pos = moves[len(moves)-1]
            moves = moves[:len(moves)-1]
            moveColor = moveColor ^ 192
            if len(moves) <= 0{
              isExplore = false
            }
          }
        case "q":   //stop exploring
          isExplore = false;
          moves = moves[:0]
        case "exit":
          isSure := YN("Are you sure want to exit the program?")
          if(isSure){
            os.Exit(0)
          }
        default:
          newMove, err := ptn.ParseMove(cmd)
          if(err != nil){
            fmt.Println("Command not recognized! cmd: ",cmd)
          }else{ // handling a move
            lastPos := pos
            pos, err = pos.Move(&newMove)
            if err != nil{
              fmt.Println("Movement error: ",err)
              continue //go back to commands ?
            }
            moves = append(moves,lastPos)
            isExplore = true
            moveColor = moveColor ^ 192 // swap color ()
          }

      }

      //check to make sure upcoming position is valid
      if !isExplore{
        pos, err = parsed.PositionAtMove(moveNum,moveColor)
        if err != nil{
          fmt.Println("position isn't valid: ",err)
          continue
        }
      }
      //moveNum = int(groups[1])
      break
    }
    //fmt.Scanln()
  }


}

func YN(q string) bool {
  var in string
  for {
    fmt.Println(q)
    fmt.Scanln(&in)
    in = strings.ToLower(in)
    if (in == "y") || (in == "yes") {
      return true
    }else if (in =="n") || (in == "no") {
      return false
    } else{
      fmt.Println("I don't recognize that response")
    }
  }
}
