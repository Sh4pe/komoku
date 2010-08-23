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
 */

package komoku

import (
    "sort"
    "container/vector"
    //"fmt"
)

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

