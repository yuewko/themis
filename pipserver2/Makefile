################################################################
# content filter interface (cfi) demo
################################################################
#
#  make         to build demo
#  make test    to test
#

PRE = /usr/local
INC = $(PRE)/include
LIB = $(PRE)/lib/64 $(PRE)/lib64 $(PRE)/lib

EXE = demo-list
EXE-adult = demo-adult
EXE-alcohol = demo-alcohol
EXE-drugs = demo-drugs
EXE-gambling = demo-gambling
EXE-guns = demo-guns
EXE-sports = demo-sports
EXE-tobacco = demo-tobacco
EXE-weapons = demo-weapons

ifeq ($(origin CC), default)
CC = gcc
endif

CFLAGS  += $(INC:%=-I%) -Wall -ansi -pedantic -D_GNU_SOURCE
LDFLAGS += $(LIB:%=-L%) -lcfi -lm

ifeq ($(shell uname),Linux)
LDFLAGS += -ldl
endif

ifeq ($(shell uname),SunOS)
CFLAGS  += -D__EXTENSIONS__
LDFLAGS += -lrt -lsocket -lnsl -ldl
endif

pathsearch  = $(firstword $(wildcard $(addsuffix /$(1),$(LIB))))
LIBCFIPATH := $(call pathsearch,libcfi.so)
ifneq (,$(findstring 64,$(LIBCFIPATH)))
    CFLAGS  += -m64
    LDFLAGS += -m64
endif

all: $(EXE) $(EXE-adult) $(EXE-alcohol) $(EXE-drugs) $(EXE-gambling) $(EXE-guns) \
     $(EXE-sports) $(EXE-tobacco) $(EXE-weapons)

$(EXE): $(EXE).o
	$(CC) $(EXE).o $(LDFLAGS) -o $(EXE)

$(EXE).o: $(EXE).c $(INC)/cfi.h

$(EXE-adult): $(EXE-adult).o
	$(CC) $(EXE-adult).o $(LDFLAGS) -o $(EXE-adult)

$(EXE-adult).o: $(EXE-adult).c $(INC)/cfi.h

$(EXE-alcohol): $(EXE-alcohol).o
	$(CC) $(EXE-alcohol).o $(LDFLAGS) -o $(EXE-alcohol)

$(EXE-alcohol).o: $(EXE-alcohol).c $(INC)/cfi.h

$(EXE-drugs): $(EXE-drugs).o
	$(CC) $(EXE-drugs).o $(LDFLAGS) -o $(EXE-drugs)

$(EXE-drugs).o: $(EXE-drugs).c $(INC)/cfi.h

$(EXE-gambling): $(EXE-gambling).o
	$(CC) $(EXE-gambling).o $(LDFLAGS) -o $(EXE-gambling)

$(EXE-gambling).o: $(EXE-gambling).c $(INC)/cfi.h

$(EXE-guns): $(EXE-guns).o
	$(CC) $(EXE-guns).o $(LDFLAGS) -o $(EXE-guns)

$(EXE-guns).o: $(EXE-guns).c $(INC)/cfi.h

$(EXE-sports): $(EXE-sports).o
	$(CC) $(EXE-sports).o $(LDFLAGS) -o $(EXE-sports)

$(EXE-sports).o: $(EXE-sports).c $(INC)/cfi.h

$(EXE-tobacco): $(EXE-tobacco).o
	$(CC) $(EXE-tobacco).o $(LDFLAGS) -o $(EXE-tobacco)

$(EXE-tobacco).o: $(EXE-tobacco).c $(INC)/cfi.h

$(EXE-weapons): $(EXE-weapons).o
	$(CC) $(EXE-weapons).o $(LDFLAGS) -o $(EXE-weapons)

$(EXE-weapons).o: $(EXE-weapons).c $(INC)/cfi.h

clean:
	rm -f $(EXE) $(EXE).o

config: $(EXE)
	@./demo.sh

test: $(EXE)
	@./ivp.sh

help:
	@echo ""
	@echo "Makefile usage:"
	@echo ""
	@echo "  make        - builds $(EXE) program"
	@echo ""
	@echo "  make config - builds a custom configuration file in"
	@echo "                conf.demo and tests some sample content"
	@echo ""
	@echo "  make test   - runs a full suite of tests"
	@echo ""
