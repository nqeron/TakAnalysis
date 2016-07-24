package analysis

//command line direct for meta-analysis call
//run stuff built
import (
  "log"
  "flag"
  "time"
  "strings"
  "os"
  "github.com/nelhage/taktician/ptn"
)

var (
  out = flag.String("out", "", "ptn file to export results to")
  depth = flag.Int("depth", 6, "minimax depth")
  debug = flag.Bool("debug", false, "debug logs for analysis")
  timeLimit = flag.Duration("limit", time.Minute, "analysis time limit")
  sensitivity = flag.Int("sensitivity", 1, "level to highlight")
)

func main() {
  flag.Parse()
  fileName := flag.Arg(0)
  f, e := os.Open(fileName)
  defer f.Close()
  if e != nil {
    log.Fatal("open: ",e)
  }

  if *out == "" {
    *out = strings.Replace(fileName,".ptn","-flagged.ptn",-1)
  }

  outfile, e := os.Create(*out)
  defer outfile.Close()
  if e != nil{
    log.Fatal("open out file: ", e)
  }

  parsed, e := ptn.ParsePTN(f)
  if e!= nil {
    log.Fatal("parse: ",e)
  }
  Meta(parsed,outfile,Config{
      Depth: *depth,
      Sensitivity: *sensitivity,
      TimeLimit: *timeLimit,
      Debug: *debug,
      Verbose: false,
      AnnotationOnly: false,
  })
}
