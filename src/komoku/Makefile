# vim: set noexpandtab :

CC=6g
TARG=komoku
PACK=gopack
LINK=6l

TESTRUNNER='package main; import ( "./common_test"; "testing"; ); func main() { testing.Main(komoku.Testsuite()) }'
TRUNNERPRE='package main; import ( "./'
TRUNNERPOST='"; "testing"; ); func main() { testing.Main(komoku.Testsuite()) }'

all: $(TARG).6
	rm -f $(TARG).a
	$(PACK) grc $(TARG).a $(TARG).6

$(TARG).6:
	$(CC) -o $(TARG).6 common.go gamestate.go treenode.go

# the stuff in the target may be copy&pasted
common_test:
	# create tmp file with adjusted testrunner in it
	echo $(TRUNNERPRE)$(@)$(TRUNNERPOST) > $(@).tmp.go 
	# if the target is foo_test, sed is used to produce the last word of the command
	# 6g -o foo_test.6 foo_test.go foo.go
	$(CC) -o $(@).6 $(@).go `echo $(@) | sed 's/_[a-zA-Z]*//g'`.go
	# compile the tmp file
	$(CC) -o $(@).tmp.6 $(@).tmp.go
	$(LINK) -o $(@) $(@).tmp.6
	# clean up the tmp files
	rm -f $(@).tmp.6
	rm -f $(@).tmp.go

clean:
	rm -f $(TARG).6 $(TARG).a common_test common_test.6
