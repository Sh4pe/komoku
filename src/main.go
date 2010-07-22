package main

import (
    "fmt"
    //"./common"
    //"./treenode"
    "./komoku/komoku"
)

func main() {
    t := komoku.NewRootTreeNode()
    fmt.Printf("%v\n", t)
    a := 0*19 + 5
    fmt.Println(a)
    fmt.Printf("%d = (%d, %d)\n", a, a/19, a%19)
}
