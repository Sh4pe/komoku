# vim: set noexpandtab :

all: main.6
	6l -o komoku main.6
	
main.6: main.go treenode.6 gamestate.6
	6g main.go

treenode.6: treenode.go gamestate.6 common.6
	6g treenode.go

gamestate.6: gamestate.go
	6g gamestate.go

common.6: common.go
	6g common.go

clean:
	rm -f komoku main.6 treenode.6 gamestate.6 common.6
