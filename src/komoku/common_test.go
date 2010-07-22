package komoku

import (
    "testing"
)

func TestXYToPos(t *testing.T) {
    sz := BoardSize
    for x := 0; x < sz; x++ {
        for y := 0; y < sz; y++ {
            rx, ry := posToXY(xyToPos(x, y))
            if rx != x || ry != y {
                t.Fatalf("expected (%d, %d), got (%d, %d)", x, y, rx, ry)
            }
        }
    }
}

func TestPosToXY(t *testing.T) {
    for i := 0; i<BoardSize*BoardSize; i++ {
        x, y := posToXY(i)
        retPos := xyToPos(x,y)
        if retPos != i {
            t.Fatalf("expected %d, got %d", i, retPos)
        }
    }
}

func Testsuite() []testing.Test {
    return []testing.Test { testing.Test{"XYToPos", TestXYToPos},
                            testing.Test{"PosToXY", TestPosToXY} }
}

