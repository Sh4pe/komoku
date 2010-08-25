/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */


// Code for the (text)UI comes here
// The output these routines generate is inspiered by the output gnugo prints.


/*
 * TODO:
 *      - move charDigit et al. to common.go
 */
package komoku

import (
    "fmt"
)

// ################################################################################
// ####################### Functions for printing boards #########################$
// ################################################################################

func PrintBoard(b *Board) {
    fmt.Printf(printBoardPrimitive(b, " ", -1, -1, nil))
}

// this does the actual work
// kind of ugly...
// returns a string
func printBoardPrimitive(b *Board,
                         leftOffset string, // printed on the very beginning of each line
                         lastX, lastY int, // Marks the last played move. Negative values indicate that no last move is provided.
                         marks []Point, // nil means "no marks". Marked points will be displayed as an !
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
    if b.BoardSize() > 9 {
        line += " "
    }
    for i := 0; i < b.BoardSize(); i++ {
        char, err := DigitToChar(i)
        if err != nil {
            fmt.Printf("Error in printBoardPrimitive. Error:\n%s\n", err)
            return
        }
        line += " " + char
    }
    s += line + "\n"
    // print the board
    //for y := 0; y < b.BoardSize(); y++ {
    for y := b.BoardSize()-1; y >= 0; y-- {
        lineNumber := fmt.Sprintf("%d", y+1)
        if b.BoardSize() > 9 {
            lineNumber = fmt.Sprintf("%2d", y+1)
        }
        line = leftOffset + lineNumber
        rightSpace := ""
        for x := 0; x < b.BoardSize(); x++ {
            //fmt.Printf("(%d,%d)\n",x,y)
            //field := b.GetField(x,y)
            empty, group := b.GetGroup(x,y)
            fieldChar := ""
            //fmt.Printf("empty: %v, group: %v\n", empty, group)
            inMarks := false
            for _, p := range marks {
                if p.X == x && p.Y == y {
                    inMarks = true
                    break
                }
            }
            if inMarks {
                fieldChar = "!"
            } else if empty {
                if isHoshi(x,y, b.BoardSize()) {
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
    if b.BoardSize() > 9 {
        line += " "
    }
    for i := 0; i < b.BoardSize(); i++ {
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

