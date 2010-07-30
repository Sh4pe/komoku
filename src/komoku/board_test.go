package komoku

import (
    "testing"
)

func TestFieldIndicesRemove(t *testing.T) {
    // Very dump at the moment.
    dummy := [...]int {0,1,2,3,4,5,6,7,8,9}
    length := len(dummy)
    test := NewFieldIndices(length, length)
    for i := 0; i < length; i++ {
        test.Set(i,i)
    }
    l1 := test.Length()
    test.Remove(2)
    l2 := test.Length()
    if l2 != l1 - 1 {
        t.Fatalf("expected %d, got %d", l1 - 1, l2)
    }
}

func Testsuite() []testing.Test {
    return []testing.Test { testing.Test{"TestFieldIndicesRemove", TestFieldIndicesRemove},
                          }
}
