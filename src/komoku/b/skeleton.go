/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */

// This is a skeleton for benchmark programs.

package main

import (
    "testing"
    "flag"
    "fmt"
    "os"
    "regexp"
    "SUBSTITUTE_THIS"
)

func main() {
    flag.Parse()
    b := komoku.Benchmarks()
    if len(os.Args) == 1 {
        // No arguments supplied, so simply print the available benchmarks.
        bn := ""
        for _, n := range b {
            bn += n.Name + " "
        }
        fmt.Println("benchmarks:")
        fmt.Println(bn)
        return
    }
    matchAlways := func(pat, str string) (bool, os.Error) {
        return regexp.MatchString(pat, str)
    }
    testing.RunBenchmarks(matchAlways, b)
}

