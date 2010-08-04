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

    fmt.Printf("Board of size %d\n", komoku.BoardSize)
    komoku.PrintBoard(b)

    fmt.Println("")

    il := komoku.NewIntList()
    fmt.Printf("len: %d\n", il.Length())
    for i := 0; i < 10; i++ {
        il.Append(i)
    }
    fmt.Printf("len: %d\n", il.Length())
    last := il.Last()
    for it := il.First(); it != last; it = it.Next() {
        fmt.Printf("%d\n", it.Value())
    }
    il.Remove(5)
    fmt.Printf("len: %d\n", il.Length())
    last = il.Last()
    for it := il.First(); it != last; it = it.Next() {
        fmt.Printf("%d\n", it.Value())
    }
}
