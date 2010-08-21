/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */
package main

import (
    //"fmt"
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
    b.PlayMove(3,3)
    b.PlayMove(3,4)
    b.PlayMove(2,4)
    b.PlayMove(8,8)
    b.PlayMove(3,5)
    b.PlayMove(9,9)
    b.PlayMove(4,4)

    komoku.PrintBoard(b)
}
