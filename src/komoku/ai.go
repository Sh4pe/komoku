/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */

package komoku

import (
    //"runtime"
    "fmt"
    "time"
)

// ################################################################################
// ########################### gobal variables and initialization #################
// ################################################################################
var averageGameLength = make(map[int]int)

func init() {
    // these numbers are taken from e/gamelength
    averageGameLength[ 5] = 36
    averageGameLength[ 6] = 50
    averageGameLength[ 7] = 67
    averageGameLength[ 8] = 85
    averageGameLength[ 9] = 109
    averageGameLength[10] = 137
    averageGameLength[11] = 163
    averageGameLength[12] = 191
    averageGameLength[13] = 221
    averageGameLength[14] = 259
    averageGameLength[15] = 292
    averageGameLength[16] = 333
    averageGameLength[17] = 374
    averageGameLength[18] = 419
    averageGameLength[19] = 469
    averageGameLength[20] = 514
    averageGameLength[21] = 571
    averageGameLength[22] = 621
    averageGameLength[23] = 675
    averageGameLength[24] = 741
    averageGameLength[25] = 791
}

// ################################################################################
// ########################### AI struct ##########################################
// ################################################################################
type AI struct {
    topNode *TreeNode // The current top node of the scoring tree
    environment *Environment
    numThinkers int // number of thinking goroutines
    runThinkers bool
    thinkerFinished[]chan bool // the thinkers answer here when they are finished
}

// ##################### AI methods ##########################


// Generate a move using the current statistics as a guide to the best move
// and play this move.
func (a *AI) GenMove(color Color) Vertex {
    // Give komoku 10 seconds to think. This behaviour will change in future.
    return a.genMove(color, 10000000000)
}

// Generate a move using the current statistics as a guide to the best move
// and play this move. Thinks for timeToThink nanoseconds
func (a *AI) genMove(color Color, timeToThink int64) Vertex {

    a.startThinking(true)
    time.Sleep(timeToThink)
    defer a.startThinking(a.stopThinking())


    // find the best move
    var bestWinPercentage float = -1.0
    bestPos := -1
    for pos, childNode := range a.topNode.children {
        simulations := childNode.NodeInfo.simulations
        var wonByColor int
        if color == Black {
            wonByColor = childNode.NodeInfo.wonByBlack
        } else {
            wonByColor = childNode.NodeInfo.wonByWhite
        }
        var winPercentage float = float(wonByColor)/float(simulations)
        if winPercentage > bestWinPercentage {
            bestWinPercentage = winPercentage
            bestPos = pos
        }
    }

    // play the best move
    bestX, bestY := a.environment.Game.Board.posToXY(bestPos)
    a.PlayMove(bestX, bestY, color)

    // TODO: this never plays pass yet
    return *NewVertexByInts(bestX, bestY, false)
}

// Debug version of AI.GenMove
func (a *AI) genMoveDbg(color Color, timeToThink int64) Vertex {
    fmt.Printf("\nBoard before simulations:\n")
    PrintBoard(a.environment.Game.Board)

    a.startThinking(true)
    time.Sleep(timeToThink)
    defer a.startThinking(a.stopThinking())


    // find the best move
    var bestWinPercentage float = -1.0
    bestPos := -1
    for pos, childNode := range a.topNode.children {
        simulations := childNode.NodeInfo.simulations
        var wonByColor int
        if color == Black {
            wonByColor = childNode.NodeInfo.wonByBlack
        } else {
            wonByColor = childNode.NodeInfo.wonByWhite
        }
        var winPercentage float = float(wonByColor)/float(simulations)
        if winPercentage > bestWinPercentage {
            bestWinPercentage = winPercentage
            bestPos = pos
        }
    }

    // play the best move
    bestX, bestY := a.environment.Game.Board.posToXY(bestPos)
    a.PlayMove(bestX, bestY, color)

    fmt.Printf("number of simulations: %d, bestWinPercentage: %f\n", a.topNode.NodeInfo.simulations, bestWinPercentage)
    sum := 0
    numNodes := 0
    for _, node := range a.topNode.children{
        sum += node.NodeInfo.simulations
        numNodes++
    }
    fmt.Printf("collected number of simulations: %d, number of nodes: %d\n", sum, numNodes)

    fmt.Printf("\nBoard after simulations:\n")
    PrintBoard(a.environment.Game.Board)

    // TODO: this never plays pass yet
    return *NewVertexByInts(bestX, bestY, false)
}

// Thinks until a.runThinker is false, and sends true to a.thinkerFinished[index] when finished
func (a *AI) makeThinker(index int) {
    for a.runThinkers {
        a.runSimulation()
    }
    a.thinkerFinished[index] <- true
}

// Returns the total number of simulations currently run
func (a *AI) NumSimulations() int {
    return a.topNode.NodeInfo.simulations
}

// Play a move on the board
func (a *AI) PlayMove(x,y int, color Color) (err Error) {
    // if komoku is thinking already, stop thinking and restart it afterwards
    /*if a.runThinkers {
        a.stopThinking()
        defer a.startThinking()
    }*/
    defer a.startThinking(a.stopThinking())
    return a.playMove(x,y,color)
}

// Does not stop the thinking (random game generation) before it does anything. If you
// need to stop before, think about calling PlayMove(...) instead.
func (a *AI) playMove(x,y int, color Color) (err Error) {
    if err := a.environment.Game.PlayMove(x,y,color); err != nil {
        return err
    }

    pos := a.environment.Game.Board.xyToPos(x,y)
    if childNode, ok := a.topNode.children[pos]; ok {
        // child node present
        a.topNode = childNode
    } else {
        // child node not present, so create a new node
        a.topNode = a.topNode.ChildNode(pos)
    }

    return
}


// Runs one simulation originating from the current state in a. This func also scores in the game tree.
func (a *AI) runSimulation() {
    board := a.environment.Game.Board.Copy()

    // play random games until both players pass in a row
    lastPass := false
    currentNode := a.topNode
    for {
        v := board.PlayRandomMove(board.ColorOfNextPlay())
        if v.Pass {
            if lastPass {
                break
            } else {
                lastPass = true
                currentNode = currentNode.ChildNode(-1)
            }
        } else {
            lastPass = false
            playedPos := board.xyToPos(v.X, v.Y)
            currentNode = currentNode.ChildNode(playedPos)
        }
    }

    // the game is finished, now calculate who won or if its a jigo
    prisonersBlack, prisonersWhite := board.numberOfPrisoners()
    stonesBlack, stonesWhite := board.numberOfStones()
    areaBlack, areaWhite := board.getArea()
    scoreBlack := float(prisonersBlack + stonesBlack + areaBlack)
    scoreWhite := float(prisonersWhite + stonesWhite + areaWhite) + a.environment.komi
    var wonBlack, wonWhite, jigo int
    if scoreBlack > scoreWhite {
        wonBlack = 1
        wonWhite = 0
        jigo = 0
    } else if scoreBlack < scoreWhite {
        wonBlack = 0
        wonWhite = 1
        jigo = 0
    } else if scoreBlack == scoreWhite {
        wonBlack = 0
        wonWhite = 0
        jigo = 1
    }

    // the game is finished, now score in the game tree
    for currentNode != nil {
        currentNode.IncrementScore(1, wonBlack, wonWhite, jigo)
        currentNode = currentNode.parent
    }
}

// Only starts thinking if think == true - so this can be used as a sort of 
// (rails-like) "around wrapper".
// If a is already thinking, this does nothing
func (a *AI) startThinking(think bool) {
    if !think || a.runThinkers {
        return
    }

    a.runThinkers = true
    for i := 0; i < a.numThinkers; i++ {
        go a.makeThinker(i)
    }
}

// Stop thinking and waits for all thinkers to finish. Does nothing if a is not thinking.
// Returns true iff a was thinking before
func (a *AI) stopThinking() bool {
    if !a.runThinkers {
        return false
    }

    // block until every thinker is ready
    a.runThinkers = false
    for i := 0; i < a.numThinkers; i++ {
        <-a.thinkerFinished[i]
    }

    return true
}

// ##################### AI helper functions ##########################
func NewAI(boardsize int) *AI {
    // numThinkers := runtime.GOMAXPROCS(0)
    // TODO: numThinkers = 1 seems to be the fastest, but why??
    numThinkers := 1
    a := &AI{
        numThinkers: numThinkers,
        topNode: NewTreeNode(nil),
        environment: NewEnvironment(boardsize),
        thinkerFinished: make([]chan bool, numThinkers),
    }
    for i := 0; i < numThinkers; i++ {
         a.thinkerFinished[i] = make(chan bool)
    }
    return a
}



