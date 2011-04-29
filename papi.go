// Copyright (C) 2011, Los Alamos National Security, LLC.
// Use of this source code is governed by a BSD-style license.

/*

This package is a wrapper for PAPI, the Performance API.  PAPI
provides access to CPU performance counters and to other low-level
system information.

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


// NumCounters is the number of hardware counters available on the
// system.  Consequently, the slice passed to functions such as
// StartCounters() should contain no more than NumCounters elements.
var NumCounters int


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


// Return the real-time counter's value in clock cycles.
func GetRealCyc() int64 {
	return int64(C.PAPI_get_real_cyc())
}


// Return the real-time counter's value in microseconds.
func GetRealUsec() int64 {
	return int64(C.PAPI_get_real_usec())
}


// Return the virtual-time counter's value in clock cycles.
func GetVirtCyc() int64 {
	return int64(C.PAPI_get_virt_cyc())
}


// Return the virtual-time counter's value in microseconds.
func GetVirtUsec() int64 {
	return int64(C.PAPI_get_virt_usec())
}


// Return the total real time, total process time, total
// floating-point instructions, and average Mflip/s since the previous
// call to PAPI.Flips().
func Flips() (rtime, ptime float32, flpins int64, mflips float32, err Errno) {
	var c_rtime, c_ptime, c_mflips C.float
	var c_flpins C.longlong
	err = Errno(C.PAPI_flips(&c_rtime, &c_ptime, &c_flpins, &c_mflips))
	if err == OK {
		rtime, ptime, flpins, mflips = float32(c_rtime), float32(c_ptime), int64(c_flpins), float32(c_mflips)
	}
	return
}


// Return the total real time, total process time, total
// floating-point operations, and average Mflop/s since the previous
// call to PAPI.Flops().
func Flops() (rtime, ptime float32, flpops int64, mflops float32, err Errno) {
	var c_rtime, c_ptime, c_mflops C.float
	var c_flpops C.longlong
	err = Errno(C.PAPI_flops(&c_rtime, &c_ptime, &c_flpops, &c_mflops))
	if err == OK {
		rtime, ptime, flpops, mflops = float32(c_rtime), float32(c_ptime), int64(c_flpops), float32(c_mflops)
	}
	return
}


// Return the total real time, total process time, total number of
// instructions, and average instructions per cycle since the previous
// call to PAPI.Ipc().
func Ipc() (rtime, ptime float32, ins int64, ipc float32, err Errno) {
	var c_rtime, c_ptime, c_ipc C.float
	var c_ins C.longlong
	err = Errno(C.PAPI_ipc(&c_rtime, &c_ptime, &c_ins, &c_ipc))
	if err == OK {
		rtime, ptime, ins, ipc = float32(c_rtime), float32(c_ptime), int64(c_ins), float32(c_ipc)
	}
	return
}


// Given a slice of event codes, start counting the corresponding events.
func StartCounters(evcodes []Event) (err Errno) {
	events := (*C.int)(&evcodes[0])
	numEvents := C.int(len(evcodes))
	err = Errno(C.PAPI_start_counters(events, numEvents))
	return
}


// Store the current event counts in a given slice and reset the
// counters to zero.
func ReadCounters(values []int64) (err Errno) {
	valuePtr := (*C.longlong)(&values[0])
	numValues := C.int(len(values))
	err = Errno(C.PAPI_read_counters(valuePtr, numValues))
	return
}


// Add the current event counts to those in a given slice and reset
// the counters to zero.
func AccumCounters(values []int64) (err Errno) {
	valuePtr := (*C.longlong)(&values[0])
	numValues := C.int(len(values))
	err = Errno(C.PAPI_accum_counters(valuePtr, numValues))
	return
}


// Store the current event counts in a given slice, reset the
// counters to zero, and stop counting the events.
func StopCounters(values []int64) (err Errno) {
	valuePtr := (*C.longlong)(&values[0])
	numValues := C.int(len(values))
	err = Errno(C.PAPI_stop_counters(valuePtr, numValues))
	return
}
