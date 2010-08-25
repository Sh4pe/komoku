/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */
package komoku

import (
    "testing"
    "os"
    //"fmt"
)

var komokuSource = []string{ // at least a part of it
    "./board_test.go",
    "./common_test.go",
    "./gtp_test.go",
    "./intlist_test.go",
    "./skeleton.go",
    "./ui_test.go",
    "../board.go",
    "../common.go",
    "../debug.go",
    "../environment.go",
    "../game.go",
    "../group.go",
    "../gtp.go",
    "../gtpcmd.go",
    "../intlist.go",
    "../treenode.go",
    "../ui.go",
    "../../main.go",
}
// tests if all the komokuSource-files exist
func TestRelPathToAbs(t *testing.T) {
    // The tests are normally invoked through a Makefile. We first check if every file exists...
    for _, rel := range komokuSource {
        abs := relPathToAbs(rel)
        fi, err := os.Stat(abs)
        if err != nil || !fi.IsRegular() {
            t.Fatalf("%s is no regular file", abs)
        }
    }
    // Then run the whole common_test executable from the system
    os.Exec("./common_test", nil, nil)
}

func Testsuite() []testing.Test {
    return []testing.Test { testing.Test{"TestRelPathToAbs", TestRelPathToAbs},
                          }
}

