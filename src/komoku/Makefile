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

OBJS = $(TARG).6 $(TARG).a common_test common_test.6 board_test board_test.6 ui_test.6 ui_test

#########################################################################################
############### Stuff needed for generating test executables ############################
#########################################################################################

# _ALL is just there for 'documentaions sake'
TESTCOMPILEPROG_ALL='package main; import ( "./blah_blah"; "testing"; ); func main() { testing.Main(komoku.Testsuite()) }'
TESTCOMPILEPROG_PRE='package main; import ( "./'
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
TESTRUN = @echo Running `echo $(@) | sed 's/_[a-zA-Z]*$$//g'`...;
TESTRUN += ./`echo $(@) | sed 's/_[a-zA-Z]*$$//g'`

#########################################################################################
####################### build rules #####################################################
#########################################################################################
all: $(KOMOKULIB)

$(KOMOKULIB): $(TARG).6
	@rm -f $(TARG).a
	@echo '  6l $(TARG).a'
	@$(PACK) grc $(TARG).a $(TARG).6
	
$(TARG).6: common.go board.go treenode.go ui.go
	$(COMPILE_QUIET)

common_test: common_test.go common.go
	$(TESTCOMPILE_QUIET)

common_test_run: common_test
	$(TESTRUN)

board_test: board_test.go common.go board.go
	$(TESTCOMPILE_QUIET)

board_test_run: board_test
	$(TESTRUN)

ui_test: ui_test.go common.go board.go ui.go
	$(TESTCOMPILE_QUIET)

ui_test_run: ui_test
	$(TESTRUN)

tests_compile: common_test board_test ui_test

tests_run: $(KOMOKULIB) common_test board_test ui_test
	@for targ in $(^); do \
		 if test x"$$targ" != x"$(KOMOKULIB)"; then \
		 	make "$$targ"_run; \
		 fi \
	done;

clean:
	@for obj in $(OBJS); do \
		echo '  rm '"$$obj"; \
		rm -f "$$obj"; \
	done;
