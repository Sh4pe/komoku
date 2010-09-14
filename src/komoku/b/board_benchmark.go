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
    "os"
    "rand"
)

const boardsize = 9


func BenchmarkRandomGameByListLegalPoints(b *testing.B) {
    b.StopTimer()
    board := NewBoard(boardsize)
    b.StartTimer()
    for i := 0; i < b.N; i++ {
        legalMoves := board.ListLegalPoints(board.ColorOfNextPlay())
        if len(legalMoves) == 0 {
            b.StopTimer()
            board = NewBoard(boardsize)
            b.StartTimer()
        } else {
            sec, nsec, _ := os.Time()
            random := rand.New(rand.NewSource(sec+nsec))
            randomMove := legalMoves[random.Intn(len(legalMoves))]
            board.TurnPlayMove(randomMove.X, randomMove.Y)
        }
    }
    DbgHistogram.PrintSorted()
}

func BenchmarkRandomGameByPlayRandomMove(b *testing.B) {
    b.StopTimer()
    board := NewBoard(boardsize)
    var color Color = Black
    b.StartTimer()
    for i := 0; i < b.N; i++ {
        v := board.PlayRandomMove(color)
        if v.Pass {
            b.StopTimer()
            board = NewBoard(boardsize)
            b.StartTimer()
        }
        color = !color
    }
    DbgHistogram.PrintSorted()
}

func Benchmarks() []testing.Benchmark {
    return []testing.Benchmark { testing.Benchmark{"BenchmarkRandomGameByListLegalPoints", BenchmarkRandomGameByListLegalPoints},
                                 testing.Benchmark{"BenchmarkRandomGameByPlayRandomMove", BenchmarkRandomGameByPlayRandomMove},
                               }
}
