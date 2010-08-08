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

func TestIntListRemove(t *testing.T) {
    dummy := [...]int {0,1,2,3,4,5,6,7,8,9}
    length := len(dummy)
    test := NewIntList()
    for i := 0; i < length; i++ {
        test.Append(i)
    }
    l1 := test.Length()
    test.Remove(2)
    l2 := test.Length()
    if l2 != l1 - 1 {
        t.Fatalf("expected %d, got %d", l1 - 1, l2)
    }
}

// Tests wheather Remove and Append yield consistend IntLists of correct length
func TestRemoveAndAppendWork(t *testing.T) {
    for testLen := 3; testLen < BoardSize*BoardSize + 10; testLen++ {
        il := NewIntList()
        for i := 0; i < testLen; i++ {
            il.Append(i)
        }
        removed := 0
        for even := 0; even < testLen; even += 2 {
            //fmt.Printf("\nfi : %s\n", il)
            //fmt.Printf("removed: %d, even: %d, il.Length(): %d\n", removed, even, il.Length())
            il.Remove(even)
            removed++
        }
        if il.Length() != testLen - removed {
            t.Fatalf("IntList has wrong lenth after removing: expected %d, got %d", testLen-removed, il.Length())
        }
        // Add the same values we removed before
        for even := 0; even < testLen; even += 2 {
            il.Append(even)
        }
        if il.Length() != testLen {
            t.Fatalf("IntList has wrong lenth after re-adding the removed indices: expected %d, got %d", testLen, il.Length())
        }
        // Test if every value occurs only once
        unique := make(map[int]bool)
        last := il.Last()
        for it := il.First(); it != last; it = it.Next() {
            unique[it.Value()] = true
        }
        if len(unique) != testLen {
            t.Fatalf("IntList seems to store non-unique values. Expected %d, got %d", testLen, len(unique))
        }
    }
}

func Testsuite() []testing.Test {
    return []testing.Test { testing.Test{"TestIntListRemove", TestIntListRemove},
                            testing.Test{"TestRemoveAndAppendWork", TestRemoveAndAppendWork},
                          }
}
