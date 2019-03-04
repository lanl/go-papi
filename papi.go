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
import "unsafe"
import "fmt"

// An Errno is the PAPI error number.
type Errno int32

// Internally to the package, we test for papi_ok even though we
// always convert this to nil when returning an error to the user.
const papi_ok = C.PAPI_OK

// Convert a PAPI error number to a string.
func (err Errno) String() (errMsg string) {
	if papiErrStr := C.PAPI_strerror(C.int(err)); papiErrStr == nil {
		errMsg = "Unknown PAPI error"
	} else {
		errMsg = C.GoString(papiErrStr)
	}
	return
}

// Make PAPI error numbers implement the error interface.
func (err Errno) Error() (errMsg string) {
	return err.String()
}

// ----------------------------------------------------------------------

// An Event is a PAPI event code, either preset or native.
type Event int32

// Convert a PAPI event code to a string.
func (ecode Event) String() (ename string) {
	cstring := (*C.char)(C.malloc(C.PAPI_MAX_STR_LEN))
	defer C.free(unsafe.Pointer(cstring))
	if Errno(C.PAPI_event_code_to_name(C.int(ecode), cstring)) == papi_ok {
		ename = C.GoString(cstring)
	}
	return
}

// Convert a string to a PAPI event code.  This is particularly useful
// for looking up the event code associated with a PAPI native event.
func StringToEvent(ename string) (ecode Event, err error) {
	cstring := C.CString(ename)
	defer C.free(unsafe.Pointer(cstring))
	var c_ecode C.int
	if errno := Errno(C.PAPI_event_name_to_code(cstring, &c_ecode)); errno == papi_ok {
		ecode = Event(c_ecode)
	} else {
		err = errno
	}
	return
}

// An EventInfo textually describes a PAPI event.
type EventInfo struct {
	EventCode  Event         // Preset (0x8xxxxxxx) or native (0x4xxxxxxx) event code
	EventType  EventModifier // Event type or category (for preset events only)
	Symbol     string        // Name of the event
	ShortDescr string        // A description suitable for use as a label, typically only implemented for preset events
	LongDescr  string        // A longer description of the event (sentence to paragraph length)
	Derived    string        // Name of the derived type (for presets, usually NOT_DERIVED; for native events, empty string)
	Postfix    string        // String containing postfix operations; only defined for preset events of derived type DERIVED_POSTFIX */
	Code       []uint32      // Array of values that further describe the event (for presets, native event_code values; for native events, register values for event programming)
	Name       []string      // Names of code terms (for presets, native event names, as in Symbol, above; for native events, descriptive strings for each register value presented in the code array)
	Note       string        // An optional developer note supplied with a preset event to delineate platform-specific anomalies or restrictions
}

// An EventModifier filters by characteristic the set of events
// returned by EnumEvents().  It can be ENUM_EVENTS to match all
// events, PRESET_ENUM_AVAIL to match the available preset events,
// NTV_* to match particular sets of native events, or a bitwise-or of
// one or more PRESET_BIT_* EventModifier constants to match by
// characteristic.
type EventModifier int32

// Convert an EventModifier to a string.  Caveat: The result is
// meaningful only for PAPI preset events.
func (emod EventModifier) String() string {
	// Handle a few non-masking events specially.
	switch emod {
	case ENUM_EVENTS:
		return "PAPI_ENUM_EVENTS"
	case ENUM_FIRST:
		return "PAPI_ENUM_FIRST"
	case PRESET_ENUM_AVAIL:
		return "PAPI_PRESET_ENUM_AVAIL"
	}

	// Handle masking events by concatenating their strings.
	result := ""
	unnamedBits := EventModifier(0)
	for b := uint32(2); b < 31; b++ {
		if maskBit := emod & (1 << b); maskBit != 0 {
			if maskName, found := presetBitToString[maskBit]; found {
				result += "|" + maskName
			} else {
				unnamedBits |= maskBit
			}
		}
	}
	if unnamedBits != 0 {
		result += fmt.Sprintf("|0x%x", uint32(unnamedBits))
	}
	return result[1:]
}

// An EventMask filters by defining group (preset or native) the set
// of events returned by EnumEvents().
type EventMask int32

// The following may be used individually or ORed together when passed
// to EnumEvents().
const (
	PRESET_MASK EventMask = C.PAPI_PRESET_MASK // Predefined events only
	NATIVE_MASK EventMask = C.PAPI_NATIVE_MASK // Native events only
)

// Native events associated with a particular PAPI component can be
// selected by EnumEvents() by ORing NATIVE_MASK with a component
// mask.
func ComponentMask(cid int) EventMask {
	return EventMask(0x3c000000 & (cid << 26))
}

// An EventSet is a handle to a PAPI-internal set of events.
type EventSet int32

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

// Describe one level of TLB and one level of cache.  Note that if
// multiple TLB page sizes are supported, this will show up as
// multiple TLBInfo values at the same memory-hierarchy level.
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

// Support for different "components" -- sets of counters -- can be
// compiled into the PAPI library.  A ComponentInfo structure
// describes a wealth of information about an individual component.
type ComponentInfo struct {
	Name                   string // Name of the substrate we're using, usually CVS RCS Id
	Version                string // Version of this substrate, usually CVS Revision
	SupportVersion         string // Version of the support library
	KernelVersion          string // Version of the kernel PMC support driver
	CmpIdx                 int    // Index into the vector array for this component; set at init time
	NumCntrs               int    // Number of hardware counters the substrate supports
	NumMpxCntrs            int    // Number of hardware counters the substrate or PAPI can multiplex supports
	NumPresetEvents        int    // Number of preset events the substrate supports
	NumNativeEvents        int    // Number of native events the substrate supports
	DefaultDomain          int    // The default domain when this substrate is used
	AvailableDomains       int    // Available domains
	DefaultGranularity     int    // The default granularity when this substrate is used
	AvailableGranularities int    // Available granularities
	ItimerSig              int    // Signal number used by the multiplex timer, 0 if not
	ItimerNum              int    // Number of the itimer used by mpx and overflow/profile emulation
	ItimerNs               int    // ns between mpx switching and overflow/profile emulation
	ItimerResNs            int    // ns of resolution of itimer
	HardwareIntrSig        int    // Signal used by hardware to deliver PMC events
	ClockTicks             int    // Clock ticks per second
	OpcodeMatchWidth       int    // Width of opcode matcher if exists, 0 if not
	OSVersion              int    // Currently running kernel version
	HardwareIntr           bool   // HW overflow intr, does not need to be emulated in software
	PreciseIntr            bool   // Performance interrupts happen precisely
	POSIX1bTimers          bool   // Using POSIX 1b interval timers (timer_create) instead of setitimer
	KernelProfile          bool   // Has kernel profiling support (buffered interrupts or sprofil-like)
	KernelMultiplex        bool   // In kernel multiplexing
	DataAddressRange       bool   // Supports data address range limiting
	InstrAddressRange      bool   // Supports instruction address range limiting
	FastCounterRead        bool   // Supports a user level PMC read instruction
	FastRealTimer          bool   // Supports a fast real timer
	FastVirtualTimer       bool   // Supports a fast virtual timer
	Attach                 bool   // Supports attach
	AttachMustPtrace       bool   // Attach must first ptrace and stop the thread/process
	CPU                    bool   // Supports specifying cpu number to use with event set
	Inherit                bool   // Supports child processes inheriting parents counters
	EdgeDetect             bool   // Supports edge detection on events
	Invert                 bool   // Supports invert detection on events
	ProfileEAR             bool   // Supports data/instr/tlb miss address sampling
	CntrGroups             bool   // Underlying hardware uses counter groups (e.g. POWER5)
	CntrUmasks             bool   // Counters have unit masks
	CntrIEAREvents         bool   // Counters support instr event addr register
	CntrDEAREvents         bool   // Counters support data event addr register
	CntrOPCMEvents         bool   // Counter events support opcode matching
}

// ----------------------------------------------------------------------

// The following debug levels can be passed to SetDebugLevel().
const (
	QUIET      = C.PAPI_QUIET      // Option to turn off automatic reporting of return codes < 0 to stderr
	VERB_ECOND = C.PAPI_VERB_ECONT // Option to automatically report any return codes < 0 to stderr and continue
	VERB_ESTOP = C.PAPI_VERB_ESTOP // Option to automatically report any return codes < 0 to stderr and exit
)

// Set the PAPI library's debug level.
func SetDebugLevel(level int) (err error) {
	if errno := Errno(C.PAPI_set_debug(C.int(level))); errno != papi_ok {
		err = errno
	}
	return
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

// Enable PAPI support for multiplexed event sets (event sets
// supporting more counters than what the underlying hardware allows
// by timesharing counters) at the cost of periodic process
// interruptions from an interval timer.  InitMultiplex() needs to be
// called only once per application.
func InitMultiplex() {
	C.PAPI_multiplex_init()
}
