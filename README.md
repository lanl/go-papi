go-papi
=======


Description
-----------

go-papi provides a Go interface to PAPI, the [Performance Application
Programming Interface](http://icl.cs.utk.edu/papi/).  PAPI provides
convenient access to hardware performance counters, primarily those
provided by the CPU but with others (e.g., various networks) also
available.

As the [PAPI(3) man
page](http://icl.cs.utk.edu/projects/papi/wiki/PAPIC:PAPI.3) explains,
the PAPI API is split into "high level" and "low level" functions,
with the former being simpler to use but less flexible and the latter
providing finer-grained control over the sets of data measured and
reported.


Installation
------------

Installation is a bit of a pain, as the new
[`go`](http://weekly.golang.org/cmd/go/) tool doesn't yet handle
custom `Makefile`s such as the one `go-papi` requires.  (This
requirement is due to some of `go-papi`'s source files being generated
automatically by a Perl script.)

First, download `go-papi` into your Go build tree without
automatically building/installing it:

<pre>
    go get -d -v github.com/losalamos/go-papi
</pre>

Set the `PAPI_INCDIR` environment variable to the directory containing
`papi.h`.  Also, ensure that the directory containing `libpapi.so` is
listed in your `LD_LIBRARY_PATH`.

Next, switch to the `go-papi` directory and build/test/install the
package:

<pre>
    cd $GOROOT/src/pkg/github.com/losalamos/go-papi
    make
    make check
    make install
</pre>

It is then safe to do a `make clean` to remove all of the byproducts
of the installation process.


Documentation
-------------

Pre-built documentation for the core part of the go-papi API is
available online at
<http://gopkgdoc.appspot.com/pkg/github.com/losalamos/go-papi>,
courtesy of [GoPkgDoc](http://gopkgdoc.appspot.com/).  Unfortunately,
the online documentation omits descriptions of all constants,
variables, etc. that are generated during the build process,
specifically the list of PAPI events (`papi-event.go`), event
modifiers (`papi-emod.go`), and error values (`papi-errno.go`).

Once you install go-papi, you can view the complete go-papi API with
[`godoc`](http://golang.org/cmd/godoc/), for example by running

<pre>
    godoc -http=:6060 -index
</pre>

to start a local Web server then viewing the documentation at
<http://localhost:6060/pkg/github.com/losalamos/go-papi/> in your
favorite browser.

For code examples, take a look at the `*_test.go` files in the go-papi
source distribution.  `papi_hl_test.go` utilizes PAPI's high-level
API; `papi_ll_test.go` utilizes PAPI's low-level API; and
`papi_test.go` utilizes a few miscellaneous functions.


License
-------

BSD-ish with a "modifications must be indicated" clause.  See
<http://github.com/losalamos/go-papi/blob/master/LICENSE> for the full
text.


Author
------

Scott Pakin, <pakin@lanl.gov>
