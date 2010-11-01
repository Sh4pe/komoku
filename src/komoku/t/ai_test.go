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
)

func TestRunSimulation(t *testing.T) {
    numTestSimulations := 500
    ai := NewAI(9)

    for i := 0; i < numTestSimulations; i++ {
        ai.runSimulation()
    }
    countedSimuls := 0
    for _, node := range ai.topNode.children {
        countedSimuls += node.NodeInfo.simulations
    }

    if countedSimuls != numTestSimulations {
        t.Fatalf("Wrong number of simulations, expected %d, got %d", numTestSimulations, countedSimuls)
    }
    if ai.topNode.simulations != numTestSimulations {
        t.Fatalf("AI.topNode has a wrong number of simulations, expected %d, got %d", numTestSimulations, ai.topNode.simulations)
    }

}

func Testsuite() []testing.Test {
    return []testing.Test {
        testing.Test{"TestRunSimulation", TestRunSimulation},
    }
}
