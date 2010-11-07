# vim: set noexpandtab :

#
# (c) 2010 by David Nies (nies.david@googlemail.com)
#     http://www.twitter.com/Sh4pe
#
# Use of this source code is governed by a license 
# that can be found in the LICENSE file.
#

# These are settings shared by all makefiles of this project

# Finds out if we have to use 6g or 8g, 5g is not supported
OBJSUFF=$(shell if [ -x "`which 6g`" ]; then echo 6; else echo 8; fi )
CCSTR=$(OBJSUFF)g
CC=$(OBJSUFF)g
LINKSTR=$(OBJSUFF)l
LINK=$(OBJSUFF)l -e
PACK=gopack
PROFILER=sudo 6prof
GOPPROF=gopprof
