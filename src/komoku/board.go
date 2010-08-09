/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */
package komoku

import (
    "fmt"
    //"os"
)

/*
 * TODO:
 *      - is FieldIndices.Sequence needed/useful?
 *      - let this talk the GTP-protocol
 */


// ################################################################################
// ########################### Board struct #######################################
// ################################################################################

// This object is responsible for recording a current state of a game.
type Board struct {
    fields [BoardSize*BoardSize]Field
    legalBlackMoves IntList // Indices of fields at which it is legal to play a black stone.
    legalWhiteMoves IntList // Indices of fields at which it is legal to play a white stone.
    emptyFields IntList // indices of empty fields
}

// ##################### Board methods ##########################


// Returns a copy of the field at (x,y)
func (b *Board) GetField(x,y int) Field {
    //fmt.Printf("(%d,%d), pos: %d\n",x,y,xyToPos(x,y))
    return b.fields[xyToPos(x,y)]
}

// Is it legal to play a stone of color 'color' at (x,y)?
func (b *Board) IsLegalMove(x, y int, color Color) bool {
    // TODO: write test for this!
    var indices *IntList = &b.legalBlackMoves
    if color == White {
        indices = &b.legalWhiteMoves
    }
    pos := xyToPos(x,y)
    last := indices.Last()
    for it := indices.First(); it != last; it = it.Next() {
        if it.Value() == pos {
            return true
        }
    }
    return false
}

// Play a stone of color 'color'  at (x,y). If an error occurs (such as that this 
// place is already occupied) this error is returned.
func (b *Board) Move(x, y int, color Color) (err Error) {
    index := xyToPos(x,y)
    if !b.fields[index].Empty() {
        return NewFieldOccupiedError(x,y)
    }
    if !b.IsLegalMove(x,y, color) {
        return NewIllegalMoveError(x,y, color)
    }
    // TODO: not yet done at all!
    // Go on here!
    return
}

// Removes a stone at (x,y), if there is any and updates b.emptyFields.
// This method does not alter legalWhiteMoves or legalBlackMoves
func (b *Board) removeStone(x, y int) {
    pos := xyToPos(x,y)
    if !b.fields[pos].Empty() {
        b.fields[pos].Clear()
        b.emptyFields.Append(pos)
    }
}

func (b* Board) updateLegalMoves() {
    // This method assumes that b.emptyFields is correctly set.
    b.legalWhiteMoves.Clear()
    b.legalBlackMoves.Clear()
    last := b.emptyFields.Last()
    //for i := 0; i < b.emptyFields.Length(); i++ {
    for it := b.emptyFields.First(); it != last; it = it.Next() {
        //pos := b.emptyFields.Get(i)
        pos := it.Value()
        x, y := posToXY(pos)
        nbours := neighbours(x,y)
        freeNBours := 0
        // how many are free?
        for _, p := range nbours {
            if b.fields[xyToPos(p.X, p.Y)].Empty() {
                freeNBours++
            }
        }
        switch {
            case freeNBours > 0: // this is a legal move for both colors
                b.legalWhiteMoves.Append(pos)
                b.legalBlackMoves.Append(pos)
            case freeNBours == 0: // this case is more difficult
                return
        }
    }
}

// ##################### Board helper functions ##########################
// Creates a new, initial board
func NewBoard() *Board {
    ret := &Board{ legalWhiteMoves: *NewIntList(),
                   legalBlackMoves: *NewIntList(),
                   emptyFields: *NewIntList(),
                 }
    for i := 0; i < BoardSize*BoardSize; i++ {
        ret.legalWhiteMoves.Append(i)
        ret.legalBlackMoves.Append(i)
        ret.emptyFields.Append(i)
    }
    return ret
}

// Creates a new FieldOccupiedError, indicating that (x,y) is alrady used.
func NewFieldOccupiedError(x, y int) (err Error) {
    return NewError(fmt.Sprintf("(%d,%d) is already occupied", x, y), ErrFieldOccupied)
}

// Creates a new FieldOccupiedError, indicating that (x,y) is alrady used.
func NewIllegalMoveError(x, y int, color Color) (err Error) {
    return NewError(fmt.Sprintf("a %v move at (%d,%d) is illegal", color, x, y), ErrIllegalMove)
}


