package komoku

import (
    "os"
)

// ################ constants ####################
const (
    BoardSize = 19 // we support only quadratic boards at the moment
)

// komokus error constants
const (
    ErrFieldOccupied = iota;
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
    return 19*y + x
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
