package main

import(
  //"fmt"
  "flag"
  //"io"
  "strings"
  "os"
  "log"

  "github.com/nqeron/TakAnalysis/utility"
)

var (
  out = flag.String("out","","outfile")
  //debug = flag.Bool("debug",false,debug)
)

func main() {
  flag.Parse()
  fileName := flag.Arg(0)
  f, e := os.Open(fileName)
  defer f.Close()

  if e != nil{
    log.Fatal("open:", e)
  }

  tposed, e := utility.FTransformPTN(f)

  if *out == ""{
    *out = strings.Replace(fileName, ".ptn", "-transposed.ptn",-1)
  }

  outfile, _ := os.Create(*out)
  defer outfile.Close()

  outfile.WriteString(tposed.Render())
}
