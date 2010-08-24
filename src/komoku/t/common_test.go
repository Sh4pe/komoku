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

func dummyTest(t *testing.T) {
}

func Testsuite() []testing.Test {
    return []testing.Test { testing.Test{"dummyTest", dummyTest},
                          }
}

