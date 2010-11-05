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

func BenchmarkRunSimulation9(b *testing.B) {
    b.StopTimer()
    ai := NewAI(9)
    b.StartTimer()
    for i := 0; i < b.N; i++ {
        ai.runSimulation()
    }
}

func BenchmarkRunSimulation19(b *testing.B) {
    b.StopTimer()
    ai := NewAI(19)
    b.StartTimer()
    for i := 0; i < b.N; i++ {
        ai.runSimulation()
    }
}

func Benchmarks() []testing.InternalBenchmark {
    return []testing.InternalBenchmark {
        testing.InternalBenchmark{"BenchmarkRunSimulation9", BenchmarkRunSimulation9},
        testing.InternalBenchmark{"BenchmarkRunSimulation19", BenchmarkRunSimulation19},
    }
}
