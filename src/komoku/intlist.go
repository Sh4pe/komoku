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
 */

// ################################################################################
// ########################### struct IntList #####################################
// ################################################################################

// This is a doubly-linked list of ints
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

// Appends the value 'v' to the end if il
func (il *IntList) Append(v int) {
    secondLast := il.last.prev
    newNode := newIntListNode(secondLast, il.last, v)
    secondLast.next = newNode
    il.last.prev = newNode
    il.length++
}

// Clears the whole 'il' entirely
func (il *IntList) Clear() {
    il.first.next = il.last
    il.last.prev = il.first
    il.length = 0
}

func (il *IntList) First() *IntListNode {
    return il.first.next
}

func (il *IntList) Last() *IntListNode {
    return il.last
}

func (il *IntList) Length() int {
    return il.length
}

// Removes the value 'val' from il, if it exists therein. This method (as well as the whole IntList) ssumes that
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
