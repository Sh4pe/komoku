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
    "container/vector"
)

const (
    vectorLength = 100;
)

func BenchmarkGenericVector(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var vec vector.Vector
        for j := 0; j < vectorLength; j++ {
            vec.Push(j)
        }
        for j := 0; j < vectorLength; j++ {
            val, _ := vec.At(j).(int)
            val++
        }
    }
}

func BenchmarkIntVector(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var vec vector.IntVector
        for j := 0; j < vectorLength; j++ {
            vec.Push(j)
        }
        for j := 0; j < vectorLength; j++ {
            val := vec.At(j)
            val++
        }
    }
}


func Benchmarks() []testing.Benchmark {
    return []testing.Benchmark { testing.Benchmark{"BenchmarkGenericVector", BenchmarkGenericVector},
                                 testing.Benchmark{"BenchmarkIntVector", BenchmarkIntVector},
                               }
}
