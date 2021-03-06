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
    "rand"
    "os"
)

/*
 * This file defines the Board object. This object handles the representation of one current
 * state in a game. It provides the methods for creating legal games.
 */

// ################################################################################
// ########################### global variables ###################################
// ################################################################################
var neighbourCache []([]([]int))

// ################################################################################
// ########################### init func ##########################################
// ################################################################################
func init() {
    // Initialize the neighbourCache. We only support boards of sizes between 3 and 25
    neighbourCache = make([]([]([]int)), 30)
    for boardSize := 3; boardSize < 26; boardSize++ {
        neighbourCache[boardSize] = make([]([]int), boardSize*boardSize)
        for pos := 0; pos < boardSize*boardSize; pos++ {
            x, y := posToXY(pos, boardSize)
            nbourPoints := calculateNeighbours(x,y,boardSize)
            l := len(nbourPoints)
            neighbourCache[boardSize][pos] = make([]int, l)
            for i, p := range nbourPoints {
                nbourPos := xyToPos(p.X, p.Y, boardSize)
                neighbourCache[boardSize][pos][i] = nbourPos
            }
        }
    }
}

// ################################################################################
// ########################### actionFunc #########################################
// ################################################################################

// Context is implicit in closures. We use this to make a part of the context needed by an 
// actionFunc explicit, because we need to change it later, when we copy a board.
type boardBoundFuncContext struct {
    enemiesInAtari, enemiesNotInAtari, adjSameColor, adjOtherColor GroupSlice
}

// This is the actual closure
type actionFuncClosure func(boardPtr *Board, c *boardBoundFuncContext) (blackUpdated, whiteUpdated bool)

// An actionFunc represents the action performed by playing at one field. Use .Call() to
// execute the action.
type actionFunc struct {
    context *boardBoundFuncContext
    boardPtr *Board
    f actionFuncClosure
}

// The return values are only relevant if we assume that the legalities for all fields have been 
// calculated before. In this case, {black,white}Updated indicates that all legalities are valid 
// afterwards.
func (a *actionFunc) Call() (blackUpdated, whiteUpdated bool) {
    return a.f(a.boardPtr, a.context)
}

func NewActionFunc(board *Board, context *boardBoundFuncContext, f actionFuncClosure) *actionFunc {
    return &actionFunc{
        context: context,
        boardPtr: board,
        f: f,
    }
}

// ################################################################################
// ########################### Auxiliary types for Board ##########################
// ################################################################################

// Used for ko. Says that a play of color 'Color' at 'Point' is forbidden by the ko rule
type koLock struct {
    Pos int
    Color
}

func NewKoLock(pos int, color Color) *koLock {
    return &koLock{ Pos: pos,
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
    actionOnNextBlackMove []*actionFunc // This stores the appropriate code which has to be run if a black move is played on a field
    actionOnNextWhiteMove []*actionFunc // see the obvious analogue
    colorOfNextPlay Color
    boardSize int
    currentSequence uint32 // Represents the state of the board. Everything flagged with a different sequence has to be updated.
    fieldSequencesBlack []uint32
    fieldSequencesWhite []uint32
    rand *rand.Rand
    prisonersBlack int // number of black prisoners
    prisonersWhite int // number of white prisoners
}

// ##################### Board methods ##########################

// Returns the board size.
func (b *Board) BoardSize() int {
    return b.boardSize
}

// Calculates if a move of 'color' at 'pos' is legal. Does not use 
// b.legal{Black,White}Moves for this. Also, this method returns the appropriate
// `actions` to be performed if this move is played using 'action'. The return value
// of these actions indicate if any legality update has been performed.
// Note that this method assumes that 'pos' is empty.
func (b *Board) calculateIfLegal(pos int, color Color) (isLegal bool, action *actionFunc) {
    /*x, y := b.posToXY(pos)
    v, _ := pointToGTPVertex(*NewPoint(x, y))
    printDbgMsgBTf(6,"entering Board.calculateIfLegal(%s, %v)\n", v, color) // <DBG>*/
    //printDbgMsgf("ko status: %v\n", b.ko)

    // Is this move prohibited because of a ko? It is not prohibited for the player who 
    // took the ko to fill it in the next move
    if b.ko != nil && b.ko.Color == color && b.ko.Pos == pos {
        //printDbgMsgf("returned from Board.calculateIfLegal(%d, %d, %v)\n", x, y, color) // <DBG/>
        return false, nil
    }
    // assemble context which is bound to the board
    nFree, context := b.getEnvironmentAndContext(pos, color)
    /*printDbgMsgf("sameColLen: %d, othColLen: %d, enemInAtariLen: %d, enemNotInAtariLen: %d\n", len(context.adjSameColor), len(context.adjOtherColor), 
                len(context.enemiesInAtari), len(context.enemiesNotInAtari)) // <DBG>*/
    removeGroups := len(context.enemiesInAtari) > 0

    if len(context.adjSameColor) == 0 {
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
                firstGroup := context.enemiesInAtari[0]
                if len(context.enemiesInAtari) == 1 && firstGroup.Fields.Length() == 1 {
                    // It's a ko, so remove the group, play the stone, update the liberties
                    // and set b.ko to the right point.
                    koPos := firstGroup.Fields.First().Value()
                    action = NewActionFunc(b, context, func(boardPtr *Board, c *boardBoundFuncContext) (blackUpToDate, whiteUpToDate bool) {
                        //printDbgMsgf("Board.calculateIfLegal: sameColLen == nFree == 0, removeGroups = true, ko case.\n") // <DBG>
                        //DbgHistogram.Score() // </DBG>

                        alreadyUpdated := make([]bool, boardPtr.boardSize*boardPtr.boardSize)
                        boardPtr.ko = nil
                        boardPtr.removeGroup(c.enemiesInAtari[0]) // Remove the enemy stone
                        boardPtr.CreateGroup(pos,color) // Create the new group
                        for _, grp := range c.enemiesNotInAtari { // Update liberties of the groups adjacent to the new stone at (x,y)
                            boardPtr.updateGroupLiberties(grp)
                        }
                        // Update liberties and legality for the groups adjacent to the removed stone
                        var adjToKoSameColor GroupSlice
                        if color == Black {
                            _, adjToKoSameColor, _ = boardPtr.GetEnvironment(koPos)
                        } else {
                            _, _, adjToKoSameColor = boardPtr.GetEnvironment(koPos)
                        }
                        for _, grp := range adjToKoSameColor {
                            grp.Liberties.AppendUnique(koPos)
                            boardPtr.updateLegalityForLibertiesOfExcept(grp, koPos, boardPtr.currentSequence + 1, alreadyUpdated)
                        }
                        // The player who took the ko may fill it, so make it legal for the player 'color' for the next round.
                        // TODO: this is sort of an evil hack... or is it?
                        if color == Black {
                            _, boardPtr.actionOnNextBlackMove[koPos] = boardPtr.calculateIfLegal(koPos, Black)
                            boardPtr.fieldSequencesBlack[koPos] = boardPtr.currentSequence + 1
                        } else {
                            _, boardPtr.actionOnNextWhiteMove[koPos] = boardPtr.calculateIfLegal(koPos, White)
                            boardPtr.fieldSequencesWhite[koPos] = boardPtr.currentSequence + 1
                        }
                        // Update the legality for the liberties of the groups adjacent to the new created stone
                        //printDbgMsg("Update the legality for the liberties of the groups adjacent to the new created stone\n") // </DBG>
                        for _, grp := range c.enemiesNotInAtari {
                            boardPtr.updateLegalityForLibertiesOf(grp, boardPtr.currentSequence + 1, alreadyUpdated)
                        }

                        boardPtr.ko = NewKoLock(koPos, !color)
                        return true, true
                    })
                } else {
                    // It's not a ko
                    action = NewActionFunc(b, context, func(boardPtr *Board, c *boardBoundFuncContext) (blackUpToDate, whiteUpToDate bool) {
                        //printDbgMsgf("Board.calculateIfLegal: sameColLen == nFree == 0, removeGroups = true, not ko case.\n") // <DBG>
                        //DbgHistogram.Score() // </DBG>
                        for _, grp := range c.enemiesInAtari {
                            boardPtr.removeGroup(grp)
                        }
                        boardPtr.CreateGroup(pos,color)
                        for _, grp := range c.enemiesNotInAtari {
                            boardPtr.updateGroupLiberties(grp)
                        }
                        boardPtr.ko = nil
                        return false, false
                    })
                }
                //printDbgMsgf("returned from Board.calculateIfLegal(%d, %d, %v)\n", x, y, color) // <DBG/>
                return true, action
            }
        } else {
            // There are no adjacent friendly groups, but free neighbour fields, so this move
            // is always legal. Remove adjacent enemy groups if necessary and create a new group.
            if removeGroups {
                action = NewActionFunc(b, context, func(boardPtr *Board, c *boardBoundFuncContext) (blackUpToDate, whiteUpToDate bool) {
                    //printDbgMsgf("Board.calculateIfLegal: sameColLen == 0, nFree > 0, removeGroups = true.\n") // <DBG>
                    //DbgHistogram.Score() // </DBG>
                    for _, grp := range c.enemiesInAtari {
                        boardPtr.removeGroup(grp)
                    }
                    boardPtr.CreateGroup(pos,color)
                    for _, grp := range c.enemiesNotInAtari {
                        boardPtr.updateGroupLiberties(grp)
                    }
                    boardPtr.ko = nil
                    return false, false
                })
            } else {
                action = NewActionFunc(b, context, func(boardPtr *Board, c *boardBoundFuncContext) (blackUpToDate, whiteUpToDate bool) {
                    // Experiments show that this case is run the most often

                    //printDbgMsgf("Board.calculateIfLegal: sameColLen == 0, nFree > 0, removeGroups = false.\n") // <DBG>
                    //DbgHistogram.Score() // </DBG>
                    alreadyUpdated := make([]bool, boardPtr.boardSize*boardPtr.boardSize)
                    boardPtr.ko = nil
                    boardPtr.CreateGroup(pos,color)
                    boardPtr.dropLibertyFromEach(pos, c.adjOtherColor)
                    boardPtr.updateLegalityForFreeNeighboursOf(pos, boardPtr.currentSequence + 1, alreadyUpdated)
                    boardPtr.updateLegalityForAdjacentGroups(c.adjOtherColor, boardPtr.currentSequence + 1, alreadyUpdated)
                    return true, true
                })
            }
            //printDbgMsgf("returned from Board.calculateIfLegal(%d, %d, %v)\n", x, y, color) // <DBG/>
            return true, action
        }
    } else {
        if nFree == 0 {
            if removeGroups {
                // This move captures stones and thus produces empty fields, so it is legal. Capture
                // the stones first and then join the adjacent groups of the same color.
                action = NewActionFunc(b, context, func(boardPtr *Board, c *boardBoundFuncContext) (blackUpToDate, whiteUpToDate bool) {
                    //printDbgMsgf("Board.calculateIfLegal: sameColLen > 0, nFree == 0, removeGroups = true.\n") // <DBG>
                    //DbgHistogram.Score() // </DBG>
                    for _, grp := range c.enemiesInAtari {
                        boardPtr.removeGroup(grp)
                    }
                    boardPtr.joinGroupsByPlayAt(pos, c.adjSameColor)
                    for _, grp := range c.enemiesNotInAtari {
                        boardPtr.updateGroupLiberties(grp)
                    }
                    boardPtr.ko = nil
                    return false, false
                })
                //printDbgMsgf("returned from Board.calculateIfLegal(%d, %d, %v)\n", x, y, color) // <DBG/>
                return true, action
            } else {
                // This move does not capture any stones, so its legality depends upon the liberties 
                // of the adjacent groups of the same color. We check if at least one of these group
                // has at least two liberties.
                oneHasTwo := false
                for _, g := range context.adjSameColor {
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
                    action = NewActionFunc(b, context, func(boardPtr *Board, c *boardBoundFuncContext) (blackUpToDate, whiteUpToDate bool) {
                        // Experiments show that this case is run the 3rd most often

                        //printDbgMsgf("Board.calculateIfLegal: sameColLen > 0, nFree == 0, removeGroups = false, oneHasTwo = true.\n") // <DBG>
                        //DbgHistogram.Score() // </DBG>
                        alreadyUpdated := make([]bool, boardPtr.boardSize*boardPtr.boardSize)
                        boardPtr.ko = nil
                        boardPtr.joinGroupsByPlayAt(pos, c.adjSameColor)
                        boardPtr.dropLibertyFromEach(pos, c.adjOtherColor)
                        boardPtr.updateLegalityForAdjacentGroups(c.adjOtherColor, boardPtr.currentSequence + 1, alreadyUpdated)
                        boardPtr.updateLegalityForLibertiesOf(boardPtr.fields[pos], boardPtr.currentSequence + 1, alreadyUpdated)

                        return true, true
                    })
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
                action = NewActionFunc(b, context, func(boardPtr *Board, c *boardBoundFuncContext) (blackUpToDate, whiteUpToDate bool) {
                    //printDbgMsgf("Board.calculateIfLegal: sameColLen > 0, nFree > 0, removeGroups = true.\n") // <DBG>
                    //DbgHistogram.Score() // </DBG>
                    for _, grp := range c.enemiesInAtari {
                        boardPtr.removeGroup(grp)
                    }
                    boardPtr.joinGroupsByPlayAt(pos, c.adjSameColor)
                    for _, grp := range c.enemiesNotInAtari {
                        boardPtr.updateGroupLiberties(grp)
                    }
                    boardPtr.ko = nil
                    return false, false
                })
            } else {
                action = NewActionFunc(b, context, func(boardPtr *Board, c *boardBoundFuncContext) (blackUpToDate, whiteUpToDate bool) {
                    // Experiments show that this case is run the 2nd most often

                    //printDbgMsgf("Board.calculateIfLegal: sameColLen > 0, nFree > 0, removeGroups = false.\n") // <DBG>
                    //DbgHistogram.Score() // </DBG>
                    alreadyUpdated := make([]bool, boardPtr.boardSize*boardPtr.boardSize)
                    boardPtr.ko = nil

                    /*printDbgMsgf("len(adjSameColor): %d\n", len(c.adjSameColor))
                    nblack, nwhite := boardPtr.numberOfGroups()
                    printDbgMsgf("before: nblack %d, nwhite %d\n", nblack, nwhite)*/

                    boardPtr.joinGroupsByPlayAt(pos, c.adjSameColor)

                    /*nblack, nwhite = boardPtr.numberOfGroups()
                    printDbgMsgf("after: nblack %d, nwhite %d\n", nblack, nwhite)*/

                    for _, grp := range c.enemiesNotInAtari {
                        boardPtr.updateGroupLiberties(grp)
                    }
                    boardPtr.updateLegalityForAdjacentGroups(c.adjOtherColor, boardPtr.currentSequence + 1, alreadyUpdated)
                    // update legality for the newly joined group, which is at pos
                    boardPtr.updateLegalityForLibertiesOf(boardPtr.fields[pos], boardPtr.currentSequence + 1, alreadyUpdated)

                    return true, true
                })
            }
            //printDbgMsgf("returned from Board.calculateIfLegal(%d, %d, %v)\n", x, y, color) // <DBG/>
            return true, action
        }
    }

    // Control flow should never reach here; if it does however, this is an error
    panic("Control reaced the end of Board.calculateIfLegal")
    return false, nil
}

// Tries to choose a favorable move out of candidatePos randomly. Favorable means that we don't fill own eyes. pos's marked 
// in alreadyConsidered are skipped. tries denotes the number of tries. This method returns true and the favorable bos if it was 
// able to find a move, false, 0 otherwise
func (b *Board) chooseRandomFavorableMove(candidatePos []int, color Color, alreadyConsidered []bool, tries int) (found bool, retPos int) {
    length := len(candidatePos)
    for i := 0; i < tries; i++ {
        randomIndex := b.rand.Intn(length)
        randomPos := candidatePos[randomIndex]
        if !alreadyConsidered[randomPos] {
            if b.IsLegalMove(randomPos, color) {
                // Do we want to play this move? This move is not played if one it only fills eyes.
                if b.isEyeFillingMove(randomPos, color) {
                    alreadyConsidered[randomPos] = true
                } else {
                    // yes, we want to play this move
                    found, retPos = true, randomPos
                    break
                }
            } else {
                alreadyConsidered[randomPos] = true
            }
        }
    }
    return
}

// Returns the color of the player who plays the next turn
func (b *Board) ColorOfNextPlay() Color {
    return b.colorOfNextPlay
}

// returns an equivalent but completlely independent copy of b
func (b *Board) Copy() *Board {
    sec, nsec, _ := os.Time()
    l := b.boardSize*b.boardSize
    cpy := &Board{
        fields: make([]*Group, l),
        colorOfNextPlay: b.colorOfNextPlay,
        boardSize: b.boardSize,
        currentSequence: b.currentSequence,
        actionOnNextBlackMove: make([]*actionFunc, l),
        actionOnNextWhiteMove: make([]*actionFunc, l),
        fieldSequencesBlack: make([]uint32, l),
        fieldSequencesWhite: make([]uint32, l),
        rand: rand.New(rand.NewSource(sec+nsec)),
        prisonersBlack: b.prisonersBlack,
        prisonersWhite: b.prisonersWhite,
    }
    if b.ko != nil {
        cpy.ko = &koLock{
            Pos: b.ko.Pos,
            Color: b.ko.Color,
        }
    }
    copy(cpy.fieldSequencesBlack, b.fieldSequencesBlack)
    copy(cpy.fieldSequencesWhite, b.fieldSequencesWhite)

    // Copy .fields in a seperate loop first. The next loop, in which we set .actionOnNext{Black,White}Move
    // depends upon correcty and already completely set .fields. Remember to copy each group only once!
    copiedGroups := make(map[*Group]bool)
    for i := 0; i < l; i++ {
        grp := b.fields[i]
        _, present := copiedGroups[grp]
        if !present && grp != nil  {
            gcpy := grp.Copy()
            last := gcpy.Fields.Last()
            for it := gcpy.Fields.First(); it != last; it = it.Next() {
                cpy.fields[it.Value()] = gcpy
            }
            // remember that we already copied grp
            copiedGroups[grp] = true
        }
    }

    for i := 0; i < l; i++ {
        // copy .actionOnNextBlackMove and adjust its context
        if f := b.actionOnNextBlackMove[i]; f != nil {
            cpy.actionOnNextBlackMove[i] = NewActionFunc(cpy, nil, f.f)
            _, cpy.actionOnNextBlackMove[i].context = cpy.getEnvironmentAndContext(i, Black)
        }
        // copy .actionOnNextWhiteMove and adjust its context
        if f := b.actionOnNextWhiteMove[i]; f != nil {
            cpy.actionOnNextWhiteMove[i] = NewActionFunc(cpy, nil, f.f)
            _, cpy.actionOnNextWhiteMove[i].context = cpy.getEnvironmentAndContext(i, White)
        }
    }

    return cpy
}

// Create a new, one-stone-group of 'color' at 'pos' and sets its liberties appropriately. 
// This method does not perform any legality checks or liberty updates for other groups.
func (b *Board) CreateGroup(pos int, color Color) {
// TODO: unexport this?
    newGroup := NewGroup(color)
    newGroup.Fields.Append(pos)
    nbours := b.neighboursByPos(pos)
    for _, npos := range nbours {
        if b.fields[npos] == nil {
            newGroup.Liberties.Append(npos)
        }
    }
    b.fields[pos] = newGroup
    b.actionOnNextBlackMove[pos] = nil
    b.actionOnNextWhiteMove[pos] = nil
}

// as CreateGroup
func (b *Board) CreateGroupByPoint(x, y int, color Color) {
    pos := b.xyToPos(x,y)
    b.CreateGroup(pos, color)
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

// This is used to cound the area surrounded by black and white.
// This func assumes that the connected, free fields all have size one.
func (b *Board) getArea() (black, white int) {
    black, white = 0, 0
    for pos := 0; pos < b.boardSize*b.boardSize; pos++ {
        if b.fields[pos] == nil {
            _, adjBlack, adjWhite := b.GetEnvironment(pos)
            blackLen, whiteLen := len(adjBlack), len(adjWhite)
            if blackLen == 0 && whiteLen > 0 {
                white++
            } else if blackLen > 0 && whiteLen == 0 {
                black++
            }
        }
    }
    return
}

// Returns the `environment` of 'pos', i.e. the number 'nFree' of free neighbours
// and GroupSlices 'adj{Black,White}' containing the adjacent {black,white} groups.
func (b *Board) GetEnvironment(pos int) (nFree int, adjBlack, adjWhite GroupSlice) {
    // TODO: unexport this?
    nbours := b.neighboursByPos(pos)
    adjBlack = NewGroupSlice()
    adjWhite = NewGroupSlice()
    for _, npos := range nbours {
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

// returns nFree and the context for an actionFunc for a move of the designated color at pos
func (b *Board) getEnvironmentAndContext(pos int, color Color) (int, *boardBoundFuncContext) {
    nFree, sameColor, otherColor := b.GetEnvironment(pos)
    if color == White {
        sameColor, otherColor = otherColor, sameColor
    }
    inAtari, notInAtari := b.determineGroupsAtariStatus(otherColor)
    context := &boardBoundFuncContext{
        enemiesInAtari: inAtari,
        enemiesNotInAtari: notInAtari,
        adjSameColor: sameColor,
        adjOtherColor: otherColor,
    }
    return nFree, context
}

// Returs a pointer to the group which occupies (x,y). Nil means that this 
// field is empty.
func (b *Board) GetGroupByPoint(x,y int) *Group {
    index := b.xyToPos(x,y)
    return b.fields[index]
}

// Returns the action which is performed when a stone of the designated color is played at pos on
// an empty board
func (b *Board) initialActionGenerator(pos int, color Color) *actionFunc {
    return NewActionFunc(b, nil, func(boardPtr *Board, c *boardBoundFuncContext) (blackUpToDate, whiteUpToDate bool) {
    //return func() (updateBlack, updateWhite bool) {
        boardPtr.CreateGroup(pos, color)
        boardPtr.colorOfNextPlay = !boardPtr.colorOfNextPlay
        return false, false
    })
}

// is playing a stone of the designated color at pos an eye filling move?
func (b *Board) isEyeFillingMove(pos int, color Color) bool {
    // TODO: Write tests for this!!
    var nFree int
    var adjSameColor, adjOtherColor GroupSlice
    if color == Black {
        nFree, adjSameColor, adjOtherColor = b.GetEnvironment(pos)
    } else {
        nFree, adjOtherColor, adjSameColor = b.GetEnvironment(pos)
    }
    return nFree == 0 && len(adjSameColor) == 1 && len(adjOtherColor) == 0
}

// Is it legal to play a stone of color 'color' at 'pos'?
func (b *Board) IsLegalMove(pos int, color Color) bool {
    /*printDbgMsgf("IsLegalMove(%d, %d, %s): fieldSeqW[pos]: %d, fieldSeqB[pos]: %d, currSeq: %d\n", x,y,color, b.fieldSequencesWhite[pos], b.fieldSequencesBlack[pos],
                b.currentSequence)*/
    if color == Black {
        if b.fieldSequencesBlack[pos] != b.currentSequence {
            b.updateLegalityForBlack(pos, b.currentSequence)
        }
        if b.actionOnNextBlackMove[pos] != nil {
            return true
        }
    } else {
        if b.fieldSequencesWhite[pos] != b.currentSequence {
            b.updateLegalityForWhite(pos, b.currentSequence)
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

// Returns a slice of legal moves for color 'color' as Points
func (b *Board) ListLegalPoints(color Color) []Point {
    legalPosses := b.listLegalPosses(color)
    ret := make([]Point, len(legalPosses))
    for i, pos := range legalPosses {
        x, y := b.posToXY(pos)
        ret[i] = *NewPoint(x,y)
    }
    return ret
}

// Returns a slice of legal moves for color 'color' as posses
func (b *Board) listLegalPosses(color Color) []int {
    b.updateLegalMoves(color)

    var actions []*actionFunc
    if color == Black {
        actions = b.actionOnNextBlackMove
    } else {
        actions = b.actionOnNextWhiteMove
    }
    ret := make([]int, b.boardSize*b.boardSize)
    index := 0
    for i := 0; i < b.boardSize*b.boardSize; i++ {
        if actions[i] != nil {
            ret[index] = i
            index++
        }
    }
    return ret[0:index]
}

// Returns the neighbours of 'pos' as a pos
func (b *Board) neighboursByPos(pos int) []int {
    // TODO: remove the ..ByPos in the name
    return neighbourCache[b.boardSize][pos]
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

// returns the number of prisoners for black and white
func (b *Board) numberOfPrisoners() (nblack, nwhite int) {
    return b.prisonersBlack, b.prisonersWhite
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
func (b *Board) PlayMove(x,y int, color Color) (err Error) {
    pos := b.xyToPos(x,y)
    return b.playMoveByPos(pos, color)
}

// Play a move of color 'color' at 'pos'
func (b *Board) playMoveByPos(pos int, color Color) (err Error) {
    if !b.IsLegalMove(pos, color) {
        x, y := b.posToXY(pos)
        return NewIllegalMoveError(x,y, color)
    }

    var action *actionFunc
    if color == White {
        action = b.actionOnNextWhiteMove[pos]
    } else {
        action = b.actionOnNextBlackMove[pos]
    }
    blackUpToDate, whiteUpToDate := action.Call()
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
    b.currentSequence++

    return nil
}

// The player of color 'color' plays a pass.
func (b *Board) PlayPass(color Color) {
    b.colorOfNextPlay = !color
    b.currentSequence++
}

// Plays a random move for player 'color' and returns the played vertex.
func (b *Board) PlayRandomMove(color Color) Vertex {
    // FIXME: what about sekis and the like????

    // Collect empty fields
    emptyPos := make([]int, b.boardSize*b.boardSize)
    alreadyConsidered := make([]bool, b.boardSize*b.boardSize)
    index := 0
    for i := 0; i < b.boardSize*b.boardSize; i++ {
        if b.fields[i] == nil {
            emptyPos[index] = i
            index++
        }
    }
    emptyPos = emptyPos[0:index]
    // Now pick a random move and play it if it is legal. randomTries how often only random
    // moves should be picked. If we pick illegal moves more often than randomTries, we determine
    // all legal moves and pick one of them
    const randomTries = 4
    if found, pos := b.chooseRandomFavorableMove(emptyPos, color, alreadyConsidered, randomTries); found {
        x, y := b.posToXY(pos)
        b.PlayMove(x,y,color)
        return *NewVertexByInts(x,y,false)
    }
    // this didn't work, so we have to look at all legal moves
    legalMoves := b.listLegalPosses(color)
    if len(legalMoves) == 0 {
        // TODO: panic here? There should be some legal moves....
        return *NewVertexByInts(0,0,true)
    }
    if found, pos := b.chooseRandomFavorableMove(legalMoves, color, alreadyConsidered, randomTries); found {
        x, y := b.posToXY(pos)
        b.PlayMove(x,y,color)
        return *NewVertexByInts(x,y,false)
    }
    // Guessing inside the legal moves didn't yield a favorable one, so we have to go through these systematically
    for _, pos := range legalMoves {
        if !alreadyConsidered[pos] {
            if b.isEyeFillingMove(pos, color) {
                alreadyConsidered[pos] = true
            } else {
                // yey! We want to play this move
                x, y := b.posToXY(pos)
                b.PlayMove(x,y,color)
                return *NewVertexByInts(x,y,false)
            }
        }
    }
    // Okay then, we don't want to play any move, so we pass
    b.PlayPass(color)
    return *NewVertexByInts(0,0,true)
}

// Plays out the sequence 'seq' of moves
func (b *Board) playSequence(seq []Move) {
    for _, m := range seq {
        if m.Vertex.Pass {
            b.PlayPass(m.Color)
        } else {
            /*v, _ := pointToGTPVertex(*NewPoint(m.Vertex.X, m.Vertex.Y))
            printDbgMsgf("playing %s stone at %s\n", m.Color, v)*/
            b.PlayMove(m.Vertex.X, m.Vertex.Y, m.Color)
        }
    }
}

func (b *Board) posToXY(pos int) (x, y int) {
    return posToXY(pos, b.boardSize)
}

// Removes the group which occupies 'pos', if there is any, and updates b.emptyFields.
// This method does not alter legalWhiteMoves or legalBlackMoves.
// TODO(david): Do I really want to export this?
func (b *Board) RemoveGroupByPos(pos int) {
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
    var prisoners *int
    if group.Color == Black {
        prisoners = &b.prisonersBlack
    } else {
        prisoners = &b.prisonersWhite
    }
    for e := group.Fields.First(); e != last; e = e.Next() {
        // Collect adjacend groups so that we can update their liberties later
        pos := e.Value()
        //x, y := b.posToXY(pos)
        nbours := b.neighboursByPos(pos)
        for _, npos := range nbours {
            if grp := b.fields[npos]; grp != nil {
                adjGroups.PushUnique(grp)
            }
        }
        b.fields[pos] = nil
        // count prisoners
        *prisoners++
    }
    // Update the liberties of the adjacent groups.
    for _, grp := range adjGroups {
        b.updateGroupLiberties(grp)
    }
}

// Resets the board. The state is the same as after creating a new instance.
func (b *Board) Reset() {
    sec, nsec, _ := os.Time()
    b.rand = rand.New(rand.NewSource(sec+nsec))
    b.ko = nil
    b.colorOfNextPlay = Black
    b.currentSequence = 0
    b.prisonersWhite = 0
    b.prisonersBlack = 0
    for i := 0; i < b.boardSize*b.boardSize; i++ {
        b.fields[i] = nil
        b.actionOnNextBlackMove[i] = b.initialActionGenerator(i, Black)
        b.actionOnNextWhiteMove[i] = b.initialActionGenerator(i, White)
        b.fieldSequencesBlack[i] = 0
        b.fieldSequencesWhite[i] = 0
    }

}

// The player whose turn it is plays a stone (x,y). If an error occurs (such as that 
// this place is already occupied) this error is returned. This method assumes that 
// b.actionOnNextMove is correcty set
func (b *Board) TurnPlayMove(x, y int) (err Error) {
    // TODO: write tests for this...
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
        //x, y := b.posToXY(it.Value())
        nbours := b.neighboursByPos(it.Value())
        for _, npos := range nbours {
            //npos := b.xyToPos(p.X, p.Y)
            if b.fields[npos] == nil {
                group.Liberties.AppendUnique(npos)
            }
        }
    }
}

// Updates the legality for posToXY(pos) and makes sure the state of the board remains correct. 
// 'whichSequence' denotes the sequence to set in b.fieldSequences{Black,White}.
func (b *Board) updateLegalityFor(pos int, whichSequence uint32) {
    //printDbgMsgBTf(4,"in updateLegalityFor(%s, %d), currentSequence: %d\n", v, whichSequence, b.currentSequence)

    // use this 'if' only for debugging!
    /*if b.fieldSequencesBlack[pos] == whichSequence || b.fieldSequencesWhite[pos] == whichSequence {
        pX, pY := b.posToXY(pos)
        v, _ := pointToGTPVertex(*NewPoint(pX, pY))
        panic(NewFieldLegalityCheckedMoreThanOnceError(fmt.Sprintf("checked more than once for %s",v)))
    }*/

    b.updateLegalityForBlack(pos, whichSequence)
    b.updateLegalityForWhite(pos, whichSequence)
}

// Checks if a black move at 'pos' is legal and makes sure the state of the board remains correct. 
// 'whichSequence' denotes the sequence to set in b.fieldSequences{Black,White}.
func (b *Board) updateLegalityForBlack(pos int, whichSequence uint32) {
    _, b.actionOnNextBlackMove[pos] = b.calculateIfLegal(pos, Black)
    b.fieldSequencesBlack[pos] = whichSequence
}

// Checks if a white move at 'pos' is legal and makes sure the state of the board remains correct. 
// 'whichSequence' denotes the sequence to set in b.fieldSequences{Black,White}.
func (b *Board) updateLegalityForWhite(pos int, whichSequence uint32) {
    _, b.actionOnNextWhiteMove[pos] = b.calculateIfLegal(pos, White)
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

            /*x, y := b.posToXY(lpos)
            v, _ := pointToGTPVertex(*NewPoint(x, y))
            printDbgMsgf("updateLegalityForAdjacentGroups walks over %s\n", v)
            if x == 3 && y == 2 {
                printDbgMsgf("we are at %s, this is the board:\n", v)
                PrintBoard(b)
            }*/

            if !alreadyUpdated[lpos] {
                b.updateLegalityFor(lpos, whichSequence)
                alreadyUpdated[lpos] = true
            }
        }
    }
}

// Helper method for Board.calculateIfLegal. Does what it name indicates and skips the fields which are mared in 
// alreadyUpdated.
func (b *Board) updateLegalityForFreeNeighboursOf(pos int, whichSequence uint32, alreadyUpdated []bool) {
    //printDbgMsg("in updateLegalityForFreeNeighboursOf\n") // <DBG>
    //defer printDbgMsgf("returned from updateLegalityForFreeNeighboursOf\n") // </DBG>
    nbours := b.neighboursByPos(pos)
    for _, npos := range nbours {
        //npos := b.xyToPos(p.X, p.Y)
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
func (b *Board) updateLegalityForLibertiesOfExcept(group *Group, exceptPos int, whichSequence uint32, alreadyUpdated []bool) {
    //printDbgMsg("in Board.updateLegalityForLibertiesOfExcept\n") // <DBG/>
    //defer printDbgMsg("returned from Board.updateLegalityForLibertiesOfExcept\n") // <DBG/>
    lastLib := group.Liberties.Last()
    for itLib := group.Liberties.First(); itLib != lastLib; itLib = itLib.Next() {
        //exceptPos := b.xyToPos(exceptX, exceptY)
        lpos := itLib.Value()
        if exceptPos != lpos && !alreadyUpdated[lpos] {
            b.updateLegalityFor(lpos, whichSequence)
            alreadyUpdated[lpos] = true
        }
    }
}

// Updates the legal moves for the color 'color'
func (b *Board) updateLegalMoves(color Color) {
// TODO: write tests for this...
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
    return xyToPos(x,y, b.boardSize)
}

// ##################### Board helper functions ##########################
// Creates a new, initial board of size 'boardsize'.
func NewBoard(boardsize int) *Board {
    ret := &Board{ 
        fields: make([]*Group, boardsize*boardsize),
        actionOnNextBlackMove: make([]*actionFunc, boardsize*boardsize),
        actionOnNextWhiteMove: make([]*actionFunc, boardsize*boardsize),
        boardSize: boardsize,
        fieldSequencesBlack: make([]uint32, boardsize*boardsize),
        fieldSequencesWhite: make([]uint32, boardsize*boardsize),
    }
    ret.Reset()
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


// ################################################################################
// ########################### common helper funcs ################################
// ################################################################################
func calculateNeighbours(x, y, boardsize int) []Point {
    ret := make([]Point, 4)
    count := 0
    switch x {
        case 0:
            ret[count] = Point{ 1, y }
            count++
        case boardsize-1:
            ret[count] = Point{ boardsize - 2 , y }
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
        case boardsize-1:
            ret[count] = Point{ x, boardsize - 2 }
            count++
        default:
            ret[count] = Point{ x, y-1 }
            count++
            ret[count] = Point{ x, y+1 }
            count++
    }
    return ret[0:count]
}

func xyToPos(x, y, boardSize int) int {
    return boardSize*y + x
}

func posToXY(pos, boardsize int) (x, y int) {
    return pos%boardsize, pos/boardsize
}



