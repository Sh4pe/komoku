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


func BenchmarkRandomGame(b *testing.B) {
    b.StopTimer()
    boardSize := 19
    board := NewBoard(boardSize)
    b.StartTimer()
    for i := 0; i < b.N; i++ {
        legalMoves := board.ListLegalPoints(board.ColorOfNextPlay())
        if len(legalMoves) == 0 {
            b.StopTimer()
            board = NewBoard(boardSize)
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

func Benchmarks() []testing.Benchmark {
    return []testing.Benchmark { testing.Benchmark{"BenchmarkRandomGame", BenchmarkRandomGame},
                               }
}
