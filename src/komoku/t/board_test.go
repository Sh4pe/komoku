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
    "os"
    "rand"
    "time"
    "fmt"
)

/*
 * TODO:
 *      - This becomes ugly and hardly maintainable. Refactor it!
 */

func TestCreateGroup(t *testing.T) {
    testname := "TestCreateGroup"
    fmt.Printf("entering %s\n", testname)
    defer fmt.Printf("leaving %s\n", testname)

    for x := 0; x < DefaultBoardSize; x++ {
        for y := 0; y < DefaultBoardSize; y++ {
            b := NewBoard(DefaultBoardSize)
            b.CreateGroup(x,y, Black)
            g := b.GetGroup(x,y)
            if g == nil {
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
            if len(b.ListEmptyFields()) != expect {
                t.Fatalf("emptyFields has wrong length afterwards, expected %d, got %d", expect, len(b.ListEmptyFields()))
            }
            if len(b.ListLegalPoints(Black)) != expect {
                t.Fatalf("wrong number of legal black moves afterwards, expected %d, got %d", expect, len(b.ListLegalPoints(Black)))
            }
            if len(b.ListLegalPoints(White)) != expect {
                t.Fatalf("wrong number of legal white moves afterwards, expected %d, got %d", expect, len(b.ListLegalPoints(White)))
            }
        }
    }
}

func TestUpdateGroupLiberties(t *testing.T) {
    testname := "TestUpdateGroupLiberties"
    fmt.Printf("entering %s\n", testname)
    defer fmt.Printf("leaving %s\n", testname)

    b := NewBoard(DefaultBoardSize)
    // single stones
    p := NewPoint(1,1)
    pos := b.xyToPos(p.X, p.Y)
    b.CreateGroup(p.X, p.Y, Black)
    g := b.GetGroup(p.X, p.Y)
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
    testname := "TestJoinGroups"
    fmt.Printf("entering %s\n", testname)
    defer fmt.Printf("leaving %s\n", testname)

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
    nblack, nwhite := b.numberOfGroups()
    nall := nblack + nwhite
    if nall != 1 {
        t.Fatalf("(Created 1st group) Wrong number of groups, got %d, want 1", nall)
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

    nblack, nwhite = b.numberOfGroups()
    nall = nblack + nwhite
    if nall != 2 {
        t.Fatalf("(Created 2nd group) Wrong number of groups, got %d, want 2", nall)
    }

    b.joinGroups(b.fields[t1pos], b.fields[t2pos])
    nblack, nwhite = b.numberOfGroups()
    nall = nblack + nwhite
    if nall != 1 {
        t.Fatalf("(Joined the two groups) Wrong number of groups, got %d, want 1", nall)
    }
    g := b.GetGroup(t1[0].X, t1[0].Y)
    b.updateGroupLiberties(g)
    if g.Liberties.Length() != 11 {
        t.Fatalf("Wrong number of liberties, got %d, wanted 11", g.Liberties.Length())
    }
}

type testGetEnvironmentCase struct {
    sequence []Move
    where Point
    enFree, eadjBlackLen, eadjWhiteLen int
}

var testGetEnvironment9 = []testGetEnvironmentCase{
    testGetEnvironmentCase {
        sequence: []Move{
            Move{ Color: Black, Vertex: *NewVertex(Point{1,1}, false) },
            Move{ Color: Black, Vertex: *NewVertex(Point{2,2}, false) },
            Move{ Color: Black, Vertex: *NewVertex(Point{0,2}, false) },
            Move{ Color: Black, Vertex: *NewVertex(Point{1,3}, false) },
        },
        where: Point{1,2},
        enFree: 0,
        eadjBlackLen: 4,
        eadjWhiteLen: 0,
    },
    testGetEnvironmentCase {
        sequence: []Move{
            Move{ Color: White, Vertex: *NewVertex(Point{0,3}, false) },
            Move{ Color: White, Vertex: *NewVertex(Point{0,4}, false) },
            Move{ Color: Black, Vertex: *NewVertex(Point{8,3}, false) },
        },
        where: Point{8,4},
        enFree: 2,
        eadjBlackLen: 1,
        eadjWhiteLen: 0,
    },
}

func TestGetEnvironment(t *testing.T) {
    testname := "TestGetEnvironment"
    fmt.Printf("entering %s\n", testname)
    defer fmt.Printf("leaving %s\n", testname)

    for i, tc := range testGetEnvironment9 {
        b := NewBoard(9)
        b.playSequence(tc.sequence)
        nFree, adjBlack, adjWhite := b.GetEnvironment(tc.where.X, tc.where.Y)
        if nFree != tc.enFree || len(adjBlack) != tc.eadjBlackLen || len(adjWhite) != tc.eadjWhiteLen {
            t.Fatalf("Board.GetEnvironment fails test case %d", i)
        }
    }
}

func TestRemoveGroup(t *testing.T) {
    testname := "TestRemoveGroup"
    fmt.Printf("entering %s\n", testname)
    defer fmt.Printf("leaving %s\n", testname)

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
    b.RemoveGroupByPos(1,1)
    if len(b.ListEmptyFields()) != b.BoardSize()*b.BoardSize() {
        t.Fatalf("Board.RemoveGroup does not completely remove the group")
    }
}

func TestXYToPos(t *testing.T) {
    testname := "TestXYToPos"
    fmt.Printf("entering %s\n", testname)
    defer fmt.Printf("leaving %s\n", testname)

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
    testname := "TestPosToXY"
    fmt.Printf("entering %s\n", testname)
    defer fmt.Printf("leaving %s\n", testname)

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
    testname := "TestNeighbours"
    fmt.Printf("entering %s\n", testname)
    defer fmt.Printf("leaving %s\n", testname)

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


type testNumGroupsCase struct {
    sequence []Move
    eNumBlack, eNumWhite int
}

var testNumGroups9 = []testNumGroupsCase {
    testNumGroupsCase {
        sequence: []Move{
            Move{ Color: Black, Vertex: *NewVertex(Point{0,3}, false) },
            Move{ Color: Black, Vertex: *NewVertex(Point{0,4}, false) },
        },
        eNumBlack: 1,
        eNumWhite: 0,
    },
    testNumGroupsCase {
        sequence: []Move{
            Move{ Color: White, Vertex: *NewVertex(Point{0,3}, false) },
            Move{ Color: White, Vertex: *NewVertex(Point{0,4}, false) },
            Move{ Color: Black, Vertex: *NewVertex(Point{8,3}, false) },
            Move{ Color: Black, Vertex: *NewVertex(Point{8,4}, false) },
        },
        eNumBlack: 1,
        eNumWhite: 1,
    },
    testNumGroupsCase {
        sequence: []Move{
            Move{ Color: Black, Vertex: *NewVertex(Point{8,3}, false) },
            Move{ Color: Black, Vertex: *NewVertex(Point{8,4}, false) },
            Move{ Color: White, Vertex: *NewVertex(Point{0,3}, false) },
            Move{ Color: White, Vertex: *NewVertex(Point{0,4}, false) },
        },
        eNumBlack: 1,
        eNumWhite: 1,
    },
}

func TestNumGroups(t *testing.T) {
    testname := "TestNumGroups"
    fmt.Printf("entering %s\n", testname)
    defer fmt.Printf("leaving %s\n", testname)

    // on a 9x9 board
    for i, tc := range testNumGroups9 {
        b := NewBoard(9)
        b.playSequence(tc.sequence)
        nblack, nwhite := b.numberOfGroups()
        if nblack != tc.eNumBlack || nwhite != tc.eNumWhite {
            t.Fatalf("testcase %d: expected (b%d,w%d), got (b%d,w%d)", i, tc.eNumBlack, tc.eNumWhite, nblack, nwhite)
        }
    }
}

var testNumStones9 = []testNumGroupsCase {
    testNumGroupsCase {
        sequence: []Move{
            Move{ Color: Black, Vertex: *NewVertex(Point{0,3}, false) },
            Move{ Color: Black, Vertex: *NewVertex(Point{0,4}, false) },
        },
        eNumBlack: 2,
        eNumWhite: 0,
    },
}

func TestNumStones(t *testing.T) {
    testname := "TestNumStones"
    fmt.Printf("entering %s\n", testname)
    defer fmt.Printf("leaving %s\n", testname)

    // on a 9x9 board
    for _, tc := range testNumStones9 {
        b := NewBoard(9)
        b.playSequence(tc.sequence)
        nblack, nwhite := b.numberOfStones()
        if nblack != tc.eNumBlack || nwhite != tc.eNumWhite {
            t.Fatalf("expected (b%d,w%d), got (b%d,w%d)", tc.eNumBlack, tc.eNumWhite, nblack, nwhite)
        }
    }
}

type testGroupGeometryCase struct {
    seqBlack, seqWhite []Point
    blackGroupPivot, whiteGroupPivot *Point
    eBlackPoints, eWhitePoints []Point
}

var testGroupGeometry9 = []testGroupGeometryCase {
    testGroupGeometryCase {
        seqBlack: []Point {
            Point{0,3}, Point{0,4},
        },
        seqWhite: nil,
        blackGroupPivot: &Point{0,4},
        whiteGroupPivot: nil,
        eBlackPoints: []Point{ Point{0,3}, Point{0,4} },
        eWhitePoints: nil,
    },
}

func TestGroupGeometry(t *testing.T) {
    testname := "TestGroupGeometry"
    fmt.Printf("entering %s\n", testname)
    defer fmt.Printf("leaving %s\n", testname)


    tester := func(number int, t *testing.T, b *Board, setPoints []Point, color Color) {
        for _, p := range setPoints {
            grp := b.GetGroup(p.X,p.Y)
            if grp == nil {
                t.Fatalf("Failed testcase %d (%s), there seems to be no group at (%d,%d)", number, color, p.X, p.Y)
            } else if grp.Color != color {
                t.Fatalf("Failed testcase %d (%s), stone has the wrong color", number, color)
            }
        }
        grp := b.GetGroup(setPoints[0].X, setPoints[0].Y)
        if len(setPoints) != grp.Fields.Length() {
            t.Fatalf("Failed testcase %d (%s), group has wrong number of stones, got %d, wanted %d", number, color, grp.Fields.Length(), len(setPoints))
        }
        all := make(map[int]bool)
        for _, p := range setPoints {
            pos := b.xyToPos(p.X, p.Y)
            all[pos] = true
        }
        last := grp.Fields.Last()
        for it := grp.Fields.First(); it != last; it = it.Next() {
            all[it.Value()] = false, false
        }
        if len(all) != 0 {
            t.Fatalf("Failed testcase %d (%s), expected black fields and white fields seem to differ", number, color)
        }

    }

    // on a 9x9 board
    for i, tc := range testGroupGeometry9 {
        b := NewBoard(9)
        for _, p := range tc.seqBlack {
            b.PlayMove(p.X, p.Y, Black)
        }
        for _, p := range tc.seqWhite {
            b.PlayMove(p.X, p.Y, White)
        }
        if tc.seqBlack != nil {
            tester(i, t, b, tc.seqBlack, Black)
        }
        if tc.seqWhite != nil {
            tester(i, t, b, tc.seqWhite, White)
        }
    }
}

type writeStringer interface {
    WriteString(s string) (ret int, err os.Error)
}

// Tests a simple ko and tenuki situation
func TestKo(t *testing.T) {
    testname := "TestKo"
    fmt.Printf("entering %s\n", testname)
    defer fmt.Printf("leaving %s\n", testname)

    board := NewBoard(9)
    sequence := []Move{
        Move{ Color: Black, Vertex: *NewVertex(Point{3,4}, false) },
        Move{ Color: Black, Vertex: *NewVertex(Point{4,5}, false) },
        Move{ Color: Black, Vertex: *NewVertex(Point{5,4}, false) },
        Move{ Color: Black, Vertex: *NewVertex(Point{4,3}, false) },
        Move{ Color: White, Vertex: *NewVertex(Point{3,5}, false) },
        Move{ Color: White, Vertex: *NewVertex(Point{4,6}, false) },
        Move{ Color: White, Vertex: *NewVertex(Point{5,5}, false) },
    }
    board.playSequence(sequence)
    board.PlayMove(4,4,White)
    legalWhite, _ := board.calculateIfLegal(4,5, White)
    legalBlack, _ := board.calculateIfLegal(4,5, Black)
    if !legalWhite {
        t.Fatalf("expected legal white move at (%d,%d)", 4,5)
    }
    if legalBlack {
        t.Fatalf("expected that a black move at (%d,%d) is illegal", 4,5)
    }
    board.PlayMove(1,1,Black)
    legalWhite, _ = board.calculateIfLegal(4,5, White)
    legalBlack, _ = board.calculateIfLegal(4,5, Black)
    if !legalWhite {
        t.Fatalf("after tennuki: expected legal white move at (%d,%d)", 4,5)
    }
    if !legalBlack {
        t.Fatalf("after tennuki: expected that a black move at (%d,%d) is now illegal", 4,5)
    }

}

// Tests two kos at one board
func TestDoubleKo(t *testing.T) {
    testname := "TestDoubleKo"
    fmt.Printf("entering %s\n", testname)
    defer fmt.Printf("leaving %s\n", testname)

    board := NewBoard(9)
    sequence := []Move{
        Move{ Color: Black, Vertex: *NewVertex(Point{3,4}, false) },
        Move{ Color: Black, Vertex: *NewVertex(Point{4,5}, false) },
        Move{ Color: Black, Vertex: *NewVertex(Point{5,4}, false) },
        Move{ Color: Black, Vertex: *NewVertex(Point{4,3}, false) },
        Move{ Color: White, Vertex: *NewVertex(Point{3,5}, false) },
        Move{ Color: White, Vertex: *NewVertex(Point{4,6}, false) },
        Move{ Color: White, Vertex: *NewVertex(Point{5,5}, false) },
        Move{ Color: White, Vertex: *NewVertex(Point{2,3}, false) },
        Move{ Color: White, Vertex: *NewVertex(Point{3,2}, false) },
        Move{ Color: White, Vertex: *NewVertex(Point{1,2}, false) },
        Move{ Color: White, Vertex: *NewVertex(Point{2,1}, false) },
        Move{ Color: Black, Vertex: *NewVertex(Point{3,1}, false) },
        Move{ Color: Black, Vertex: *NewVertex(Point{1,1}, false) },
        Move{ Color: Black, Vertex: *NewVertex(Point{2,0}, false) },
    }
    board.playSequence(sequence)
    legalBlack, _ := board.calculateIfLegal(4,4,Black)
    legalWhite, _ := board.calculateIfLegal(4,4,White)
    if !legalBlack || !legalWhite {
        t.Fatalf("wrong legality status at (%d,%d)", 4,4)
    }
    legalBlack, _ = board.calculateIfLegal(2,2,Black)
    legalWhite, _ = board.calculateIfLegal(2,2,White)
    if !legalBlack || !legalWhite {
        t.Fatalf("wrong legality status at (%d,%d)", 2,2)
    }
    board.PlayMove(4,4,White)
    legalBlack, _ = board.calculateIfLegal(4,5,Black)
    legalWhite, _ = board.calculateIfLegal(4,5,White)
    if legalBlack || !legalWhite {
        t.Fatalf("wrong legality status at (%d,%d)", 4,5)
    }
    board.PlayMove(2,2,Black)
    legalBlack, _ = board.calculateIfLegal(2,1,Black)
    legalWhite, _ = board.calculateIfLegal(2,1,White)
    if !legalBlack || legalWhite {
        t.Fatalf("wrong legality status at (%d,%d)", 2,1)
    }
    legalBlack, _ = board.calculateIfLegal(4,5,Black)
    legalWhite, _ = board.calculateIfLegal(4,5,White)
    if !legalBlack || !legalWhite {
        t.Fatalf("after 2nd ko: wrong legality status at (%d,%d)", 4,5)
    }
}

func gameStateCheck(game *Game,
                    dumpFile string,
                    t *testing.T,
                    nGame, nMove int,
                    legalBlack, legalWhite []Point,
                    ) {

    dumpSequence := func(w writeStringer) {
        for _, mv := range game.sequence {
            m, _ := mv.(Move)
            vertex, _ := pointToGTPVertex(*NewPoint(m.Vertex.X, m.Vertex.Y))
            line := fmt.Sprintf("  play %s %s\n", colorToGTPColor(m.Color), vertex)
            if _, werr := w.WriteString(line); werr != nil {
                t.Fatalf("dumpSequence does not work, check the writeStringer")
            }
        }
    }

    failBecausePositionIsIllegal := func(w writeStringer,
                                         failMessage, infoString string,
                                         ) {
        time := time.LocalTime()
        w.WriteString(fmt.Sprintf("# These moves lead to an illegal position. %s\n", time))
        _, werr := w.WriteString(fmt.Sprintf("  boardsize %d\n", game.Board.BoardSize()))
        if werr != nil {
            failMessage = fmt.Sprintf("The sequence should have been dumped into %s, but this file could not be opened.\n", dumpFile)
            failMessage += fmt.Sprintf("The error was: %s\n", werr)
        } else {
            dumpSequence(w)
            line := fmt.Sprintf("# %s\n", infoString)
            w.WriteString(line)

        }
        t.Fatalf(failMessage)
    }

    fileFail := func(failMessage, infoString string) {
        os.Remove(dumpFile)
        if file, err := os.Open(dumpFile, os.O_CREATE | os.O_RDWR, 0666); err == nil {
            failBecausePositionIsIllegal(file, failMessage, infoString)
        } else {
            t.Fatalf("Tried to create the error output file %s, but this error occured: %s", dumpFile, err)
        }
    }

    for _, p := range legalBlack {
        if gptr := game.Board.GetGroup(p.X, p.Y); gptr != nil {
            illegalVertex, _ := pointToGTPVertex(*NewPoint(p.X,p.Y))
            failMessage := fmt.Sprintf("Game %d, move %d: the point %s was legal for black but is already occupied\n", nGame, nMove, illegalVertex)
            failMessage += fmt.Sprintf("The sequence was dumped into %s", dumpFile)
            infoString := fmt.Sprintf("the vertex %s is legal for black but already occupied", illegalVertex)
            fileFail(failMessage, infoString)
        }
    }
    for _, p := range legalWhite {
        if gptr := game.Board.GetGroup(p.X, p.Y); gptr != nil {
            illegalVertex, _ := pointToGTPVertex(*NewPoint(p.X,p.Y))
            failMessage := fmt.Sprintf("Game %d, move %d: the point %s was legal for white but is already occupied\n", nGame, nMove, illegalVertex)
            failMessage += fmt.Sprintf("The sequence was dumped into %s", dumpFile)
            infoString := fmt.Sprintf("the vertex %s is legal for white but already occupied", illegalVertex)
            fileFail(failMessage, infoString)
        }
    }
    // Assemble the set of all groups in game.Board
    groupMap := make(map[*Group]bool)
    for _, grpPtr := range game.Board.fields {
        if grpPtr != nil {
            groupMap[grpPtr] = true
        }
    }
    // check that no groups have 0 liberties
    gmEach := func(grp *Group) {
        if grp.Liberties.Length() == 0 {
            failMessage := fmt.Sprintf("in game %d, there were groups with 0 liberties after %d moves.\nSeq dumped into %s", nGame, nMove, dumpFile)
            gX, gY := game.Board.posToXY(grp.Fields.First().Value())
            vertex, _ := pointToGTPVertex(*NewPoint(gX, gY))
            infoString := fmt.Sprintf("# the group with 0 libs is around %s\n", vertex)
            fileFail(failMessage, infoString)
        }
    }
    for grp, _ := range groupMap {
        gmEach(grp)
    }
    // check if the legal moves and the occupied fields completely exhaust the whole board
    allmap := make(map[int]bool)
    for i := 0; i < game.Board.BoardSize()*game.Board.BoardSize(); i++ {
        allmap[i] = true
    }
    for i := 0; i < game.Board.BoardSize()*game.Board.BoardSize(); i++ {
        //if !game.Board.fields[i].Empty() {
        if game.Board.fields[i] != nil {
            allmap[i] = false, false
        }
    }
    for _, p := range legalBlack {
        pos := game.Board.xyToPos(p.X, p.Y)
        allmap[pos] = false, false
    }
    for _, p := range legalWhite {
        pos := game.Board.xyToPos(p.X, p.Y)
        allmap[pos] = false, false
    }
    if len(allmap) != 0 {
        failMessage := fmt.Sprintf("in game %d, there were fields which were neigher empty nor legal for any color.\nSeq dumped into %s", nGame, dumpFile)
        infoString := "# the fields which are neither occupied nor legal are:\n# "
        for fpos, _ := range allmap {
            fx, fy := game.Board.posToXY(fpos)
            vertex, _ := pointToGTPVertex(*NewPoint(fx, fy))
            infoString += fmt.Sprintf(" %s ", vertex)
        }
        fileFail(failMessage, infoString)
    }
}

// used for some panic situations that occur while debugging komoku
func panicChecks(t *testing.T) {
    if e := recover(); e != nil {
        jjjjjjj
    }
}

// Generates random games and checks if the []Points returned by Board.ListLegalPoints do not intersect
// already occupied points
func TestListLegalPoints(t *testing.T) {
    testname := "TestListLegalPoints"
    fmt.Printf("entering %s\n", testname)
    defer fmt.Printf("leaving %s\n", testname)

    numGames := 500 // Number of games this test should play
    gamesLen := 100 // Number of random moves to play
    boardsize := 9
    dumpFile := relPathToAbs("../../../data/tmp/TestListLegalPoints.GTPsequence.tmp")
    lastMovePass := false
    for nGame := 0; nGame < numGames; nGame++ {
        fmt.Printf("Game %d\n", nGame)
        game := NewGame(boardsize)
        var currentColor Color = Black
        for nMove := 0; nMove < gamesLen; nMove++ {

            legalBlack := game.Board.ListLegalPoints(Black)
            legalWhite := game.Board.ListLegalPoints(White)

            gameStateCheck(game, dumpFile, t, nGame, nMove, legalBlack, legalWhite)

            var legal []Point
            if currentColor == Black {
                legal = legalBlack
            } else {
                legal = legalWhite
            }
            // If there is no legal move for the side whos turn it is, this side passes.
            // If there are two passes in a row, finish the game.
            if len(legal) == 0 {
                if lastMovePass {
                    break
                } else {
                    game.PlayPass(currentColor)
                    lastMovePass = true
                }
            } else {
                // Play a random move
                sec, nsec, _ := os.Time()
                random := rand.New(rand.NewSource(sec+nsec))
                randomMove := legal[random.Intn(len(legal))]
                game.PlayMove(randomMove.X, randomMove.Y, currentColor)
                currentColor = !currentColor
                lastMovePass = false
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
                            testing.Test{"TestNumGroups", TestNumGroups},
                            testing.Test{"TestNumStones", TestNumStones},
                            testing.Test{"TestGroupGeometry", TestGroupGeometry},
                            testing.Test{"TestKo", TestKo},
                            testing.Test{"TestDoubleKo", TestDoubleKo},
                            testing.Test{"TestListLegalPoints", TestListLegalPoints},
                         }
}
