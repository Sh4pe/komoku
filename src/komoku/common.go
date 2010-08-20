/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */

/*
 * TODO:
 *      - This whole file is not very tidy. Clean it up.
 */
package komoku

import (
    "os"
)

// ################ constants ####################
const (
    BoardSize = 19 // ...says that we are playing on a (BoardSize * BoardSize) - board
                   // We support only quadratic boards at the moment.
                   // This should be less than 25, because ui.go.PrintBoard will have problems otherwise....
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

// ###########################################
// ################ types ####################
// ###########################################

type Point struct {
    X, Y int
}

func NewPoint(x, y int) *Point {
    return &Point{ X: x, Y: y }
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


// ################## type color #####################
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

// ###########################################
// ################ helper functions #########
// ###########################################

// TODO: use point!
func posToXY(pos int) (x, y int) {
    return pos%BoardSize, pos/BoardSize
}

// TODO: use point!
func xyToPos(x, y int) int {
    return BoardSize*y + x
}

// Returns the neighbours of a field (x,y) as a slice of Points.
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

// Returns true if (x,y) is a hoshi point. 
// Currently, the only supported board sizes are 9, 13 and 19. 
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
