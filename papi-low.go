// Copyright (C) 2011, Los Alamos National Security, LLC.
// Use of this source code is governed by a BSD-style license.

// This file provides an interface to PAPI's low-level functions.

package papi


// #include <papi.h>
import "C"


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
