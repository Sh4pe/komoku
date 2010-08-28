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

func BenchmarkCreationPointerAssign(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var il *IntList
        // Does this copy? Thats the question of this benchmark...
        il = NewIntList()
        // Dummy op to avoid "declared and not used"-error
        il.Length()
    }
}

func BenchmarkCreationValueAssign(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var il IntList
        // Does this copy? Thats the question of this benchmark...
        il = *NewIntList()
        // dummy op to avoid "declared and not used"-error
        il.Length()
    }
}

func BenchmarkIntAssignment(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var a int = 100
        var b int = 200

        b = a
        b++
    }
}

func BenchmarkIntCast(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var a int = 100
        var b uint32 = 200

        b = uint32(a)
        b++
    }
}

func BenchmarkIntListIteratorLoop(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        il := NewIntList()
        for k := 0; k < 10; k++ {
            il.Append(k)
        }
        b.StartTimer()
        last := il.Last()
        for it := il.First(); it != last; it = it.Next() {
            v := it.Value()
            v++
        }
    }
}

func BenchmarkIntListDoLoop(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        il := NewIntList()
        for k := 0; k < 10; k++ {
            il.Append(k)
        }
        b.StartTimer()
        eachFunc := func(val int) {
            local := val
            local++
        }
        il.Do(eachFunc)
    }
}

func Benchmarks() []testing.Benchmark {
    return []testing.Benchmark { testing.Benchmark{"BenchmarkGenericVector", BenchmarkGenericVector},
                                 testing.Benchmark{"BenchmarkIntVector", BenchmarkIntVector},
                                 testing.Benchmark{"BenchmarkCreationPointerAssign", BenchmarkCreationPointerAssign},
                                 testing.Benchmark{"BenchmarkCreationValueAssign", BenchmarkCreationValueAssign},
                                 testing.Benchmark{"BenchmarkIntAssignment", BenchmarkIntAssignment},
                                 testing.Benchmark{"BenchmarkIntCast", BenchmarkIntCast},
                                 testing.Benchmark{"BenchmarkIntListIteratorLoop", BenchmarkIntListIteratorLoop},
                                 testing.Benchmark{"BenchmarkIntListDoLoop", BenchmarkIntListDoLoop},
                               }
}
