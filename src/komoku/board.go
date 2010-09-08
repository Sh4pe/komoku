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
 *      - think about the cast-mess between int and GroupIndexType. Use interfaces for this!
 */


// ######################## Auxiliary type for Board struct #############################
// If the {black,white} legal moves are up to date after this actions, {black,white}UpToDate 
// indicate this.
type actionFunc func() (blackUpToDate, whiteUpToDate bool)

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
    legalBlackMoves *IntList // Indices of fields at which it is legal to play a black stone.
    legalWhiteMoves *IntList // Indices of fields at which it is legal to play a white stone.
    emptyFields *IntList // indices of empty fields
    ko *koLock // nil means that there is no ko
    actionOnNextBlackMove []actionFunc // This stores the appropriate code which has to be run if a black move is played on a field
    actionOnNextWhiteMove []actionFunc // see the obvious analogue
    colorOfNextPlay Color
    boardSize int
    acBlackMoveUpToDate, acWhiteMoveUpToDate bool
}

// ##################### Board methods ##########################

// Returns the board size.
func (b *Board) BoardSize() int {
    return b.boardSize
}

// Calculates if a move of 'color' at (x,y) is legal. Does not use 
// b.legal{Black,White}Moves for this. Also, this method returnsthe appropriate
// `actions` to be performed if this move is played using 'action'. This 'action' does not include any updates
// of legal{Black,White}Moves... 
// Note that this method assumes that (x,y) is empty.
func (b *Board) calculateIfLegal(x,y int, color Color) (isLegal bool, action actionFunc) {
    //printDbgMsgf("entering Board.calculateIfLegal(%d, %d, %v)\n", x,y,color) // <DBG>
    //printDbgMsgf("ko status: %v\n", b.ko)

    // Is this move prohibited because of a ko?
    if b.ko != nil {
        // It is not prohibited for the player who took the ko to fill it in the next move
        if b.ko.Point.X == x && b.ko.Point.Y == y && b.ko.Color == color {
            //printDbgMsgf("returned from Board.calculateIfLegal(%d, %d, %v)\n", x, y, color) // <DBG/>
            return false, nil
        }
    }
    pos := b.xyToPos(x,y)
    nFree, adjSameColor, adjOtherColor := b.GetEnvironment(x,y)
    if color == White {
        adjSameColor, adjOtherColor = adjOtherColor, adjSameColor
    }
    enemiesInAtari, enemiesNotInAtari := b.determineGroupsAtariStatus(adjOtherColor)
    // func for capturing groups, if there are any.
    removeGroupsFunc := func() {}
    removeGroups := false // true if there are groups to remove
    if len(enemiesInAtari) > 0 {
        removeGroups = true
        removeGroupsFunc = func() {
            for _, grp := range enemiesInAtari {
                b.removeGroup(grp)
            }
        }
    }
    // func for updating liberties for each adjacent enemy group which is not captured.
    // This func is used if some stones are captured so that we have to walk over the entire groups.
    // This func is safe to call even if there are no groups whose liberties need to be
    // updated.
    updateLiberyFunc := func() {
        for _, grp := range enemiesNotInAtari {
            b.updateGroupLiberties(grp)
        }
    }

    //if adjSameColor.Length() == 0 {
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

                        b.clearKo()
                        // Remove the enemy stone
                        b.removeGroup(enemiesInAtari[0]) // TODO: use firstGroup instead!
                        // Update liberties and legality for the groups adjacent to the removed stone
                        var adjToKoSameColor GroupSlice
                        if color == Black {
                            _, adjToKoSameColor, _ = b.GetEnvironment(koX, koY)
                        } else {
                            _, _, adjToKoSameColor = b.GetEnvironment(koX, koY)
                        }
                        for _, grp := range adjToKoSameColor {
                            //grpX, grpY := b.posToXY(grp.Fields.First().Value()) // <DBG>
                            //printDbgMsgf("adjko around (%d,%d)\n", grpX, grpY) // </DBG>
                            grp.Liberties.AppendUnique(koPos)
                            b.updateLegalityForLibertiesOfExcept(grp, koX, koY)
                        }
                        // Create the new group
                        b.CreateGroup(x,y,color)
                        // Update liberties of the groups adjacent to the new stone at (x,y)
                        updateLiberyFunc()
                        // The player who took the ko may fill it, so put koPos into the appropriate legal moves array
                        if color == White {
                            b.legalWhiteMoves.AppendUnique(koPos)
                        } else {
                            b.legalBlackMoves.AppendUnique(koPos)
                        }
                        // Update the legality for the liberties of the groups adjacent to the new created stone
                        //printDbgMsg("Update the legality for the liberties of the groups adjacent to the new created stone\n") // </DBG>
                        for _, grp := range enemiesNotInAtari {
                            b.updateLegalityForLibertiesOf(grp)
                        }

                        b.ko = NewKoLock(koX, koY, !color)
                        return true, true
                    }
                } else {
                    // It's not a ko
                    action = func() (blackUpToDate, whiteUpToDate bool) {
                        //printDbgMsgf("Board.calculateIfLegal: sameColLen == nFree == 0, removeGroups = true, not ko case.\n") // <DBG>
                        //DbgHistogram.Score() // </DBG>
                        removeGroupsFunc()
                        b.CreateGroup(x,y,color)
                        updateLiberyFunc()
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
                    removeGroupsFunc()
                    b.CreateGroup(x,y,color)
                    updateLiberyFunc()
                    b.ko = nil
                    return false, false
                }
            } else {
                action = func() (blackUpToDate, whiteUpToDate bool) {
                    // Experiments show that this case is run the most often

                    //printDbgMsgf("Board.calculateIfLegal: sameColLen == 0, nFree > 0, removeGroups = false.\n") // <DBG>
                    //DbgHistogram.Score() // </DBG>
                    b.clearKo()
                    b.CreateGroup(x,y,color)
                    b.dropLibertyFromEach(pos, adjOtherColor)
                    b.updateLegalityForFreeNeighboursOf(x,y)
                    b.updateLegalityForAdjacentGroups(adjOtherColor)
                    //b.ko = nil
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
                    removeGroupsFunc()
                    b.joinGroupsByPlayAt(pos, adjSameColor)
                    updateLiberyFunc()
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
                        b.clearKo()
                        b.joinGroupsByPlayAt(pos, adjSameColor)
                        b.dropLibertyFromEach(pos, adjOtherColor)
                        b.updateLegalityForAdjacentGroups(adjOtherColor)
                        b.updateLegalityForLibertiesOf(b.fields[pos])

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
                    removeGroupsFunc()
                    b.joinGroupsByPlayAt(pos, adjSameColor)
                    updateLiberyFunc()
                    b.ko = nil
                    return false, false
                }
            } else {
                action = func() (blackUpToDate, whiteUpToDate bool) {
                    // Experiments show that this case is run the 2nd most often

                    //printDbgMsgf("Board.calculateIfLegal: sameColLen > 0, nFree > 0, removeGroups = false.\n") // <DBG>
                    //DbgHistogram.Score() // </DBG>
                    b.clearKo()
                    b.joinGroupsByPlayAt(pos, adjSameColor)
                    updateLiberyFunc()
                    b.updateLegalityForFreeNeighboursOf(x,y)
                    b.updateLegalityForAdjacentGroups(adjOtherColor)
                    // update legality for the newly joined group, which is at pos
                    b.updateLegalityForLibertiesOf(b.fields[pos])

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

// If there was a ko, make this point legal again and remove the ko lock
func (b *Board) clearKo() {
    if b.ko != nil {
        koPos := b.xyToPos(b.ko.Point.X, b.ko.Point.Y)
        b.ko = nil
        b.updateLegalityFor(koPos)
    }
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
    b.legalBlackMoves.Remove(pos)
    b.legalWhiteMoves.Remove(pos)
    b.emptyFields.Remove(pos)
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
// TODO: write tests for this...
func (b *Board) IsLegalMove(x, y int, color Color) bool {
    // TODO: write test for this!
    var indices *IntList = b.legalBlackMoves
    if color == White {
        indices = b.legalWhiteMoves
    }
    pos := b.xyToPos(x,y)
    last := indices.Last()
    for it := indices.First(); it != last; it = it.Next() {
        if it.Value() == pos {
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
// groups in 'adjSameColor' by playing a stone of the same color at posToXY(playPos)
// This func is _not_safe to call if there are no groups to join (i.e. adjSameColor empty)
func (b *Board) joinGroupsByPlayAt(playPos int, adjSameColor GroupSlice) {
    //printDbgMsg("in joinGroupsByPlayAt\n") // <DBG>
    //defer printDbgMsgf("returned from joinGroupsByPlayAt\n") // </DBG>

    firstGroup := adjSameColor[0]
    // Add the stone at posToXY(playPos) to the first group.
    firstGroup.Fields.Append(playPos)
    b.fields[playPos] = firstGroup
    b.emptyFields.Remove(playPos)
    b.legalBlackMoves.Remove(playPos)
    b.legalWhiteMoves.Remove(playPos)

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

// Returns a slice of legal moves of color 'color'
func (b *Board) ListLegalPoints(color Color) []Point {
    upToDate := b.acBlackMoveUpToDate
    legalMoves := b.legalBlackMoves
    if color == White {
        upToDate = b.acWhiteMoveUpToDate
        legalMoves = b.legalWhiteMoves
    }
    if !upToDate {
        b.updateLegalMoves(color)
    }
    last := legalMoves.Last()
    var ret []Point = make([]Point, legalMoves.Length())
    i := 0
    for it := legalMoves.First(); it != last; it = it.Next() {
        x, y := b.posToXY(it.Value())
        ret[i].X, ret[i].Y = x, y
        i++
    }
    return ret
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
    // It this is not the expected move sequence, we have to update the legal move arrays
    if color == Black && !b.acBlackMoveUpToDate {
        b.updateLegalMoves(Black)
    }
    if color == White && !b.acWhiteMoveUpToDate {
        b.updateLegalMoves(White)
    }

    pos := b.xyToPos(x,y)
    if !b.IsLegalMove(x,y, color) {
        return NewIllegalMoveError(x,y, color)
    }

    if color == White {
        b.acBlackMoveUpToDate, b.acWhiteMoveUpToDate = b.actionOnNextWhiteMove[pos]()
    } else {
        b.acBlackMoveUpToDate, b.acWhiteMoveUpToDate = b.actionOnNextBlackMove[pos]()
    }
    // Clear the appropriate actionOnNextMove array. 
    b.colorOfNextPlay = !color
    if b.colorOfNextPlay == White && !b.acWhiteMoveUpToDate {
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
    }

    return nil
}

// The player of color 'color' plays a pass.
func (b *Board) PlayPass(color Color) {
    nextActions := b.actionOnNextBlackMove
    if color == Black {
        nextActions = b.actionOnNextWhiteMove
    }
    for i := 0; i < b.BoardSize()*b.BoardSize(); i++ {
        nextActions[i] = nil
    }
    b.colorOfNextPlay = !color
    b.updateLegalMoves(b.colorOfNextPlay)
}

// Plays out the sequence 'seq' of moves
func (b *Board) playSequence(seq []Move) {
    for _, m := range seq {
        if m.Vertex.Pass {
            b.PlayPass(m.Color)
        } else {
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
// because of the removal of 'gid'.
// This method does not alter legalWhiteMoves or legalBlackMoves.
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
        b.emptyFields.Append(e.Value())
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
func (b *Board) updateLegalityFor(pos int) {
    pX, pY := b.posToXY(pos)
    var legalBlack, legalWhite bool
    legalBlack, b.actionOnNextBlackMove[pos] = b.calculateIfLegal(pX, pY, Black)
    legalWhite, b.actionOnNextWhiteMove[pos] = b.calculateIfLegal(pX, pY, White)
    //printDbgMsgf("at liberty (%d,%d); legalBlack: %v, legalWhite: %v\n", pX, pY, legalBlack, legalWhite) // <DBG/>
    if legalBlack {
        b.legalBlackMoves.AppendUnique(pos)
    } else {
        b.legalBlackMoves.Remove(pos)
    }
    if legalWhite {
        b.legalWhiteMoves.AppendUnique(pos)
    } else {
        b.legalWhiteMoves.Remove(pos)
    }
}

// Helper method for Board.calculateIfLegal. Does what it name indicates.
func (b *Board) updateLegalityForAdjacentGroups(adjGroups GroupSlice) {
    //printDbgMsg("in updateLegalityForAdjacentGroups\n") // <DBG>
    //defer printDbgMsgf("returned from updateLegalityForAdjacentGroups\n") // </DBG>

    for _, grp := range adjGroups {
        lastLib := grp.Liberties.Last()
        for itLib := grp.Liberties.First(); itLib != lastLib; itLib = itLib.Next() {
            lpos := itLib.Value()
            b.updateLegalityFor(lpos)
        }
    }
}

// Helper method for Board.calculateIfLegal. Does what it name indicates.
func (b *Board) updateLegalityForFreeNeighboursOf(x,y int) {
    //printDbgMsg("in updateLegalityForFreeNeighboursOf\n") // <DBG>
    //defer printDbgMsgf("returned from updateLegalityForFreeNeighboursOf\n") // </DBG>
    nbours := b.neighbours(x,y)
    for _, p := range nbours {
        npos := b.xyToPos(p.X, p.Y)
        if b.fields[npos] == nil {
            b.updateLegalityFor(npos)
        }
    }
}

// Helper method for Board.calculateIfLegal. Does what it name indicates.
func (b *Board) updateLegalityForLibertiesOf(group *Group) {
    //printDbgMsg("in Board.updateLegalityForLibertiesOf\n") // <DBG/>
    //defer printDbgMsg("returned from Board.updateLegalityForLibertiesOf\n") // <DBG/>
    lastLib := group.Liberties.Last()
    for itLib := group.Liberties.First(); itLib != lastLib; itLib = itLib.Next() {
        lpos := itLib.Value()
        b.updateLegalityFor(lpos)
    }
}

// Helper method for Board.calculateIfLegal, analogous to updateLegalityForLibertiesOf, except that (exceptX, exceptY) will be 
// left out
func (b *Board) updateLegalityForLibertiesOfExcept(group *Group, exceptX, exceptY int) {
    //printDbgMsg("in Board.updateLegalityForLibertiesOfExcept\n") // <DBG/>
    //defer printDbgMsg("returned from Board.updateLegalityForLibertiesOfExcept\n") // <DBG/>
    lastLib := group.Liberties.Last()
    for itLib := group.Liberties.First(); itLib != lastLib; itLib = itLib.Next() {
        exceptPos := b.xyToPos(exceptX, exceptY)
        lpos := itLib.Value()
        if exceptPos != lpos {
            b.updateLegalityFor(lpos)
        }
    }
}

// Updates the legal moves for the color 'color'
// TODO: write tests for this...
func (b *Board) updateLegalMoves(color Color) {
    // This method assumes that b.emptyFields is correctly set.
    legalMoves := b.legalBlackMoves
    actions := b.actionOnNextBlackMove
    flag := &b.acBlackMoveUpToDate
    if color == White {
        legalMoves = b.legalWhiteMoves
        actions = b.actionOnNextWhiteMove
        flag = &b.acWhiteMoveUpToDate
    }
    legalMoves.Clear()
    last := b.emptyFields.Last()
    for it := b.emptyFields.First(); it != last; it = it.Next() {
        pos := it.Value()
        x, y := b.posToXY(pos)
        isLegal, action := b.calculateIfLegal(x,y, color)
        if isLegal {
            legalMoves.Append(pos)
        }
        actions[pos] = action
    }
    *flag = true
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
                   legalWhiteMoves: NewIntList(),
                   legalBlackMoves: NewIntList(),
                   emptyFields: NewIntList(),
                   colorOfNextPlay: Black,
                   boardSize: boardsize,
                   acBlackMoveUpToDate: true,
                   acWhiteMoveUpToDate: true,
                 }
    for i := 0; i < boardsize*boardsize; i++ {
        ret.legalWhiteMoves.Append(i)
        ret.legalBlackMoves.Append(i)
        ret.emptyFields.Append(i)
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


