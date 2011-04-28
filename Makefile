# Build the PAPI Go package.
# By Scott Pakin <pakin@lanl.gov>

include $(GOROOT)/src/Make.inc

VERSION=1.0

TARG=papi

CGOFILES=\
	papi.go\
	papi-errno.go\
	papi-event.go\

DISTFILES=\
	papi.go\
	consts2code\
	Makefile

include $(GOROOT)/src/Make.pkg

# ---------------------------------------------------------------------------

# We use a helper Perl script, consts2code, to generate papi-errno.go
# and papi-event.go.

PERL=perl
PAPI_INCDIR:=$(dir $(shell $(PERL) consts2code papi.h))

papi-errno.go: consts2code $(PAPI_INCDIR)/papi.h
	$(PERL) consts2code \
	  papi.h \
	  Errno \
	  "The following constants can be returned as Errno values from PAPI functions." \
	  'PAPI_E.*-\d|PAPI_OK' > papi-errno.go

papi-event.go: consts2code $(PAPI_INCDIR)/papiStdEventDefs.h
	$(PERL) consts2code \
	  papiStdEventDefs.h \
	  Event \
	  "The following constants represent PAPI's standard event types." \
	  '_idx' | grep -v PAPI_END > papi-event.go

CLEANFILES += papi-errno.go papi-event.go

# ---------------------------------------------------------------------------

FULLNAME=gopapi-$(VERSION)

dist: $(DISTFILES)
	mkdir $(FULLNAME)
	cp $(DISTFILES) $(FULLNAME)
	tar -czf $(FULLNAME).tar.gz $(FULLNAME)
	$(RM) -r $(FULLNAME)
	tar -tzvf $(FULLNAME).tar.gz
