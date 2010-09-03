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

func Benchmark10Appends(b *testing.B) {
    for i := 0; i < b.N; i++ {
        il := NewIntList()
        for l := 0; l < 10; l++ {
            il.Append(l)
        }
    }
}

func Benchmark10Removes(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        il := NewIntList()
        for l := 0; l < 10; l++ {
            il.Append(l)
        }
        b.StartTimer()
        for l := 0; l < 10; l++ {
            il.Remove(l)
        }
    }
}

func Benchmark10AppendUniques(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        il := NewIntList()
        for l := 0; l < 10; l++ {
            il.Append(l)
        }
        b.StartTimer()
        for l := 5; l < 15; l++ {
            il.AppendUnique(l)
        }
    }
}

func BenchmarkJoinUnique(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        il1 := NewIntList()
        il2 := NewIntList()
        for k := 0; k < 100; k++ {
            il1.Append(k)
        }
        for k := 50; k < 150; k++ {
            il2.Append(k)
        }
        b.StartTimer()
        il1.JoinUnique(il2)
    }
}

func Benchmarks() []testing.Benchmark {
    return []testing.Benchmark { testing.Benchmark{"BenchmarkIntListCreation", BenchmarkIntListCreation},
                                 testing.Benchmark{"Benchmark10Appends", Benchmark10Appends},
                                 testing.Benchmark{"Benchmark10Removes", Benchmark10Removes},
                                 testing.Benchmark{"Benchmark10AppendUniques", Benchmark10AppendUniques},
                                 testing.Benchmark{"BenchmarkJoinUnique", BenchmarkJoinUnique},
                               }
}
