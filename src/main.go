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
    "strings"
    //"container/list"
    "./komoku/komoku"
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
func (i *fimpl) Wee() string {
    return "Father wees"
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

func main() {
    var f *fimpl
    var c *cimpl
    childAsFather(c)
    fmt.Println()
    fatherDowntoChild(f)

    for i := 0; i < 5; i++ { fmt.Println("") }

    b := komoku.NewBoard()
    b.TurnPlayMove(3,3)
    b.TurnPlayMove(3,4)
    b.TurnPlayMove(2,4)
    b.TurnPlayMove(8,8)
    b.TurnPlayMove(3,5)
    b.TurnPlayMove(9,9)
    b.TurnPlayMove(4,4)
    komoku.PrintBoard(b)

    for i := 0; i < 5; i++ { fmt.Println("") }

    GTP := komoku.NewGTPObject()
    //line := "       This is     \t\ttext # and this is behind the hash :-) \n"
    line := "2 command_name arguments additional_arguments\n"
    res, _, _ := GTP.ExecuteCommand(line)
    fmt.Printf("'%v'\n\n", res)
    resSlice := strings.Split(res, " ", -1)
    for _, s := range resSlice {
        fmt.Printf("'%v'\n", s)
    }
}

