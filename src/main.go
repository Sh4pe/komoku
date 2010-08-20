/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */
package main

import (
    "fmt"
    //"container/list"
    "./komoku/komoku"
)

const (
    a = iota-1;
    b;
    c;
)

func main() {
    b := komoku.NewBoard()
    b.IsLegalMove(1,1,komoku.White)
    b.CreateGroup(3,3,komoku.White)
    b.CreateGroup(2,2,komoku.Black)

    x,y := 8,8
    nFree, adjBlack, adjWhite := b.GetEnvironment(x,y)

    fmt.Printf("Board of size %d\n", komoku.BoardSize)
    komoku.PrintBoard(b)

    fmt.Println()
    fmt.Printf("At (%d,%d):\n\n", x,y)
    fmt.Printf("nFree: %d,\nadjBlack: %v\nadjWhite: %v\n", nFree, adjBlack, adjWhite)

}
