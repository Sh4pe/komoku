/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */
package main

/*
 * Plan for this branch:
 *      Get rid of the GroupIndexType in Board.fields. Use *Group instead, nil can denote empty.
 *      Get rid of most (if not all) [probably only some] the use of IntList, implement a GroupSlice instead....
 *      Board.GetEnvironment and .determineGroupsAtariStatus should be refactored in this new context.
 *      Use vector.IntVector for legal{Black,White}Moves and emptyFields
 */


import (
    "fmt"
    "./komoku/komoku"
)

func testMain() {
    fmt.Println("Testmain")
    slc := make([]int, 10)
    slc2 := make([]int, 10, 20)
    fmt.Printf("slc: %d, slc2: %d\n", len(slc), len(slc2))

    for i := 0; i < 5; i++ { fmt.Println("") }

    iPtrs := make([]*int, 10)
    for i := 0; i < len(iPtrs); i++ {
        iPtrs[i] = new(int)
        *iPtrs[i] = i
    }

    iPtrs = iPtrs[0:0]
    fmt.Println("start")
    for i := 0; i < len(iPtrs); i++ {
        fmt.Printf("%v\n", *iPtrs[i])
    }
    fmt.Println("stop")
    iPtrs = iPtrs[0:5]
    fmt.Println("start")
    for i := 0; i < len(iPtrs); i++ {
        fmt.Printf("%v\n", *iPtrs[i])
    }
    fmt.Println("stop")
    iPtrs = iPtrs[0:10]
    fmt.Println("start")
    for i := 0; i < len(iPtrs); i++ {
        fmt.Printf("%v\n", *iPtrs[i])
    }
    fmt.Println("stop")

    for i := 0; i < 5; i++ { fmt.Println("") }

    gs := komoku.NewGroupSlice()
    const max = 10
    for i := 0; i < max; i++ {
        gs.Push(komoku.NewGroup(komoku.Black))
    }
    for i, g := range gs {
        fmt.Printf("%d, %v\n", i, g)
    }
    i := 10
    j := i 
    i++
    fmt.Printf("%d\n", j)
}

func normalMain() {
    komoku.RunGTPMode()
}

func main() {
    //now start implementing the AI
    //testMain()
    komoku.RunGTPMode()
}
