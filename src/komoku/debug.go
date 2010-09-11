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
    "container/vector"
    "sort"
    "runtime/pprof"
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
    fmt.Fprintf(os.Stderr, "[%s:%d] %s", callerFile, callerLine, msg)
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

// like printDbgMsgf, but adds dips one step further into the backtrace
func printDbgMsgCallerf(format string, a ...interface{}) {
    if !printDebugOutput {
        return
    }
    _, callerFilePath, callerLine, _ := runtime.Caller(1)
    _, callerFilePath2, callerLine2, _ := runtime.Caller(2)
    splitPath := strings.Split(callerFilePath, "/", -1)
    splitPath2 := strings.Split(callerFilePath2, "/", -1)
    callerFile := splitPath[len(splitPath) - 1]
    callerFile2 := splitPath2[len(splitPath2)-1]
    prefix := fmt.Sprintf("[%s:%d <- %s:%d] ", callerFile, callerLine, callerFile2, callerLine2)
    fmt.Fprintf(os.Stderr, prefix+format, a)
}

// like printDbgMsgCallerf, but with a backtrace of depth 'depth'
func printDbgMsgBTf(depth int, format string, a ...interface{}) {
    if !printDebugOutput {
        return
    }
    /*_, callerFilePath, callerLine, _ := runtime.Caller(1)
    _, callerFilePath2, callerLine2, _ := runtime.Caller(2)
    splitPath := strings.Split(callerFilePath, "/", -1)
    splitPath2 := strings.Split(callerFilePath2, "/", -1)
    callerFile := splitPath[len(splitPath) - 1]
    callerFile2 := splitPath2[len(splitPath2)-1]
    prefix := fmt.Sprintf("[%s:%d <- %s:%d] ", callerFile, callerLine, callerFile2, callerLine2)
    fmt.Fprintf(os.Stderr, prefix+format, a)*/
    pc := make([]uintptr, depth)
    d := runtime.Callers(2, pc)
    _, callerFilePath, callerLine, _ := runtime.Caller(1)
    splitPath := strings.Split(callerFilePath, "/", -1)
    callerFile := splitPath[len(splitPath) - 1]
    prefix := fmt.Sprintf("%s:%d", callerFile, callerLine)
    for i := 2; i < d; i++ {
        /*p, _, _, _ := runtime.Caller(i)
        fmt.Printf("%v, %v\n", pc[i], p)*/
        _, callerFilePath, callerLine, _ = runtime.Caller(i)
        splitPath = strings.Split(callerFilePath, "/", -1)
        callerFile = splitPath[len(splitPath) - 1]
        prefix += fmt.Sprintf(" <- %s:%d", callerFile, callerLine)
    }
    prefix = fmt.Sprintf("[%s] ", prefix)
    fmt.Fprintf(os.Stderr, prefix+format, a)
}

func ProfileInfoToFile(profFile string) {
    file, err := os.Open(profFile, os.O_CREATE | os.O_RDWR, 0666)
    if err != nil {
        fmt.Printf("unable to open .prof file %s\nerror: %s\n", profFile, err)
    }
    if err := pprof.WriteHeapProfile(file); err != nil {
        fmt.Printf("error in WriteHeapProfile: %s", err)
    }
}

func ProfileInfoToFileByRelPath(profFile string) {
    profFile = relPathToAbs(profFile)
    os.Remove(profFile)
    file, err := os.Open(profFile, os.O_CREATE | os.O_RDWR, 0666)
    if err != nil {
        fmt.Printf("unable to open .prof file %s\nerror: %s\n", profFile, err)
    }
    if err := pprof.WriteHeapProfile(file); err != nil {
        fmt.Printf("error in WriteHeapProfile: %s", err)
    }
}

// ################################################################################
// ########################### debugHistogram struct ##############################
// ################################################################################
type debugHistogram struct {
    mapping map[string]int
}

// ######################## debugHistogram methods ####################

func (d *debugHistogram) Print() {
    if !printDebugOutput {
        return
    }
    if len(d.mapping) == 0 {
        return
    }
    fmt.Fprintf(os.Stderr, "Debug histogram\n")
    for k, v := range d.mapping {
        fmt.Fprintf(os.Stderr, "%s: %d\n", k, v)
    }
}

func (d *debugHistogram) PrintSorted() {
    if !printDebugOutput {
        return
    }
    if len(d.mapping) == 0 {
        return
    }
    var scores vector.IntVector
    inverseMapping := make(map[int]string)
    sum := 0
    for key, value := range d.mapping {
        inverseMapping[value] = key
        scores.Push(value)
        sum += value
    }
    scoresArray := sort.IntArray(scores)
    scoresArray.Sort()

    fmt.Fprintf(os.Stderr, "Debug histogram - sorted:\n")
    length := len(scoresArray)
    for i := 0; i < length; i++ {
        current := scoresArray[length - i-1]
        fmt.Fprintf(os.Stderr, "%s: %d (%2.1f%%)\n", inverseMapping[current], current, float(current)/float(sum)*100)
    }
}

func (d *debugHistogram) Score() {
    _, callerFilePath, callerLine, _ := runtime.Caller(1)
    splitPath := strings.Split(callerFilePath, "/", -1)
    callerFile := splitPath[len(splitPath) - 1]
    str := fmt.Sprintf("%s:%d", callerFile, callerLine)
    if v, ok := d.mapping[str]; !ok {
        d.mapping[str] = 1
    } else {
        d.mapping[str] = v+1
    }
}

func (d *debugHistogram) ScoreTagged(tag string) {
    _, callerFilePath, callerLine, _ := runtime.Caller(1)
    splitPath := strings.Split(callerFilePath, "/", -1)
    callerFile := splitPath[len(splitPath) - 1]
    str := fmt.Sprintf("%s:%d '%s'", callerFile, callerLine, tag)
    if v, ok := d.mapping[str]; !ok {
        d.mapping[str] = 1
    } else {
        d.mapping[str] = v+1
    }
}

// ######################## debugHistogram helpers ####################
func newDebugHistogram() *debugHistogram {
    return &debugHistogram{ mapping: make(map[string]int),
                          }
}


// ######################## debugHistogram helpers ####################
var DbgHistogram = newDebugHistogram()

