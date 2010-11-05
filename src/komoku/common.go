/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */

package komoku

import (
    "os"
    "path"
    "strings"
    "fmt"
    "strconv"
)

// ################ constants ####################
const (
    DefaultBoardSize =  9 // This should be less than 25, because ui.go.PrintBoard and the GTP protocol 
                          // will have problems otherwise....
    DefaultKomi = 6.5
    komokuVersion = "0.1a"
    komokuProgramName = "komoku"
)

// komokus error constants
const (
    ErrFieldOccupied = iota;
    ErrIllegalMove;
    ErrInvalidCoordinateChar;
    ErrInvalidCoordinateDigit;
    ErrUnacceptableBoardSize;
    ErrGTPSyntaxError;
    ErrIOError;
    ErrFieldLegalityCheckedMoreThanOnce;
    ErrGTPNotImplemented;
    ErrGTPIllegalCommand;
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

type Vertex struct {
    X, Y int
    Pass bool
}

func NewVertex(p Point, pass bool) *Vertex {
    return &Vertex{ X: p.X, Y: p.Y, Pass: pass }
}

func NewVertexByInts(x, y int , pass bool) *Vertex {
    return &Vertex{ X: x, Y: y, Pass: pass }
}

type Move struct {
    Color
    Vertex
}

func NewMove(p Point, c Color, pass bool) *Move {
    return &Move{ Vertex: *NewVertex(p, pass), Color: c }
}

func NewMoveByVertex(v *Vertex, c Color) *Move {
    return &Move{ Vertex: *v, Color: c }
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

func NewIOError(err os.Error) Error {
    return NewError(err.String(), ErrIOError)
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
// ################ own conversions ##########
// ###########################################
// TODO: get rid of this or make it less ugly at least
var charDigit = map[string]int {
    "A": 0, "B": 1, "C": 2, "D": 3, "E": 4, "F": 5, "G": 6, "H": 7, "J": 8, "K": 9, "L": 10,
    "M": 11, "N": 12, "O": 13, "P": 14, "Q": 15, "R": 16, "S": 17, "T": 18, "U": 19, "V": 20, "W": 21,
    "X": 22, "Y": 23, "Z": 24,
}

var coordinateChars string = "ABCDEFGHJKLMNOPQRSTUVWXYZ"

// TODO: there must be a better way to do this...
func CharToDigit(c string) (digit int, err Error) {
    digit, ok := charDigit[c]
    if !ok {
        return -1, NewInvalidCoordinateCharError(c)
    }
    return
}

// TODO: ...and for this also...
func DigitToChar(digit int) (char string, err Error) {
    if digit < 0 || digit > len(coordinateChars) {
        return "", NewInvalidCoordinateDigitError(digit)
    }
    return coordinateChars[digit:digit+1], err
}

func NewInvalidCoordinateCharError(char string) (err Error) {
    return NewError(fmt.Sprintf("'%s' is not a valid character for a coordinate", char), ErrInvalidCoordinateChar)
}

func NewInvalidCoordinateDigitError(digit int) (err Error) {
    return NewError(fmt.Sprintf("%d is not a valid coordinate digit", digit), ErrInvalidCoordinateDigit)
}

// ###########################################
// ################ helper functions #########
// ###########################################


// Returns true if (x,y) is a hoshi point. 
// Currently, the only supported board sizes are 9, 13 and 19. 
// For other sizes, this function returns false for all (x,y).
// Right now, this function is only used by printing functions in ui.go, so this
// is not a 'very important' function
func isHoshi(x, y, boardsize int) bool {
    switch boardsize {
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

// If c is one of {white,w}, color is White. If c is one of {black,b}, color is Black. In both cases, ok is true
// c is not treated case sensitive. If c is something else, color is meaningless and ok is false
func gtpColorToColor(c string) (color Color, ok bool) {
    c = strings.ToLower(c)
    if c == "w" || c == "white" {
        return White, true
    }
    if c == "b" || c == "black" {
        return Black, true
    }
    return color, false
}

func colorToGTPColor(c Color) string {
    if c == Black {
        return "B"
    }
    return "W"
}

// c is case insensitive. A vertex in the GTP spec is something like "B13" or "a2" or "pass". If it is "pass",
// 'pass' and 'ok' are true and 'point' is meaningless. If it is a coordinate, 'point' points there and 'ok' is
// true. If it is something else, 'ok' is false and 'point' and 'pass' are meaningless.

// TODO: charDigit from ui.go should be in this file...
// TODO: write tests for this
func gtpVertexToPoint(c string) (point Point, okay, pass bool) {
    c = strings.ToUpper(c)
    if c == "PASS" {
        return Point{0,0}, true, true
    }
    if len(c) > 1 {
        if x, ok := charDigit[c[0:1]]; ok {
            if y, err := strconv.Atoi(c[1:len(c)]); err == nil {
                //fmt.Printf("gtpVertexToPoint: coords: (%d,%d)\n", x,y)
                return Point{ X: x, Y: y-1 }, true, false
            }
        }
    }
    return Point{0,0}, false, false
}

// Returns strings as "E5"...
func pointToGTPVertex(p Point) (ret string, ok bool) {
    ret = ""
    digit, err := DigitToChar(p.X)
    if err != nil {
        return "", false
    }
    ret += digit
    ret += fmt.Sprintf("%d",p.Y + 1)
    return ret, true
}

// Rel is a path relative to the executable which is run. This func returns the absolute path.
func relPathToAbs(rel string) string {
    wd, _ := os.Getwd()
    //fmt.Println("wdir", wd)
    execDir := os.Getenv("_")
    //fmt.Println("execDir", execDir)
    base := path.Base(execDir)
    execDir = execDir[0:len(execDir)-len(base)]
    fname := wd + "/" + execDir + "/" + rel
    // zsh tweak (possibly for other shells too?)
    if strings.Index(execDir, wd) != -1 {
        fname = execDir + "/" + rel
    }
    return path.Clean(fname)
}

