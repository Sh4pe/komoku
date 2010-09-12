/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */

/*
 * Here the Game struct is defined. It stores all the information which komoku associates
 * with one game
 */

package komoku

import (
    "container/vector"
)

// ################################################################################
// ########################### Game struct ########################################
// ################################################################################
type Game struct {
    Board *Board // The current board
    komi float
    sequence vector.Vector // the sequence of moves
}

// ##################### Game methods ##########################

func (g *Game) GetKomi() float {
    return g.komi
}

// Returns the last move.
func (g *Game) LastMove() *Move {
    if g.sequence.Len() == 0 {
        return nil
    }
    ret, _ := g.sequence.At(g.sequence.Len()-1).(Move)
    return &ret
}

func (g *Game) PlayMove(x, y int, color Color) (err Error) {
    g.sequence.Push(*NewMove(*NewPoint(x,y), color, false))
    return g.Board.PlayMove(x,y,color)
}

func (g *Game) PlayPass(color Color) {
    g.sequence.Push(*NewMove(*NewPoint(0,0), color, true))
    g.Board.PlayPass(color)
}

// Plays out the sequence 'seq' of moves
func (g *Game) PlaySequence(seq []Move) {
    for _, m := range seq {
        if m.Vertex.Pass {
            g.PlayPass(m.Color)
        } else {
            g.PlayMove(m.Vertex.X, m.Vertex.Y, m.Color)
        }
    }
}

func (g *Game) SetKomi(newKomi float) {
    g.komi = newKomi
}


// ##################### Game helper functions ##########################
func NewGame(boardsize int) *Game {
    return &Game{ Board: NewBoard(boardsize),
                  komi: defaultKomi,
                }
}


