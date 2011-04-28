# Build the PAPI Go package.
# By Scott Pakin <pakin@lanl.gov>

include $(GOROOT)/src/Make.inc

VERSION=1.0

TARG=papi

CGOFILES=\
	papi.go\
	papi-errno.go\
	papi-event.go\

EXTRA_DIST=\
	consts2code\
	Makefile

include $(GOROOT)/src/Make.pkg

# ---------------------------------------------------------------------------

# If you need to regenerate papi-errno.go and papi-event.go, here's
# how to do it.  Note that we add the generated .go files to
# CLEANFILES only if we believe we can recreate them.

PAPI_INCDIR=/usr/include
PERL=perl

ifeq ($(wildcard $(PAPI_INCDIR)/papi.h), )

papi-errno.go papi-event.go:
	$(error Please define PAPI_INCDIR as the directory containing papi.h and papiStdEventDefs.h)

else

papi-errno.go: consts2code $(PAPI_INCDIR)/papi.h
	$(PERL) consts2code \
	  $(PAPI_INCDIR)/papi.h \
	  Errno \
	  "The following constants can be returned as Errno values from PAPI functions." \
	  'PAPI_E.*-\d|PAPI_OK' > papi-errno.go

papi-event.go: consts2code $(PAPI_INCDIR)/papiStdEventDefs.h
	$(PERL) consts2code \
	  $(PAPI_INCDIR)/papiStdEventDefs.h \
	  Event \
	  "The following constants represent PAPI's standard event types." \
	  '_idx' | grep -v PAPI_END > papi-event.go

CLEANFILES += papi-errno.go papi-event.go

endif

# ---------------------------------------------------------------------------

FULLNAME=gopapi-$(VERSION)

dist: $(CGOFILES) $(EXTRA_DIST)
	mkdir $(FULLNAME)
	cp $(CGOFILES) $(EXTRA_DIST) $(FULLNAME)
	tar -czf $(FULLNAME).tar.gz $(FULLNAME)
	$(RM) -r $(FULLNAME)
	tar -tzvf $(FULLNAME).tar.gz
