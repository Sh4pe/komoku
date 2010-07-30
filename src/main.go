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
    b.LegalMove(1,1,komoku.White)

    fmt.Printf("Board of size %d\n", komoku.BoardSize)
    komoku.PrintBoard(b)

}
