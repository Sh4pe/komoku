/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */

package komoku

// ################################################################################
// ########################### Group struct #######################################
// ################################################################################

type Group struct {
    Color Color
    Fields *IntList
    Liberties *IntList
}

// ##################### Group methods ##########################
func (g *Group) NumLiberties() int {
    return g.Liberties.Length()
}

func (g *Group) NumStones() int {
    return g.Fields.Length()
}

// ##################### Group helper functions ##########################

// Creates a new empty group of colof 'c'.
func NewGroup(c Color) *Group {
    return &Group{ Color: c,
                   Fields: NewIntList(),
                   Liberties: NewIntList(),
                 }
}

// ################################################################################
// ########################### GroupMap struct ####################################
// ################################################################################

type GroupMap struct {
    mapping map[uint32]*Group
    topIndex uint32
}

// ##################### GroupMap methods ##########################

// 'Appends' 'group' to the GroupMap, returns its index
func (gm *GroupMap) Append(group *Group) uint32 {
    ret := gm.topIndex
    gm.mapping[gm.topIndex] = group
    gm.topIndex++
    return ret
}

func (gm *GroupMap) Get(index uint32) (gropu *Group) {
    return
}

// ##################### GroupMap helper functions ##########################

func NewGroupMap() *GroupMap {
    return &GroupMap{ mapping: make(map[uint32]*Group),
                      topIndex: 1, // 0 is reserved to represent empty fields
                    }
}

