package komoku

import (
    "testing"
    //"fmt"
)

func TestFieldIndicesRemove(t *testing.T) {
    dummy := [...]int {0,1,2,3,4,5,6,7,8,9}
    length := len(dummy)
    test := NewFieldIndices(length)
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

// Tests wheather Remove and Append yield consistend FieldIndices of correct length
func TestRemoveAndAppendWork(t *testing.T) {
    for testLen := 3; testLen < BoardSize*BoardSize + 10; testLen++ {
        fi := NewFieldIndices(testLen)
        for i := 0; i < testLen; i++ {
            fi.Append(i)
        }
        removed := 0
        for even := 0; even < testLen; even += 2 {
            //fmt.Printf("\nfi : %s\n", fi)
            //fmt.Printf("removed: %d, even: %d, fi.Length(): %d\n", removed, even, fi.Length())
            fi.Remove(even)
            removed++
        }
        if fi.Length() != testLen - removed {
            t.Fatalf("FieldIndices has wrong lenth after removing: expected %d, got %d", testLen-removed, fi.Length())
        }
        // Add the same values we removed before
        for even := 0; even < testLen; even += 2 {
            fi.Append(even)
        }
        if fi.Length() != testLen {
            t.Fatalf("FieldIndices has wrong lenth after re-adding the removed indices: expected %d, got %d", testLen, fi.Length())
        }
        // Test if every value occurs only once
        unique := make(map[int]bool)
        for i := 0; i < fi.Length(); i++ {
            unique[fi.Get(i)] = true
        }
        if len(unique) != testLen {
            t.Fatalf("FieldIndices seems to store non-unique values")
        }
    }
}

func Testsuite() []testing.Test {
    return []testing.Test { testing.Test{"TestFieldIndicesRemove", TestFieldIndicesRemove},
                            testing.Test{"TestRemoveAndAppendWork", TestRemoveAndAppendWork},
                          }
}
