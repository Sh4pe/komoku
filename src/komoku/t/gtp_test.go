/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */
package komoku

import (
    "testing"
    //"fmt"
)

// Stuff for TestParseLine
type testParseLineCase struct {
    Line string
    Empty, HasId bool
    Id uint
    CommandName string
    Args []string
}

func TestParseLine(t *testing.T) {
    cases := []testParseLineCase {
        testParseLineCase{
            Line: "33 command_name argument1\n",
            Empty: false,
            HasId: true,
            Id: 33,
            CommandName: "command_name",
            Args: []string {
                "argument1",
            },
        },
        testParseLineCase{
            Line: "33 command_name argument1 argument2\n",
            Empty: false,
            HasId: true,
            Id: 33,
            CommandName: "command_name",
            Args: []string {
                "argument1",
                "argument2",
            },
        },
        testParseLineCase{
            Line: "33 command_name\n",
            Empty: false,
            HasId: true,
            Id: 33,
            CommandName: "command_name",
            Args: []string { },
        },
        testParseLineCase{
            Line: "command_name argument1 argument2\n",
            Empty: false,
            HasId: false,
            Id: 0,
            CommandName: "command_name",
            Args: []string {
                "argument1",
                "argument2",
            },
        },
        testParseLineCase{
            Line: "command_name argument1 \n",
            Empty: false,
            HasId: false,
            Id: 0,
            CommandName: "command_name",
            Args: []string {
                "argument1",
            },
        },
        testParseLineCase{
            Line: "command_name \n",
            Empty: false,
            HasId: false,
            Id: 0,
            CommandName: "command_name",
            Args: []string { },
        },
        testParseLineCase{
            Line: "\n",
            Empty: true,
            HasId: false,
            Id: 0,
            CommandName: "",
            Args: []string { },
        },
        testParseLineCase{
            Line: "23\n",
            Empty: true,
            HasId: false,
            Id: 0,
            CommandName: "",
            Args: []string { },
        },
    }

    obj := NewGTPObject(nil)
    for _, testCase := range cases {
        failed := false
        empty, hasId, id, commandName, args := obj.parseLine(testCase.Line)
        if empty != testCase.Empty { failed = true }
        if hasId != testCase.HasId { failed = true }
        if id != testCase.Id { failed = true }
        if commandName != testCase.CommandName { failed = true }
        if len(args) != len(testCase.Args) {
            failed = true
        } else {
            for i, a := range args {
                if testCase.Args[i] != a {
                    failed = true
                    break
                }
            }
        }
        if failed {
            // fmt.Printf("empty: %v\nhasId: %v\nid: %v\ncommandName: %v\nargs: %v\n",
                // testCase.Empty, testCase.HasId, testCase.Id, testCase.CommandName, testCase.Args)
            t.Fatalf("parseLine seems not to parse this line well:\n%s", testCase.Line)
        }
    }

    return
}

func Testsuite() []testing.Test {
    return []testing.Test { testing.Test{"TestParseLine", TestParseLine},
                          }
}
