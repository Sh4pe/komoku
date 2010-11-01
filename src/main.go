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
    "runtime"
    "./komoku/komoku"
    //"time"
)

func testMain() {
    fmt.Printf("runtime.GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))

    for i := 0; i < 5; i++ {
        fmt.Println()
    }


}

func normalMain() {
    komoku.RunGTPMode()
}

func main() {
    //testMain()
    komoku.RunGTPMode()
}
