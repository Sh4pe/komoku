# vim: set noexpandtab :

#  *(1)
#
#  # create tmp file with adjusted testrunner in it
#  echo $(TRUNNERPRE)$(@)$(TRUNNERPOST) > $(@).tmp.go 
#  # if the target is foo_test, sed is used to produce the last word of the command
#  # 6g -o foo_test.6 foo_test.go foo.go
#  $(CC) -o $(@).6 $(^) `echo $(@) | sed 's/_[a-zA-Z]*//g'`.go
#  # compile the tmp file
#  $(CC) -o $(@).tmp.6 $(@).tmp.go
#  $(LINK) -o $(@) $(@).tmp.6
#  # clean up the tmp files
#  rm -f $(@).tmp.6
#  rm -f $(@).tmp.go
#
#  Notes: the sed part is not used any more...

CC=6g
TARG=komoku
PACK=gopack
LINK=6l

TESTRUNNER='package main; import ( "./common_test"; "testing"; ); func main() { testing.Main(komoku.Testsuite()) }'
TRUNNERPRE='package main; import ( "./'
TRUNNERPOST='"; "testing"; ); func main() { testing.Main(komoku.Testsuite()) }'

all: $(TARG).a

$(TARG).a: $(TARG).6
	rm -f $(TARG).a
	$(PACK) grc $(TARG).a $(TARG).6
	
$(TARG).6: common.go board.go treenode.go ui.go
	$(CC) -o $(TARG).6 common.go board.go treenode.go ui.go

# for documentation see *(1)
common_test: common_test.go common.go
	echo $(TRUNNERPRE)$(@)$(TRUNNERPOST) > $(@).tmp.go 
	$(CC) -o $(@).6 $(^)
	$(CC) -o $(@).tmp.6 $(@).tmp.go
	$(LINK) -o $(@) $(@).tmp.6
	rm -f $(@).tmp.6
	rm -f $(@).tmp.go

common_test_run: common_test
	echo 'Running common_test'
	./common_test

board_test: board_test.go common.go board.go
	echo $(TRUNNERPRE)$(@)$(TRUNNERPOST) > $(@).tmp.go 
	$(CC) -o $(@).6 $(^) 
	$(CC) -o $(@).tmp.6 $(@).tmp.go
	$(LINK) -o $(@) $(@).tmp.6
	rm -f $(@).tmp.6
	rm -f $(@).tmp.go

board_test_run: board_test
	echo 'Running board_test'
	./board_test

ui_test: ui_test.go common.go board.go ui.go
	echo $(TRUNNERPRE)$(@)$(TRUNNERPOST) > $(@).tmp.go 
	$(CC) -o $(@).6 $(^)
	$(CC) -o $(@).tmp.6 $(@).tmp.go
	$(LINK) -o $(@) $(@).tmp.6
	rm -f $(@).tmp.6
	rm -f $(@).tmp.go

ui_test_run: ui_test
	echo 'Running ui_test'
	./ui_test

tests_compile: common_test board_test ui_test

tests_run: $(TARG).a tests_compile
	echo
	echo 'running tests....'
	./common_test 
	./board_test 
	./ui_test

clean:
	rm -f $(TARG).6 $(TARG).a common_test common_test.6 board_test board_test.6 ui_test.6 ui_test
