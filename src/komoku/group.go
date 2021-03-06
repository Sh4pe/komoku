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
    Color
    Fields *IntList // The fields which are occupied by the group
    Liberties *IntList
}

// ##################### Group methods ##########################

// returns an equivalent but completlely independent copy of g
func (g *Group) Copy() *Group {
    cpy := NewGroup(g.Color)
    // copy Fields
    last := g.Fields.Last()
    var it *IntListNode
    for it = g.Fields.First(); it != last; it = it.Next() {
        cpy.Fields.Append(it.Value())
    }
    // copy Liberties
    last = g.Liberties.Last()
    for it = g.Liberties.First(); it != last; it = it.Next() {
        cpy.Liberties.Append(it.Value())
    }
    return cpy
}

func (g *Group) NumLiberties() int {
    return g.Liberties.Length()
}

func (g *Group) NumStones() int {
    return g.Fields.Length()
}

// ##################### Group helper functions ##########################

// Creates a new empty group of color 'c'.
func NewGroup(c Color) *Group {
    return &Group{
        Color: c,
        Fields: NewIntList(),
        Liberties: NewIntList(),
    }
}

// ################################################################################
// ########################### GroupIndexType ######################################
// ################################################################################

type GroupIndexType int

// Makes 'g' emypt, i.e. sets it to 0
func (g *GroupIndexType) Clear() {
    *g = 0
}

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

// `Appends` 'group' to the GroupMap, returns its index
func (gm *GroupMap) Append(group *Group) GroupIndexType {
    ret := gm.topIndex
    gm.mapping[gm.topIndex] = group
    gm.topIndex++
    return ret
}

// Do calls f for each key, value pair of the GroupMap
func (gm *GroupMap) Do(f func(key GroupIndexType, group *Group)) {
    for key, value := range gm.mapping {
        f(key, value)
    }
}

func (gm *GroupMap) Get(index GroupIndexType) (group *Group) {
    return gm.mapping[index]
}

// Returns the number of elements in the GroupMap
func (gm *GroupMap) Length() int {
    return len(gm.mapping)
}

// Removes 'index' from 'gm' if it was contained in it.
func (gm *GroupMap) Remove(index GroupIndexType) {
    gm.mapping[index] = nil, false
}

// ##################### GroupMap helper functions ##########################

func NewGroupMap() *GroupMap {
    return &GroupMap{
        mapping: make(map[GroupIndexType]*Group),
        topIndex: 1, // 0 is reserved to represent empty fields
    }
}

// ################################################################################
// ########################### GroupSlice struct ##################################
// ################################################################################

// ##################### constants for GroupSlice ##########################
const (
    initialGroupSliceSize = 4
    groupSliceGrow = 4
)

// GroupSlice is a reduced version of a vector of *Groups
type GroupSlice []*Group

// ##################### GroupSlice methods ##########################

// Increases the length of the slice by 1, reallocates if necessary.
func (g *GroupSlice) grow() {
    if length := len(*g); length + 1 > cap(*g) {
        newSlice := make([]*Group, length + 1, length + groupSliceGrow)
        copy(newSlice, *g)
        *g = newSlice
    } else {
        *g = (*g)[0:length+1]
    }
}

// Appends group to the GroupSlice
func (g *GroupSlice) Push(group *Group) {
    g.grow()
    (*g)[len(*g)-1] = group
}

// Appends group to the GroupSlice only if group is not already contained.
func (g *GroupSlice) PushUnique(group *Group) {
    for _, grp := range *g {
        if grp == group {
            return
        }
    }
    g.Push(group)
}

// ##################### GroupSlice helper functions ##########################

func NewGroupSlice() GroupSlice {
    gs := make([]*Group, 0, initialGroupSliceSize)
    return GroupSlice(gs)
}

