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
    //"container/list"
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


// ################################################################################
// ########################### Board struct #######################################
// ################################################################################

// This object is responsible for recording a current state of a game.
type Board struct {
    fields [BoardSize*BoardSize]GroupIndexType // Stores the indices for groupMap
    groupMap *GroupMap // Maps fields to groups
    legalBlackMoves *IntList // Indices of fields at which it is legal to play a black stone.
    legalWhiteMoves *IntList // Indices of fields at which it is legal to play a white stone.
    emptyFields *IntList // indices of empty fields
    ko *Point // if not nil, this points to where you can't play because of the ko rule
    actionOnNextMove [BoardSize*BoardSize]func() // This stores the appropriate code which has to be run if a move is played on a field
    colorOfNextPlay Color
}

// ##################### Board methods ##########################


// Calculates if a move of 'color' at (x,y) is legal. Does not use 
// b.legal{Black,White}Moves for this. Also, this method returnsthe appropriate
// `actions` to be performed if this move is played using 'action'. This 'action' does not include any updates
// of legal{Black,White}Moves... 
// Note that this method assumes that (x,y) is empty.
func (b *Board) calculateIfLegal(x,y int, color Color) (isLegal bool, action func()) {
    //fmt.Printf("Board.calculateIfLegal(%d, %d, %v)\n", x,y,color)
    // Is this move prohibited because of a ko?
    if b.ko != nil {
        if b.ko.X == x && b.ko.Y == y {
            return false, nil
        }
    }
    pos := xyToPos(x,y)
    nFree, adjSameColor, adjOtherColor := b.GetEnvironment(x,y)
    //fmt.Printf("Board.calculateIfLegal: environment: nFree: %d, adjSameColor: %v, adjOtherColor: %v\n", nFree, adjSameColor, adjOtherColor)
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
    // func for updating liberties for each adjacent enemy group which is not catpured.
    // This func is safe to call even if there are no groups whose liberties need to be
    // updated.
    updateLiberyFunc := func() {
        last := enemiesNotInAtari.Last()
        for it := enemiesNotInAtari.First(); it != last; it = it.Next() {
            b.updateGroupLiberties(b.groupMap.Get( GroupIndexType(it.Value()) ))
        }
    }
    // func to join all the adjacent groups of the same color by playing a stone of color 'color'
    // at (x,y). This func is _not_safe to call if there are no groups to join!
    joinGroupsFunc := func() {
        // Add the stone at (x,y) to the first group.
        first := adjSameColor.First()
        firstIndex := GroupIndexType(first.Value())
        firstGroup := b.groupMap.Get(firstIndex)
        firstGroup.Fields.Append(pos)
        // Then join the other groups into the first group
        last := adjSameColor.Last()
        for it := first.Next(); it != last; it = it.Next() {
            b.joinGroups(firstIndex, GroupIndexType(it.Value()))
        }
        b.updateGroupLiberties(firstGroup)
    }

    if adjSameColor.Length() == 0 {
        if nFree == 0 {
            if !removeGroups {
                // There are no adjacent friendly groups, no neighbour is a free field and 
                // this move does not capture enemy stones. So it is illegal.
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
                    koX, koY := posToXY(firstGroup.Fields.First().Value())
                    action = func() {
                        //fmt.Printf("Board.calculateIfLegal: sameColLen == nFree == 0, removeGroups = true, ko case\n")
                        removeGroupsFunc()
                        b.CreateGroup(x,y,color)
                        updateLiberyFunc()
                        b.ko.X, b.ko.Y = koX, koY
                    }
                } else {
                    // It's not a ko
                    action = func() {
                        //fmt.Printf("Board.calculateIfLegal: sameColLen == nFree == 0, removeGroups = true, not ko case\n")
                        removeGroupsFunc()
                        b.CreateGroup(x,y,color)
                        updateLiberyFunc()
                        b.ko = nil
                    }
                }
                return true, action
            }
        } else {
            // There are no adjacent friendly groups, but free neighbour fields, so this move
            // is always legal. Remove adjacent enemy groups if necessary and create a new group.
            if removeGroups {
               action = func() {
                    //fmt.Printf("Board.calculateIfLegal: sameColLen == 0, nFree > 0, removeGroups = true\n")
                    removeGroupsFunc()
                    b.CreateGroup(x,y,color)
                    updateLiberyFunc()
                    b.ko = nil
                }
            } else {
                action = func() {
                    //fmt.Printf("Board.calculateIfLegal: sameColLen == 0, nFree > 0, removeGroups = false\n")
                    b.CreateGroup(x,y,color)
                    updateLiberyFunc()
                    b.ko = nil
                }
            }
            return true, action
        }
    } else {
        if nFree == 0 {
            if removeGroups {
                // This move captures stones and thus produces empty fields, so it is legal. Capture
                // the stones first and then join the adjacent groups of the same color.
                action = func() {
                    //fmt.Printf("Board.calculateIfLegal: sameColLen > 0, nFree == 0, removeGroups = true\n")
                    removeGroupsFunc()
                    joinGroupsFunc()
                    updateLiberyFunc()
                    b.ko = nil
                }
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
                    action = func() {
                        //fmt.Printf("Board.calculateIfLegal: sameColLen > 0, nFree == 0, removeGroups = false, oneHasTwo = true\n")
                        joinGroupsFunc()
                        updateLiberyFunc()
                        b.ko = nil
                    }
                    return true, action
                } else {
                    // There is no groups to remove and every adjacient group of the same color has only one liberty,
                    // which must be the field (x,y) we want to play at, so this move is illegal.
                    return false, nil
                }
            }
        } else {
            // There are free neighbour fields, so this move is always legal. Capture adjacent enemy groups if necessary, 
            // then join groups and update liberties
            if removeGroups {
                action = func() {
                    //fmt.Printf("Board.calculateIfLegal: sameColLen > 0, nFree > 0, removeGroups = true\n")
                    removeGroupsFunc()
                    joinGroupsFunc()
                    updateLiberyFunc()
                    b.ko = nil
                }
            } else {
                action = func() {
                    //fmt.Printf("Board.calculateIfLegal: sameColLen > 0, nFree > 0, removeGroups = false\n")
                    joinGroupsFunc()
                    updateLiberyFunc()
                    b.ko = nil
                }
            }
            return true, action
        }
    }

    // Control flow should never reach here; if it does however, this is an error
    panic("Control reaced the end of Board.calculateIfLegal")
    return false, nil
}

// Create a new, one-stone-group of 'color' at (x,y) and sets its liberties appropriately. 
// This method does not perform any legality checks or liberty updates for other groups.
// TODO: unexport this?
func (b *Board) CreateGroup(x, y int, color Color) {
    pos := xyToPos(x,y)
    newGroup := NewGroup(color)
    newGroup.Fields.Append(pos)
    nbours := neighbours(x,y)
    for _, p := range nbours {
        npos := xyToPos(p.X, p.Y)
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

// Returns the `environment` of (x,y), i.e. the number 'nFree' of free neighbours
// and IntLists 'adj{Black,White}' containing instances of GroupIndexTypes of adjacent {black,white} groups.
// TODO: unexport this?
func (b *Board) GetEnvironment(x,y int) (nFree int, adjBlack, adjWhite *IntList) {
    nbours := neighbours(x,y)
    adjBlack = NewIntList()
    adjWhite = NewIntList()
    for _, p := range nbours {
        npos := xyToPos(p.X, p.Y)
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
    x, y := posToXY(pos)
    return b.GetEnvironment(x,y)
}

// If the field (x,y) is empty, this method returns (false, nil).
// If the field is not empty, it returns (true, 'pointer to group')...
func (b *Board) GetGroup(x,y int) (empty bool, group *Group) {
    index := xyToPos(x,y)
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
    pos := xyToPos(x,y)
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


// Removes the group which occupies (x,y), if there is any, and updates b.emptyFields.
// This method does not alter legalWhiteMoves or legalBlackMoves.
// TODO(david): Do I really want to export this?
func (b *Board) RemoveGroup(x,y int) {
    pos := xyToPos(x,y)
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
        x, y := posToXY(pos)
        nbours := neighbours(x,y)
        for _, p := range nbours {
            npos := xyToPos(p.X, p.Y)
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
    //fmt.Printf("Board.Play(%d, %d)\n", x, y)
    pos := xyToPos(x,y)
    if !b.fields[pos].Empty() {
        return NewFieldOccupiedError(x,y)
    }
    if !b.IsLegalMove(x,y, b.colorOfNextPlay) {
        return NewIllegalMoveError(x,y, b.colorOfNextPlay)
    }

    b.actionOnNextMove[pos]()
    // Clear the actionOnNextMove array. 
    // David: This is really important so that the GC can free the closures with all the
    // associated contexts. Am I right?
    for i := 0; i < BoardSize*BoardSize; i++ {
        b.actionOnNextMove[i] = nil
    }
    b.colorOfNextPlay = !b.colorOfNextPlay
    b.updateLegalMoves()

    return
}

// The player, whose turn it is, plays a pass.
func (b *Board) TurnPlayPass() {
    for i := 0; i < BoardSize*BoardSize; i++ {
        b.actionOnNextMove[i] = nil
    }
    b.colorOfNextPlay = !b.colorOfNextPlay
    b.updateLegalMoves()
}

// Recomputes the liberties of 'group' based only on the currently occupied and empty
// fields.
func (b *Board) updateGroupLiberties(group *Group) {
    // this is expensive...
    group.Liberties.Clear()
    last := group.Fields.Last()
    for it := group.Fields.First(); it != last; it = it.Next() {
        x, y := posToXY(it.Value())
        nbours := neighbours(x,y)
        for _, p := range nbours {
            npos := xyToPos(p.X, p.Y)
            if b.fields[npos].Empty() {
                group.Liberties.AppendUnique(npos)
            }
        }
    }
}

// TODO: write tests for this...
func (b *Board) updateLegalMoves() {
    // This method assumes that b.emptyFields is correctly set.
    b.legalWhiteMoves.Clear()
    b.legalBlackMoves.Clear()
    last := b.emptyFields.Last()
    for it := b.emptyFields.First(); it != last; it = it.Next() {
        pos := it.Value()
        x, y := posToXY(pos)
        isLegal, action := b.calculateIfLegal(x,y, b.colorOfNextPlay)
        if isLegal {
            if b.colorOfNextPlay == Black {
                b.legalBlackMoves.Append(pos)
            } else {
                b.legalWhiteMoves.Append(pos)
            }
        }
        b.actionOnNextMove[pos] = action
    }
}

// ##################### Board helper functions ##########################
// Creates a new, initial board
func NewBoard() *Board {
    ret := &Board{ legalWhiteMoves: NewIntList(),
                   legalBlackMoves: NewIntList(),
                   emptyFields: NewIntList(),
                   groupMap: NewGroupMap(),
                   colorOfNextPlay: Black,
                 }
    for i := 0; i < BoardSize*BoardSize; i++ {
        ret.legalWhiteMoves.Append(i)
        ret.legalBlackMoves.Append(i)
        ret.emptyFields.Append(i)
        // set the initial actions for playing a field
        x, y := posToXY(i)
        ret.actionOnNextMove[i] = func() {
            ret.CreateGroup(x, y, ret.colorOfNextPlay)
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


