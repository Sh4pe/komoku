/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */

/*
 * Here are some small debug helpers
 */
package komoku

import (
    "runtime"
    "strings"
    "fmt"
    "os"
)

const (
    printDebugOutput = true // To omit debug output, you should remove the reporting intline code. This should always
                            // be set to true to make sure that the reporting inline code is all properly commented out.
)

func printDbgMsg(msg string) {
    if !printDebugOutput {
        return
    }
    _, callerFilePath, callerLine, _ := runtime.Caller(1)
    splitPath := strings.Split(callerFilePath, "/", -1)
    callerFile := splitPath[len(splitPath) - 1]
    fmt.Fprintf(os.Stderr, "[%s:%d] %s\n", callerFile, callerLine, msg)
}

func printDbgMsgf(format string, a ...interface{}) {
    if !printDebugOutput {
        return
    }
    _, callerFilePath, callerLine, _ := runtime.Caller(1)
    splitPath := strings.Split(callerFilePath, "/", -1)
    callerFile := splitPath[len(splitPath) - 1]
    prefix := fmt.Sprintf("[%s:%d] ", callerFile, callerLine)
    fmt.Fprintf(os.Stderr, prefix+format, a)
}
