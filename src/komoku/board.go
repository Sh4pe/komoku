/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */
package komoku

import (
    "fmt"
    //"runtime"
    //"os"
)

/*
 * This file defines the Board object. This object handles the representation of one current
 * state in a game. It provides the methods for creating legal games.
 */

/*
 * TODO:
 *      - make this all more Go-ideomatic. Can the use of interfaces make IntList obsolete?
 *      - make updateLegalityFor take a color parameter which specifies the color for which legality
 *        should be checked.
 *      - make the usage of color more generic. A lot of 'ifs' might be removed by this. A bad example
 *        might be the code around line 171...
 *      - clean up the mess between (x,y) vs. pos...
 *      - investigate on the strange panics inside board_test...
 */


// ######################## Auxiliary type for Board struct #############################
// The return values are only relevant if we assume that the legalities for all fields have been 
// calculated before. In this case, {black,white}Updated indicates that all legalities are valid 
// afterwards.
type actionFunc func() (blackUpdated, whiteUpdated bool)

// Used for ko. Says that a play of color 'Color' at 'Point' is forbidden by the ko rule
type koLock struct {
    Point Point
    Color Color
}

func NewKoLock(x,y int, color Color) *koLock {
    return &koLock{ Point: *NewPoint(x,y),
                    Color: color,
                  }
}

// ################################################################################
// ########################### Board struct #######################################
// ################################################################################

// This object is responsible for recording a current state of a game.
type Board struct {
    fields []*Group // Stores pointers to the groups. nil denotes an empty field
    ko *koLock // nil means that there is no ko
    actionOnNextBlackMove []actionFunc // This stores the appropriate code which has to be run if a black move is played on a field
    actionOnNextWhiteMove []actionFunc // see the obvious analogue
    colorOfNextPlay Color
    boardSize int
    acBlackMoveUpToDate, acWhiteMoveUpToDate bool
    currentSequence uint32 // Represents the state of the board. Everything flagged with a different sequence has to be updated.
    fieldSequencesBlack []uint32
    fieldSequencesWhite []uint32
}

// ##################### Board methods ##########################

// Returns the board size.
func (b *Board) BoardSize() int {
    return b.boardSize
}

// Calculates if a move of 'color' at (x,y) is legal. Does not use 
// b.legal{Black,White}Moves for this. Also, this method returnsthe appropriate
// `actions` to be performed if this move is played using 'action'. The return value
// of these actions indicate if any legality update has been performed.
// Note that this method assumes that (x,y) is empty.
func (b *Board) calculateIfLegal(x,y int, color Color) (isLegal bool, action actionFunc) {
    /*v, _ := pointToGTPVertex(*NewPoint(x, y))
    printDbgMsgBTf(4,"entering Board.calculateIfLegal(%s, %v)\n", v, color) // <DBG>*/
    //printDbgMsgf("ko status: %v\n", b.ko)

    // Is this move prohibited because of a ko? It is not prohibited for the player who 
    // took the ko to fill it in the next move
    if b.ko != nil && b.ko.Color == color&& b.ko.Point.X == x && b.ko.Point.Y == y {
        //printDbgMsgf("returned from Board.calculateIfLegal(%d, %d, %v)\n", x, y, color) // <DBG/>
        return false, nil
    }
    pos := b.xyToPos(x,y)
    nFree, adjSameColor, adjOtherColor := b.GetEnvironment(x,y)
    if color == White {
        adjSameColor, adjOtherColor = adjOtherColor, adjSameColor
    }
    enemiesInAtari, enemiesNotInAtari := b.determineGroupsAtariStatus(adjOtherColor)
    /*printDbgMsgf("sameColLen: %d, othColLen: %d, enemInAtariLen: %d, enemNotInAtariLen: %d\n", len(adjSameColor), len(adjOtherColor), 
                len(enemiesInAtari), len(enemiesNotInAtari)) // <DBG>*/
    removeGroups := len(enemiesInAtari) > 0

    if len(adjSameColor) == 0 {
        if nFree == 0 {
            if !removeGroups {
                // There are no adjacent friendly groups, no neighbour is a free field and 
                // this move does not capture enemy stones. So it is illegal.

                //printDbgMsgf("returned from Board.calculateIfLegal(%d, %d, %v)\n", x, y, color) // <DBG/>
                return false, nil
            } else {
                // There are no adjacent friendly groups and no free neighbours, but this 
                // move captures, so it is legal. 

                // But first we have to check if this move is a ko play. If we capture exactly one
                // group consisting of exactly one stone, than it's a ko.
                firstGroup := enemiesInAtari[0]
                if len(enemiesInAtari) == 1 && firstGroup.Fields.Length() == 1 {
                    // It's a ko, so remove the group, play the stone, update the liberties
                    // and set b.ko to the right point.
                    koPos := firstGroup.Fields.First().Value()
                    koX, koY := b.posToXY(koPos)
                    action = func() (blackUpToDate, whiteUpToDate bool) {
                        //printDbgMsgf("Board.calculateIfLegal: sameColLen == nFree == 0, removeGroups = true, ko case.\n") // <DBG>
                        //DbgHistogram.Score() // </DBG>

                        alreadyUpdated := make([]bool, b.boardSize*b.boardSize)
                        b.ko = nil
                        b.removeGroup(enemiesInAtari[0]) // Remove the enemy stone
                        b.CreateGroup(x,y,color) // Create the new group
                        for _, grp := range enemiesNotInAtari { // Update liberties of the groups adjacent to the new stone at (x,y)
                            b.updateGroupLiberties(grp)
                        }
                        // Update liberties and legality for the groups adjacent to the removed stone
                        var adjToKoSameColor GroupSlice
                        if color == Black {
                            _, adjToKoSameColor, _ = b.GetEnvironment(koX, koY)
                        } else {
                            _, _, adjToKoSameColor = b.GetEnvironment(koX, koY)
                        }
                        for _, grp := range adjToKoSameColor {
                            grp.Liberties.AppendUnique(koPos)
                            b.updateLegalityForLibertiesOfExcept(grp, koX, koY, b.currentSequence + 1, alreadyUpdated)
                        }
                        // The player who took the ko may fill it, so make it legal for the player 'color' for the next round.
                        // TODO: this is sort of an evil hack... or is it?
                        if color == Black {
                            _, b.actionOnNextBlackMove[koPos] = b.calculateIfLegal(koX, koY, Black)
                            b.fieldSequencesBlack[koPos] = b.currentSequence + 1
                        } else {
                            _, b.actionOnNextWhiteMove[koPos] = b.calculateIfLegal(koX, koY, White)
                            b.fieldSequencesWhite[koPos] = b.currentSequence + 1
                        }
                        // Update the legality for the liberties of the groups adjacent to the new created stone
                        //printDbgMsg("Update the legality for the liberties of the groups adjacent to the new created stone\n") // </DBG>
                        for _, grp := range enemiesNotInAtari {
                            b.updateLegalityForLibertiesOf(grp, b.currentSequence + 1, alreadyUpdated)
                        }

                        b.ko = NewKoLock(koX, koY, !color)
                        return true, true
                    }
                } else {
                    // It's not a ko
                    action = func() (blackUpToDate, whiteUpToDate bool) {
                        //printDbgMsgf("Board.calculateIfLegal: sameColLen == nFree == 0, removeGroups = true, not ko case.\n") // <DBG>
                        //DbgHistogram.Score() // </DBG>
                        for _, grp := range enemiesInAtari {
                            b.removeGroup(grp)
                        }
                        b.CreateGroup(x,y,color)
                        for _, grp := range enemiesNotInAtari {
                            b.updateGroupLiberties(grp)
                        }
                        b.ko = nil
                        return false, false
                    }
                }
                //printDbgMsgf("returned from Board.calculateIfLegal(%d, %d, %v)\n", x, y, color) // <DBG/>
                return true, action
            }
        } else {
            // There are no adjacent friendly groups, but free neighbour fields, so this move
            // is always legal. Remove adjacent enemy groups if necessary and create a new group.
            if removeGroups {
               action = func() (blackUpToDate, whiteUpToDate bool) {
                    //printDbgMsgf("Board.calculateIfLegal: sameColLen == 0, nFree > 0, removeGroups = true.\n") // <DBG>
                    //DbgHistogram.Score() // </DBG>
                    for _, grp := range enemiesInAtari {
                        b.removeGroup(grp)
                    }
                    b.CreateGroup(x,y,color)
                    for _, grp := range enemiesNotInAtari {
                        b.updateGroupLiberties(grp)
                    }
                    b.ko = nil
                    return false, false
                }
            } else {
                action = func() (blackUpToDate, whiteUpToDate bool) {
                    // Experiments show that this case is run the most often

                    //printDbgMsgf("Board.calculateIfLegal: sameColLen == 0, nFree > 0, removeGroups = false.\n") // <DBG>
                    //DbgHistogram.Score() // </DBG>
                    alreadyUpdated := make([]bool, b.boardSize*b.boardSize)
                    b.ko = nil
                    b.CreateGroup(x,y,color)
                    b.dropLibertyFromEach(pos, adjOtherColor)
                    b.updateLegalityForFreeNeighboursOf(x,y, b.currentSequence + 1, alreadyUpdated)
                    b.updateLegalityForAdjacentGroups(adjOtherColor, b.currentSequence + 1, alreadyUpdated)
                    return true, true
                }
            }
            //printDbgMsgf("returned from Board.calculateIfLegal(%d, %d, %v)\n", x, y, color) // <DBG/>
            return true, action
        }
    } else {
        if nFree == 0 {
            if removeGroups {
                // This move captures stones and thus produces empty fields, so it is legal. Capture
                // the stones first and then join the adjacent groups of the same color.
                action = func() (blackUpToDate, whiteUpToDate bool) {
                    //printDbgMsgf("Board.calculateIfLegal: sameColLen > 0, nFree == 0, removeGroups = true.\n") // <DBG>
                    //DbgHistogram.Score() // </DBG>
                    for _, grp := range enemiesInAtari {
                        b.removeGroup(grp)
                    }
                    b.joinGroupsByPlayAt(pos, adjSameColor)
                    for _, grp := range enemiesNotInAtari {
                        b.updateGroupLiberties(grp)
                    }
                    b.ko = nil
                    return false, false
                }
                //printDbgMsgf("returned from Board.calculateIfLegal(%d, %d, %v)\n", x, y, color) // <DBG/>
                return true, action
            } else {
                // This move does not capture any stones, so its legality depends upon the liberties 
                // of the adjacent groups of the same color. We check if at least one of these group
                // has at least two liberties.
                oneHasTwo := false
                for _, g := range adjSameColor {
                    if g.Liberties.Length() > 1 {
                        /*fX, fY := b.posToXY(g.Fields.First().Value())
                        grpVertex, _ := pointToGTPVertex(*NewPoint(fX, fY))
                        Libs := g.Liberties
                        libs := ""
                        libEnd := Libs.Last()
                        for it := Libs.First(); it != libEnd; it = it.Next() {
                            lX, lY := b.posToXY(it.Value())
                            ver, _ := pointToGTPVertex(*NewPoint(lX, lY))
                            libs += " " + ver
                        }
                        printDbgMsgf("group around %s has many libs: %s\n", grpVertex, libs)*/
                        oneHasTwo = true
                        break
                    }
                }
                if oneHasTwo {
                    // If we join the groups, the resulting group has at least one liberty, so this move is legal.
                    // Since there are no groups to capture, simply join the adjacient groups of color 'color'.
                    action = func() (blackUpToDate, whiteUpToDate bool) {
                        // Experiments show that this case is run the 3rd most often

                        //printDbgMsgf("Board.calculateIfLegal: sameColLen > 0, nFree == 0, removeGroups = false, oneHasTwo = true.\n") // <DBG>
                        //DbgHistogram.Score() // </DBG>
                        alreadyUpdated := make([]bool, b.boardSize*b.boardSize)
                        b.ko = nil
                        b.joinGroupsByPlayAt(pos, adjSameColor)
                        b.dropLibertyFromEach(pos, adjOtherColor)
                        b.updateLegalityForAdjacentGroups(adjOtherColor, b.currentSequence + 1, alreadyUpdated)
                        b.updateLegalityForLibertiesOf(b.fields[pos], b.currentSequence + 1, alreadyUpdated)

                        return true, true
                    }
                    //printDbgMsgf("returned from Board.calculateIfLegal(%d, %d, %v)\n", x, y, color) // <DBG/>
                    return true, action
                } else {
                    // There is no groups to remove and every adjacient group of the same color has only one liberty,
                    // which must be the field (x,y) we want to play at, so this move is illegal.

                    //printDbgMsgf("returned from Board.calculateIfLegal(%d, %d, %v)\n", x, y, color) // <DBG/>
                    return false, nil
                }
            }
        } else {
            // There are free neighbour fields, so this move is always legal. Capture adjacent enemy groups if necessary, 
            // then join groups and update liberties
            if removeGroups {
                action = func() (blackUpToDate, whiteUpToDate bool) {
                    //printDbgMsgf("Board.calculateIfLegal: sameColLen > 0, nFree > 0, removeGroups = true.\n") // <DBG>
                    //DbgHistogram.Score() // </DBG>
                    for _, grp := range enemiesInAtari {
                        b.removeGroup(grp)
                    }
                    b.joinGroupsByPlayAt(pos, adjSameColor)
                    for _, grp := range enemiesNotInAtari {
                        b.updateGroupLiberties(grp)
                    }
                    b.ko = nil
                    return false, false
                }
            } else {
                action = func() (blackUpToDate, whiteUpToDate bool) {
                    // Experiments show that this case is run the 2nd most often

                    //printDbgMsgf("Board.calculateIfLegal: sameColLen > 0, nFree > 0, removeGroups = false.\n") // <DBG>
                    //DbgHistogram.Score() // </DBG>
                    alreadyUpdated := make([]bool, b.boardSize*b.boardSize)
                    b.ko = nil
                    b.joinGroupsByPlayAt(pos, adjSameColor)
                    for _, grp := range enemiesNotInAtari {
                        b.updateGroupLiberties(grp)
                    }
                    b.updateLegalityForAdjacentGroups(adjOtherColor, b.currentSequence + 1, alreadyUpdated)
                    // update legality for the newly joined group, which is at pos
                    b.updateLegalityForLibertiesOf(b.fields[pos], b.currentSequence + 1, alreadyUpdated)

                    return true, true
                }
            }
            //printDbgMsgf("returned from Board.calculateIfLegal(%d, %d, %v)\n", x, y, color) // <DBG/>
            return true, action
        }
    }

    // Control flow should never reach here; if it does however, this is an error
    panic("Control reaced the end of Board.calculateIfLegal")
    return false, nil
}

// Returns the color of the player who plays the next turn
func (b *Board) ColorOfNextPlay() Color {
    return b.colorOfNextPlay
}

// Create a new, one-stone-group of 'color' at (x,y) and sets its liberties appropriately. 
// This method does not perform any legality checks or liberty updates for other groups.
// TODO: unexport this?
func (b *Board) CreateGroup(x, y int, color Color) {
    pos := b.xyToPos(x,y)
    newGroup := NewGroup(color)
    newGroup.Fields.Append(pos)
    nbours := b.neighbours(x,y)
    for _, p := range nbours {
        npos := b.xyToPos(p.X, p.Y)
        if b.fields[npos] == nil {
            newGroup.Liberties.Append(npos)
        }
    }
    b.fields[pos] = newGroup
    b.actionOnNextBlackMove[pos] = nil
    b.actionOnNextWhiteMove[pos] = nil
}

// Returns the sublists '{not,}inAtari' of 'groups' which are {not,} in atari. Assumes that the groups in
// 'groups' are pairwise different.
// TODO: write tests for this...
func (b *Board) determineGroupsAtariStatus(groups GroupSlice) (inAtari, notinAtari GroupSlice) {
    inAtari = NewGroupSlice()
    notinAtari = NewGroupSlice()
    for _, group := range groups {
        if group.Liberties.Length() == 1 {
            inAtari.Push(group)
        } else {
            notinAtari.Push(group)
        }
    }
    return
}

// Helper for Board.calculateIfLegal. Takes each group (identified by one index) in adjGroups and drops
// the liberty posToXY(libertyPos) from it
func (b *Board) dropLibertyFromEach(libertyPos int, adjGroups GroupSlice) {
    for _, grp := range adjGroups {
        grp.Liberties.Remove(libertyPos)
    }
}

// Returns the `environment` of (x,y), i.e. the number 'nFree' of free neighbours
// and IntLists 'adj{Black,White}' containing instances of GroupIndexTypes of adjacent {black,white} groups.
// TODO: unexport this?
func (b *Board) GetEnvironment(x,y int) (nFree int, adjBlack, adjWhite GroupSlice) {
    nbours := b.neighbours(x,y)
    adjBlack = NewGroupSlice()
    adjWhite = NewGroupSlice()
    for _, p := range nbours {
        npos := b.xyToPos(p.X, p.Y)
        if b.fields[npos] == nil {
            nFree++
        } else {
            group := b.fields[npos]
            // Don't append the same group more than once!
            if group.Color == Black {
                adjBlack.PushUnique(group)
            } else {
                adjWhite.PushUnique(group)
            }
        }
    }
    return
}

// Analogous to GetEnvironment.
// TODO: write tests for this...
// TODO: unexport this?
func (b *Board) GetEnvironmentByPos(pos int) (nFree int, adjBlack, adjWhite GroupSlice) {
    x, y := b.posToXY(pos)
    return b.GetEnvironment(x,y)
}

// Returs a pointer to the group which occupies (x,y). Nil means that this 
// field is empty.
func (b *Board) GetGroup(x,y int) *Group {
    index := b.xyToPos(x,y)
    return b.fields[index]
}

// Is it legal to play a stone of color 'color' at (x,y)?
func (b *Board) IsLegalMove(x, y int, color Color) bool {
    pos := b.xyToPos(x,y)
    /*printDbgMsgf("IsLegalMove(%d, %d, %s): fieldSeqW[pos]: %d, fieldSeqB[pos]: %d, currSeq: %d\n", x,y,color, b.fieldSequencesWhite[pos], b.fieldSequencesBlack[pos],
                b.currentSequence)*/
    if color == Black {
        if b.fieldSequencesBlack[pos] != b.currentSequence {
            b.updateLegalityFor(pos, b.currentSequence)
        }
        if b.actionOnNextBlackMove[pos] != nil {
            return true
        }
    } else {
        if b.fieldSequencesWhite[pos] != b.currentSequence {
            b.updateLegalityFor(pos, b.currentSequence)
        }
        if b.actionOnNextWhiteMove[pos] != nil {
            return true
        }
    }
    return false
}

// Joins the groups 'into' and 'from'. The stones from 'from' become stones of
// 'into' and 'into' is deleted. One could write the effect of this method 
// as `` into += from ''. Note that this method does not update 'into's liberties
// afterwards, you'll have to do this manually if you want it.
func (b *Board) joinGroups(into, from *Group) {
    lastFrom := from.Fields.Last()
    for fromIt := from.Fields.First(); fromIt != lastFrom; fromIt = fromIt.Next() {
        fpos := fromIt.Value()
        b.fields[fpos] = into
    }
    into.Fields.JoinUnique(from.Fields)
}

// This is a helper function for Board.calculateIfLegal. It joint the adjacent
// groups in 'adjSameColor' by playing a stone of the same color at posToXY(playPos),
// but does not update legalities.
// This func is _not_safe to call if there are no groups to join (i.e. adjSameColor empty)
func (b *Board) joinGroupsByPlayAt(playPos int, adjSameColor GroupSlice) {
    //printDbgMsg("in joinGroupsByPlayAt\n") // <DBG>
    //defer printDbgMsgf("returned from joinGroupsByPlayAt\n") // </DBG>

    firstGroup := adjSameColor[0]
    // Add the stone at posToXY(playPos) to the first group.
    firstGroup.Fields.Append(playPos)
    b.fields[playPos] = firstGroup
    b.actionOnNextBlackMove[playPos] = nil
    b.actionOnNextWhiteMove[playPos] = nil

    // Then join the other groups into the first group
    for i := 1; i < len(adjSameColor); i++ {
        b.joinGroups(firstGroup, adjSameColor[i])
    }
    // Finally update the liberties for the joined group.
    b.updateGroupLiberties(firstGroup)
}

func (b *Board) legalMovesNeedUpdate() {
    b.acBlackMoveUpToDate = false
    b.acWhiteMoveUpToDate = false
}

// Returns a slice containing the empty fields of b.
func (b *Board) ListEmptyFields() []*Point {
    ret := make([]*Point, b.boardSize*b.boardSize)
    index := 0
    for i := 0; i < b.boardSize*b.boardSize; i++ {
        if b.fields[i] == nil {
            x, y := b.posToXY(i)
            ret[index] = NewPoint(x,y)
            index++
        }
    }
    return ret[0:index]
}

// Returns a slice of legal moves of color 'color'
func (b *Board) ListLegalPoints(color Color) []Point {
    b.updateLegalMoves(color)

    var actions []actionFunc
    if color == Black {
        actions = b.actionOnNextBlackMove
    } else {
        actions = b.actionOnNextWhiteMove
    }
    ret := make([]Point, b.boardSize*b.boardSize)
    index := 0
    for i := 0; i < b.boardSize*b.boardSize; i++ {
        if actions[i] != nil {
            x, y := b.posToXY(i)
            //v, _ := pointToGTPVertex(*NewPoint(x,y))
            //printDbgMsgf("actions[%s] != nil\n", v)
            ret[index] = *NewPoint(x,y)
            index++
        }
    }
    return ret[0:index]
}

// Returns the neighbours of a field (x,y) as a slice of Points.
func (b *Board) neighbours(x, y int) []Point {
    // TODO: can this be implemented better?
    ret := make([]Point, 4)
    count := 0
    switch x {
        case 0:
            ret[count] = Point{ 1, y }
            count++
        case b.BoardSize()-1:
            ret[count] = Point{ b.BoardSize()-2, y }
            count++
        default:
            ret[count] = Point{ x-1, y }
            count++
            ret[count] = Point{ x+1, y }
            count++
    }
    switch y {
        case 0:
            ret[count] = Point{ x, 1 }
            count++
        case b.BoardSize()-1:
            ret[count] = Point{ x, b.BoardSize()-2 }
            count++
        default:
            ret[count] = Point{ x, y-1 }
            count++
            ret[count] = Point{ x, y+1 }
            count++
    }
    return ret[0:count]
}

// Returns the number of {black,white} groups in 'n{black,white}'
func (b *Board) numberOfGroups() (nblack, nwhite int) {
    // Note: this method is only used by the GTP-command komoku-numgroups, so
    // the implementation doesn't need to be very efficient.
    uniqueBlack := make(map[*Group]bool)
    uniqueWhite := make(map[*Group]bool)
    for _, grp := range b.fields {
        if grp != nil {
            if grp.Color == Black {
                uniqueBlack[grp] = true
            } else {
                uniqueWhite[grp] = true
            }
        }
    }
    nblack, nwhite = len(uniqueBlack), len(uniqueWhite)
    return
}

// Returns the number of {black,white} stones on the board
func (b *Board) numberOfStones() (nblack, nwhite int) {
    nblack, nwhite = 0,0
    for _, grp := range b.fields {
        if grp != nil {
            if grp.Color == Black {
                nblack++
            } else {
                nwhite++
            }
        }
    }
    return
}

// Play a move of color 'color' at (x,y)
func (b *Board) PlayMove(x, y int, color Color) (err Error) {
    pos := b.xyToPos(x,y)

    // Check if this is legal if necessary

    if !b.IsLegalMove(x,y, color) {
        return NewIllegalMoveError(x,y, color)
    }

    var action actionFunc
    if color == White {
        action = b.actionOnNextWhiteMove[pos]
    } else {
        action = b.actionOnNextBlackMove[pos]
    }
    blackUpToDate, whiteUpToDate := action()
    if blackUpToDate || whiteUpToDate {
        for i := 0; i < b.BoardSize()*b.BoardSize(); i++ {
            //pX, pY := b.posToXY(i)
            //v, _ := pointToGTPVertex(*NewPoint(pX, pY))
            /*printDbgMsgf("PlayMove after action at %s, bSeq(pos): %d, wSeq(pos): %d, currSeq: %d\n", v, b.fieldSequencesBlack[i], b.fieldSequencesWhite[i],
                          b.currentSequence)*/
            if blackUpToDate && b.fieldSequencesBlack[i] <= b.currentSequence {
                //printDbgMsgf("black updated at %s\n", v)
                b.fieldSequencesBlack[i]++
            }
            if whiteUpToDate && b.fieldSequencesWhite[i] <= b.currentSequence {
                //printDbgMsgf("white updated at %s\n", v)
                b.fieldSequencesWhite[i]++
            }
        }
    }
    // Clear the appropriate actionOnNextMove array. 
    b.colorOfNextPlay = !color
    /*if b.colorOfNextPlay == White && !b.acWhiteMoveUpToDate {
        // David: This is really important so that the GC can free the closures with all the
        // associated contexts. Am I right?
        for i := 0; i < b.BoardSize()*b.BoardSize(); i++ {
            b.actionOnNextBlackMove[i] = nil
        }
        b.updateLegalMoves(Black)
        return nil
    }
    if b.colorOfNextPlay == Black && !b.acBlackMoveUpToDate {
        for i := 0; i < b.BoardSize()*b.BoardSize(); i++ {
            b.actionOnNextWhiteMove[i] = nil
        }
        b.updateLegalMoves(White)
    }*/
    b.currentSequence++

    return nil
}

// The player of color 'color' plays a pass.
func (b *Board) PlayPass(color Color) {
    /*nextActions := b.actionOnNextBlackMove
    if color == Black {
        nextActions = b.actionOnNextWhiteMove
    }
    for i := 0; i < b.BoardSize()*b.BoardSize(); i++ {
        nextActions[i] = nil
    }*/
    b.colorOfNextPlay = !color
    //b.updateLegalMoves(b.colorOfNextPlay)
    b.currentSequence++
}

// Plays out the sequence 'seq' of moves
func (b *Board) playSequence(seq []Move) {
    for _, m := range seq {
        if m.Vertex.Pass {
            b.PlayPass(m.Color)
        } else {
            //printDbgMsgf("playing %s stone at (%d,%d)\n", m.Color, m.Vertex.X, m.Vertex.Y)
            b.PlayMove(m.Vertex.X, m.Vertex.Y, m.Color)
        }
    }
}

// TODO: use point!
func (b *Board) posToXY(pos int) (x, y int) {
    return pos%b.BoardSize(), pos/b.BoardSize()
}

// Removes the group which occupies (x,y), if there is any, and updates b.emptyFields.
// This method does not alter legalWhiteMoves or legalBlackMoves.
// TODO(david): Do I really want to export this?
func (b *Board) RemoveGroupByPos(x,y int) {
    pos := b.xyToPos(x,y)
    if grp := b.fields[pos]; grp != nil {
        b.removeGroup(grp)
    }
}

// Removes the group which has the index 'gid' in .groupMap and updates b.emptyFields.
// Also, this method calls updateGroupLiberties for each group whose liberties might change
// because of the removal of 'group'.
// This method does not alter legalities for any fields.
func (b *Board) removeGroup(group *Group) {
    last := group.Fields.Last()
    adjGroups := NewGroupSlice()
    for e := group.Fields.First(); e != last; e = e.Next() {
        // Collect adjacend groups so that we can update their liberties later
        pos := e.Value()
        x, y := b.posToXY(pos)
        nbours := b.neighbours(x,y)
        for _, p := range nbours {
            npos := b.xyToPos(p.X, p.Y)
            if grp := b.fields[npos]; grp != nil {
                adjGroups.PushUnique(grp)
            }
        }
        b.fields[pos] = nil
        //b.emptyFields.Append(e.Value())
    }
    // Update the liberties of the adjacent groups.
    for _, grp := range adjGroups {
        b.updateGroupLiberties(grp)
    }
}


// The player whose turn it is plays a stone (x,y). If an error occurs (such as that 
// this place is already occupied) this error is returned. This method assumes that 
// b.actionOnNextMove is correcty set
// TODO: write tests for this...
func (b *Board) TurnPlayMove(x, y int) (err Error) {
    return b.PlayMove(x,y, b.colorOfNextPlay)
}

// The player, whose turn it is, plays a pass.
func (b *Board) TurnPlayPass() {
    b.PlayPass(b.colorOfNextPlay)
}

// Recomputes the liberties of 'group' based only on the currently occupied and empty
// fields.
func (b *Board) updateGroupLiberties(group *Group) {
    // this is expensive... TODO(David): can this be made faster?
    group.Liberties.Clear()
    last := group.Fields.Last()
    for it := group.Fields.First(); it != last; it = it.Next() {
        x, y := b.posToXY(it.Value())
        nbours := b.neighbours(x,y)
        for _, p := range nbours {
            npos := b.xyToPos(p.X, p.Y)
            if b.fields[npos] == nil {
                group.Liberties.AppendUnique(npos)
            }
        }
    }
}

// Updates the legality for posToXY(pos) and makes sure the state of the board remains correct. 
// 'whichSequence' denotes the sequence to set in b.fieldSequences{Black,White}.
func (b *Board) updateLegalityFor(pos int, whichSequence uint32) {
    pX, pY := b.posToXY(pos)
    //v, _ := pointToGTPVertex(*NewPoint(pX, pY))
    //printDbgMsgBTf(4,"in updateLegalityFor(%s, %d), currentSequence: %d\n", v, whichSequence, b.currentSequence)

    // use this 'if' only for debugging!
    /*if b.fieldSequencesBlack[pos] == whichSequence || b.fieldSequencesWhite[pos] == whichSequence {
        panic(NewFieldLegalityCheckedMoreThanOnceError(fmt.Sprintf("checked more than once for %s",v)))
    }*/

    _, b.actionOnNextBlackMove[pos] = b.calculateIfLegal(pX, pY, Black)
    _, b.actionOnNextWhiteMove[pos] = b.calculateIfLegal(pX, pY, White)
    b.fieldSequencesBlack[pos] = whichSequence
    b.fieldSequencesWhite[pos] = whichSequence
}

// Checks if a black move at 'pos' is legal and makes sure the state of the board remains correct. 
// 'whichSequence' denotes the sequence to set in b.fieldSequences{Black,White}.
func (b *Board) updateLegalityForBlack(pos int, whichSequence uint32) {
    pX, pY := b.posToXY(pos)
    _, b.actionOnNextBlackMove[pos] = b.calculateIfLegal(pX, pY, Black)
    b.fieldSequencesBlack[pos] = whichSequence
}

// Checks if a white move at 'pos' is legal and makes sure the state of the board remains correct. 
// 'whichSequence' denotes the sequence to set in b.fieldSequences{Black,White}.
func (b *Board) updateLegalityForWhite(pos int, whichSequence uint32) {
    pX, pY := b.posToXY(pos)
    _, b.actionOnNextWhiteMove[pos] = b.calculateIfLegal(pX, pY, White)
    b.fieldSequencesWhite[pos] = whichSequence
}

// Helper method for Board.calculateIfLegal. Does what it name indicates and skips the fields which are mared in 
// alreadyUpdated.
func (b *Board) updateLegalityForAdjacentGroups(adjGroups GroupSlice, whichSequence uint32, alreadyUpdated []bool) {
    //printDbgMsg("in updateLegalityForAdjacentGroups\n") // <DBG>
    //defer printDbgMsgf("returned from updateLegalityForAdjacentGroups\n") // </DBG>

    for _, grp := range adjGroups {
        lastLib := grp.Liberties.Last()
        for itLib := grp.Liberties.First(); itLib != lastLib; itLib = itLib.Next() {
            lpos := itLib.Value()
            if !alreadyUpdated[lpos] {
                b.updateLegalityFor(lpos, whichSequence)
                alreadyUpdated[lpos] = true
            }
        }
    }
}

// Helper method for Board.calculateIfLegal. Does what it name indicates and skips the fields which are mared in 
// alreadyUpdated.
func (b *Board) updateLegalityForFreeNeighboursOf(x,y int, whichSequence uint32, alreadyUpdated []bool) {
    //printDbgMsg("in updateLegalityForFreeNeighboursOf\n") // <DBG>
    //defer printDbgMsgf("returned from updateLegalityForFreeNeighboursOf\n") // </DBG>
    nbours := b.neighbours(x,y)
    for _, p := range nbours {
        npos := b.xyToPos(p.X, p.Y)
        if b.fields[npos] == nil && !alreadyUpdated[npos] {
            b.updateLegalityFor(npos, whichSequence)
            alreadyUpdated[npos] = true
        }
    }
}

// Helper method for Board.calculateIfLegal. Does what it name indicates and skips the fields which are mared in 
// alreadyUpdated.
func (b *Board) updateLegalityForLibertiesOf(group *Group, whichSequence uint32, alreadyUpdated []bool) {
    //printDbgMsg("in Board.updateLegalityForLibertiesOf\n") // <DBG/>
    //defer printDbgMsg("returned from Board.updateLegalityForLibertiesOf\n") // <DBG/>
    lastLib := group.Liberties.Last()
    for itLib := group.Liberties.First(); itLib != lastLib; itLib = itLib.Next() {
        lpos := itLib.Value()
        if !alreadyUpdated[lpos] {
            b.updateLegalityFor(lpos, whichSequence)
            alreadyUpdated[lpos] = true
        }
    }
}

// Helper method for Board.calculateIfLegal, analogous to updateLegalityForLibertiesOf, except that (exceptX, exceptY) will be 
// left also.
func (b *Board) updateLegalityForLibertiesOfExcept(group *Group, exceptX, exceptY int, whichSequence uint32, alreadyUpdated []bool) {
    //printDbgMsg("in Board.updateLegalityForLibertiesOfExcept\n") // <DBG/>
    //defer printDbgMsg("returned from Board.updateLegalityForLibertiesOfExcept\n") // <DBG/>
    lastLib := group.Liberties.Last()
    for itLib := group.Liberties.First(); itLib != lastLib; itLib = itLib.Next() {
        exceptPos := b.xyToPos(exceptX, exceptY)
        lpos := itLib.Value()
        if exceptPos != lpos && !alreadyUpdated[lpos] {
            b.updateLegalityFor(lpos, whichSequence)
            alreadyUpdated[lpos] = true
        }
    }
}

// Updates the legal moves for the color 'color'
// TODO: write tests for this...
func (b *Board) updateLegalMoves(color Color) {
    if color == Black {
        for i := 0; i < b.BoardSize()*b.BoardSize(); i++ {
            if b.fields[i] == nil && b.fieldSequencesBlack[i] != b.currentSequence {
                b.updateLegalityForBlack(i, b.currentSequence)
            }
        }
    } else {
        for i := 0; i < b.BoardSize()*b.BoardSize(); i++ {
            if b.fields[i] == nil && b.fieldSequencesWhite[i] != b.currentSequence {
                b.updateLegalityForWhite(i, b.currentSequence)
            }
        }
    }
}

func (b *Board) xyToPos(x, y int) int {
    return b.BoardSize()*y + x
}

// ##################### Board helper functions ##########################
// Creates a new, initial board of size 'boardsize'.
func NewBoard(boardsize int) *Board {
    ret := &Board{ fields: make([]*Group, boardsize*boardsize),
                   actionOnNextBlackMove: make([]actionFunc, boardsize*boardsize),
                   actionOnNextWhiteMove: make([]actionFunc, boardsize*boardsize),
                   //legalWhiteMoves: NewIntList(),
                   //legalBlackMoves: NewIntList(),
                   //emptyFields: NewIntList(),
                   colorOfNextPlay: Black,
                   boardSize: boardsize,
                   acBlackMoveUpToDate: true,
                   acWhiteMoveUpToDate: true,
                   fieldSequencesBlack: make([]uint32, boardsize*boardsize),
                   fieldSequencesWhite: make([]uint32, boardsize*boardsize),
                 }
    for i := 0; i < boardsize*boardsize; i++ {
        //ret.legalWhiteMoves.Append(i)
        //ret.legalBlackMoves.Append(i)
        //ret.emptyFields.Append(i)
        // set the initial actions for playing a field
        x, y := ret.posToXY(i)
        ret.actionOnNextBlackMove[i] = func() (updateBlack, updateWhite bool) {
            ret.CreateGroup(x, y, Black)
            ret.colorOfNextPlay = !ret.colorOfNextPlay
            ret.legalMovesNeedUpdate()
            return false, false
        }
        ret.actionOnNextWhiteMove[i] = func() (updateBlack, updateWhite bool) {
            ret.CreateGroup(x, y, White)
            ret.colorOfNextPlay = !ret.colorOfNextPlay
            ret.legalMovesNeedUpdate()
            return false, false
        }
    }
    return ret
}

// Creates a new FieldOccupiedError, indicating that (x,y) is alrady used.
func NewFieldOccupiedError(x, y int) (err Error) {
    return NewError(fmt.Sprintf("(%d,%d) is already occupied", x, y), ErrFieldOccupied)
}

// Creates a new FieldOccupiedError, indicating that (x,y) is alrady used.
func NewIllegalMoveError(x, y int, color Color) (err Error) {
    return NewError(fmt.Sprintf("a %v move at (%d,%d) is illegal", color, x, y), ErrIllegalMove)
}

func NewFieldLegalityCheckedMoreThanOnceError(msg string) Error {
    return NewError(msg, ErrFieldLegalityCheckedMoreThanOnce)
}

