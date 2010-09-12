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
    for x := 0; x < DefaultBoardSize; x++ {
        for y := 0; y < DefaultBoardSize; y++ {
            b := NewBoard(DefaultBoardSize)
            pos := b.xyToPos(x,y)
            b.CreateGroup(pos, Black)
            g := b.GetGroup(x,y)
            if g == nil {
                t.Fatalf("Field empty after CreateGroup created a group on it")
            }
            nbours := b.neighboursByPos(pos)
            if len(nbours) != g.Liberties.Length() {
                t.Fatalf("Different number of liberties (%d) and neighbours (%d)", g.Liberties.Length(), len(nbours))
            }
            // are the neighbours of (x,y) exactly the liberties of the new group?
            last := g.Liberties.Last()
            for it := g.Liberties.First(); it != last; it = it.Next() {
                found := false
                for _, npos := range nbours {
                    if npos == it.Value() {
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
    b := NewBoard(DefaultBoardSize)
    // single stones
    p := NewPoint(1,1)
    pos := b.xyToPos(p.X, p.Y)
    b.CreateGroup(pos, Black)
    g := b.GetGroup(p.X, p.Y)
    b.updateGroupLiberties(g)
    if g.Liberties.Length() != 4 {
        t.Fatalf("Wrong # of liberties, got %d, want 4", g.Liberties.Length())
    }
    // two stones in a row
    p2 := NewPoint(1,2)
    pos2 := b.xyToPos(p2.X, p2.Y)
    b.CreateGroup(pos2, Black)
    b.joinGroups(b.fields[pos], b.fields[pos2])
    b.updateGroupLiberties(g)
    if g.Liberties.Length() != 6 {
        t.Fatalf("(2 stones): Wrong # of liberties, got %d, want 6", g.Liberties.Length())
    }
    // empty triangle - ewww!
    p3 := NewPoint(2,2)
    pos3 := b.xyToPos(p3.X, p3.Y)
    b.CreateGroup(pos3, Black)
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
    b.CreateGroup(refpos, Black)
    for i := 1; i < len(ref); i++ {
        p := ref[i]
        b.CreateGroupByPoint(p.X, p.Y, Black)
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
    b.CreateGroup(refpos, Black)
    for i := 1; i < len(ref); i++ {
        p := ref[i]
        b.CreateGroupByPoint(p.X, p.Y, Black)
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
    b := NewBoard(DefaultBoardSize)

    t1 := []*Point{ NewPoint(1,1), NewPoint(1,2), NewPoint(2,2) }
    ref := t1
    refpos := b.xyToPos(ref[0].X, ref[0].Y)
    b.CreateGroup(refpos, Black)
    for i := 1; i < len(ref); i++ {
        p := ref[i]
        b.CreateGroupByPoint(p.X, p.Y, Black)
        b.joinGroups(b.fields[refpos], b.fields[b.xyToPos(p.X, p.Y)])
    }
    b.RemoveGroupByPos(1,1)
    if len(b.ListEmptyFields()) != b.BoardSize()*b.BoardSize() {
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
            pos := b.xyToPos(col, row)
            n := b.neighboursByPos(pos)
            if len(n) != expectedLen {
                t.Fatalf("expected %d, got %d", expectedLen, len(n))
            }
            for _, npos := range n {
                // Each neighbour on same row or same column?
                niX, niY := b.posToXY(npos)
                if (niX != col) && (niY != row) {
                    t.Fatalf("(%d,%d)'s neighbour (%d,%d) not on same row/column", col, row, niX, niY)
                }
                // Each neighbour has the right distance?
                dx := (niX-col)*(niX-col)
                dy := (niY-row)*(niY-row)
                // one of dx, dy has to be 1, the other 0
                if (dx - dy)*(dx - dy) != 1 {
                    //t.Logf("dx^2: %d", dx)
                    //t.Logf("dy^2: %d", dy)
                    t.Fatalf("(%d,%d)'s neighbour (%d,%d) has the wrong distance", col, row, niX, niY)
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
    /*testName := "TestKo"
    fmt.Printf("entering %s\n", testName)
    defer fmt.Printf("leaving %s\n", testName)*/

    legalityCheckedTooOften := func(err Error) {
        printDbgMsgBTf(4, "legalityCheckedTooOften not yet initialized!\n")
        t.Fatalf("legalityCheckedTooOften not yet initialized!\n")
    }
    inPanicLegalityCheckedTooOften := func(err Error) {
        legalityCheckedTooOften(err)
    }
    defer panicChecks(inPanicLegalityCheckedTooOften)

    game := NewGame(9)
    dumpFile := relPathToAbs("../../../data/tmp/TestKo.GTPsequence.tmp")
    sequence := []Move{
        Move{ Color: Black, Vertex: *NewVertex(Point{3,4}, false) },
        Move{ Color: Black, Vertex: *NewVertex(Point{4,5}, false) },
        Move{ Color: Black, Vertex: *NewVertex(Point{5,4}, false) },
        Move{ Color: Black, Vertex: *NewVertex(Point{4,3}, false) },
        Move{ Color: White, Vertex: *NewVertex(Point{3,5}, false) },
        Move{ Color: White, Vertex: *NewVertex(Point{4,6}, false) },
        Move{ Color: White, Vertex: *NewVertex(Point{5,5}, false) },
    }
    // set the error handling for panic checks to the appropriate context
    legalityCheckedTooOften = func(err Error) {
        failLegalityCheckedTooOften(dumpFile, game, err, t)
    }
    game.PlaySequence(sequence)
    game.PlayMove(4,4,White)
    legalWhite, _ := game.Board.calculateIfLegal(4,5, White)
    legalBlack, _ := game.Board.calculateIfLegal(4,5, Black)
    if !legalWhite {
        t.Fatalf("expected legal white move at (%d,%d)", 4,5)
    }
    if legalBlack {
        t.Fatalf("expected that a black move at (%d,%d) is illegal", 4,5)
    }
    game.PlayMove(1,1,Black)
    legalWhite, _ = game.Board.calculateIfLegal(4,5, White)
    legalBlack, _ = game.Board.calculateIfLegal(4,5, Black)
    if !legalWhite {
        t.Fatalf("after tennuki: expected legal white move at (%d,%d)", 4,5)
    }
    if !legalBlack {
        t.Fatalf("after tennuki: expected that a black move at (%d,%d) is now illegal", 4,5)
    }

}

// Tests two kos at one board
func TestDoubleKo(t *testing.T) {
    /*testName := "TestDoubleKo"
    fmt.Printf("entering %s\n", testName)
    defer fmt.Printf("leaving %s\n", testName)*/

    legalityCheckedTooOften := func(err Error) {
        printDbgMsgBTf(4, "legalityCheckedTooOften not yet initialized!\n")
        t.Fatalf("legalityCheckedTooOften not yet initialized!\n")
    }
    inPanicLegalityCheckedTooOften := func(err Error) {
        legalityCheckedTooOften(err)
    }
    defer panicChecks(inPanicLegalityCheckedTooOften)

    game := NewGame(9)
    dumpFile := relPathToAbs("../../../data/tmp/TestDoubleKo.GTPsequence.tmp")
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
    // set the error handling for panic checks to the appropriate context
    legalityCheckedTooOften = func(err Error) {
        failLegalityCheckedTooOften(dumpFile, game, err, t)
    }
    game.PlaySequence(sequence)
    legalBlack, _ := game.Board.calculateIfLegal(4,4,Black)
    legalWhite, _ := game.Board.calculateIfLegal(4,4,White)
    if !legalBlack || !legalWhite {
        t.Fatalf("wrong legality status at (%d,%d)", 4,4)
    }
    legalBlack, _ = game.Board.calculateIfLegal(2,2,Black)
    legalWhite, _ = game.Board.calculateIfLegal(2,2,White)
    if !legalBlack || !legalWhite {
        t.Fatalf("wrong legality status at (%d,%d)", 2,2)
    }
    game.PlayMove(4,4,White)
    legalBlack, _ = game.Board.calculateIfLegal(4,5,Black)
    legalWhite, _ = game.Board.calculateIfLegal(4,5,White)
    if legalBlack || !legalWhite {
        t.Fatalf("wrong legality status at (%d,%d)", 4,5)
    }
    game.PlayMove(2,2,Black)
    legalBlack, _ = game.Board.calculateIfLegal(2,1,Black)
    legalWhite, _ = game.Board.calculateIfLegal(2,1,White)
    if !legalBlack || legalWhite {
        t.Fatalf("wrong legality status at (%d,%d)", 2,1)
    }
    legalBlack, _ = game.Board.calculateIfLegal(4,5,Black)
    legalWhite, _ = game.Board.calculateIfLegal(4,5,White)
    if !legalBlack || !legalWhite {
        t.Fatalf("after 2nd ko: wrong legality status at (%d,%d)", 4,5)
    }
}

func dumpSequence(w writeStringer, game *Game, t *testing.T) {
    for _, mv := range game.sequence {
        m, _ := mv.(Move)
        vertex, _ := pointToGTPVertex(*NewPoint(m.Vertex.X, m.Vertex.Y))
        line := fmt.Sprintf("  play %s %s\n", colorToGTPColor(m.Color), vertex)
        if _, werr := w.WriteString(line); werr != nil {
            t.Fatalf("dumpSequence does not work, check the writeStringer")
        }
    }
}

func failBecausePositionIsIllegal(w writeStringer, failMessage, infoString string, game *Game, t *testing.T, dumpFile string) {
    time := time.LocalTime()
    w.WriteString(fmt.Sprintf("# These moves lead to an illegal position. %s\n", time))
    _, werr := w.WriteString(fmt.Sprintf("  boardsize %d\n", game.Board.BoardSize()))
    if werr != nil {
        failMessage = fmt.Sprintf("The sequence should have been dumped into %s, but this file could not be opened.\n", dumpFile)
        failMessage += fmt.Sprintf("The error was: %s\n", werr)
    } else {
        dumpSequence(w, game, t)
        line := fmt.Sprintf("# %s\n", infoString)
        w.WriteString(line)

    }
    t.Fatalf(failMessage)
}

func fileFail(failMessage, infoString, dumpFile string, game *Game, t *testing.T) {
    os.Remove(dumpFile)
    if file, err := os.Open(dumpFile, os.O_CREATE | os.O_RDWR, 0666); err == nil {
        failBecausePositionIsIllegal(file, failMessage, infoString, game, t, dumpFile)
    } else {
        t.Fatalf("Tried to create the error output file %s, but this error occured: %s", dumpFile, err)
    }
}

func gameStateCheck(game *Game,
                    dumpFile string,
                    t *testing.T,
                    nGame, nMove int,
                    legalBlack, legalWhite []Point,
                    ) {


    for _, p := range legalBlack {
        if gptr := game.Board.GetGroup(p.X, p.Y); gptr != nil {
            illegalVertex, _ := pointToGTPVertex(*NewPoint(p.X,p.Y))
            failMessage := fmt.Sprintf("Game %d, move %d: the point %s was legal for black but is already occupied\n", nGame, nMove, illegalVertex)
            failMessage += fmt.Sprintf("The sequence was dumped into %s", dumpFile)
            infoString := fmt.Sprintf("the vertex %s is legal for black but already occupied", illegalVertex)
            fileFail(failMessage, infoString, dumpFile, game, t)
        }
    }
    for _, p := range legalWhite {
        if gptr := game.Board.GetGroup(p.X, p.Y); gptr != nil {
            illegalVertex, _ := pointToGTPVertex(*NewPoint(p.X,p.Y))
            failMessage := fmt.Sprintf("Game %d, move %d: the point %s was legal for white but is already occupied\n", nGame, nMove, illegalVertex)
            failMessage += fmt.Sprintf("The sequence was dumped into %s", dumpFile)
            infoString := fmt.Sprintf("the vertex %s is legal for white but already occupied", illegalVertex)
            fileFail(failMessage, infoString, dumpFile, game, t)
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
            infoString := fmt.Sprintf("# the group with 0 libs is around %s", vertex)
            fileFail(failMessage, infoString, dumpFile, game, t)
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
        fileFail(failMessage, infoString, dumpFile, game, t)
    }
}

// used for some panic situations that occur while debugging komoku
func panicChecks(legalityCheckedTooOften func(err Error)) {
    //fmt.Printf("entering panicChecks\n")
    if e := recover(); e != nil {
        //fmt.Printf("after first if\n")
        // catch only komoku errors, repanic otherwise
        //printDbgMsgBTf(4, "its an error != nil\n")
        if err, ok := e.(Error); ok {
            //fmt.Printf("after successful type assertion, errno: %d\n", err.Errno())
            //printDbgMsgBTf(4, "its a komoku error\n")
            switch err.Errno() {
                case ErrFieldLegalityCheckedMoreThanOnce:
                    //printDbgMsgBTf(4, "calling legalityCheckedTooOften\n")
                    legalityCheckedTooOften(err)
                default:
                    // repanic
                    panic(e)
            }
        } else {
            panic(e)
        }
    }
}

func failLegalityCheckedTooOften(dumpFile string, game *Game, err Error, t *testing.T) {
    backtrace := sPrintDbgMsgBTf(13,"")
    failMessage := fmt.Sprintf("'%s', bt: '%s'.\nSeq dumped into %s", err, backtrace, dumpFile)
    infoString := fmt.Sprintf("'%s', bt: '%s'.", err, backtrace)
    fileFail(failMessage, infoString, dumpFile, game, t)
}

// Generates random games and checks if the []Points returned by Board.ListLegalPoints do not intersect
// already occupied points
func TestListLegalPoints(t *testing.T) {
    /*testName := "TestListLegalPoints"
    fmt.Printf("entering %s\n", testName)
    defer fmt.Printf("leaving %s\n", testName)*/

    legalityCheckedTooOften := func(err Error) {
        printDbgMsgBTf(4, "legalityCheckedTooOften not yet initialized!\n")
        t.Fatalf("legalityCheckedTooOften not yet initialized!\n")
    }
    inPanicLegalityCheckedTooOften := func(err Error) {
        legalityCheckedTooOften(err)
    }
    defer panicChecks(inPanicLegalityCheckedTooOften)

    numGames := 500 // Number of games this test should play
    gamesLen := 100 // Number of random moves to play
    boardsize := 9
    dumpFile := relPathToAbs("../../../data/tmp/TestListLegalPoints.GTPsequence.tmp")
    lastMovePass := false
    for nGame := 0; nGame < numGames; nGame++ {
        //fmt.Printf("Game %d\n", nGame)
        game := NewGame(boardsize)
        var currentColor Color = Black
        for nMove := 0; nMove < gamesLen; nMove++ {

            // set the error handling for panic checks to the appropriate context
            legalityCheckedTooOften = func(err Error) {
                failLegalityCheckedTooOften(dumpFile, game, err, t)
            }
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
