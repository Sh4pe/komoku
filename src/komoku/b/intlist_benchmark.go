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

func BenchmarkIntListCreation(b *testing.B) {
    for i := 0; i < b.N; i++ {
        il := NewIntList()
        il.Append(3)
    }
}

func Benchmarks() []testing.Benchmark {
    return []testing.Benchmark { testing.Benchmark{"BenchmarkIntListCreation",BenchmarkIntListCreation},
                               }
}
