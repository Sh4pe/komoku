/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */

/*
 * TODO:
 *      - use an io.Writer for ExecuteCommand and RunGTPMode..
 */

package komoku

import (
    "os"
    "bufio"
    "strings"
    "strconv"
    "fmt"
)

const (
    GTPBool = iota;
    GTPColor
    GTPFloat
    GTPInt
    GTPVertex
    GTPString
)

// ################################################################################
// ########################### GTPCommand et al ###################################
// ################################################################################

// The type of an executable GTP command primitive. 'quit' == true means that komoku has to 
// quit.
type GTPCommandFunc func(object *GTPObject, params []interface{}) (result string, quit bool, err Error)

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
    ai *AI // pointer to the current game AI
}

// ##################### GTPObject methods ##########################

// This is the entry point for processing input lines. It takes the raw, yet unparsed or
// line 'input', performs everything which has to be done and returns the response string in
// 'result'. If komoku has to quit after this command (e.g. if the command is "quit"), 'quit'
// will be true, otherwise false. err is != nil if an error occurs.

// TODO: Write tests for the arg checking
func (obj *GTPObject) ExecuteCommand(input string) (result string, quit bool, err Error) {
    empty, hasId, id, commandName, args := obj.parseLine(input)
    if empty {
        return "", false, nil
    }

    gtpCmd, ok := obj.commands[commandName];
    if !ok { // This command is not found
        return obj.formatErrorResponse(hasId, id, "unknown command"), false, nil
    }
    // Check the arguments
    signatureLen := len(gtpCmd.Signature)
    if signatureLen != len(args) {
        return obj.formatErrorResponse(hasId, id, fmt.Sprintf("wrong number of arguments, %d argument(s) expected", signatureLen)), false, nil
    }
    argsToPass := make([]interface{}, len(args))
    for i := 0; i < len(args); i++ {
        // TODO: refactor this! 
        // TODO: Do the type conversion (e.g. gtpVertexToPoint) here
        switch gtpCmd.Signature[i] {
            case GTPBool:
                if args[i] != "true" || args[i] != "false" {
                    errmsg := fmt.Sprintf("argument %d has to be a boolean", i)
                    return obj.formatErrorResponse(hasId, id, errmsg), false, nil
                } else {
                    val := true
                    if args[i] == "false" {
                        val = false
                    }
                    argsToPass[i] = val
                }
            case GTPColor:
                color, ok := gtpColorToColor(args[i])
                if !ok {
                    errmsg := fmt.Sprintf("argument %d has to be a color", i)
                    return obj.formatErrorResponse(hasId, id, errmsg), false, nil
                } else {
                    argsToPass[i] = color
                }
            case GTPFloat:
                fval, err := strconv.Atof(args[i])
                if err != nil {
                    errmsg := fmt.Sprintf("argument %d has to be a float", i)
                    return obj.formatErrorResponse(hasId, id, errmsg), false, nil
                } else {
                    argsToPass[i] = fval
                }
            case GTPInt:
                ival, err := strconv.Atoui(args[i])
                if err != nil {
                    errmsg := fmt.Sprintf("argument %d has to be an unsigned int", i)
                    return obj.formatErrorResponse(hasId, id, errmsg), false, nil
                } else {
                    argsToPass[i] = ival
                }
            case GTPVertex:
                point, ok, pass := gtpVertexToPoint(args[i])
                if !ok {
                    errmsg := fmt.Sprintf("argument %d has to be a vertex", i)
                    return obj.formatErrorResponse(hasId, id, errmsg), false, nil
                } else {
                    //fmt.Printf("ExecuteCommand: point: %v\n", point)
                    argsToPass[i] = *NewVertex(point, pass)
                }
            case GTPString:
                argsToPass[i] = args[i]
            default:
                // This should never happen
                panic("\n\nThe signature of " + commandName + " is set erroneous.\n\n")
        }
    }
    cmdResponse, retQuit, err := gtpCmd.Func(obj, argsToPass)
    if err != nil {
        return obj.formatErrorResponse(hasId, id, cmdResponse), retQuit, nil
    }
    return obj.formatSuccessResponse(hasId, id, cmdResponse), retQuit, nil
}

// Returns the error response.
func (obj *GTPObject) formatErrorResponse(hasId bool, id uint, msg string) string {
    ret := "?"
    if hasId {
        ret += fmt.Sprintf("%d ", id)
    } else {
        ret += " "
    }
    ret += msg + "\n\n"
    return ret
}

// Returns the success response.
func (obj *GTPObject) formatSuccessResponse(hasId bool, id uint, msg string) string {
    ret := "="
    if hasId {
        ret += fmt.Sprintf("%d ", id)
    } else {
        ret += " "
    }
    ret += msg + "\n\n"
    return ret
}

// Parses a line. A command has one of the following syntaxes:
// (1) id command_name arguments\n
// (2) id command_name\n
// (3) command_name arguments\n
// (3) command_name\n
// If the line does not contain any command, 'empty' is true, false otherwise.
// The returned 'id' is only meaningful if 'hasId' == true. 'args' is a slice of arguments

// TODO: should args be []interface{} ?
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
                       ai: NewAI(DefaultBoardSize),
                     }

    // GTP commands
    ret.commands["boardsize"] = gtpboardsize(ret)
    ret.commands["clear_board"] = gtpclear_board(ret)
    ret.commands["genmove"] = gtpgenmove(ret)
    ret.commands["known_command"] = gtpknown_command(ret)
    ret.commands["komi"] = gtpkomi(ret)
    ret.commands["list_commands"] = gtplist_commands(ret)
    ret.commands["name"] = gtpname(ret)
    ret.commands["play"] = gtpplay(ret)
    ret.commands["protocol_version"] = gtpprotocol_version(ret)
    ret.commands["quit"] = gtpquit(ret)
    ret.commands["showboard"] = gtpshowboard(ret)
    ret.commands["version"] = gtpversion(ret)

    // Private extensions
    ret.commands["komoku-alllegal"] = gtpkomoku_alllegal(ret)
    ret.commands["komoku-genmovedbg"] = gtpkomoku_genmovedbg(ret)
    ret.commands["komoku-getenv"] = gtpkomoku_getenv(ret)
    ret.commands["komoku-getgroup"] = gtpkomoku_getgroup(ret)
    ret.commands["komoku-infocmd"] = gtpkomoku_infocmd(ret)
    ret.commands["komoku-numgroups"] = gtpkomoku_numgroups(ret)
    ret.commands["komoku-numstones"] = gtpkomoku_numstones(ret)
    ret.commands["komoku-playfork"] = gtpkomoku_playfork(ret)
    ret.commands["komoku-placehandi"] = gtpkomoku_placehandi(ret)
    ret.commands["komoku-showliberties"] = gtpkomoku_showliberties(ret)
    ret.commands["komoku-source"] = gtpkomoku_source(ret)
    ret.commands["komoku-sourceforkn"] = gtpkomoku_sourceforkn(ret)
    ret.commands["komoku-sourcen"] = gtpkomoku_sourcen(ret)

    return ret
}

func NewUnacceptableBoardSizeError() (err Error) {
    return NewError(fmt.Sprintf("unacceptable size"), ErrUnacceptableBoardSize)
}

func NewGTPSyntaxError(msg string) (err Error) {
    return NewError(fmt.Sprintf(msg), ErrGTPSyntaxError)
}

func NewGTPNotImplementedError(msg string) (err Error) {
    return NewError(fmt.Sprintf(msg), ErrGTPNotImplemented)
}

func NewGTPIllegalCommand(msg string) (err Error) {
    return NewError(fmt.Sprintf(msg), ErrGTPIllegalCommand)
}

// ################################################################################
// #################### Function for running the GTP-mode #########################
// ################################################################################

func RunGTPMode() {
    // Create the GTPObject and start the input loop
    gtpObject := NewGTPObject()
    in := bufio.NewReader(os.Stdin)
    for {
        line, err := in.ReadString('\n')
        switch err {
            case nil:
                result, quit, execErr := gtpObject.ExecuteCommand(line)
                if execErr != nil {
                    fmt.Printf("Error in GTPObject.ExecuteCommand:\n%s\n", execErr)
                }

                fmt.Printf(result)
                if quit {
                    return
                }
            case os.EOF:
                break
            default:
                panic("\n\nUnexpected case in RunGTPMode.\n\n")
        }
    }
}

