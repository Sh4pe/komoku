# vim: set noexpandtab :

#
# (c) 2010 by David Nies (nies.david@googlemail.com)
#     http://www.twitter.com/Sh4pe
#
# Use of this source code is governed by a license 
# that can be found in the LICENSE file.
#

# TODO:
#	

#########################################################################################
####################### most common defintions ##########################################
#########################################################################################
CC=6g
TARG=komoku
PACK=gopack
LINK=6l

KOMOKULIB=$(TARG).a

COMPILE  = $(CC) -o $(@) $(^)
COMPILE_QUIET = @echo '  6g $(@)'; $(COMPILE)

ALLSOURCE  = board.go
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

#########################################################################################
############### Stuff needed for generating test executables ############################
#########################################################################################
TESTDIR = t/

### the command for compiling tests
#
TESTCOMPILE  = $(CC) -o $(@).6 $(^); # compile all dependencies to $(@).6
# create tmp file with adjusted testrunner code in it
TESTCOMPILE += cat $(TESTDIR)/skeleton.go | sed 's/SUBSTITUTE_THIS/$(notdir $(@))/g' > $(TESTDIR)$(notdir $(@)).tmp.go; 
TESTCOMPILE += $(CC) -I$(TESTDIR) -o $(@).tmp.6 $(@).tmp.go; # compiles the generated program with all dependencies as .go-file arguments
TESTCOMPILE += $(LINK) -L$(TESTDIR) -o $(@) $(@).tmp.6; # link the program to generate an executable
TESTCOMPILE += rm -f $(@).tmp.6; # remove the .6 file for the temorarily created file
TESTCOMPILE += rm -f $(@).tmp.go; # and finally remove the temorarily created file

# the command for doing this quietly with a nice output
TESTCOMPILE_QUIET = @echo '  6l $(@)'; $(TESTCOMPILE)

ALLTESTS_TARGS = board_test common_test gtp_test intlist_test ui_test 
ALLTESTS = $(patsubst %,$(TESTDIR)%,$(ALLTESTS_TARGS))


TESTOBJS += board_test
TESTOBJS += board_test.6
TESTOBJS += common_test
TESTOBJS += common_test.6 
TESTOBJS += gtp_test
TESTOBJS += gtp_test.6 
TESTOBJS += intlist_test
TESTOBJS += intlist_test.6
TESTOBJS += ui_test
TESTOBJS += ui_test.6

#########################################################################################
############### Stuff needed for generating benchmark executables #######################
#########################################################################################
BENCHMARKDIR = b/

# The command for compiling benchmarks. Mostly analogous to $(TESTCOMPILE)
BENCHMARKCOMPILE  = $(CC) -o $(@).6 $(^);
BENCHMARKCOMPILE += cat $(BENCHMARKDIR)/skeleton.go | sed 's/SUBSTITUTE_THIS/$(notdir $(@))/g' > $(BENCHMARKDIR)$(notdir $(@)).tmp.go; 
BENCHMARKCOMPILE += $(CC) -I$(BENCHMARKDIR) -o $(@).tmp.6 $(@).tmp.go;
BENCHMARKCOMPILE += $(LINK) -L$(BENCHMARKDIR) -o $(@) $(@).tmp.6;
BENCHMARKCOMPILE += rm -f $(@).tmp.6;
BENCHMARKCOMPILE += rm -f $(@).tmp.go;

BENCHMARKCOMPILE_QUIET = @echo '  6l $(@)'; $(BENCHMARKCOMPILE)

# This is the standard regexp supplied to -benchmarks= in the benchmark executables
BENCHMARKREGEX = [a-zA-Z]*

# Command for running benchmarks
BENCHMARKRUN  = @echo Running $(patsubst %_run,%, $(@));
BENCHMARKRUN += ./$(patsubst %_run,%, $(@)) -benchmarks='$(BENCHMARKREGEX)';

ALLBENCHMARKS_TARGS = design_decision_benchmark
ALLBENCHMARKS = $(patsubst %,$(BENCHMARKDIR)%,$(ALLBENCHMARKS_TARGS))


BENCHMARKOBJS += design_decision_benchmark
BENCHMARKOBJS += design_decision_benchmark.6


#########################################################################################
####################### 'final' global vars #############################################
#########################################################################################

OBJS  = $(TARG).6 $(TARG).a $(patsubst %,$(TESTDIR)%,$(TESTOBJS)) 
OBJS += $(patsubst %,$(BENCHMARKDIR)%,$(BENCHMARKOBJS))

#########################################################################################
####################### build rules #####################################################
#########################################################################################
all: $(KOMOKULIB)

$(KOMOKULIB): $(TARG).6
	@rm -f $(TARG).a
	@echo '  6l $(TARG).a'
	@$(PACK) grc $(TARG).a $(TARG).6
	
$(TARG).6: $(ALLSOURCE)
	$(COMPILE_QUIET)

#################### tests ################

$(TESTDIR)board_test: $(TESTDIR)board_test.go board.go common.go debug.go group.go intlist.go ui.go
	$(TESTCOMPILE_QUIET)

$(TESTDIR)common_test: $(TESTDIR)common_test.go common.go
	$(TESTCOMPILE_QUIET)

$(TESTDIR)gtp_test: $(TESTDIR)gtp_test.go board.go common.go debug.go environment.go game.go gtp.go gtpcmd.go group.go intlist.go ui.go
	$(TESTCOMPILE_QUIET)

$(TESTDIR)intlist_test: $(TESTDIR)intlist_test.go common.go intlist.go 
	$(TESTCOMPILE_QUIET)

$(TESTDIR)ui_test: $(TESTDIR)ui_test.go board.go common.go debug.go group.go intlist.go ui.go 
	$(TESTCOMPILE_QUIET)

.PHONY: tests_compile
tests_compile: $(ALLTESTS)

tests_run: $(KOMOKULIB) $(ALLTESTS)
	@for targ in $(^); do \
		 if test x"$$targ" != x"$(KOMOKULIB)"; then \
			$$targ; \
		 fi \
	done;

#################### benchmarks ################
$(BENCHMARKDIR)design_decision_benchmark: $(BENCHMARKDIR)design_decision_benchmark.go intlist.go
	$(BENCHMARKCOMPILE_QUIET)

$(BENCHMARKDIR)design_decision_benchmark_run: $(BENCHMARKDIR)design_decision_benchmark
	$(BENCHMARKRUN)


.PHONY: benchmarks_compile
benchmarks_compile: $(ALLBENCHMARKS)
	
benchmarks_run: $(KOMOKULIB) $(ALLBENCHMARKS)
	@for targ in $(^); do \
		 if test x"$$targ" != x"$(KOMOKULIB)"; then \
		 	make "$$targ"_run; \
		 fi \
	done;

#################### phony targets ################
.PHONY: clean
clean:
	@for obj in $(OBJS); do \
		echo '  rm '"$$obj"; \
		rm -f "$$obj"; \
	done;
