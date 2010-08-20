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
 * TODO:
 *      - let this talk the GTP-protocol
 *      - make this all more Go-ideomatic. Can the use of interfaces make IntList obsolete?
 *      - think about the cast-mess between int and GroupIndexType
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
}

// ##################### Board methods ##########################


// Calculates if a move of 'color' at (x,y) is legal. Does not use 
// b.legal{Black,White}Moves for this. Also, this stores the appropriate
// `actions` in actionOnNextMove. This `action` does not include any update 
// of legal{Black,White}Moves...

// TODO(david); this method might need some serious refactoring...
func (b *Board) calculateIfLegal(x,y int, color Color) bool {
    // Is this move prohibited because of a ko?
    if b.ko != nil {
        if b.ko.X == x && b.ko.Y == y {
            return false
        }
    }
    pos := xyToPos(x,y)
    nFree, adjSameColor, adjOtherColor := b.GetEnvironment(x,y)
    if color == White {
        adjSameColor, adjOtherColor = adjOtherColor, adjSameColor
    }

    groupsToRemove := b.filterGroupsInAtari(adjOtherColor, pos) // these groups would be captures by the move

    // func for capturing groups, if there are any.
    removeGroupsFunc := func() {}
    removeGroups := false // true if there are groups to remove
    if groupsToRemove.Length() > 0 {
        removeGroups = true
        removeGroupsFunc = func() {
            rmLast := groupsToRemove.Last()
            for rmIt := groupsToRemove.First(); rmIt != rmLast; rmIt = rmIt.Next() {
                b.removeGroupByGroupIndex( GroupIndexType(rmIt.Value()) )
            }
        }
    }

    // TODO: updateGroupLiberties #######################################

    if adjSameColor.Length() == 0 {
        if nFree == 0 {
            if !removeGroups {
                // There are no adjacent friendly groups, no neighbour is a free field and 
                // this move does not capture enemy stones. So it is illegal.
                return false
            } else {
                // There are no adjacent friendly groups and no free neighbours, but this 
                // move captures, so it is legal. Do the capture and create a new group afterwards.
                b.actionOnNextMove[pos] = func() {
                    removeGroupsFunc()
                    b.CreateGroup(x,y,color)
                }
                return true
            }
        } else {
            // There are no adjacent friendly groups, but free neighbour fields, so this move
            // is always legal. Remove adjacent enemy groups if necessary and create a new group.
            if removeGroups {
                b.actionOnNextMove[pos] = func() {
                    removeGroupsFunc()
                    b.CreateGroup(x,y,color)
                }
            } else {
                b.actionOnNextMove[pos] = func() {
                    b.CreateGroup(x,y,color)
                }
            }
            return true
        }
    } else {
        if nFree == 0 {
            if removeGroups {
                // This move captures stones and thus produces empty fields, so it is legal. Capture
                // the stones first and then join the adjacent groups of the same color.
                b.actionOnNextMove[pos] = func() {
                    removeGroupsFunc()
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
                }
                return true
            } else {
                // go on here
            }
        }
    }

    return false
}

// Create a new, one-stone-group of 'color' at (x,y). This method does not perform
// any legality checks.
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

// Returns the (possibly emyty) sublist of 'groups' which are in atari and whose only
// liberty is at posToXY(pos).
// TODO: write tests for this...
func (b *Board) filterGroupsInAtari(groups *IntList, pos int) (*IntList) {
    groupsToRemove := NewIntList()
    last := groups.Last()
    for it := groups.First(); it != last; it = it.Next() {
        groupIndex := GroupIndexType(it.Value())
        group := b.groupMap.Get(groupIndex)
        if group.Liberties.Length() == 1 && group.Liberties.First().Value() == pos {
            groupsToRemove.Append(int(groupIndex))
        }
    }
    return groupsToRemove
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
// as `` into += from '' ...
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

// Play a stone of color 'color' at (x,y). If an error occurs (such as that this 
// place is already occupied) this error is returned.
// TODO: write tests for this...
func (b *Board) Move(x, y int, color Color) (err Error) {
    index := xyToPos(x,y)
    if !b.fields[index].Empty() {
        return NewFieldOccupiedError(x,y)
    }
    if !b.IsLegalMove(x,y, color) {
        return NewIllegalMoveError(x,y, color)
    }
    // TODO: not yet done at all!
    // Go on here!
    return
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
// This method does not alter legalWhiteMoves or legalBlackMoves.
func (b *Board) removeGroupByGroupIndex(gid GroupIndexType) {
    g := b.groupMap.Get(gid)
    last := g.Fields.Last()
    for e := g.Fields.First(); e != last; e = e.Next() {
        b.fields[e.Value()].Clear()
        b.emptyFields.Append(e.Value())
    }
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
    // FIXME: use calculateIfLegal... this is now obsolete!
    // This method assumes that b.emptyFields is correctly set.
    b.legalWhiteMoves.Clear()
    b.legalBlackMoves.Clear()
    last := b.emptyFields.Last()
    for it := b.emptyFields.First(); it != last; it = it.Next() {
        //pos := b.emptyFields.Get(i)
        pos := it.Value()
        x, y := posToXY(pos)
        nbours := neighbours(x,y)
        freeNBours := 0
        // how many are free?
        for _, p := range nbours {
            if b.fields[xyToPos(p.X, p.Y)].Empty() {
                freeNBours++
            }
        }
        switch {
            case freeNBours > 0: // this is a legal move for both colors
                b.legalWhiteMoves.Append(pos)
                b.legalBlackMoves.Append(pos)
            case freeNBours == 0: // this case is more difficult
                // TODO: implement this!
                return
        }
    }
}

// ##################### Board helper functions ##########################
// Creates a new, initial board
func NewBoard() *Board {
    ret := &Board{ legalWhiteMoves: NewIntList(),
                   legalBlackMoves: NewIntList(),
                   emptyFields: NewIntList(),
                   groupMap: NewGroupMap(),
                 }
    for i := 0; i < BoardSize*BoardSize; i++ {
        ret.legalWhiteMoves.Append(i)
        ret.legalBlackMoves.Append(i)
        ret.emptyFields.Append(i)
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


