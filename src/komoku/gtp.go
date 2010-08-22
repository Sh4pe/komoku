/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */

package komoku

import (
    //"os"
    "strings"
    "strconv"
    "fmt"
)

const (
    GTPBool = iota;
    GTPColor
    GTPFloat
    GTPInt
    GTPPoint
    GTPString
)

// ################################################################################
// ########################### GTPCommand et al ###################################
// ################################################################################

// The type of an executable GTP command primitive
type GTPCommandFunc func(object *GTPObject, params ...string) (result string, err Error)

// Everything to completely describe a GTP command
type GTPCommand struct {
    Signature []int // Signature of expected arguments, such as GTPBool etc..
    Func GTPCommandFunc // The actual code
}

// ################################################################################
// ########################### GTPObject struct ###################################
// ################################################################################

// This object provides the functionality for comunicating via the GTP protocol.
type GTPObject struct {
    commands map[string]*GTPCommand
}

// ##################### GTPObject methods ##########################

// This is the entry point for processing input lines. It takes the raw, yet unparsed or
// line 'input', performs everything which has to be done and returns the response string in
// 'result'. If komoku has to quit after this command (e.g. if the command is "quit"), 'quit'
// will be true, otherwise false. err is != nil if an error occurs.
func (obj *GTPObject) ExecuteCommand(input string) (result string, quit bool, err Error) {
    empty, hasId, id, commandName, args := obj.parseLine(input)
    fmt.Printf("input: '%v'\n", input)
    fmt.Printf("empty: %v\nhasId: %v\nid: %v\ncommandName: %v\nargs: %v\n", empty, hasId, id, commandName, args)
    return "not implemented", false, nil
}

// Parses a line. A command has one of the following syntaxes:
// (1) id command_name arguments\n
// (2) id command_name\n
// (3) command_name arguments\n
// (3) command_name\n
// If the line does not contain any command, 'empty' is true, false otherwise.
// The returned 'id' is only meaningful if 'hasId' == true. 'args' is a slice of arguments
func (obj *GTPObject) parseLine(line string) (empty, hasId bool, id uint, commandName string, args []string) {
    line = obj.preprocessLine(line)
    empty = false
    if line == "" {
        return true, false, 0, "", nil
    }
    split := strings.Split(line, " ", -1)
    // Check if the first part is an id
    hasId = false
    //fmt.Printf("%v\n", split)
    if i, err := strconv.Atoui(split[0]); err == nil {
        // fmt.Printf("Atoui(%s) = %d\n", split[0], i)
        // fmt.Printf("err: %s\n", err)
        hasId = true
        id = i
        split = split[1:len(split)]
    }
    //fmt.Printf("%v\n", split)
    // If there is nothing after the id, the line is treated as if it were empty.
    if len(split) == 0 {
        return true, false, 0, "", nil
    }
    commandName = split[0]
    split = split[1:len(split)]
    args = make([]string, len(split))
    for index, arg := range split {
        args[index] = arg
    }

    return
}

func (obj *GTPObject) preprocessLine(input string) (result string) {
    if hashPos := strings.Index(input, "#"); hashPos != -1 {
        result = input[0:hashPos]
    } else {
        result = input
    }
    result = strings.Replace(result, "\n", "", -1)
    result = strings.Replace(result, "\t", " ", -1)
    // Modify the string so that at contains at most one consecutive whitespace
    dblWhitespacePos := strings.Index(result, "  ")
    for dblWhitespacePos != -1 {
        result = strings.Replace(result, "  ", " ", -1)
        dblWhitespacePos = strings.Index(result, "  ")
    }
    // Remove leading and trailing whitespaces
    if strings.HasPrefix(result, " ") { result = result[1:len(result)] }
    if strings.HasSuffix(result, " ") { result = result[0:len(result)-1] }

    return
}

// ##################### GTPObject helper functions ##########################
func NewGTPObject() *GTPObject {
    ret := &GTPObject{ commands: make(map[string]*GTPCommand),
                     }

    return ret
}

