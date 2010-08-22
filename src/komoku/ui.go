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

// ################################################################################
// ####################### Functions for printing boards #########################$
// ################################################################################

func PrintBoard(b *Board) {
    fmt.Printf(printBoardPrimitive(b, " ", -1, -1))
}

// this does the actual work
// kind of ugly...
// returns a string
// TODO: cache the printed lines so that error output is never mixed with regular output.
// TODO: make this GTP-compatible
func printBoardPrimitive(b *Board,
                         leftOffset string, // printed on the very beginning of each line
                         lastX, lastY int, // Marks the last played move. Negative values indicate that no last move is provided.
                        ) (s string) {
    // so that the values of lastX, lastY don't interfere with our algorithm...
    // Now we can assume that lastX, lastY are valid or 'too small'
    if lastX < 0 || lastY < 0 {
        lastX = -10
        lastY = -10
    }

    // print coordinates at the header
    line := leftOffset
    line += " "
    if BoardSize > 9 {
        line += " "
    }
    for i := 0; i < BoardSize; i++ {
        char, err := DigitToChar(i)
        if err != nil {
            fmt.Printf("Error in printBoardPrimitive. Error:\n%s\n", err)
            return
        }
        line += " " + char
    }
    s += line + "\n"
    // print the board
    //for y := 0; y < BoardSize; y++ {
    for y := BoardSize-1; y >= 0; y-- {
        lineNumber := fmt.Sprintf("%d", y+1)
        if BoardSize > 9 {
            lineNumber = fmt.Sprintf("%2d", y+1)
        }
        line = leftOffset + lineNumber
        rightSpace := ""
        for x := 0; x < BoardSize; x++ {
            //fmt.Printf("(%d,%d)\n",x,y)
            //field := b.GetField(x,y)
            empty, group := b.GetGroup(x,y)
            fieldChar := ""
            //fmt.Printf("empty: %v, group: %v\n", empty, group)
            if empty {
                // is this a hoshi?
                if isHoshi(x,y) {
                    fieldChar = "+"
                } else {
                    fieldChar = "."
                }
            } else if group.Color == Black {
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
        line += lineNumber
        s += line + "\n"
    }
    // and print the footer
    line = leftOffset
    line += " "
    if BoardSize > 9 {
        line += " "
    }
    for i := 0; i < BoardSize; i++ {
        char, err := DigitToChar(i)
        if err != nil {
            fmt.Printf("Error in printBoardPrimitive. Error:\n%s\n", err)
            return
        }
        line += " " + char
    }
    s += line + "\n"
    return s
}

// ###############################################################################
// ####################### The interactive mode ##################################
// ###############################################################################

func InteractiveMode() {
    return
}

