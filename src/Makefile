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
LINK=6l

COMPILEFIRST  = $(CC) -o $(@) $(<)
COMPILEFIRST_QUIET = @echo '  6g $(@)'; $(COMPILE)

OBJS= main main.6

#########################################################################################
####################### build rules #####################################################
#########################################################################################
all: $(TARG).6
	@echo '  6l $(TARG)'; 6l -o $(TARG) $(TARG).6
	
$(TARG).6: main.go komoku/komoku.a
	@echo '  6g $(@)'; $(CC) -o $(@) main.go

komoku/komoku.a: komoku/common.go komoku/board.go komoku/treenode.go komoku/ui.go
	@echo 'Building in subdirectory komoku...'
	@cd komoku && make
	@echo 'Leaving komoku...'

tests_run:
	@cd komoku && make tests_run

benchmarks_run:
	@cd komoku && make benchmarks_run

clean:
	@for obj in $(OBJS); do \
		echo '  rm '"$$obj"; \
		rm -f "$$obj"; \
	done;
	@echo 
	@echo 'Cleaning subdirectory komoku...'
	@cd komoku && make clean
