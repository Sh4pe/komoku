/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */


// Code for the (text)UI comes here
// The output these routines generate is inspiered by the output gnugo prints.
package komoku

import (
    "fmt"
)

var charDigit = map[string]int {
    "a": 0, "b": 1, "c": 2, "d": 3, "e": 4, "f": 5, "g": 6, "h": 7, "j": 8, "k": 9, "l": 10,
    "m": 11, "n": 12, "o": 13, "p": 14, "q": 15, "r": 16, "s": 17, "t": 18, "u": 19, "v": 20, "w": 21,
    "x": 22, "y": 23, "z": 24,
}

var coordinateChars string = "abcdefghjklmnopqrstuvwxyz"

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

// ################################################################################
// ####################### Functions for printing boards #########################$
// ################################################################################

func PrintBoard(b *Board) {
    printBoardPrimitive(b, "", -1, -1)
}

// this does the actual work
// kind of ugly...
// TODO: cache the printed lines so that error output is never mixed with regular output.
func printBoardPrimitive(b *Board,
                         leftOffset string, // printed on the very beginning of each line
                         lastX, lastY int, // Marks the last played move. Negative values indicate that no last move is provided.
                        ) {
    // so that the values of lastX, lastY don't interfere with our algorithm...
    // Now we can assume that lastX, lastY are valid of 'too small'
    if lastX < 0 || lastY < 0 {
        lastX = -10
        lastY = -10
    }

    // print coordinates at the header
    line := leftOffset
    line += " "
    for i := 0; i < BoardSize; i++ {
        char, err := DigitToChar(i)
        if err != nil {
            fmt.Printf("Error in printBoardPrimitive. Error:\n%s\n", err)
            return
        }
        line += " " + char
    }
    fmt.Println(line)
    // print the board
    for y := 0; y < BoardSize; y++ {
        char, err := DigitToChar(y)
        if err != nil {
            fmt.Printf("Error in printBoardPrimitive. Error:\n%s\n", err)
            return
        }
        line = leftOffset + char
        rightSpace := ""
        for x := 0; x < BoardSize; x++ {
            //fmt.Printf("(%d,%d)\n",x,y)
            field := b.GetField(x,y)
            fieldChar := ""
            if field.Empty() {
                // is this a hoshi?
                if isHoshi(x,y) {
                    fieldChar = "+"
                } else {
                    fieldChar = "."
                }
            } else if field.Black() {
                fieldChar = "X"
            } else {
                // so the field must be occupied by a white stone...
                fieldChar = "O"
            }
            leftSpace := " "
            rightSpace = ""
            if y == lastY {
                if x == lastX {
                    leftSpace = "("
                    rightSpace = ")"
                } else if x-1 == lastX {
                    leftSpace = ""
                }
            }
            line += leftSpace + fieldChar + rightSpace
        }
        if rightSpace != ")" {
            line += " "
        }
        line += char
        fmt.Println(line)
    }
    // and print the footer
    line = leftOffset
    line += " "
    for i := 0; i < BoardSize; i++ {
        char, err := DigitToChar(i)
        if err != nil {
            fmt.Printf("Error in printBoardPrimitive. Error:\n%s\n", err)
            return
        }
        line += " " + char
    }
    fmt.Println(line)
}
