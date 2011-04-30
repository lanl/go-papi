// Copyright (C) 2011, Los Alamos National Security, LLC.
// Use of this source code is governed by a BSD-style license.

// This file defines various PAPI datatypes and methods on those types
// and initializes the PAPI library.

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

// ----------------------------------------------------------------------

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

// ----------------------------------------------------------------------

// An AddressMap stores information about the currently running program.
type AddressMap struct {
	Name      string
	TextStart uintptr // Start address of program text segment 
	TextEnd   uintptr // End address of program text segment 
	DataStart uintptr // Start address of program data segment 
	DataEnd   uintptr // End address of program data segment 
	BssStart  uintptr // Start address of program bss segment 
	BssEnd    uintptr // End address of program bss segment 
}


// A ProgramInfo is just like an AddressMap but additionally stores
// the full patname of the executable.
type ProgramInfo struct {
	FullName    string     // Program path+name
	AddressInfo AddressMap // Other program information
}

// ----------------------------------------------------------------------

// Attributes of each level of the memory hierarchy are described by
// MHAttrs.
type MHAttrs int32


// Map MHAttrs cache-type primitives to strings.
var ctype2string map[MHAttrs]string = map[MHAttrs]string{
	MH_TYPE_INST:    "INST",
	MH_TYPE_DATA:    "DATA",
	MH_TYPE_VECTOR:  "VECTOR",
	MH_TYPE_TRACE:   "TRACE",
	MH_TYPE_UNIFIED: "UNIFIED"}

// Map MHAttrs write-policy primitives to strings.
var wpol2string map[MHAttrs]string = map[MHAttrs]string{
	MH_TYPE_WT: "WT",
	MH_TYPE_WB: "WB"}

// Map MHAttrs replacement-policy primitives to strings.
var repl2string map[MHAttrs]string = map[MHAttrs]string{
	MH_TYPE_LRU:        "LRU",
	MH_TYPE_PSEUDO_LRU: "PSEUDO_LRU"}


// Output memory-hierarchy attributes as a user-friendly string.
func (a MHAttrs) String() string {
	var str string // String to return
	if ctype, ok := ctype2string[a&0xF]; ok {
		str = ctype
	}
	if wpol, ok := wpol2string[a&0xF0]; ok {
		str += "|" + wpol
	}
	if repl, ok := repl2string[a&0xF00]; ok {
		str += "|" + repl
	}
	if a&MH_TYPE_TLB == MH_TYPE_TLB {
		str += "|" + "TLB"
	}
	if a&MH_TYPE_PREF == MH_TYPE_PREF {
		str += "|" + "PREF"
	}
	if str[0] == '|' {
		return str[1:]
	}
	return str
}


// A fully associative cache or TLB is defined to have associativity
// FullyAssociative.
const FullyAssociative = C.SHRT_MAX


// Describe a translation lookaside buffer's characteristics.
type TLBInfo struct {
	Type          MHAttrs // Cache attributes of the TLB
	NumEntries    int32   // Number of entries in the TLB
	PageSize      int32   // Page size in bytes
	Associativity int32   // TLB associativity (0=unknown)
}


// Return the cache type of a TLB (MH_TYPE_EMPTY, MH_TYPE_INST,
// MH_TYPE_DATA, MH_TYPE_VECTOR, or MH_TYPE_UNIFIED, but not
// MH_TYPE_TRACE).
func (ci *TLBInfo) CacheType() MHAttrs {
	return ci.Type & 0xF
}


// Return the TLB write policy (MH_TYPE_WT or MH_TYPE_WB).
func (ci *TLBInfo) CacheWritePolicy() MHAttrs {
	return ci.Type & 0xF0
}


// Return the TLB replacement policy (MH_TYPE_UNKNOWN, MH_TYPE_LRU, or
// MH_TYPE_PSEUDO_LRU).
func (ci *TLBInfo) CacheReplacementPolicy() MHAttrs {
	return ci.Type & 0xF00
}


// Return the TLB usage type (MH_TYPE_TLB, MH_TYPE_PREF, or neither).
func (ci *TLBInfo) CacheUsage() MHAttrs {
	return ci.Type & 0xF000
}


// Describe a cache's characteristics.
type CacheInfo struct {
	Type          MHAttrs // Cache attributes
	Size          int32   // Cache size in bytes
	LineSize      int32   // Line size in bytes
	NumLines      int32   // Number of cache lines
	Associativity int32   // Cache associativity (0=unknown)
}


// Return the cache type of a cache (MH_TYPE_EMPTY, MH_TYPE_INST,
// MH_TYPE_DATA, MH_TYPE_VECTOR, MH_TYPE_TRACE, or MH_TYPE_UNIFIED).
func (ci *CacheInfo) CacheType() MHAttrs {
	return ci.Type & 0xF
}


// Return the cache write policy (MH_TYPE_WT or MH_TYPE_WB).
func (ci *CacheInfo) CacheWritePolicy() MHAttrs {
	return ci.Type & 0xF0
}


// Return the cache replacement policy (MH_TYPE_UNKNOWN, MH_TYPE_LRU,
// or MH_TYPE_PSEUDO_LRU).
func (ci *CacheInfo) CacheReplacementPolicy() MHAttrs {
	return ci.Type & 0xF00
}


// Return the cache usage type (MH_TYPE_TLB, MH_TYPE_PREF, or neither).
func (ci *CacheInfo) CacheUsage() MHAttrs {
	return ci.Type & 0xF000
}


// Describe one level of TLB and one level of cache.
type MHLevelInfo struct {
	TLB   []TLBInfo   // Information about all TLBs at the current level of the memory hierarchy
	Cache []CacheInfo // Information about all caches at the current level of the memory hierarchy
}


// Describe all of the hardware that PAPI knows about.
type HardwareInfo struct {
	CPUs          int32         // Number of CPUs per NUMA node
	Threads       int32         // Number of hardware threads per core
	Cores         int32         // Number of cores per socket
	Sockets       int32         // Number of sockets
	NUMANodes     int32         // Total Number of NUMA Nodes
	TotalCPUs     int32         // Total number of CPUs in the entire system
	Vendor        int32         // Vendor number of the CPU
	VendorName    string        // Vendor name of the CPU
	Model         int32         // Model number of the CPU
	ModelName     string        // Model name of the CPU
	Revision      float32       // Revision number of the CPU
	CPUIDFamily   int32         // CPUID family
	CPUIDModel    int32         // CPUID model
	CPUIDStepping int32         // CPUID stepping
	MHz           float32       // CPU's current clock rate in megahertz
	ClockMHz      int32         // CPUs cycle counter's current clock rate in megahertz
	MemHierarchy  []MHLevelInfo // Information about each level of the memory hierarchy
}

// ----------------------------------------------------------------------

// DynMemInfo represents the dynamic memory usage of the current
// program.  According to the PAPI documentation, this function is
// currently implemented only for the Linux operating system.
type DynMemInfo struct {
	Peak          int64 // Peak size of process image, may be 0 on older Linux systems
	Size          int64 // Size of process image
	Resident      int64 // Resident set size
	HighWaterMark int64 // High-water memory usage
	Shared        int64 // Shared memory
	Text          int64 // Memory allocated to code
	Library       int64 // Memory allocated to libraries
	Heap          int64 // Size of the heap
	Locked        int64 // Locked memory
	Stack         int64 // Size of the stack
	PageSize      int64 // Size of a page
	PTE           int64 // Size  of page table entries, may be 0 on older Linux systems
}

// ----------------------------------------------------------------------

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
