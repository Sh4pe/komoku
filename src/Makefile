# vim: set noexpandtab :

CC=6g
TARG=main
LINK=6l

all: main.6 
	6l -o main main.6
	
main.6: main.go komoku/komoku.a
	6g main.go

komoku/komoku.a: komoku/common.go komoku/board.go komoku/treenode.go
	cd komoku && make

clean:
	cd komoku && make clean
	rm -f main main.6
