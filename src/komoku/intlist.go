/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */
package komoku


/* 
 * TODO:
 *      - .Sequence for IntList?
 *      - Write test for .Clear()
 */

// ################################################################################
// ########################### struct IntList #####################################
// ################################################################################

// This is a doubly-linked list of ints. It is assumed that the entries are 
// pairwise different.
type IntList struct {
    first *IntListNode
    last *IntListNode
    length int
}

type IntListNode struct {
    prev *IntListNode
    value int
    next *IntListNode
}

// ################### methods of IntList ###########################

// Appends the value 'v' to the end of 'il'. This method does not check if
// 'v' is already contained in 'il'.
func (il *IntList) Append(v int) {
    secondLast := il.last.prev
    newNode := newIntListNode(secondLast, il.last, v)
    secondLast.next = newNode
    il.last.prev = newNode
    il.length++
}

// Appends 'v' only if it is not already contained in 'il'
func (il *IntList) AppendUnique(v int) {
    last := il.Last()
    for it := il.First(); it != last; it = it.Next() {
        if it.Value() == v {
            return
        }
    }
    il.Append(v)
}

// Clears the whole 'il' entirely
func (il *IntList) Clear() {
    il.first.next = il.last
    il.last.prev = il.first
    il.length = 0
}

// Do calls f for each entry in 'il'
func (il *IntList) Do(f func(val int)) {
    last := il.last
    for it := il.first.next; it != last; it = it.next {
        f(it.value)
    }
}

func (il *IntList) First() *IntListNode {
    return il.first.next
}

// This method joins the IntLists 'other' into 'il', in such
// a way that the entries in the resulting 'il' are pairwise different.
// 'other' is not changed by this method, you will have to delete it for yourself.
func (il *IntList) JoinUnique(other *IntList) {
    // The algorithm here is stupid. Are there better ones?
    otherLast := other.Last()
    for otherIt := other.First(); otherIt != otherLast; otherIt = otherIt.Next() {
        contained := false
        value := otherIt.Value()
        last := il.Last()
        for it := il.First(); it != last; it = it.Next() {
            if it.Value() == value {
                contained = true
                break
            }
        }
        if !contained {
            il.Append(value)
        }
    }
}

func (il *IntList) Last() *IntListNode {
    return il.last
}

func (il *IntList) Length() int {
    return il.length
}

// Removes the value 'val' from il, if it exists therein. This method (as well as the whole IntList) assumes that
// each value occurs at most once.
// Returns true if 'val' was contained in 'il'
func (il *IntList) Remove(val int) bool {
    last := il.Last()
    for it := il.First(); it != last; it = it.Next() {
        if it.Value() == val {
            it.prev.next = it.next
            it.next.prev = it.prev
            // TODO(david): are the next two lines needed (for the GC)?
            it.next = nil
            it.prev = nil

            il.length--
            return true
        }
    }
    return false
}

// ################### methods of IntListNode ###########################

func (iln *IntListNode) Next() *IntListNode {
    return iln.next
}

func (iln *IntListNode) Value() int {
    return iln.value
}

// ################### helper functions ###############################

// Creates an empty IntList
func NewIntList() *IntList {
    // -111 for no good reason. FIXME(david)?
    f := newIntListNode(nil, nil, -111)
    l := newIntListNode(f, nil, -111)
    f.next = l
    return &IntList{ first: f,
                     last: l,
                   }
}

// prevNode and nextNode are... (guess what!), val is the value
func newIntListNode(prevNode, nextNode *IntListNode, val int) *IntListNode {
    return &IntListNode{ next: nextNode,
                         prev: prevNode,
                         value: val,
                       }
}
