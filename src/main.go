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

    var a komoku.FieldIndices = make(komoku.FieldIndices,10,10)
    for i := 0; i < len(a); i++ {
        a[i] = i
    }
    fmt.Printf("%v\n", a)
    a.Remove(3)
    fmt.Printf("%v\n", a)
}
