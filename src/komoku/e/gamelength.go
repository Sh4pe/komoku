/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */

// Here are experiments to determine the average simulated game lengths

package main

import (
    "fmt"
    "./komoku"
)

const (
    numTestGames = 100;
)

func averageGameLength(boardsize int) {
    board := komoku.NewBoard(boardsize)
    sum := 1
    lastPass := false
    for i := 0; i < numTestGames; i++ {
        for {
            v := board.PlayRandomMove(board.ColorOfNextPlay())
            if v.Pass {
                if lastPass {
                    break
                }
                lastPass = true
            } else {
                lastPass = false
            }
            sum++
        }
        board.Reset()
    }
    //fmt.Printf("played %d simulations on boardsize %d, average number of moves per simulation: %2.2f\n", numTestGames, boardsize, float(sum)/float(numTestGames))
    fmt.Printf("averageGameLength[%2d] = %2.0f\n", boardsize, float(sum)/float(numTestGames))
}

func main() {
    for i := 5; i < 26; i++ {
        averageGameLength(i)
    }
}

