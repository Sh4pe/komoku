/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */

// This is a skeleton for test programs.
package main
import (
    "SUBSTITUTE_THIS"
    "testing"
    "fmt"
    "os"
)

func main() {
    testsuite := komoku.Testsuite()
    numTests := len(testsuite)
    if numTests > 1 || numTests == 0 {
        fmt.Printf("Running %d tests in SUBSTITUTE_THIS...\n", numTests)
    } else {
        fmt.Printf("Running %d test in SUBSTITUTE_THIS...\n", numTests)
    }
    matchAlways := func(pat, str string) (bool, os.Error) {
        return true, nil
    }
    testing.Main(matchAlways, testsuite)
}
