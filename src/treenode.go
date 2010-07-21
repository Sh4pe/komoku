package komoku

import (
    "container/vector"
    //"./gamestate"
    //"./common"
)

/*
 * ############### NodeInfo struct ############
 * This represents information attached to a node
 */
type NodeInfo struct {
}

/*
 * ############# TreeNode struct ##################
 */
type TreeNode struct {
    parent      *TreeNode // nil means that node is at the root of the tree
    children    vector.Vector // children of this node
    nodeInfo    NodeInfo
}

/*
 * ############# methods of TreeNode ##################
 */
func (tn *TreeNode) ChildAt(index int) *TreeNode {
    // we don't need to check the type, because we only add TreeNodes to children
    val, _ := tn.children.At(index).(TreeNode)
    return &val
}

func (tn *TreeNode) HasChildren() bool {
    return tn.children.Len() == 0
}

func (tn *TreeNode) IsRoot() bool {
    return tn.parent == nil
}

func (tn *TreeNode) NumChildren() int {
    return tn.children.Len()
}

func (tn *TreeNode) Parent() *TreeNode {
    return tn.parent
}

/*
 * ############# helper functions ##################
 */
// creates a root node for a tree
func NewRootTreeNode() *TreeNode {
    return &TreeNode{ parent: nil }
}
