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

// Tests if CharToDigit and DigitToChar are inverse to each other.
func TestCharDigitConversionInverse(t *testing.T) {
    l := len(coordinateChars)

    for i := 0; i < l; i++ {
        c := coordinateChars[i:i+1]
        digit, err := CharToDigit(c)
        if err != nil {
            t.Fatalf("CharToDigit returns an error: ``%s''", err)
        }
        if digit != i {
            t.Fatalf("expected %d, got %d", i, digit)
        }
    }

    for key, value := range charDigit {
        char, err := DigitToChar(value)
        if err != nil {
            t.Fatalf("DigitToChar returns an error: ``%s''", err)
        }
        if char != key {
            t.Fatalf("expected '%s', got '%s'", key, char)
        }
    }
}

func Testsuite() []testing.Test {
    return []testing.Test { testing.Test{"TestCharDigitConversionInverse", TestCharDigitConversionInverse},
                          }
}
