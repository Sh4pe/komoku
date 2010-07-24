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

    fmt.Printf("%d %d %d\n\n", a,b,c)
    a := [...]int{1,2,3,4,5,6}
    for k, v := range a {
        fmt.Printf("k=%d v=%d\n", k, v)
    }

    fmt.Println("")
    b := a[0:3]
    for k, v := range b {
        fmt.Printf("k=%d v=%d\n", k, v)
    }
}
