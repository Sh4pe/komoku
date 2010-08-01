package komoku

import (
    "fmt"
    //"os"
)

/*
 * TODO:
 *      - is FieldIndices.Sequence needed/useful?
 */


type Color bool
const (
    White = true
    Black = false
)

func (c Color) String() string {
    if c == White {
        return "white"
    }
    return "black"
}

// ############### Field struct ###############
type Field struct {
    value int8
}

const (
    fieldWhite = iota - 1;
    fieldEmpty;
    fieldBlack
)

// Is this field occupied by a black stone?
func (f *Field) Black() bool {
    return f.value == fieldBlack
}

// makes f Empty()
func (f *Field) Clear() {
    f.value = fieldEmpty
}

// Is this field empty?
func (f *Field) Empty() bool {
    return f.value == fieldEmpty
}


// Is this field occupied by a white stone?
func (f *Field) White() bool {
    return f.value == fieldWhite
}

// ################### helpers ################

// Create empty field.
func NewField() (ret *Field) {
    return
}

func NewFieldBlack() *Field {
    return &Field{ fieldBlack }
}

func NewFieldWhite() *Field {
    return &Field{ fieldWhite }
}

// ################################################################################
// ########################### FieldIndices struct ################################
// ################################################################################

// TODO: Do we want to export this?
type SequenceType uint64

// Storage type for indices of fields. It is assumed that each index occures at most
// once in a FieldIndices.
// TODO: do we want to export this?
type FieldIndices struct {
    indices []int
    Sequence SequenceType // Used to check if the indices are up to date
    topIndex int // points to the last element
}

// ######################## methods ###############################

// Appends i to the FieldIndices
func (fi *FieldIndices) Append(i int) {
    // is fi.indices big enough?
    if fi.topIndex >= cap(fi.indices) {
        fi.grow()
    }
    // adjust length of the slice
    fi.indices = fi.indices[0:fi.topIndex+1]
    fi.indices[fi.topIndex] = i
    fi.topIndex++
    fi.Sequence++
}

func (fi *FieldIndices) Capacity() int {
    return cap(fi.indices)
}

// Empties a FieldIndices entirely
func (fi *FieldIndices) Clear() {
    fi.indices = fi.indices[0:0]
    fi.topIndex = 0
    fi.Sequence++
}

// Returns index-th element...
func (fi *FieldIndices) Get(index int) int {
    return fi.indices[index]
}

// If the slice fi.indices gets too small, this function lets it grow.
func (fi *FieldIndices) grow() {
    const newElements = 10
    newIndices := make([]int, fi.topIndex, fi.topIndex + newElements)
    copy(newIndices, fi.indices)
    fi.indices = newIndices
    // TODO: fi.Sequence++ here?
}

// Returns the length of a FieldIndices.
func (fi *FieldIndices) Length() int {
    //return len(fi.indices)
    return fi.topIndex
}

// Removes 'index' from a FieldIndices. The rest remains unchanged.
func (fi *FieldIndices) Remove(index int) {
    // This whole method doesn't work if an index might occur more than once.
    jump := 0
    length := fi.Length()
    for i := 0; i < length - jump; i++ {
        if fi.indices[i] == index {
            jump++
            if i+jump >= length {
                break
            }
        }
        //fmt.Printf("i: %d, jump: %d, length: %d\n", i, jump, length)
        fi.indices[i] = fi.indices[i+jump]
    }
    fi.indices = fi.indices[0:length-jump]
    fi.topIndex = length-jump
    fi.Sequence++
}

// Implemented so that FieldIndices implements Stringer interface. 
func (fi *FieldIndices) String() string {
    return fmt.Sprintf("%v", fi.indices)
}

// ########################### helper functions ###################################

// c is the capacity of the 'FieldIndices'
func NewFieldIndices(c int) *FieldIndices {
    return &FieldIndices{ indices: make([]int, 0, c),
                        }
}

// ################################################################################
// ########################### Board struct #######################################
// ################################################################################

// This object is responsible for recording a current state of a game.
type Board struct {
    fields [BoardSize*BoardSize]Field
    legalBlackMoves FieldIndices // Indices of fields at which it is legal to play a black stone.
    legalWhiteMoves FieldIndices // Indices of fields at which it is legal to play a white stone.
    emptyFields FieldIndices // indices of empty fields
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
    var indices *FieldIndices = &b.legalBlackMoves
    if color == White {
        indices = &b.legalWhiteMoves
    }
    pos := xyToPos(x,y)
    for i := 0; i < indices.Length(); i++ {
        if indices.Get(i) == pos {
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
    for i := 0; i < b.emptyFields.Length(); i++ {
        pos := b.emptyFields.Get(i)
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
    ret := &Board{ legalWhiteMoves: *NewFieldIndices(BoardSize*BoardSize),
                   legalBlackMoves: *NewFieldIndices(BoardSize*BoardSize),
                   emptyFields: *NewFieldIndices(BoardSize*BoardSize),
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


