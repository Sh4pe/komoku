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
    fmt.Printf("%d\n", komoku.BoardSize)
    fmt.Printf("%d %d\n\n\n", komoku.ErrFieldOccupied, komoku.ErrIllegalMove)

    b := komoku.NewBoard()
    b.LegalMove(1,1,komoku.White)

    err := komoku.NewIllegalMoveError(1,1,komoku.White)
    fmt.Printf("%v\n", err)
}
