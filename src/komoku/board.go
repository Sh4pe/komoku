package komoku

import (
    "fmt"
    //"os"
)

type Color bool
const (
    White = true
    Black = false
)

// ############### Field struct ###############
type Field struct {
    value int8
}

const (
    fieldWhite = iota - 1;
    fieldEmpty;
    fieldBlack
)

// Is this field empty?
func (f *Field) Empty() bool {
    return f.value == fieldEmpty
}

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
}

// Returns the length of a FieldIndices.
func (fi *FieldIndices) Length() int {
    return len(fi.indices)
}

// Removes 'index' from a FieldIndices. The rest remains unchanged.
func (fi *FieldIndices) Remove(index int) {
    // This whole method doesn't work if an index might occur more than once.
    jump := 0
    length := fi.Length()
    for i := 0; i < length - jump; i++ {
        //fmt.Printf("i: %d, jump: %d, length: %d\n", i, jump, length)
        if fi.indices[i] == index {
            jump++
        }
        fi.indices[i] = fi.indices[i+jump]
    }
    fi.indices = fi.indices[0:length-jump]
}

// Returns index-th element...
func (fi *FieldIndices) Get(index int) int {
    return fi.indices[index]
}

// Returns index-th element...
func (fi *FieldIndices) Set(index, value int) {
    fi.Sequence++
    fi.indices[index] = value
}

// Implemented so that FieldIndices implements Stringer interface. 
func (fi *FieldIndices) String() string {
    return fmt.Sprintf("%v", fi.indices)
}

// ########################### helper functions ###################################

// a, b as in make([]int, a, b)...
func NewFieldIndices(a, b int) *FieldIndices {
    return &FieldIndices{ indices: make([]int, a, b),
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

// Play a stone of color 'color'  at (x,y). If an error occurs (such as that this 
// place is already occupied) this error is returned.
func (b *Board) Move(x, y int, color Color) (err Error) {
    index := xyToPos(x,y)
    if !b.fields[index].Empty() {
        return FieldOccupiedError(x,y)
    }
    // TODO: not yet done at all!
    return
}

// Is it legal to play a stone of color 'color' at (x,y)?
func (b *Board) LegalMove(x, y int, color Color) bool {
    return true
}

// ##################### Board helper functions ##########################
// Creates a new, initial board
func NewBoard() *Board {
    ret := &Board{ legalWhiteMoves: *NewFieldIndices(BoardSize*BoardSize, BoardSize*BoardSize),
                   legalBlackMoves: *NewFieldIndices(BoardSize*BoardSize, BoardSize*BoardSize),
                   emptyFields: *NewFieldIndices(BoardSize*BoardSize, BoardSize*BoardSize),
                 }
    for i := 0; i < BoardSize*BoardSize; i++ {
        ret.legalWhiteMoves.Set(i,i)
        ret.legalBlackMoves.Set(i,i)
        ret.emptyFields.Set(i,i)
    }
    return ret
}

// Creates a new FieldOccupiedError, indicating that (x,y) is alrady used.
func FieldOccupiedError(x, y int) (err Error) {
    return NewError(fmt.Sprintf("(%d,%d) is already occupied", x, y), ErrFieldOccupied)
}

// methods
