# vim: set noexpandtab :

#
# (c) 2010 by David Nies (nies.david@googlemail.com)
#     http://www.twitter.com/Sh4pe
#
# Use of this source code is governed by a license 
# that can be found in the LICENSE file.
#

include common.mk

####################### most common defintions ##########################################
TARG=main

# prefix of this directory
THISDIR=

OBJS= main main.$(OBJSUFF)

#########################################################################################
####################### build rules #####################################################
#########################################################################################
all: $(TARG)


$(TARG): $(TARG).$(OBJSUFF)
	@echo '  $(LINKSTR) $(THISDIR)$(TARG)'; $(LINKSTR) -o $(TARG) $(TARG).$(OBJSUFF)

$(TARG).$(OBJSUFF): main.go komoku/komoku.a
	@echo '  $(CCSTR) $(@)'; $(CC) -o $(@) main.go

komoku/komoku.a: force
	@cd komoku && $(MAKE)

tests_run:
	@cd komoku && $(MAKE) tests_run

benchmarks_run:
	@cd komoku && $(MAKE) benchmarks_run

.PHONY: clean
clean:
	@for obj in $(OBJS); do \
		echo '  rm '"$(THISDIR)$$obj"; \
		rm -f "$$obj"; \
	done;
	@cd komoku && $(MAKE) clean

force: ;
