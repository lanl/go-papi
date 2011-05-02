// Copyright (C) 2011, Los Alamos National Security, LLC.
// Use of this source code is governed by a BSD-style license.

// This file provides an interface to PAPI's high-level functions.

package papi


// #include <papi.h>
import "C"
import "os"


// NumCounters is the number of hardware counters available on the
// system.  Consequently, the slice passed to functions such as
// StartCounters() should contain no more than NumCounters elements.
var NumCounters int


// Return the total real time, total process time, total
// floating-point instructions, and average Mflip/s since the previous
// call to PAPI.Flips().
func Flips() (rtime, ptime float32, flpins int64, mflips float32, err os.Error) {
	var c_rtime, c_ptime, c_mflips C.float
	var c_flpins C.longlong
	errno := Errno(C.PAPI_flips(&c_rtime, &c_ptime, &c_flpins, &c_mflips))
	if errno == papi_ok {
		rtime, ptime, flpins, mflips = float32(c_rtime), float32(c_ptime), int64(c_flpins), float32(c_mflips)
	} else {
		err = errno
	}
	return
}

// Return the total real time, total process time, total
// floating-point operations, and average Mflop/s since the previous
// call to PAPI.Flops().
func Flops() (rtime, ptime float32, flpops int64, mflops float32, err os.Error) {
	var c_rtime, c_ptime, c_mflops C.float
	var c_flpops C.longlong
	errno := Errno(C.PAPI_flops(&c_rtime, &c_ptime, &c_flpops, &c_mflops))
	if errno == papi_ok {
		rtime, ptime, flpops, mflops = float32(c_rtime), float32(c_ptime), int64(c_flpops), float32(c_mflops)
	} else {
		err = errno
	}
	return
}


// Return the total real time, total process time, total number of
// instructions, and average instructions per cycle since the previous
// call to PAPI.Ipc().
func Ipc() (rtime, ptime float32, ins int64, ipc float32, err os.Error) {
	var c_rtime, c_ptime, c_ipc C.float
	var c_ins C.longlong
	errno := Errno(C.PAPI_ipc(&c_rtime, &c_ptime, &c_ins, &c_ipc))
	if errno == papi_ok {
		rtime, ptime, ins, ipc = float32(c_rtime), float32(c_ptime), int64(c_ins), float32(c_ipc)
	} else {
		err = errno
	}
	return
}


// Given a slice of event codes, start counting the corresponding events.
func StartCounters(evcodes []Event) (err os.Error) {
	events := (*C.int)(&evcodes[0])
	numEvents := C.int(len(evcodes))
	if errno := Errno(C.PAPI_start_counters(events, numEvents)); errno != papi_ok {
		err = errno
	}
	return
}


// Store the current event counts in a given slice and reset the
// counters to zero.
func ReadCounters(values []int64) (err os.Error) {
	valuePtr := (*C.longlong)(&values[0])
	numValues := C.int(len(values))
	if errno := Errno(C.PAPI_read_counters(valuePtr, numValues)); errno != papi_ok {
		err = errno
	}
	return
}


// Add the current event counts to those in a given slice and reset
// the counters to zero.
func AccumCounters(values []int64) (err os.Error) {
	valuePtr := (*C.longlong)(&values[0])
	numValues := C.int(len(values))
	if errno := Errno(C.PAPI_accum_counters(valuePtr, numValues)); errno != papi_ok {
		err = errno
	}
	return
}


// Store the current event counts in a given slice, reset the
// counters to zero, and stop counting the events.
func StopCounters(values []int64) (err os.Error) {
	valuePtr := (*C.longlong)(&values[0])
	numValues := C.int(len(values))
	if errno := Errno(C.PAPI_stop_counters(valuePtr, numValues)); errno != papi_ok {
		err = errno
	}
	return
}
