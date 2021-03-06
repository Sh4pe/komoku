# vim: set noexpandtab :

#
# (c) 2010 by David Nies (nies.david@googlemail.com)
#     http://www.twitter.com/Sh4pe
#
# Use of this source code is governed by a license 
# that can be found in the LICENSE file.
#

include ../common.mk
####################### most common defintions ##########################################
TARG=komoku

# prefix of this directory
THISDIR=komoku/

KOMOKULIB=$(TARG).a

COMPILE  = $(CC) -o $(@) $(^)
COMPILE_QUIET = @echo '  $(CCSTR) $(THISDIR)$(@)'; $(COMPILE)

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

#########################################################################################
############### Stuff needed for generating test executables ############################
#########################################################################################
TESTDIR = t/

### the command for compiling tests
#
TESTCOMPILE  = $(CC) -o $(@).$(OBJSUFF) $(^); # compile all dependencies to $(@).$(OBJSUFF)
# create tmp file with adjusted testrunner code in it
TESTCOMPILE += cat $(TESTDIR)/skeleton.go | sed 's/SUBSTITUTE_THIS/$(notdir $(@))/g' > $(TESTDIR)$(notdir $(@)).tmp.go; 
TESTCOMPILE += $(CC) -I$(TESTDIR) -o $(@).tmp.$(OBJSUFF) $(@).tmp.go; # compiles the generated program with all dependencies as .go-file arguments
TESTCOMPILE += $(LINK) -L$(TESTDIR) -o $(@) $(@).tmp.$(OBJSUFF); # link the program to generate an executable
TESTCOMPILE += rm -f $(@).tmp.$(OBJSUFF); # remove the .$(OBJSUFF) file for the temorarily created file
TESTCOMPILE += rm -f $(@).tmp.go; # and finally remove the temorarily created file

# the command for doing this quietly with a nice output
TESTCOMPILE_QUIET = @echo '  $(LINKSTR) $(THISDIR)$(@)'; $(TESTCOMPILE)

ALLTESTS_TARGS = ai_test common_test group_test gtp_test intlist_test ui_test board_test
ALLTESTS = $(patsubst %,$(TESTDIR)%,$(ALLTESTS_TARGS))


TESTOBJS += board_test
TESTOBJS += board_test.$(OBJSUFF)
TESTOBJS += common_test
TESTOBJS += common_test.$(OBJSUFF) 
TESTOBJS += gtp_test
TESTOBJS += gtp_test.$(OBJSUFF) 
TESTOBJS += intlist_test
TESTOBJS += intlist_test.$(OBJSUFF)
TESTOBJS += ui_test
TESTOBJS += ui_test.$(OBJSUFF)
TESTOBJS += group_test
TESTOBJS += group_test.$(OBJSUFF)
TESTOBJS += ai_test
TESTOBJS += ai_test.$(OBJSUFF)

#########################################################################################
############### Stuff needed for generating benchmark executables #######################
#########################################################################################
BENCHMARKDIR = b/
PROFDIR = ../../data/prof/

# The command for compiling benchmarks. Mostly analogous to $(TESTCOMPILE)
BENCHMARKCOMPILE  = $(CC) -o $(@).$(OBJSUFF) $(^);
BENCHMARKCOMPILE += cat $(BENCHMARKDIR)/skeleton.go | sed 's/SUBSTITUTE_THIS/$(notdir $(@))/g' > $(BENCHMARKDIR)$(notdir $(@)).tmp.go; 
BENCHMARKCOMPILE += $(CC) -I$(BENCHMARKDIR) -o $(@).tmp.$(OBJSUFF) $(@).tmp.go;
BENCHMARKCOMPILE += $(LINK) -L$(BENCHMARKDIR) -o $(@) $(@).tmp.$(OBJSUFF);
BENCHMARKCOMPILE += rm -f $(@).tmp.$(OBJSUFF);
BENCHMARKCOMPILE += rm -f $(@).tmp.go;

BENCHMARKCOMPILE_QUIET = @echo '  $(LINKSTR) $(THISDIR)$(@)'; $(BENCHMARKCOMPILE)

# This is the standard regexp supplied to -benchmarks= in the benchmark executables
BENCHMARKREGEX = [a-zA-Z]*

# Command for running benchmarks
BENCHMARKRUN  = @echo Running $(patsubst %_run,%, $(@));
BENCHMARKRUN += ./$(patsubst %_run,%, $(@)) -benchmarks='$(BENCHMARKREGEX)';
# Commands for profiling benchmarks
BENCHMARKPROFILE  = $(eval proffile=$(PROFDIR)$(patsubst $(BENCHMARKDIR)%_profile,%,$(@)).prof)
BENCHMARKPROFILE += $(eval profexec=./$(patsubst %_profile,%, $(@)))
BENCHMARKPROFILE += @echo Profiling $(profexec);
BENCHMARKPROFILE += echo "OSX needs root rights for this";
BENCHMARKPROFILE += $(PROFILER) -d 10 -P $(proffile) $(profexec) -benchmarks='$(BENCHMARKREGEX)';
BENCHMARKPROFILE += $(GOPPROF) $(profexec) $(proffile);

# This comand runs only one benchmark in the profiler
BENCHMARKPROFILEONLY  = $(eval benchname=$(lastword $(subst _, ,$(@))))
BENCHMARKPROFILEONLY += $(eval execname=$(patsubst $(BENCHMARKDIR)%,%,$(patsubst %_profile_,%,$(subst $(benchname),,$(@)))))
BENCHMARKPROFILEONLY += $(eval proffile=$(PROFDIR)$(execname).prof)
BENCHMARKPROFILEONLY += $(eval profexec=$(patsubst %_profile_,%,$(subst $(benchname),,$(@))))
BENCHMARKPROFILEONLY += @echo Profiling benchmark \'$(benchname)\' in \'$(execname)\'...;
BENCHMARKPROFILEONLY += echo "OSX needs root rights for this";
BENCHMARKPROFILEONLY += $(PROFILER) -d 10 -P $(proffile) $(profexec) -benchmarks='$(benchname)';
BENCHMARKPROFILEONLY += $(GOPPROF) $(profexec) $(proffile);

ALLBENCHMARKS_TARGS = ai_benchmark design_decision_benchmark board_benchmark intlist_benchmark
ALLBENCHMARKS = $(patsubst %,$(BENCHMARKDIR)%,$(ALLBENCHMARKS_TARGS))


BENCHMARKOBJS += ai_benchmark
BENCHMARKOBJS += ai_benchmark.$(OBJSUFF)
BENCHMARKOBJS += board_benchmark
BENCHMARKOBJS += board_benchmark.$(OBJSUFF)
BENCHMARKOBJS += design_decision_benchmark
BENCHMARKOBJS += design_decision_benchmark.$(OBJSUFF)
BENCHMARKOBJS += intlist_benchmark
BENCHMARKOBJS += intlist_benchmark.$(OBJSUFF)

#########################################################################################
############### Stuff needed for generating experiment executables ######################
#########################################################################################
EXPERIMENTDIR = e/

EXPERIMENTCOMPILE += $(eval sourcefile=$(wordlist 2, $(words $(^)), $(^)))
EXPERIMENTCOMPILE += $(CC) -o $(@).$(OBJSUFF) $(sourcefile);
EXPERIMENTCOMPILE += $(LINK) -o $(@) $(@).$(OBJSUFF)

EXPERIMENTCOMPILE_QUIET = @echo '  $(LINKSTR) $(@)'; $(EXPERIMENTCOMPILE)

EXPERIMENTRUN  = @echo Performing the experiment $(patsubst %_run,%, $(@))...;
EXPERIMENTRUN += ./$(patsubst %_run,%, $(@));

EXPERIMENTOBJS += gamelength
EXPERIMENTOBJS += gamelength.$(OBJSUFF)


#########################################################################################
####################### global vars #####################################################
#########################################################################################

OBJS  = $(TARG).$(OBJSUFF) $(TARG).a $(patsubst %,$(TESTDIR)%,$(TESTOBJS)) 
OBJS += $(patsubst %,$(BENCHMARKDIR)%,$(BENCHMARKOBJS))
OBJS += $(patsubst %,$(EXPERIMENTDIR)%,$(EXPERIMENTOBJS))
OBJS += $(PROFDIR)*.prof

#########################################################################################
####################### build rules #####################################################
#########################################################################################
all: $(KOMOKULIB)

$(KOMOKULIB): $(TARG).$(OBJSUFF)
	@rm -f $(TARG).a
	@echo '  $(LINKSTR) $(THISDIR)$(TARG).a'
	@$(PACK) grc $(TARG).a $(TARG).$(OBJSUFF)
	
$(TARG).$(OBJSUFF): $(ALLSOURCE)
	$(COMPILE_QUIET)

#################### tests ################

$(TESTDIR)ai_test: $(TESTDIR)ai_test.go ai.go board.go common.go environment.go game.go group.go intlist.go treenode.go ui.go
	$(TESTCOMPILE_QUIET)

$(TESTDIR)board_test: $(TESTDIR)board_test.go board.go common.go debug.go game.go group.go intlist.go ui.go
	$(TESTCOMPILE_QUIET)

$(TESTDIR)common_test: $(TESTDIR)common_test.go common.go
	$(TESTCOMPILE_QUIET)

$(TESTDIR)gtp_test: $(TESTDIR)gtp_test.go ai.go board.go common.go debug.go environment.go game.go gtp.go gtpcmd.go group.go intlist.go ui.go treenode.go
	$(TESTCOMPILE_QUIET)

$(TESTDIR)group_test: $(TESTDIR)group_test.go common.go group.go intlist.go
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
$(BENCHMARKDIR)design_decision_benchmark: $(BENCHMARKDIR)design_decision_benchmark.go common.go debug.go intlist.go
	$(BENCHMARKCOMPILE_QUIET)

.PHONY: $(BENCHMARKDIR)design_decision_benchmark_run
$(BENCHMARKDIR)design_decision_benchmark_run: $(BENCHMARKDIR)design_decision_benchmark
	$(BENCHMARKRUN)

.PHONY: $(BENCHMARKDIR)design_decision_benchmark_profile_GenericVector
$(BENCHMARKDIR)design_decision_benchmark_profile_GenericVector: $(BENCHMARKDIR)design_decision_benchmark
	$(BENCHMARKPROFILEONLY)

$(BENCHMARKDIR)board_benchmark: $(BENCHMARKDIR)board_benchmark.go board.go common.go debug.go group.go intlist.go
	$(BENCHMARKCOMPILE_QUIET)

.PHONY: $(BENCHMARKDIR)board_benchmark_run
$(BENCHMARKDIR)board_benchmark_run: $(BENCHMARKDIR)board_benchmark
	$(BENCHMARKRUN)

.PHONY: $(BENCHMARKDIR)board_benchmark_profile
$(BENCHMARKDIR)board_benchmark_profile: $(BENCHMARKDIR)board_benchmark
	$(BENCHMARKPROFILE)

$(BENCHMARKDIR)intlist_benchmark: $(BENCHMARKDIR)intlist_benchmark.go intlist.go
	$(BENCHMARKCOMPILE_QUIET)

.PHONY: $(BENCHMARKDIR)intlist_benchmark_run
$(BENCHMARKDIR)intlist_benchmark_run: $(BENCHMARKDIR)intlist_benchmark
	$(BENCHMARKRUN)

$(BENCHMARKDIR)ai_benchmark: $(BENCHMARKDIR)ai_benchmark.go ai.go board.go common.go environment.go game.go group.go intlist.go treenode.go ui.go
	$(BENCHMARKCOMPILE_QUIET)

.PHONY: $(BENCHMARKDIR)ai_benchmark_run
$(BENCHMARKDIR)ai_benchmark_run: $(BENCHMARKDIR)ai_benchmark
	$(BENCHMARKRUN)

.PHONY: benchmarks_compile
benchmarks_compile: $(ALLBENCHMARKS)
	
benchmarks_run: $(KOMOKULIB) $(ALLBENCHMARKS)
	@for targ in $(^); do \
		 if test x"$$targ" != x"$(KOMOKULIB)"; then \
		 	$(MAKE) "$$targ"_run; \
		 fi \
	done;

#################### experiments ################
$(EXPERIMENTDIR)gamelength: $(KOMOKULIB) $(EXPERIMENTDIR)gamelength.go
	$(EXPERIMENTCOMPILE_QUIET)

.PHONY: $(EXPERIMENTDIR)gamelength_run
$(EXPERIMENTDIR)gamelength_run: $(EXPERIMENTDIR)gamelength
	$(EXPERIMENTRUN)

#################### pure phony targets ################
.PHONY: clean
clean:
	@for obj in $(OBJS); do \
		echo '  rm '"$(THISDIR)$$obj"; \
		rm -f "$$obj"; \
	done;
