package utility

import (
  "io"
  "strconv"
  "math"
  "log"

  "github.com/nelhage/taktician/ptn"
  "github.com/nelhage/taktician/tak"
)

const(
  DEBUG = false
)

func FTransformPTN(r io.Reader) (*ptn.PTN,error){
  //pass any source
  parsed,err := ptn.ParsePTN(r)
  if err != nil{
    log.Fatal("Parse error", err)
  }

  err = TransformPTN(parsed)
  return parsed,err
}

func TransformPTN(parsed *ptn.PTN) error{
  l := len(parsed.Ops)/3
  moves := make([]tak.Move,0,l)

  for _,o := range parsed.Ops{
    m, ok := o.(*ptn.Move)
    if !ok { continue } // skip over non-move elements

    moves = append(moves,m.Move) // add move to list
  }

  sizeTag := parsed.FindTag("Size")
	size, e := strconv.Atoi(sizeTag)
	if e != nil {
		return nil
	}

  err := TransformMoves(moves,size)
  if err != nil {
    return err
  }

  moveInd := 0
  for i,_ := range parsed.Ops{
    z := &parsed.Ops[i]
    _, ok := (*z).(*ptn.Move)
    if !ok {
      continue
    }
    o,_ := (*z).(*ptn.Move)

    o.Move = moves[moveInd]
    moveInd++

  }

  return nil

}

func TransformMoves(moves []tak.Move, size int) error{
  mid := int(math.Ceil(float64(size)/2)) - 1
  size = size - 1 //tak positions 0-N, so
  m := moves[0]
  iT := mid -1 // inner Triangle distance
  if DEBUG {log.Println("TransformMoves, mid: ", mid, ", iT:", iT, " ,m: ", m)}
  if DEBUG {log.Println("mDist: ",mDist(m.X,m.Y,mid,0)) }
  if (m.X > mid) || ( mDist(m.X,m.Y,mid,0) > iT) { // if it's not in a1 (c1)- c3 triangle
    //mirror board
    //f := func (f func([]tak.Move), m []tak.Move) { f(m) }
    if DEBUG {log.Println("not initial triangle")}
    if m.X > mid {
      mirrorX(moves,size)
    }
    if m.Y > mid {
      mirrorY(moves,size)
    }

    m = moves[0] //update to new first move

    if DEBUG {log.Println("mDist from n,0: ",mDist(m.X,m.Y,size,0)) }
    if mDist(m.X,m.Y,size,0) > size { //if move is still above diagonal
      mirrorDiag(moves,size)
      return nil
    }
  }

  if (size %2 == 0) && (m.X == mid){
    var offC tak.Move
    for _,i := range moves{
      if (i.X != mid){
        offC = i
        break
      }
    }
    if ( (offC.X > mid) && ( offC.Y < mid) ) ||
       ( (offC.X < mid) && ( offC.Y > mid) ){
      mirrorX(moves,size)
    }
    return nil
  }

  var offM tak.Move
  for _,i := range moves{
    if(i.X != i.Y){
      offM = i
      break
    }
  }

  if DEBUG { log.Printf("offM: %s, dist to n,0: %d \n",offM,mDist(offM.X,offM.Y,size,0)) }
  if mDist(offM.X,offM.Y,size,0) > size {
    mirrorDiag(moves,size)
  }

  return nil
}

//takes a sequence of tak moves and mirrors them across the x axis
func mirrorX(moves []tak.Move,size int) {
  if DEBUG { log.Println("Xtransform, size: ",size); }
  for i := range moves{
    //t := moves[i]
    moves[i].X = size - moves[i].X
    if moves[i].Type == tak.SlideRight{ moves[i].Type = tak.SlideLeft
    }else if moves[i].Type == tak.SlideLeft{ moves[i].Type = tak.SlideRight;}
    //if DEBUG { log.Printf("%s --> %s",t,moves[i]); }
  }

}

//takes a sequence of tak moves and mirrors them across the y axis
func mirrorY(moves []tak.Move,size int) {
  if DEBUG { log.Println("Ytransform, size: ",size); }
  for i := range moves{
    //t := moves[i]
    moves[i].Y = size - moves[i].Y
    if moves[i].Type == tak.SlideUp{ moves[i].Type = tak.SlideDown
    }else if moves[i].Type == tak.SlideDown{ moves[i].Type = tak.SlideUp;}
    //if DEBUG { log.Printf("%s --> %s",t,moves[i]); }
  }
}

//takes a sequence of tak moves and mirrors them across the a1-nn diagonal
func mirrorDiag(moves []tak.Move,size int) {
  if DEBUG { log.Println("Diagtransform, size: ",size); }
  for i := range moves{
    //t := moves[i]
    if moves[i].X != moves[i].Y {
      y := moves[i].Y
      moves[i].Y = moves[i].X //was size-X , -- may revert to this for a different setup?
      moves[i].X = y
    } // swap moves only for non-diagonals?
    if moves[i].Type == tak.SlideUp{ moves[i].Type = tak.SlideRight
    }else if moves[i].Type == tak.SlideDown{ moves[i].Type = tak.SlideLeft
    }else if moves[i].Type == tak.SlideRight{ moves[i].Type = tak.SlideUp
    }else if moves[i].Type == tak.SlideLeft{ moves[i].Type = tak.SlideDown;}
    //if DEBUG { log.Printf("%s --> %s",t,moves[i]); }
  }
}

//calculates manhattan distance from point (a,b) to (x,y)
func mDist(a,b,x,y int) int{
  return int( math.Abs(float64(x-a)) + math.Abs(float64(y-b) ) )
}
