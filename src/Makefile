# vim: set noexpandtab :

#
# (c) 2010 by David Nies (nies.david@googlemail.com)
#     http://www.twitter.com/Sh4pe
#
# Use of this source code is governed by a license 
# that can be found in the LICENSE file.
#

#########################################################################################
####################### most common defintions ##########################################
#########################################################################################
CC=6g
TARG=main
LINK=6l -e

COMPILEFIRST  = $(CC) -o $(@) $(<)
COMPILEFIRST_QUIET = @echo '  6g $(@)'; $(COMPILE)

OBJS= main main.6

ALLSOURCE  = ai.go
ALLSOURCE += board.go
ALLSOURCE += common.go
ALLSOURCE += debug.go
ALLSOURCE += environment.go
ALLSOURCE += game.go
ALLSOURCE += group.go 
ALLSOURCE += gtp.go 
ALLSOURCE += gtpcmd.go 
ALLSOURCE += intlist.go 
ALLSOURCE += treenode.go 
ALLSOURCE += ui.go 

KOMOKUDEPS = $(patsubst %,komoku/%, $(ALLSOURCE))

#########################################################################################
####################### build rules #####################################################
#########################################################################################
all: $(TARG).6
	@echo '  6l $(TARG)'; 6l -o $(TARG) $(TARG).6

$(TARG).6: main.go komoku/komoku.a
	@echo '  6g $(@)'; $(CC) -o $(@) main.go

komoku/komoku.a: $(KOMOKUDEPS)
	@echo 'Building in subdirectory komoku...'
	@cd komoku && $(MAKE)
	@echo 'Leaving komoku...'

tests_run:
	@cd komoku && $(MAKE) tests_run

benchmarks_run:
	@cd komoku && $(MAKE) benchmarks_run

clean:
	@for obj in $(OBJS); do \
		echo '  rm '"$$obj"; \
		rm -f "$$obj"; \
	done;
	@echo 
	@echo 'Cleaning subdirectory komoku...'
	@cd komoku && $(MAKE) clean
