package komoku

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


// ############### Board struct ###################

// This object is responsible for recording a current state of a game.
type Board struct {
    fields [BoardSize*BoardSize]Field
    playableFields []int // indices of currently playable fields
    emptyFields []int // indices of empty fields
}

// ##################### Board methods ##########################



// ##################### Board helper functions ##########################
// Creates a new, initial board
func NewBoard() *Board {
    ret := &Board{ playableFields: make([]int, BoardSize*BoardSize, BoardSize*BoardSize),
                   emptyFields: make([]int, BoardSize*BoardSize, BoardSize*BoardSize),
                 }
    for i := 0; i < BoardSize*BoardSize; i++ {
        ret.playableFields[i] = i
        ret.emptyFields[i] = i
    }
    return ret
}

// methods
