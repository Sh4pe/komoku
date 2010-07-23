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
    fmt.Println(komoku.BoardSize)

    fmt.Printf("%d %d %d\n", a,b,c)
}
