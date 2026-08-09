OUTDIR := ../shell
CC     ?= gcc
CFLAGS ?= -O2 -Wall -Wextra
LDFLAGS ?= -s

.PHONY: build build-all clean

build:
	$(CC) $(CFLAGS) $(LDFLAGS) -o $(OUTDIR)/sleep sleep.c

# Cross-compilation is not wired up for the C tool; native only.
build-all: build

clean:
	rm -f $(OUTDIR)/sleep $(OUTDIR)/sleep.exe
