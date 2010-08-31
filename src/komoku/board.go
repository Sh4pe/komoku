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
 *      - let this talk the GTP-protocol
 *      - make this all more Go-ideomatic. Can the use of interfaces make IntList obsolete?
 *      - think about the cast-mess between int and GroupIndexType. Use interfaces for this!
 *      - implement IntList.Contains(int) 
 *      - replace legal{Black,White}Moves by legalNextMove?
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
    fields []GroupIndexType // Stores the indices for groupMap
    groupMap *GroupMap // Maps fields to groups
    legalBlackMoves *IntList // Indices of fields at which it is legal to play a black stone.
    legalWhiteMoves *IntList // Indices of fields at which it is legal to play a white stone.
    emptyFields *IntList // indices of empty fields
    //ko *Point // if not nil, this points to where you can't play because of the ko rule
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
    // use this for changing IntList - maybe.
    /*if adjSameColor.Length() == 0 {
        DbgHistogram.ScoreTagged("sameColorLen == 0")
    } else {
        DbgHistogram.ScoreTagged("sameColorLen != 0")
    }
    if adjOtherColor.Length() == 0 {
        DbgHistogram.ScoreTagged("otherColorLen == 0")
    } else {
        DbgHistogram.ScoreTagged("otherColorLen != 0")
    }*/
    if color == White {
        adjSameColor, adjOtherColor = adjOtherColor, adjSameColor
    }
    enemiesInAtari, enemiesNotInAtari := b.determineGroupsAtariStatus(adjOtherColor)
    // func for capturing groups, if there are any.
    removeGroupsFunc := func() {}
    removeGroups := false // true if there are groups to remove
    if enemiesInAtari.Length() > 0 {
        removeGroups = true
        removeGroupsFunc = func() {
            rmLast := enemiesInAtari.Last()
            for rmIt := enemiesInAtari.First(); rmIt != rmLast; rmIt = rmIt.Next() {
                b.removeGroupByGroupIndex( GroupIndexType(rmIt.Value()) )
            }
        }
    }
    // func for updating liberties for each adjacent enemy group which is not captured.
    // This func is used if some stones are captured so that we have to walk over the entire groups.
    // This func is safe to call even if there are no groups whose liberties need to be
    // updated.
    updateLiberyFunc := func() {
        last := enemiesNotInAtari.Last()
        for it := enemiesNotInAtari.First(); it != last; it = it.Next() {
            b.updateGroupLiberties(b.groupMap.Get( GroupIndexType(it.Value()) ))
        }
    }

    if adjSameColor.Length() == 0 {
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
                firstGroup := b.groupMap.Get( GroupIndexType(enemiesInAtari.First().Value()) )
                if enemiesInAtari.Length() == 1 && firstGroup.Fields.Length() == 1 {
                    // It's a ko, so remove the group, play the stone, update the liberties
                    // and set b.ko to the right point.
                    koPos := firstGroup.Fields.First().Value()
                    koX, koY := b.posToXY(koPos)
                    action = func() (blackUpToDate, whiteUpToDate bool) {
                        //printDbgMsgf("Board.calculateIfLegal: sameColLen == nFree == 0, removeGroups = true, ko case.\n") // <DBG>
                        //DbgHistogram.Score() // </DBG>

                        b.clearKo()
                        // Remove the enemy stone
                        b.removeGroupByGroupIndex( GroupIndexType(enemiesInAtari.First().Value()) )
                        // Update liberties and legality for the groups adjacent to the removed stone
                        var adjToKoSameColor *IntList
                        if color == Black {
                            _, adjToKoSameColor, _ = b.GetEnvironment(koX, koY)
                        } else {
                            _, _, adjToKoSameColor = b.GetEnvironment(koX, koY)
                        }
                        adjkoLast := adjToKoSameColor.Last()
                        for adjkoIt := adjToKoSameColor.First(); adjkoIt != adjkoLast; adjkoIt = adjkoIt.Next() {
                            grp := b.groupMap.Get( GroupIndexType(adjkoIt.Value()) )
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
                        lastEnemy := enemiesNotInAtari.Last()
                        for enemyIt := enemiesNotInAtari.First(); enemyIt != lastEnemy; enemyIt = enemyIt.Next() {
                            b.updateLegalityForLibertiesOf(b.groupMap.Get( GroupIndexType(enemyIt.Value()) ))
                        }

                        // TODO: b.ko always nil after clearKo
                        if b.ko == nil {
                            //b.ko = NewPoint(koX, koY)
                            b.ko = NewKoLock(koX, koY, !color)
                        } else {
                            b.ko.Point.X, b.ko.Point.Y = koX, koY
                            b.ko.Color = !color
                        }
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
                last := adjSameColor.Last()
                oneHasTwo := false
                for it := adjSameColor.First(); it != last; it = it.Next() {
                    g := b.groupMap.Get( GroupIndexType(it.Value()) )
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
                        //b.ko = nil
                        b.clearKo()
                        b.joinGroupsByPlayAt(pos, adjSameColor)
                        b.dropLibertyFromEach(pos, adjOtherColor)
                        b.updateLegalityForAdjacentGroups(adjOtherColor)
                        // update legality for the newly joined group, which is at pos
                        b.updateLegalityForLibertiesOf(b.groupMap.Get(b.fields[pos]))

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
                    //joinGroupsFunc()
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
                    b.updateLegalityForLibertiesOf(b.groupMap.Get(b.fields[pos]))

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
        if b.fields[npos].Empty() {
            newGroup.Liberties.Append(npos)
        }
    }
    newIndex := b.groupMap.Append(newGroup)
    b.fields[pos] = newIndex
    b.legalBlackMoves.Remove(pos)
    b.legalWhiteMoves.Remove(pos)
    b.emptyFields.Remove(pos)
}

// Returns the sublists '{not,}inAtari' of 'groups' which are {not,} in atari.
// TODO: write tests for this...
func (b *Board) determineGroupsAtariStatus(groups *IntList) (inAtari, notinAtari *IntList) {
    last := groups.Last()
    inAtari = NewIntList()
    notinAtari = NewIntList()
    for it := groups.First(); it != last; it = it.Next() {
        groupIndex := GroupIndexType(it.Value())
        group := b.groupMap.Get(groupIndex)
        if group.Liberties.Length() == 1 {
            inAtari.Append(int(groupIndex))
        } else {
            notinAtari.Append(int(groupIndex))
        }
    }
    return
}

// Helper for Board.calculateIfLegal. Takes each group (identified by one index) in adjGroups and drops
// the liberty posToXY(libertyPos) from it
func (b *Board) dropLibertyFromEach(libertyPos int, adjGroups *IntList) {
    last := adjGroups.Last()
    for it := adjGroups.First(); it != last; it = it.Next() {
        grp := b.groupMap.Get( GroupIndexType(it.Value()) )
        grp.Liberties.Remove(libertyPos)
    }
}

// Returns the `environment` of (x,y), i.e. the number 'nFree' of free neighbours
// and IntLists 'adj{Black,White}' containing instances of GroupIndexTypes of adjacent {black,white} groups.
// TODO: unexport this?
func (b *Board) GetEnvironment(x,y int) (nFree int, adjBlack, adjWhite *IntList) {
    nbours := b.neighbours(x,y)
    adjBlack = NewIntList()
    adjWhite = NewIntList()
    for _, p := range nbours {
        npos := b.xyToPos(p.X, p.Y)
        if b.fields[npos].Empty() {
            nFree++
        } else {
            gindex := b.fields[npos]
            group := b.groupMap.Get(gindex)
            // Don't append the same group more than once!
            if group.Color == Black {
                adjBlack.AppendUnique(int(gindex))
            } else {
                adjWhite.AppendUnique(int(gindex))
            }
        }
    }
    return
}

// Analogous to GetEnvironment.
// TODO: write tests for this...
// TODO: unexport this?
func (b *Board) GetEnvironmentByPos(pos int) (nFree int, adjBlack, adjWhite *IntList) {
    x, y := b.posToXY(pos)
    return b.GetEnvironment(x,y)
}

// If the field (x,y) is empty, this method returns (false, nil).
// If the field is not empty, it returns (true, 'pointer to group')...
// TODO: change the return values to (group,empty)?
func (b *Board) GetGroup(x,y int) (empty bool, group *Group) {
    index := b.xyToPos(x,y)
    gindex := b.fields[index]
    if gindex.Empty() {
        return true, nil
    }
    return false, b.groupMap.Get(gindex)
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
func (b *Board) joinGroups(into, from GroupIndexType) {
    ginto := b.groupMap.Get(into)
    gfrom := b.groupMap.Get(from)
    lastFrom := gfrom.Fields.Last()
    for fromIt := gfrom.Fields.First(); fromIt != lastFrom; fromIt = fromIt.Next() {
        fpos := fromIt.Value()
        b.fields[fpos] = into
    }
    ginto.Fields.JoinUnique(gfrom.Fields)
    b.groupMap.Remove(from)
}

// This is a helper function for Board.calculateIfLegal. It joint the adjacent
// groups in 'adjSameColor' by playing a stone of the same color at posToXY(playPos)
// This func is _not_safe to call if there are no groups to join (i.e. adjSameColor empty)
func (b *Board) joinGroupsByPlayAt(playPos int, adjSameColor *IntList) {
    //printDbgMsg("in joinGroupsByPlayAt\n") // <DBG>
    //defer printDbgMsgf("returned from joinGroupsByPlayAt\n") // </DBG>
    // Add the stone at posToXY(playPos) to the first group.
    first := adjSameColor.First()
    firstIndex := GroupIndexType(first.Value())
    firstGroup := b.groupMap.Get(firstIndex)
    /*func() { // <DBG>
        dbgX, dbgY := b.posToXY(playPos)
        printDbgMsgf("firstGroup.Fields.Append at (%d,%d)\n", dbgX, dbgY)
        dbgMsg := "Fields before: "
        dbgLast := firstGroup.Fields.Last()
        for di := firstGroup.Fields.First(); di != dbgLast; di=di.Next() {
            dbgXX, dbgYY := b.posToXY(di.Value())
            dbgMsg += fmt.Sprintf("(%d,%d), ", dbgXX, dbgYY)
        }
        printDbgMsg(dbgMsg)
    }() // <DBG>*/
    firstGroup.Fields.Append(playPos)
    b.fields[playPos] = firstIndex
    b.emptyFields.Remove(playPos)
    b.legalBlackMoves.Remove(playPos)
    b.legalWhiteMoves.Remove(playPos)
    /*func() { // <DBG>
        dbgMsg := "Fields after: "
        dbgLast := firstGroup.Fields.Last()
        for di := firstGroup.Fields.First(); di != dbgLast; di=di.Next() {
            dbgXX, dbgYY := b.posToXY(di.Value())
            dbgMsg += fmt.Sprintf("(%d,%d), ", dbgXX, dbgYY)
        }
        printDbgMsg(dbgMsg)
    }() // <DBG>*/

    // Then join the other groups into the first group
    last := adjSameColor.Last()
    for it := first.Next(); it != last; it = it.Next() {
        b.joinGroups(firstIndex, GroupIndexType(it.Value()))
    }
    /*func() { // <DBG>
        dbgMsg := "Fields after every other friendly group is joined in: "
        dbgLast := firstGroup.Fields.Last()
        for di := firstGroup.Fields.First(); di != dbgLast; di=di.Next() {
            dbgXX, dbgYY := b.posToXY(di.Value())
            dbgMsg += fmt.Sprintf("(%d,%d), ", dbgXX, dbgYY)
        }
        printDbgMsg(dbgMsg)
    }() // <DBG>*/
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
    nblack, nwhite = 0,0
    // Note to self: I love closures!
    f := func(k GroupIndexType, grp *Group) {
        if grp.Color == Black {
            nblack++
        } else {
            nwhite++
        }
    }
    b.groupMap.Do(f)
    return
}

// Returns the number of {black,white} stones on the board
func (b *Board) numberOfStones() (nblack, nwhite int) {
    nblack, nwhite = 0,0
    f := func(k GroupIndexType, grp *Group) {
        if grp.Color == Black {
            nblack += grp.Fields.Length()
        } else {
            nwhite += grp.Fields.Length()
        }
    }
    b.groupMap.Do(f)
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

// Helper method for Board.calculateIfLegal. Removes all groups in 'enemiesInAtari', if
// there are any
/*func (b *Board) removeInAtari(enemiesInAtari *IntList) {
    if enemiesInAtari.Length() == 0 {
        return
    }
    rmLast := enemiesInAtari.Last()
    for rmIt := enemiesInAtari.First(); rmIt != rmLast; rmIt = rmIt.Next() {
        b.removeGroupByGroupIndex( GroupIndexType(rmIt.Value()) )
    }
}*/

// Removes the group which occupies (x,y), if there is any, and updates b.emptyFields.
// This method does not alter legalWhiteMoves or legalBlackMoves.
// TODO(david): Do I really want to export this?
func (b *Board) RemoveGroup(x,y int) {
    pos := b.xyToPos(x,y)
    if !b.fields[pos].Empty() {
        b.removeGroupByGroupIndex(b.fields[pos])
    }
}

// Removes the group which has the index 'gid' in .groupMap and updates b.emptyFields.
// Also, this method calles updateGroupLiberties for each group whose liberties might change
// because of the removal of 'gid'.
// This method does not alter legalWhiteMoves or legalBlackMoves.
func (b *Board) removeGroupByGroupIndex(gid GroupIndexType) {
    g := b.groupMap.Get(gid)
    last := g.Fields.Last()
    adjGroups := NewIntList()
    for e := g.Fields.First(); e != last; e = e.Next() {
        // Collect adjacend groups so that we can update their liberties later
        pos := e.Value()
        x, y := b.posToXY(pos)
        nbours := b.neighbours(x,y)
        for _, p := range nbours {
            npos := b.xyToPos(p.X, p.Y)
            if !b.fields[npos].Empty() && b.fields[npos] != gid {
                adjGroups.AppendUnique( int(b.fields[npos]) )
            }
        }
        b.fields[pos].Clear()
        b.emptyFields.Append(e.Value())
    }
    // Update the liberties of the adjacent groups.
    adjLast := adjGroups.Last()
    for adjIt := adjGroups.First(); adjIt != adjLast; adjIt = adjIt.Next() {
        b.updateGroupLiberties(b.groupMap.Get( GroupIndexType(adjIt.Value()) ))
    }
    b.groupMap.Remove(gid)
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
    // this is expensive...
    group.Liberties.Clear()
    last := group.Fields.Last()
    for it := group.Fields.First(); it != last; it = it.Next() {
        x, y := b.posToXY(it.Value())
        nbours := b.neighbours(x,y)
        for _, p := range nbours {
            npos := b.xyToPos(p.X, p.Y)
            if b.fields[npos].Empty() {
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
func (b *Board) updateLegalityForAdjacentGroups(adjGroups *IntList) {
    //printDbgMsg("in updateLegalityForAdjacentGroups\n") // <DBG>
    //defer printDbgMsgf("returned from updateLegalityForAdjacentGroups\n") // </DBG>
    adjLast := adjGroups.Last()
    for adjIt := adjGroups.First(); adjIt != adjLast; adjIt = adjIt.Next() {
        grp := b.groupMap.Get( GroupIndexType(adjIt.Value()) )
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
        if b.fields[npos].Empty() {
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

// TODO: use point!
func (b *Board) xyToPos(x, y int) int {
    return b.BoardSize()*y + x
}

// ##################### Board helper functions ##########################
// Creates a new, initial board of size 'boardsize'.
func NewBoard(boardsize int) *Board {
    ret := &Board{ fields: make([]GroupIndexType, boardsize*boardsize),
                   actionOnNextBlackMove: make([]actionFunc, boardsize*boardsize),
                   actionOnNextWhiteMove: make([]actionFunc, boardsize*boardsize),
                   legalWhiteMoves: NewIntList(),
                   legalBlackMoves: NewIntList(),
                   emptyFields: NewIntList(),
                   groupMap: NewGroupMap(),
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
            //fmt.Printf("Initial black move at (%d,%d)\n", x,y)
            ret.CreateGroup(x, y, Black)
            ret.colorOfNextPlay = !ret.colorOfNextPlay
            ret.legalMovesNeedUpdate()
            return false, false
        }
        ret.actionOnNextWhiteMove[i] = func() (updateBlack, updateWhite bool) {
            //fmt.Printf("Initial white move at (%d,%d)\n", x,y)
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


