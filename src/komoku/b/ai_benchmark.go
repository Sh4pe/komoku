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

func BenchmarkRunSimulation(b *testing.B) {
    b.StopTimer()
    ai := NewAI(9)
    b.StartTimer()
    for i := 0; i < b.N; i++ {
        ai.RunSimulation()
    }
}

func Benchmarks() []testing.Benchmark {
    return []testing.Benchmark {
        testing.Benchmark{"BenchmarkRunSimulation", BenchmarkRunSimulation},
    }
}
