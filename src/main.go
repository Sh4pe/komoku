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

    fi := komoku.NewFieldIndices(3)
    fi.Append(5)
    fmt.Printf("%d\n", fi.Length())

    tmp := make([]int, 10, 50)
    fmt.Printf("len: %d, cap: %d\n", len(tmp), cap(tmp))
}
