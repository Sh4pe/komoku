/* 
 * (c) 2010 by David Nies (nies.david@googlemail.com)
 *     http://www.twitter.com/Sh4pe
 *
 * Use of this source code is governed by a license 
 * that can be found in the LICENSE file.
 */

/*
 * This file defines the Environment struct. In this struct all the information for
 * running an instance of a komoku program is stored
 */

package komoku

// ################################################################################
// ########################### Environment struct #################################
// ################################################################################
type Environment struct {
    *Game
    komi float
}

// ##################### Environment methods ##########################

func (e *Environment) SetKomi(newKomi float) {
    e.komi = newKomi
}

// ##################### Environment helper functions ##########################

func NewEnvironment(boardsize int) *Environment {
    return &Environment{
        Game: NewGame(boardsize),
        komi: DefaultKomi,
    }
}


