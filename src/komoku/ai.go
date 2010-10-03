/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */

package komoku

// ################################################################################
// ########################### AI struct ##########################################
// ################################################################################
type AI struct {
    topNode *TreeNode // The current top node of the scoring tree
    environment *Environment
}

// ##################### AI methods ##########################

// Runs one simulation originating from the current state in a. This func also scores in the game tree.
func (a *AI) RunSimulation() {
    board := a.environment.CurrentGame.Board.Copy()

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

// ##################### AI helper functions ##########################
func NewAI(boardsize int) *AI {
    return &AI{
        topNode: NewTreeNode(nil),
        environment: NewEnvironment(boardsize),
    }
}



