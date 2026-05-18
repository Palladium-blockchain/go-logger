GO ?= go
PKGS ?= ./...
BENCH ?= .
BENCHTIME ?= 1s
COUNT ?= 1
PROFILE_PKG ?= ./pkg/logger

.PHONY: test bench benchmark bench-stable bench-profile

test:
	$(GO) test ./...

bench benchmark:
	$(GO) test $(PKGS) -run='^$$' -bench='$(BENCH)' -benchmem -benchtime=$(BENCHTIME) -count=$(COUNT)

bench-stable:
	$(MAKE) bench COUNT=10 BENCHTIME=3s

bench-profile:
	$(GO) test $(PROFILE_PKG) -run='^$$' -bench='$(BENCH)' -benchmem -benchtime=$(BENCHTIME) -count=1 -cpuprofile=cpu.out -memprofile=mem.out
