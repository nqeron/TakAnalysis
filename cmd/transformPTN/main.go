package main

import(
  "fmt"
  "flag"
  //"io"
  "strings"
  "os"
  "log"
  "path/filepath"

  "github.com/nqeron/TakAnalysis/utility"
)

var (
  out = flag.String("out","","outfile")
  recurs = flag.Bool("r",false,"recurs over a directory")
  verbose = flag.Bool("v",false,"verbose")
  //debug = flag.Bool("debug",false,debug)
)

func main() {
  flag.Parse()

  if *recurs {
    dir := flag.Arg(0)

    if *out == ""{
      *out = dir[:len(dir)]+"-out/"

    }
    _, err := os.Stat(*out)
    if os.IsNotExist(err){ //if the directory doesn't exist, create it
      err = os.Mkdir(*out,os.ModePerm)
      if err != nil{
        log.Fatal("create dir: ",err)
      }
    }

    fmt.Println("Outputting to: ", *out)

    e:= filepath.Walk(dir, func (path string, info os.FileInfo, err error ) error {
        if *verbose { fmt.Println(path); }
        if strings.HasSuffix(path,".ptn"){
          outfile := filepath.Join(*out, filepath.Base(path))
          if *verbose { fmt.Println("---> ", outfile); }
          err := transposeFile(path,outfile)
          return err
        }
        return err
      })

    if e != nil{
      log.Fatal("directory walk: ",e)
    }

  }else {
    fileName := flag.Arg(0)
    if *out == ""{
      *out = strings.Replace(fileName, ".ptn", "-transposed.ptn",-1)
    }
    transposeFile(fileName, *out)

  }
}

func transposeFile(fileName string, out string) error{

  f, e := os.Open(fileName)
  defer f.Close()

  if e != nil{
    return e
  }

  tposed, e := utility.FTransformPTN(f)
  outfile, _ := os.Create(out)
  defer outfile.Close()
  outfile.WriteString(tposed.Render())

  return nil
}
