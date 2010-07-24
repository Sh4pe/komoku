package komoku

import (
    //"fmt" // for debugging
)

const (
    fieldWhite = iota - 1;
    fieldEmpty;
    fieldBlack
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

// Create empty field.
func (f *Field) NewField() (ret *Field) {
    return
}

func (f *Field) NewFieldBlack() *Field {
    return &Field{ fieldBlack }
}

func (f *Field) NewFieldWhite() *Field {
    return &Field{ fieldWhite }
}

// ############### FieldIndices ##############

// Storage type for indices of fields. It is assumed that each index occures at most
// once in a FieldIndices.

// TODO: do we want to export this?
type FieldIndices []int

// Removes 'index' from a FieldIndices. The rest remains unchanged.
func (fi *FieldIndices) Remove(index int) {
    // This whole method doesn't work if an index might occur more than once.
    jump := 0
    length := len(*fi)
    for i := 0; i < length - jump; i++ {
        //fmt.Printf("i: %d, jump: %d, length: %d\n", i, jump, length)
        if (*fi)[i] == index {
            jump++
        }
        (*fi)[i] = (*fi)[i+jump]
    }
    *fi = (*fi)[0:length-jump]
}

// ############### Board struct ###################

// This object is responsible for recording a current state of a game.
type Board struct {
    fields [BoardSize*BoardSize]Field
    playableFields FieldIndices // indices of currently playable fields
    emptyFields FieldIndices // indices of empty fields
}

// ##################### Board methods ##########################

// ##################### Board helper functions ##########################
// Creates a new, initial board
func NewBoard() *Board {
    ret := &Board{ playableFields: make(FieldIndices, BoardSize*BoardSize, BoardSize*BoardSize),
                   emptyFields: make(FieldIndices, BoardSize*BoardSize, BoardSize*BoardSize),
                 }
    for i := 0; i < BoardSize*BoardSize; i++ {
        ret.playableFields[i] = i
        ret.emptyFields[i] = i
    }
    return ret
}

// methods
