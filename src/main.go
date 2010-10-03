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
    "./komoku/komoku"
)

type testerContext struct {
    string
}

type testerFunc struct {
    context *testerContext
    f func(c *testerContext)
}

func (t *testerFunc) Call() {
    t.f(t.context)
}

func NewTesterFunc() *testerFunc {
    return &testerFunc{
        context: &testerContext {
            string: "init",
        },
        f: func(c *testerContext) {
            fmt.Printf("closure; context.string: '%s'\n", c.string)
        },
    }
}

type tester struct {
    fu *testerFunc
}

func (t *tester) Copy() *tester {
    f := NewTesterFunc()
    f.context.string = t.fu.context.string
    return &tester{
        fu: f,
    }
}

func NewTester() *tester {
    return &tester{
        fu: NewTesterFunc(),
    }
}

func testMain() {
    t1 := NewTester()
    t2 := t1.Copy()

    t1.fu.Call()
    t2.fu.Call()

    t2.fu.context.string = "modified"

    t1.fu.Call()
    t2.fu.Call()

    t1.fu.context.string = "modified original"

    t1.fu.Call()
    t2.fu.Call()

    for i := 0; i < 5; i++ {
        fmt.Println()
    }

    b1 := komoku.NewBoard(19)
    b2 := b1.Copy()
    _ = b2
}

func normalMain() {
    komoku.RunGTPMode()
}

func main() {
    //testMain()
    komoku.RunGTPMode()
}
