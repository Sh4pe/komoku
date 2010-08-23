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
)

func main() {
    testsuite := komoku.Testsuite()
    fmt.Printf("Running %d tests in SUBSTITUTE_THIS...\n", len(testsuite))
    testing.Main(testsuite)
}
