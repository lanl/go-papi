// Copyright (C) 2011, Los Alamos National Security, LLC.
// Use of this source code is governed by a BSD-style license.

// This file defines some basic PAPI datatypes and methods on those
// types and initializes the PAPI library.

/*
This package presents a Go interface to PAPI, the Performance API.
PAPI provides access to CPU performance counters and to other
low-level system information.
*/
package papi


/*
#cgo LDFLAGS: -lpapi -lpthread
#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>
#include <papi.h>

// Wrap PAPI_thread_init() to simplify passing pthread_self() around.
int initialize_papi_threading (void)
{
  return PAPI_thread_init(pthread_self);
}
*/
import "C"
import "fmt"
import "unsafe"


// An Error can represent any printable error condition.
type Error interface {
	String() string
}


// An Errno is the PAPI error number.
type Errno int32


// Convert a PAPI error number to a string.
func (err Errno) String() (errMsg string) {
	if papiErrStr := C.PAPI_strerror(C.int(err)); papiErrStr == nil {
		errMsg = "Unknown PAPI error"
	} else {
		errMsg = C.GoString(papiErrStr)
	}
	return
}


// An Event is a PAPI event code, either preset or native.
type Event int32


// Convert a PAPI event code to a string.
func (ecode Event) String() (ename string) {
	cstring := (*C.char)(C.malloc(C.PAPI_MAX_STR_LEN))
	defer C.free(unsafe.Pointer(cstring))
	if Errno(C.PAPI_event_code_to_name(C.int(ecode), cstring)) == OK {
		ename = C.GoString(cstring)
	}
	return
}


// Convert a string to a PAPI event code.
func StringToEvent(ename string) (ecode Event, err Errno) {
	cstring := C.CString(ename)
	defer C.free(unsafe.Pointer(cstring))
	var c_ecode C.int
	if err = Errno(C.PAPI_event_name_to_code(cstring, &c_ecode)); err == OK {
		ecode = Event(c_ecode)
	}
	return
}


// Before we do anything else we need to initialize the PAPI library.
func init() {
	// Initialize the library proper.
	switch initval := C.PAPI_library_init(C.PAPI_VER_CURRENT); {
	case initval == C.PAPI_VER_CURRENT:
		{
		}
	case initval > 0:
		panic(fmt.Sprintf("PAPI library version mismatch: expected %d but saw %d",
			C.PAPI_VER_CURRENT, initval))
	case initval < 0:
		panic(Errno(initval).String())
	}

	// Initialize the library's thread support.
	threadval := C.initialize_papi_threading()
	if threadval != C.PAPI_OK {
		panic(Errno(threadval).String())
	}

	// Initialize the high-level counter support.
	if nc := C.PAPI_num_counters(); nc >= 0 {
		NumCounters = int(nc)
	} else {
		panic(Errno(nc).String())
	}
}
