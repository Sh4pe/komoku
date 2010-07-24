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

func TestNeighbours(t *testing.T) {
    for row := 0; row < BoardSize; row++ {
        for col := 0; col < BoardSize; col++ {
            expectedLen := 4
            if row == 0 || row == BoardSize-1 {
                expectedLen--
            }
            if col == 0 || col == BoardSize-1 {
                expectedLen--
            }
            n := neighbours(col, row)
            if len(n) != expectedLen {
                t.Fatalf("expected %d, got %d", expectedLen, len(n))
            }
            for _, ni := range n {
                // Each neighbour on same row or same column?
                if (ni.x != col) && (ni.y != row) {
                    t.Fatalf("(%d,%d)'s neighbour (%d,%d) not on same row/column", col, row, ni.x, ni.y)
                }
                // Each neighbour has the right distance?
                dx := (ni.x-col)*(ni.x-col)
                dy := (ni.y-row)*(ni.y-row)
                // one of dx, dy has to be 1, the other 0
                if (dx - dy)*(dx - dy) != 1 {
                    //t.Logf("dx^2: %d", dx)
                    //t.Logf("dy^2: %d", dy)
                    t.Fatalf("(%d,%d)'s neighbour (%d,%d) has the wrong distance", col, row, ni.x, ni.y)
                }
            }
        }
    }
}

func Testsuite() []testing.Test {
    return []testing.Test { testing.Test{"TestXYToPos", TestXYToPos},
                            testing.Test{"TestPosToXY", TestPosToXY},
                            testing.Test{"TestNeighbours", TestNeighbours}, 
                          }
}

