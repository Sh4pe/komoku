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
)

func TestGroupSlice(t *testing.T) {
    gs := NewGroupSlice()
    const max = 10
    for i := 0; i < max; i++ {
        gs.Push(NewGroup(Black))
    }
    unique := make(map[*Group]bool)
    for _, g := range gs {
        unique[g] = true
    }
    if len(unique) != max {
        t.Fatalf("GroupSlice.Push produces GroupSlices with wrong length, expected %d, got %d", max, len(unique))
    }
}

func Testsuite() []testing.Test {
    return []testing.Test {
        testing.Test{"TestGroupSlice", TestGroupSlice},
    }
}
