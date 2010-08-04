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
TARG=komoku
PACK=gopack
LINK=6l

KOMOKULIB=$(TARG).a

COMPILE  = $(CC) -o $(@) $(^)
COMPILE_QUIET = @echo '  6g $(@)'; $(COMPILE)

OBJS = $(TARG).6 $(TARG).a common_test common_test.6 board_test board_test.6 ui_test.6 ui_test \
	design_decision_benchmark design_decision_benchmark.6

ALLTESTS = common_test board_test ui_test intlist_test

#########################################################################################
############### Stuff needed for generating test executables ############################
#########################################################################################

# _ALL is here just for documentation
TESTCOMPILEPROG_ALL ='package main; import ( "./blah_blah"; "testing"; ); func main() { testing.Main(komoku.Testsuite()) }'
TESTCOMPILEPROG_PRE	='package main; import ( "./'
TESTCOMPILEPROG_POST='"; "testing"; ); func main() { testing.Main(komoku.Testsuite()) }'

# the command for compiling tests
TESTCOMPILE  = $(CC) -o $(@).6 $(^); # compile all dependencies to $(@).6
TESTCOMPILE += echo $(TESTCOMPILEPROG_PRE)$(@)$(TESTCOMPILEPROG_POST) > $(@).tmp.go; # create tmp file with adjusted testrunner code in it
TESTCOMPILE += $(CC) -o $(@).tmp.6 $(@).tmp.go; # compiles the generated program with all dependencies as .go-file arguments
TESTCOMPILE += $(LINK) -o $(@) $(@).tmp.6; # link the program to generate an executable
TESTCOMPILE += rm -f $(@).tmp.6; # remove the .6 file for the temorarily created file
TESTCOMPILE += rm -f $(@).tmp.go; # and finally remove the temorarily created file

# the command for doing this quietly with a nice output
TESTCOMPILE_QUIET = @echo '  6l $(@)'; $(TESTCOMPILE)

# targets look like 'foo_test_run', and sed chops of the trailing '_foo'
TESTRUN  = @echo Running `echo $(@) | sed 's/_[a-zA-Z]*$$//g'`...;
TESTRUN += ./`echo $(@) | sed 's/_[a-zA-Z]*$$//g'`;

#########################################################################################
############### Stuff needed for generating benchmark executables #######################
#########################################################################################

BENCHMARKCOMPILEPROG_PRE   = 'package main; import ( "testing"; "flag"; "fmt"; "os"; "./'
BENCHMARKCOMPILEPROG_POST  = '"; ); func main() { flag.Parse(); b := komoku.Benchmarks(); if len(os.Args) == 1 { bn := "";'
BENCHMARKCOMPILEPROG_POST += 'for _, n := range b { bn += n.Name + " "; }; fmt.Println("benchmarks:"); fmt.Println(bn); return; }; '
BENCHMARKCOMPILEPROG_POST += 'testing.RunBenchmarks(b); }'

# The command for compiling benchmarks. Mostly analogous to $(TESTCOMPILE)
BENCHMARKCOMPILE  = $(CC) -o $(@).6 $(^);
BENCHMARKCOMPILE += echo $(BENCHMARKCOMPILEPROG_PRE)$(@)$(BENCHMARKCOMPILEPROG_POST) > $(@).tmp.go;
BENCHMARKCOMPILE += $(CC) -o $(@).tmp.6 $(@).tmp.go;
BENCHMARKCOMPILE += $(LINK) -o $(@) $(@).tmp.6;
BENCHMARKCOMPILE += rm -f $(@).tmp.6;
BENCHMARKCOMPILE += rm -f $(@).tmp.go;

BENCHMARKCOMPILE_QUIET = @echo '  6l $(@)'; $(BENCHMARKCOMPILE)

# This is the standard regexp supplied to -benchmarks= in the benchmark executables
BENCHMARKREGEX = [a-zA-Z]*

# Command for running benchmarks
BENCHMARKRUN  = @echo Running `echo $(@) | sed 's/_[a-zA-Z]*$$//g'`...;
BENCHMARKRUN += ./`echo $(@) | sed 's/_[a-zA-Z]*$$//g'` -benchmarks='$(BENCHMARKREGEX)';

#########################################################################################
####################### build rules #####################################################
#########################################################################################
all: $(KOMOKULIB)

$(KOMOKULIB): $(TARG).6
	@rm -f $(TARG).a
	@echo '  6l $(TARG).a'
	@$(PACK) grc $(TARG).a $(TARG).6
	
$(TARG).6: common.go board.go treenode.go ui.go intlist.go
	$(COMPILE_QUIET)


#################### tests ################

board_test: board_test.go common.go board.go intlist.go
	$(TESTCOMPILE_QUIET)

board_test_run: board_test
	$(TESTRUN)

common_test: common_test.go common.go
	$(TESTCOMPILE_QUIET)

common_test_run: common_test
	$(TESTRUN)

intlist_test: intlist_test.go intlist.go common.go
	$(TESTCOMPILE_QUIET)

intlist_test_run: intlist_test
	$(TESTRUN)

ui_test: ui_test.go common.go board.go intlist.go ui.go 
	$(TESTCOMPILE_QUIET)

ui_test_run: ui_test
	$(TESTRUN)

tests_compile: $(ALLTESTS)

tests_run: $(KOMOKULIB) $(ALLTESTS)
	@for targ in $(^); do \
		 if test x"$$targ" != x"$(KOMOKULIB)"; then \
		 	make "$$targ"_run; \
		 fi \
	done;


#################### benchmarks ################

design_decision_benchmark: design_decision_benchmark.go
	$(BENCHMARKCOMPILE_QUIET)

design_decision_benchmark_run: design_decision_benchmark
	$(BENCHMARKRUN)
	
benchmarks_run: $(KOMOKULIB) design_decision_benchmark
	@for targ in $(^); do \
		 if test x"$$targ" != x"$(KOMOKULIB)"; then \
		 	make "$$targ"_run; \
		 fi \
	done;

#################### maintenance ################
clean:
	@for obj in $(OBJS); do \
		echo '  rm '"$$obj"; \
		rm -f "$$obj"; \
	done;
