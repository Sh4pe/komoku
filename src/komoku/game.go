/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */

/*
 * Here the Game struct is defined. It stores all the information which komoku associates
 * with one game
 */


package komoku

// ################################################################################
// ########################### Game struct ########################################
// ################################################################################
type Game struct {
    B *Board // The current board
    Komi float
}


