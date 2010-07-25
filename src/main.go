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
    fmt.Printf("%d\n\n\n", komoku.BoardSize)

    a := komoku.NewFieldIndices(10,10)
    for i := 0; i < 10; i++ {
        a.Set(i,i)
    }
    fmt.Printf("%v\n", a)

    for i := 0; i < 10; i++ {
        a.Set(i,i*i)
    }
    fmt.Printf("%v\n", a)
}
