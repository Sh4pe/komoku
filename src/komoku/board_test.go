package komoku

import (
    "testing"
)

func TestFieldIndicesRemove(t *testing.T) {
    // Very dump at the moment.
    dummy := [...]int {0,1,2,3,4,5,6,7,8,9}
    var test FieldIndices = dummy[0:10]
    l1 := len(test)
    test.Remove(2)
    l2 := len(test)
    if l2 != l1 - 1 {
        t.Fatalf("expected %d, got %d", l1 - 1, l2)
    }
}

func Testsuite() []testing.Test {
    return []testing.Test { testing.Test{"TestFieldIndicesRemove", TestFieldIndicesRemove},
                          }
}
