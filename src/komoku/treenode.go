/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */
package komoku

import (
    //"container/vector"
    //"fmt"
)

/*
 * ############### NodeInfo struct ############
 * This represents information attached to a node
 */
type NodeInfo struct {
    simulations int // total number of simulations that begin with this move
    wonByBlack, wonByWhite int // number of games won by {black,white}
    jigo int // number of jigos
}

/*
 * ############# TreeNode struct ##################
 */
type TreeNode struct {
    parent *TreeNode // nil means that node is at the root of the tree
    children map[int]*TreeNode // maps pos onto childnodes. The key -1 denotes a pass
    isLeaf bool // true iff this node is a leaf, i.e. if it has no children
    *NodeInfo
}

/*
 * ############# methods of TreeNode ##################
 */

// returns a pointer to the childnode with the given pos and creates it if necessary
func (t *TreeNode) ChildNode(pos int) *TreeNode {
    child, ok := t.children[pos];
    if !ok {
        newNode := NewTreeNode(t)
        t.children[pos] = newNode
        child = newNode
    }
    return child
}

// Deletes the t.NodeInfo and all its children
func (t *TreeNode) Clear() {
    t.NodeInfo = nil
    for _, child := range t.children {
        child.Clear()
    }
}

// Increments the denoted scores
func (t *TreeNode) IncrementScore(simuls, wonBlack, wonWhite, jigo int) {
    t.NodeInfo.simulations += simuls
    t.NodeInfo.wonByBlack += wonBlack
    t.NodeInfo.wonByWhite += wonWhite
    t.NodeInfo.jigo += jigo
}

/*
 * ############# helper functions ##################
 */
// Creates a new TreeNode that is a child of parent
func NewTreeNode(parent *TreeNode) *TreeNode {
    return &TreeNode{
        NodeInfo: &NodeInfo{},
        parent: parent,
        children: make(map[int]*TreeNode),
    }
}


