/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */

/*
 * This file contains the GTP commands. It consists of functions of the type
 * gtpcommand_name(*GTPObject) *GTPCommand which return the GTP command
 * "command_name". The command name should always be lower case!
 *
 * Portions of the comments are copied word by word from the GTP version 2 
 * specification, to be found at (in August 2010):
 * http://www.lysator.liu.se/~gunnar/gtp/    
 * This specification was written by Gunnar Farnebäck in Oct. 2002.
 */

package komoku

import (
    "sort"
    "container/vector"
    "fmt"
    "rand"
    "os"
)

// The board size is changed. The board configuration, number of captured stones, and move history become arbitrary.
// TODO: not yet implemented completely
func gtpboardsize(obj *GTPObject) *GTPCommand {
    signature := []int { GTPInt }
    f := func(object *GTPObject, params []interface{}) (result string, quit bool, err Error) {
        boardsize, ok := params[0].(uint)
        if !ok {
            panic("\n\nType assertion for first parameter of boardsize failed.\n\n")
        }
        if boardsize > 25 {
            return "unacceptable size", false, NewUnacceptableBoardSizeError()
        }

        // TODO: get rid of this cast
        object.env.CurrentGame.B = NewBoard(int(boardsize))
        return result, false, nil
    }
    return &GTPCommand{ Signature: signature,
                        Func: f,
                      }
}

// The board is cleared, the number of captured stones is reset to zero for both colors and the move history is reset to empty.
// TODO: not yet implemented completely
func gtpclear_board(obj *GTPObject) *GTPCommand {
    signature := []int { }
    f := func(object *GTPObject, params []interface{}) (result string, quit bool, err Error) {
        curSize := object.env.CurrentGame.B.BoardSize()
        object.env.CurrentGame.B = NewBoard(curSize)
        return result, false, nil
    }
    return &GTPCommand{ Signature: signature,
                        Func: f,
                      }
}

// Generate a move of the requested color. This is where the AI kicks in.
func gtpgenmove(obj *GTPObject) *GTPCommand {
    signature := []int { GTPColor }
    f := func(object *GTPObject, params []interface{}) (result string, quit bool, err Error) {
        color, _ := params[0].(Color)
        legalMoves := object.env.CurrentGame.B.ListLegalPoints(color)
        sec, nsec, _ := os.Time()
        random := rand.New(rand.NewSource(sec+nsec))
        randomMove := legalMoves[random.Intn(len(legalMoves))]
        obj.env.CurrentGame.B.PlayMove(randomMove.X, randomMove.Y, color)
        r, ok := pointToGTPVertex(randomMove)
        if !ok {
            panic("\n\nThe random move is a malformed coordinate.\n\n")
        }
        return r, false, nil
    }
    return &GTPCommand{ Signature: signature,
                        Func: f,
                      }
}

// Expexts one string argument, called 'cmdName'. Prints "true" if the command is known, "false" otherwise.
func gtpknown_command(obj *GTPObject) *GTPCommand {
    signature := []int { GTPString }
    f := func(object *GTPObject, params []interface{}) (result string, quit bool, err Error) {
        result = "true"
        cmdName, _ := params[0].(string) // type checking should have been done before, we assume that this works.
        if _, ok1 := object.commands[cmdName]; !ok1 {
            result = "false"
        }
        return result, false, nil
    }
    return &GTPCommand{ Signature: signature,
                        Func: f,
                      }
}

// The komi is changed. 
func gtpkomi(obj *GTPObject) *GTPCommand {
    signature := []int { GTPFloat }
    f := func(object *GTPObject, params []interface{}) (result string, quit bool, err Error) {
        newKomi, ok := params[0].(float)
        if !ok {
            panic("\n\nType assertion for first parameter of komi failed.\n\n")
        }
        obj.env.CurrentGame.Komi = newKomi
        return result, false, nil
    }
    return &GTPCommand{ Signature: signature,
                        Func: f,
                      }
}

// Shows all legal moves of the specified color
func gtpkomoku_alllegal(obj *GTPObject) *GTPCommand {
    signature := []int { GTPColor }
    f := func(object *GTPObject, params []interface{}) (result string, quit bool, err Error) {
        color, _ := params[0].(Color)
        b := obj.env.CurrentGame.B
        legalPoints := b.ListLegalPoints(color)
        return printBoardPrimitive(b, "", -1, -1, legalPoints) , false, nil
    }
    return &GTPCommand{ Signature: signature,
                        Func: f,
                      }
}

// Calls Board.GetEnvironment for the specified point
func gtpkomoku_getenv(obj *GTPObject) *GTPCommand {
    signature := []int { GTPVertex }
    f := func(object *GTPObject, params []interface{}) (result string, quit bool, err Error) {
        vertex := params[0].(Vertex)
        if vertex.Pass {
            emsg := "argument 0 has to be a vertex other than pass"
            return emsg, false, NewGTPSyntaxError(emsg)
        }
        nFree, adjBlack, adjWhite := obj.env.CurrentGame.B.GetEnvironment(vertex.X, vertex.Y)
        return fmt.Sprintf("nFree: %d, len(adjBlack): %d, len(adjWhite): %d", nFree, adjBlack.Length(), adjWhite.Length()), false, nil
    }
    return &GTPCommand{ Signature: signature,
                        Func: f,
                      }
}

// Expexts one string argument, called 'cmdName'. Prints the arguments of this command
func gtpkomoku_infocmd(obj *GTPObject) *GTPCommand {
    signature := []int { GTPString }
    f := func(object *GTPObject, params []interface{}) (result string, quit bool, err Error) {
        result = "komoku-infocmd: "
        cmdName, _ := params[0].(string) // type checking should have been done before, we assume that this works.
        if gtpCmd, ok1 := object.commands[cmdName]; !ok1 {
            result += "unknown command: '" + cmdName + "'"
        } else {
            result += cmdName + " "
            if len(gtpCmd.Signature) == 0 {
                result += "has 0 arguments"
            } else {
                for _, t := range gtpCmd.Signature {
                    switch t {
                        case GTPBool:
                            result += "bool "
                        case GTPColor:
                            result += "color "
                        case GTPFloat:
                            result += "float "
                        case GTPInt:
                            result += "int "
                        case GTPVertex:
                            result += "vertex "
                        case GTPString:
                            result += "string "
                        default:
                            panic("\n\nThe signature of " + cmdName + " is set erroneous.\n\n")
                    }
                }
            }
        }
        return result, false, nil
    }
    return &GTPCommand{ Signature: signature,
                        Func: f,
                      }
}

// Prints the number of groups in this format: "#black: <number>, #white: <number>"
func gtpkomoku_numgroups(obj *GTPObject) *GTPCommand {
    signature := []int {}
    f := func(object *GTPObject, params []interface{}) (result string, quit bool, err Error) {
        nblack, nwhite := obj.env.CurrentGame.B.numberOfGroups()
        return fmt.Sprintf("#black: %d, #white:%d", nblack, nwhite), false, nil
    }
    return &GTPCommand{ Signature: signature,
                        Func: f,
                      }
}

// Prints the number of stones in this format: "#black: <number>, #white: <number>"
func gtpkomoku_numstones(obj *GTPObject) *GTPCommand {
    signature := []int {}
    f := func(object *GTPObject, params []interface{}) (result string, quit bool, err Error) {
        nblack, nwhite := obj.env.CurrentGame.B.numberOfStones()
        return fmt.Sprintf("#black: %d, #white:%d", nblack, nwhite), false, nil
    }
    return &GTPCommand{ Signature: signature,
                        Func: f,
                      }
}

// List all commands, one by each line, sorted alphabetically
func gtplist_commands(obj *GTPObject) *GTPCommand {
    signature := []int {}
    f := func(object *GTPObject, params []interface{}) (result string, quit bool, err Error) {
        result = ""
        var cmdVector vector.StringVector
        for cmdName, _ := range obj.commands {
            cmdVector.Push(cmdName)
        }
        sort.SortStrings(sort.StringArray(cmdVector))
        result += cmdVector[0]
        for i := 1; i < cmdVector.Len(); i++ {
            result += "\n" + cmdVector.At(i)
        }
        return result, false, nil
    }
    return &GTPCommand{ Signature: signature,
                        Func: f,
                      }
}

// Print the name of this program, i.e. "komoku"
func gtpname(obj *GTPObject) *GTPCommand {
    signature := []int {}
    f := func(object *GTPObject, params []interface{}) (result string, quit bool, err Error) {
        return komokuProgramName, false, nil
    }
    return &GTPCommand{ Signature: signature,
                        Func: f,
                      }
}

// Arguments: color vertex. A stone of the requested color is played at the requested vertex and
// every action which has to be done is performed.
func gtpplay(obj *GTPObject) *GTPCommand {
    signature := []int { GTPColor, GTPVertex }
    f := func(object *GTPObject, params []interface{}) (result string, quit bool, err Error) {
        color, _ := params[0].(Color)
        vertex, _ := params[1].(Vertex)
        if vertex.Pass {
            obj.env.CurrentGame.B.PlayPass(color)
            return "", false, nil
        }
        //fmt.Printf("gtpplay: coords: (%d,%d)\n", vertex.X, vertex.Y)
        //fmt.Printf("gtpplay: vertex: %v\n", vertex)
        if er := obj.env.CurrentGame.B.PlayMove(vertex.X, vertex.Y, color); er != nil {
            if er.Errno() == ErrIllegalMove {
                return "illegal move", false, er
            } else {
                panic("\n\nBoard.PlayMove returned an error != ErrIllegalMove.\n\n")
            }
        }
        // Everything went fine
        return "", false, nil
    }
    return &GTPCommand{ Signature: signature,
                        Func: f,
                      }
}

// Print the protocol version. This implementation supports only version 2
func gtpprotocol_version(obj *GTPObject) *GTPCommand {
    signature := []int {}
    f := func(object *GTPObject, params []interface{}) (result string, quit bool, err Error) {
        return "2", false, nil
    }
    return &GTPCommand{ Signature: signature,
                        Func: f,
                      }
}

// Quit komoku
func gtpquit(obj *GTPObject) *GTPCommand {
    signature := []int {}
    f := func(object *GTPObject, params []interface{}) (result string, quit bool, err Error) {
        return "", true, nil
    }
    return &GTPCommand{ Signature: signature,
                        Func: f,
                      }
}

// Print the board
func gtpshowboard(obj *GTPObject) *GTPCommand {
    signature := []int { }
    f := func(object *GTPObject, params []interface{}) (result string, quit bool, err Error) {
        b := object.env.CurrentGame.B
        result = "\n" + printBoardPrimitive(b, "", -1, -1, nil)
        return result, false, nil
    }
    return &GTPCommand{ Signature: signature,
                        Func: f,
                      }
}

// Print the version of komoku
func gtpversion(obj *GTPObject) *GTPCommand {
    signature := []int {}
    f := func(object *GTPObject, params []interface{}) (result string, quit bool, err Error) {
        return komokuVersion, false, nil
    }
    return &GTPCommand{ Signature: signature,
                        Func: f,
                      }
}

