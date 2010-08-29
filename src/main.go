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
    "sort"
    "container/vector"
    "./komoku/komoku"
    "path"
    "os"
)

const (
    a = iota-1;
    b;
    c;
)

type father interface{
    Talk() string
}

type child interface{
    father
    Wee() string
}

type fimpl struct {
    dummyf int
}
func (i *fimpl) Talk() string {
    return "Father talks"
}

type cimpl struct {
    dummyc float
}
func (i *cimpl) Wee() string {
    return "Child wees"
}
func (i *cimpl) Talk() string {
    return "Child talks"
}

func eFather(f father) {
    fmt.Println(f.Talk())
}

func eChild(c child) {
    fmt.Println(c.Wee())
}

func childAsFather(c child) {
    fmt.Println("in childAsFather")
    if f, ok := c.(father); !ok {
        fmt.Println("type assertion failed")
    } else {
        fmt.Println(f.Talk())
        fatherDowntoChild(f)
    }
}

func fatherDowntoChild(f father) {
    if c, ok := f.(child); ok {
        fmt.Println("father to child worked")
        fmt.Println(c.Wee())
    } else {
        fmt.Println("father to child didn't work")
    }
}

func generalFunc(a interface{}) {
    if c, ok := a.(child); ok {
        fmt.Printf("type assertion okay, %s\n", c.Wee())
    } else {
        fmt.Println("type assertion not okay")
    }
}

func testMain() {
    var f *fimpl
    var c *cimpl
    childAsFather(c)
    fmt.Println()
    fatherDowntoChild(f)
    generalFunc(c)

    for i := 0; i < 5; i++ { fmt.Println("") }

    b := komoku.NewBoard(komoku.DefaultBoardSize)
    b.TurnPlayMove(3,3)
    b.TurnPlayMove(3,4)
    b.TurnPlayMove(2,4)
    b.TurnPlayMove(8,8)
    b.TurnPlayMove(3,5)
    b.TurnPlayMove(9,9)
    b.TurnPlayMove(4,4)
    komoku.PrintBoard(b)

    for i := 0; i < 5; i++ { fmt.Println("") }

    var v vector.StringVector
    v.Push("Hans")
    v.Push("Wurst")
    v.Push("Käse")
    v.Push("Adalbert")
    v.Push("dieter ")
    v.Do(func (elem string) {
        fmt.Println(elem)
    })
    fmt.Printf("\n\n")
    sort.SortStrings(sort.StringArray(v))
    v.Do(func (elem string) {
        fmt.Println(elem)
    })

    for i := 0; i < 5; i++ { fmt.Println("") }

    m := make(map[string]bool)
    m["test"] = true
    m["hur"] = false
    for a, b := range m {
        fmt.Printf("a: %v, b: %v\n", a,b)
    }
    kette := "Teststring"
    fmt.Printf("%s\n", kette[0:1])
    fmt.Printf("%s\n", kette[1:len(kette)])

    for i := 0; i < 5; i++ { fmt.Println("") }

    wd, _ := os.Getwd()
    fname := wd + "../../../data/tmp/TestListLegalPoints.GTPsequence.tmp"
    fname = path.Clean(fname)
    fmt.Println(fname)

    for i := 0; i < 5; i++ { fmt.Println("") }

    var slice = []string{ "hello", "this", "is", "slice" }
    for s := range slice {
        fmt.Printf("%s\n", s)
    }

    for i := 0; i < 5; i++ { fmt.Println("") }

    testIl := komoku.NewIntList()
    for k := 0; k < 100; k++ {
        testIl.Append(k)
    }
    eachFunc := func(val int) {
        fmt.Println(val)
    }
    testIl.Do(eachFunc)
}

func main() {
    //testMain()
    komoku.RunGTPMode()
    //see goGUI, if only 1 field is free and all others are of the same color, this move is illegal (says komoku) - thats not true.
    //write a test that tests if the whole board is exhausted by stones and legal fields.
}

