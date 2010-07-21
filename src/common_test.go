package komoku

import (
    "testing"
)

func TestXYToPos(t *testing.T) {
    sz := komoku.BoardSize
    for x := 0; x < sz; x++ {
        for y := 0; y < sz; y++ {
            rx, ry := posToXY( xyToPos(x, y))
        }
    }
}
