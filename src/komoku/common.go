package komoku

import (
    "os"
)

// ################ constants ####################
const (
    BoardSize = 13 // ...says that we are playing on a (BoardSize * BoardSize) - board
                   // We support only quadratic boards at the moment.
                   // This should be less than 26, because ui.go.PrintBoard will have problems otherwise....
)

// komokus error constants
const (
    ErrFieldOccupied = iota;
    ErrIllegalMove;
    ErrInvalidCoordinateChar;
    ErrInvalidCoordinateDigit;
)

// ################ interfaces ##############
type Error interface {
    os.Error
    Errno() int
}

// ################ types ####################

type Point struct {
    x, y int
}

type komokuError struct {
    msg string
    errno int
}

// komokuError has to implement os.Error
func (e *komokuError) String() string {
    return e.msg
}

// ...and it has to implement the Error interface
func (e *komokuError) Errno() int {
    return e.errno
}

// Create a new Error
func NewError(s string, errno int) Error {
    return &komokuError{ s, errno }
}

// ################ helper functions ####################

// TODO: use point!
func posToXY(pos int) (x, y int) {
    return pos%BoardSize, pos/BoardSize
}

// TODO: use point!
func xyToPos(x, y int) int {
    return BoardSize*y + x
}

// Returns the neighbours of a field (x,y)
func neighbours(x, y int) []Point {
    // TODO: can this be implemented better?
    ret := make([]Point, 4)
    count := 0
    switch x {
        case 0:
            ret[count] = Point{ 1, y }
            count++
        case BoardSize-1:
            ret[count] = Point{ BoardSize-2, y }
            count++
        default:
            ret[count] = Point{ x-1, y }
            count++
            ret[count] = Point{ x+1, y }
            count++
    }
    switch y {
        case 0:
            ret[count] = Point{ x, 1 }
            count++
        case BoardSize-1:
            ret[count] = Point{ x, BoardSize-2 }
            count++
        default:
            ret[count] = Point{ x, y-1 }
            count++
            ret[count] = Point{ x, y+1 }
            count++
    }
    return ret[0:count]
}

// returns true if (x,y) is a hoshi point. 
// Currently, the only supported board size are 19. 
// For other sizes, this function returns false for all (x,y).
// Right now, this function is only used by printing functions in ui.go, so this
// is not a 'very important' function
func isHoshi(x, y int) bool {
    switch BoardSize {
        case 19:
            if (x == 3) || (x == 9) || (x == 15)  {
                return (y == 3) || (y == 9) || (y == 15)
            }
            return false
        case 13:
            if (x == 3) || (x  == 9) {
                return (y == 3) || (y == 9)
            }
            if (x == 6) && (y == 6) {
                return true
            }
        case 9:
            if (x == 2) || (x == 6) {
                return (y == 2) || (y == 6)
            }
    }
    return false
}
