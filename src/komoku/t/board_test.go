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
    //"fmt"
)

/*
 * TODO:
 *      - This becomes ugly and hardly maintainable. Refactor it!
 */

func TestCreateGroup(t *testing.T) {
    for x := 0; x < DefaultBoardSize; x++ {
        for y := 0; y < DefaultBoardSize; y++ {
            b := NewBoard(DefaultBoardSize)
            b.CreateGroup(x,y, Black)
            empty, g := b.GetGroup(x,y)
            if empty {
                t.Fatalf("Field empty after CreateGroup created a group on it")
            }
            nbours := b.neighbours(x,y)
            if len(nbours) != g.Liberties.Length() {
                t.Fatalf("Different number of liberties (%d) and neighbours (%d)", g.Liberties.Length(), len(nbours))
            }
            // are the neighbours of (x,y) exactly the liberties of the new group?
            last := g.Liberties.Last()
            for it := g.Liberties.First(); it != last; it = it.Next() {
                found := false
                for _, p := range nbours {
                    if b.xyToPos(p.X, p.Y) == it.Value() {
                        found = true
                        break
                    }
                }
                if !found {
                    t.Fatalf("Neighbours and .Liberties differ")
                }
            }

            expect := b.BoardSize()*b.BoardSize() - 1
            if b.emptyFields.Length() != expect {
                t.Fatalf("emptyFields has wrong length afterwards, expected %d, got %d", expect, b.emptyFields.Length())
            }
            if b.legalBlackMoves.Length() != expect {
                t.Fatalf("legalBlackMoves has wrong length afterwards, expected %d, got %d", expect, b.legalBlackMoves.Length())
            }
            if b.legalWhiteMoves.Length() != expect {
                t.Fatalf("legalWhiteMoves has wrong length afterwards, expected %d, got %d", expect, b.legalWhiteMoves.Length())
            }
        }
    }
}

func TestUpdateGroupLiberties(t *testing.T) {
    b := NewBoard(DefaultBoardSize)
    // single stones
    p := NewPoint(1,1)
    pos := b.xyToPos(p.X, p.Y)
    b.CreateGroup(p.X, p.Y, Black)
    _, g := b.GetGroup(p.X, p.Y)
    b.updateGroupLiberties(g)
    if g.Liberties.Length() != 4 {
        t.Fatalf("Wrong # of liberties, got %d, want 4", g.Liberties.Length())
    }
    // two stones in a row
    p2 := NewPoint(1,2)
    pos2 := b.xyToPos(p2.X, p2.Y)
    b.CreateGroup(p2.X, p2.Y, Black)
    b.joinGroups(b.fields[pos], b.fields[pos2])
    b.updateGroupLiberties(g)
    if g.Liberties.Length() != 6 {
        t.Fatalf("(2 stones): Wrong # of liberties, got %d, want 6", g.Liberties.Length())
    }
    // empty triangle - ewww!
    p3 := NewPoint(2,2)
    pos3 := b.xyToPos(p3.X, p3.Y)
    b.CreateGroup(p3.X, p3.Y, Black)
    b.joinGroups(b.fields[pos], b.fields[pos3])
    b.updateGroupLiberties(g)
    if g.Liberties.Length() != 7 {
        t.Fatalf("(empty triangle): Wrong # of liberties, got %d, want 7", g.Liberties.Length())
    }
}

func TestJoinGroups(t *testing.T) {
    b := NewBoard(DefaultBoardSize)
    t1 := []*Point{ NewPoint(1,1), NewPoint(1,2), NewPoint(2,2) }
    t2 := []*Point{ NewPoint(1,3), NewPoint(1,4), NewPoint(2,4) }

    ref := t1
    refpos := b.xyToPos(ref[0].X, ref[0].Y)
    t1pos := refpos
    b.CreateGroup(ref[0].X, ref[0].Y, Black)
    for i := 1; i < len(ref); i++ {
        p := ref[i]
        b.CreateGroup(p.X, p.Y, Black)
        b.joinGroups(b.fields[refpos], b.fields[b.xyToPos(p.X, p.Y)])
    }
    if b.groupMap.Length() != 1 {
        t.Fatalf("(Created 1st group) Wrong number of groups, got %d, want 1", b.groupMap.Length())
    }

    ref = t2
    refpos = b.xyToPos(ref[0].X, ref[0].Y)
    t2pos := refpos
    b.CreateGroup(ref[0].X, ref[0].Y, Black)
    for i := 1; i < len(ref); i++ {
        p := ref[i]
        b.CreateGroup(p.X, p.Y, Black)
        b.joinGroups(b.fields[refpos], b.fields[b.xyToPos(p.X, p.Y)])
    }
    if b.groupMap.Length() != 2 {
        t.Fatalf("(Created 2nd group) Wrong number of groups, got %d, want 2", b.groupMap.Length())
    }

    b.joinGroups(b.fields[t1pos], b.fields[t2pos])
    if b.groupMap.Length() != 1 {
        t.Fatalf("(Joined the two groups) Wrong number of groups, got %d, want 1", b.groupMap.Length())
    }
    _, g := b.GetGroup(t1[0].X, t1[0].Y)
    b.updateGroupLiberties(g)
    if g.Liberties.Length() != 11 {
        t.Fatalf("Wrong number of liberties, got %d, wanted 11", g.Liberties.Length())
    }
}

func TestGetEnvironment(t *testing.T) {
    // Produce a ponnuki and get the environment inside of it
    b := NewBoard(DefaultBoardSize)
    points := []*Point{ NewPoint(1,1), NewPoint(2,2), NewPoint(0,2), NewPoint(1,3) }
    for _, p := range points {
        b.CreateGroup(p.X, p.Y, Black)
    }
    nFree, adjBlack, adjWhite := b.GetEnvironment(1,2)
    if nFree != 0 || adjBlack.Length() != 4 || adjWhite.Length() != 0 {
        t.Fatalf("Board.GetEnvironment returns wront results inside a ponnuki")
    }
}

func TestRemoveGroup(t *testing.T) {
    b := NewBoard(DefaultBoardSize)

    t1 := []*Point{ NewPoint(1,1), NewPoint(1,2), NewPoint(2,2) }
    ref := t1
    refpos := b.xyToPos(ref[0].X, ref[0].Y)
    b.CreateGroup(ref[0].X, ref[0].Y, Black)
    for i := 1; i < len(ref); i++ {
        p := ref[i]
        b.CreateGroup(p.X, p.Y, Black)
        b.joinGroups(b.fields[refpos], b.fields[b.xyToPos(p.X, p.Y)])
    }
    b.RemoveGroup(1,1)
    if b.emptyFields.Length() != b.BoardSize()*b.BoardSize() {
        t.Fatalf("Board.RemoveGroup does not completely remove the group")
    }
}

func TestXYToPos(t *testing.T) {
    b := NewBoard(DefaultBoardSize)
    sz := b.BoardSize()
    for x := 0; x < sz; x++ {
        for y := 0; y < sz; y++ {
            rx, ry := b.posToXY(b.xyToPos(x, y))
            if rx != x || ry != y {
                t.Fatalf("expected (%d, %d), got (%d, %d)", x, y, rx, ry)
            }
        }
    }
}

func TestPosToXY(t *testing.T) {
    b := NewBoard(DefaultBoardSize)
    for i := 0; i<b.BoardSize()*b.BoardSize(); i++ {
        x, y := b.posToXY(i)
        retPos := b.xyToPos(x,y)
        if retPos != i {
            t.Fatalf("expected %d, got %d", i, retPos)
        }
    }
}

func TestNeighbours(t *testing.T) {
    b := NewBoard(DefaultBoardSize)
    for row := 0; row < b.BoardSize(); row++ {
        for col := 0; col < b.BoardSize(); col++ {
            expectedLen := 4
            if row == 0 || row == b.BoardSize()-1 {
                expectedLen--
            }
            if col == 0 || col == b.BoardSize()-1 {
                expectedLen--
            }
            n := b.neighbours(col, row)
            if len(n) != expectedLen {
                t.Fatalf("expected %d, got %d", expectedLen, len(n))
            }
            for _, ni := range n {
                // Each neighbour on same row or same column?
                if (ni.X != col) && (ni.Y != row) {
                    t.Fatalf("(%d,%d)'s neighbour (%d,%d) not on same row/column", col, row, ni.X, ni.Y)
                }
                // Each neighbour has the right distance?
                dx := (ni.X-col)*(ni.X-col)
                dy := (ni.Y-row)*(ni.Y-row)
                // one of dx, dy has to be 1, the other 0
                if (dx - dy)*(dx - dy) != 1 {
                    //t.Logf("dx^2: %d", dx)
                    //t.Logf("dy^2: %d", dy)
                    t.Fatalf("(%d,%d)'s neighbour (%d,%d) has the wrong distance", col, row, ni.X, ni.Y)
                }
            }
        }
    }
}

func Testsuite() []testing.Test {
    return []testing.Test { testing.Test{"TestCreateGroup", TestCreateGroup},
                            testing.Test{"TestUpdateGroupLiberties", TestUpdateGroupLiberties},
                            testing.Test{"TestJoinGroups", TestJoinGroups},
                            testing.Test{"TestGetEnvironment", TestGetEnvironment},
                            testing.Test{"TestRemoveGroup", TestRemoveGroup},
                            testing.Test{"TestXYToPos", TestXYToPos},
                            testing.Test{"TestPosToXY", TestPosToXY},
                            testing.Test{"TestNeighbours", TestNeighbours},
                         }
}
