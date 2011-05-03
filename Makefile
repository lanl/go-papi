# Build the PAPI Go package.
# By Scott Pakin <pakin@lanl.gov>

include $(GOROOT)/src/Make.inc

VERSION=1.0

TARG=papi

CGOFILES=\
	papi.go\
	papi-high.go\
	papi-low.go\
	papi-mh.go\
	papi-errno.go\
	papi-event.go\
	papi-emod.go\

DISTFILES=\
	papi.go\
	papi-high.go\
	papi-low.go\
	papi-mh.go\
	consts2code\
	Makefile\
	papi_test.go\
	papi_hl_test.go\
	papi_ll_test.go\

include $(GOROOT)/src/Make.pkg

# ---------------------------------------------------------------------------

# We use a helper Perl script, consts2code, to generate papi-errno.go
# and papi-event.go.

PERL=perl
PAPI_INCDIR:=$(dir $(shell $(PERL) consts2code papi.h))

papi-errno.go: consts2code $(PAPI_INCDIR)/papi.h
	$(PERL) consts2code \
	  papi.h \
	  --format='%s os.Error = Errno(C.PAPI_%s)' \
	  --comment="The following constants can be returned as Errno values from PAPI functions." \
	  --keep='#define' \
	  --keep='PAPI_E.*-\d' | \
	  awk '{print} /import/ {print "import \"os\""}' | \
	  sed 's/const /var /' > papi-errno.go

papi-event.go: consts2code $(PAPI_INCDIR)/papiStdEventDefs.h
	$(PERL) consts2code \
	  papiStdEventDefs.h \
	  --format='%s Event = C.PAPI_%s' \
	  --comment="The following constants represent PAPI's standard event types." \
	  --keep='#define' \
	  --keep='_idx' | grep -v PAPI_END > papi-event.go

papi-emod.go: consts2code $(PAPI_INCDIR)/papi.h
	$(PERL) consts2code \
	  papi.h \
	  --format='%s EventModifier = C.PAPI_%s' \
	  --comment="An EventModifier filters the set of events returned by EnumEvents." \
	  --keep='PAPI_\w*ENUM\w+' \
	  --ignore='<<' \
	  --no-ifdef > papi-emod.go

CLEANFILES += papi-errno.go papi-event.go papi-emod.go

# ---------------------------------------------------------------------------

FULLNAME=gopapi-$(VERSION)

dist: $(DISTFILES)
	mkdir $(FULLNAME)
	cp $(DISTFILES) $(FULLNAME)
	tar -czf $(FULLNAME).tar.gz $(FULLNAME)
	$(RM) -r $(FULLNAME)
	tar -tzvf $(FULLNAME).tar.gz
