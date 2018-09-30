package analysis

import (
	"math"
	"time"

	tai "github.com/nelhage/taktician/ai"
	"github.com/nelhage/taktician/tak"
	"golang.org/x/net/context"
)

//HasTinue - decides whether a given position has tinue, uses ai to do so
// returns - next (predicted) move for tinue, isTinue (bool), value of position
func HasTinue(pos *tak.Position, ai *tai.MinimaxAI) (move *tak.Move, ok int, depth int) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))
	defer cancel()
	m, v, _ := ai.Analyze(ctx, pos)
	//is this best way to check for tinuê?
	if (v >= tai.WinThreshold || v <= -tai.WinThreshold) && m != nil && len(m) > 0 {
		x := float64(v)
		return &m[0], int(x / math.Abs(x)), len(m)
	}
	return nil, 0, 0
}

//IsTak - returns whether given position is Tak or not
func IsTak(pos *tak.Position) (move *tak.Move, ok bool) {
	p, _ := pos.Move(tak.Move{Type: tak.Pass}) //pass and see if there's a winning Move
	ai := MakeAI(pos, 1)                        //ai that only searches at a depth of 1

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))
	defer cancel()
	m, v, _ := ai.Analyze(ctx, p)

	if v >= tai.WinThreshold && m != nil && len(m) > 0 {
		return &m[0], true
	}
	return nil, false

}
