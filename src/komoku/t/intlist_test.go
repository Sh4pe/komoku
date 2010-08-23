/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */
package komoku

import (
    //"fmt"
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
    for testLen := 3; testLen < DefaultBoardSize*DefaultBoardSize + 10; testLen++ {
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

// Tests if JoinUnique yields IntLists with correct lengths.
func TestJoinUniqueLength(t *testing.T) {
    // First: disjoint IntLists
    for top := 0; top < 100; top++ {
        evenIntList := NewIntList()
        for even := 0; even < top; even += 2 {
            evenIntList.Append(even)
        }
        oddList := NewIntList()
        for odd := 1; odd < top; odd += 2 {
            oddList.Append(odd)
        }
        evenIntList.JoinUnique(oddList)
        if evenIntList.Length() != top {
            t.Fatalf("IntList has wrong length after .JoinUnique (disjoint lists): expected %d, got %d", top, evenIntList.Length())
        }
    }
    // Second: setwise identical IntLists
    for top := 0; top < 100; top++ {
        listOne := NewIntList()
        listTwo := NewIntList()
        for i := 0; i < top; i++ {
            listOne.Append(i)
            listTwo.Append(top-1-i)
        }
        listOne.JoinUnique(listTwo)
        if listOne.Length() != top {
            t.Fatalf("IntList has wrong length after .JoinUnique (setwise identical lists): expected %d, got %d", top, listOne.Length())
        }
    }
    // Third: non-identical list with nonempty intersection
    for length := 50; length < 100; length++ {
        for offset := 10; offset < length/2; offset++ {
            listOne := NewIntList()
            listTwo := NewIntList()
            for i := 0; i < length; i++ {
                listOne.Append(i)
                listTwo.Append(i+offset)
            }
            listOne.JoinUnique(listTwo)
            if listOne.Length() != length + offset {
                t.Fatalf("IntList has wrong length after .JoinUnique (nonempty intersection): expected %d, got %d", length + offset, listOne.Length())
            }
        }
    }
}

// Tests if JoinUnique yields 'unique' lists...
func TestJoinUniqueSetwise(t *testing.T) {
    for top := 50; top < 500; top++ {
        l1 := NewIntList()
        l2 := NewIntList()
        for i := 0; i < top - 10; i++ {
            l1.Append(i)
            l2.Append(i)
        }
        for i := top - 10; i < top; i += 2 {
            l1.Append(i)
            l2.Append(i+1)
        }
        l1.JoinUnique(l2)
        uniq := make(map[int]bool)
        last := l1.last;
        for it := l1.First(); it != last; it = it.Next() {
            uniq[it.Value()] = true
        }
        if len(uniq) != top {
            //fmt.Printf("uniq:\n%v\ntop: %d\n", uniq, top)
            t.Fatalf("JoinUnique yields non-unique lists. Expected length: %d, got %d", top, len(uniq))
        }
    }
}

func Testsuite() []testing.Test {
    return []testing.Test { testing.Test{"TestIntListRemove", TestIntListRemove},
                            testing.Test{"TestRemoveAndAppendWork", TestRemoveAndAppendWork},
                            testing.Test{"TestJoinUniqueLength", TestJoinUniqueLength},
                            testing.Test{"TestJoinUniqueSetwise", TestJoinUniqueSetwise},
                          }
}
