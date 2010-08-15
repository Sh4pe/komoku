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

// Creates a new empty group of color 'c'.
func NewGroup(c Color) *Group {
    return &Group{ Color: c,
                   Fields: NewIntList(),
                   Liberties: NewIntList(),
                 }
}

// ################################################################################
// ########################### GrouIndexType ######################################
// ################################################################################

type GroupIndexType uint32

// 0 is reserved to represent an 'empty' group
func (g *GroupIndexType) Empty() bool {
    return *g == 0
}

// ################################################################################
// ########################### GroupMap struct ####################################
// ################################################################################

type GroupMap struct {
    mapping map[GroupIndexType]*Group
    topIndex GroupIndexType
}

// ##################### GroupMap methods ##########################

// 'Appends' 'group' to the GroupMap, returns its index
func (gm *GroupMap) Append(group *Group) GroupIndexType {
    ret := gm.topIndex
    gm.mapping[gm.topIndex] = group
    gm.topIndex++
    return ret
}

func (gm *GroupMap) Get(index GroupIndexType) (group *Group) {
    return gm.mapping[index]
}

// ##################### GroupMap helper functions ##########################

func NewGroupMap() *GroupMap {
    return &GroupMap{ mapping: make(map[GroupIndexType]*Group),
                      topIndex: 1, // 0 is reserved to represent empty fields
                    }
}

