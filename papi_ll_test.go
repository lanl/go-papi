// This file tests the low-level features of the PAPI library.

package papi

import "testing"
import "time"

// Ensure that the real-time cycle counter is strictly increasing.
func TestRealCyc(t *testing.T) {
	time1 := GetRealCyc()
	time2 := GetRealCyc()
	if time1 >= time2 {
		t.Fatalf("Real-time cycles are not strictly increasing: %d >= %d",
			time1, time2)
	}
}

// Ensure that the real-time microsecond counter is increasing.
func TestRealUsec(t *testing.T) {
	const sleep_usecs = 10000
	time1 := GetRealUsec()
	time.Sleep(sleep_usecs * 1000)
	time2 := GetRealUsec()
	if time1 > time2 {
		t.Fatalf("Real-time microseconds are decreasing: %d >= %d",
			time1, time2)
	} else {
		if time2-time1 < sleep_usecs {
			t.Fatalf("Real-time microseconds increase too slowly: %d vs. %d",
				time2-time1, sleep_usecs)
		}
	}
}

// Ensure that the virtual-time cycle counter is increasing.  Ideally,
// it should be strictly increasing, but this doesn't seem to be the
// case on all systems.
func TestVirtCyc(t *testing.T) {
	const maxTimings = 1000000000
	time1 := GetVirtCyc()
	time2 := time1
	for i := 0; time1 == time2 && i < maxTimings; i++ {
		time2 = GetVirtCyc()
	}
	if time1 > time2 {
		t.Fatalf("Virtual-time cycles are decreasing: %d >= %d",
			time1, time2)
	}
}

// Ensure that the virtual-time microsecond counter is increasing.
func TestVirtUsec(t *testing.T) {
	const maxTimings = 1000000000
	time1 := GetVirtUsec()
	time2 := time1
	for i := 0; time1 == time2 && i < maxTimings; i++ {
		time2 = GetVirtUsec()
	}
	if time1 > time2 {
		t.Fatalf("Virtual-time microseconds are decreasing: %d >= %d",
			time1, time2)
	}
}

// Ensure that GetExecutableInfo() at least does *something*.  Not
// every value is populated on every system, however.
func TestExeInfo(t *testing.T) {
	exeInfo := GetExecutableInfo()
	addrInfo := exeInfo.AddressInfo
	if addrInfo.Name == "" || exeInfo.FullName == "" {
		t.Fatal("GetExecutableInfo() returned empty program names")
	}
	if addrInfo.TextStart == 0 || addrInfo.TextEnd == 0 {
		t.Fatal("GetExecutableInfo() returned zeroes for the text-segment boundaries")
	}
}

// Ensure that PAPI supports at least one counting component (the CPU).
func TestNumComponents(t *testing.T) {
	if nc := GetNumComponents(); nc < 1 {
		t.Fatalf("Expected to see at least 1 counting component but saw only %d", nc)
	}
}

// Ensure that selected pieces of hardware information are valid.
func TestHardwareInfo(t *testing.T) {
	hw := GetHardwareInfo()
	if hw.TotalCPUs == 0 {
		t.Fatal("TotalCPUs == 0")
	}
	if hw.VendorName == "" {
		t.Fatal("VendorName == \"\"")
	}
	if hw.MHz == 0.0 {
		t.Fatal("MHz == 0.0")
	}
	if hw.ClockMHz == 0.0 {
		t.Fatal("ClockMHz == 0.0")
	}
	if len(hw.MemHierarchy) == 0 {
		t.Fatal("MemHierarchy == []")
	}
}

// Do all of the work for TestEventSet and TestMultiplex.
func useEventSet(t *testing.T, events EventSet) {
	const flops = 1000

	// Check that we have at least two CPU counters.
	if GetNumCounters(0) < 2 {
		t.Fatal("Testing event sets requires a CPU with at least two counters")
	}

	// Start counting a few events.
	var err error
	if err = events.AddEvents([]Event{TOT_INS, TOT_CYC}); err != nil {
		t.Fatal(err)
	}
	if numEvents, err := events.NumEvents(); err != nil {
		t.Fatal(err)
	} else if numEvents != 2 {
		t.Fatalf("Expected 2 events, saw %d events", numEvents)
	}
	if err = events.Start(); err != nil {
		t.Fatal(err)
	}

	// Kill some time.
	performWork(flops)

	// Ensure we counted at least a minimum number of flops.
	values := make([]int64, 2)
	if err = events.Stop(values); err != nil {
		t.Fatal(err)
	}
	if err = events.RemoveEvents([]Event{TOT_CYC, TOT_INS}); err != nil {
		t.Fatal(err)
	}
	if err = events.DestroyEventSet(); err != nil {
		t.Fatal(err)
	}
}

// Test various events in an EventSet's lifetime.  This test is
// derived from examples/PAPI_add_remove_events.c in the PAPI
// distribution.
func TestEventSet(t *testing.T) {
	if events, err := CreateEventSet(); err != nil {
		t.Fatal(err)
	} else {
		useEventSet(t, events)
	}
}

// Test multiplexed event sets.
func TestMultiplex(t *testing.T) {
	InitMultiplex()
	var err error
	var events EventSet
	if events, err = CreateEventSet(); err != nil {
		t.Fatal(err)
	}
	if isMulti, err := events.GetMultiplex(); err != nil {
		t.Fatal(err)
	} else if isMulti {
		t.Fatal("Expected a non-multiplexed event set but got a multiplexed one")
	}
	if err = events.AssignComponent(0); err != nil {
		t.Fatal(err)
	}
	if err = events.SetMultiplex(); err != nil {
		t.Fatal(err)
	}
	useEventSet(t, events)
}

// Ensure that GetEventInfo returns non-empty data.
func TestGetEventInfo(t *testing.T) {
	info, err := GetEventInfo(TOT_INS)
	if err != nil {
		t.Fatal(err)
	}
	if info.Symbol != "PAPI_TOT_INS" {
		t.Fatal("Expected PAPI_TOT_INS; saw %s", info.Symbol)
	}
	if info.ShortDescr == "" || info.LongDescr == "" {
		t.Fatal("Event description is empty")
	}
}

// Ensure that enumerating events gives us at least one event.
func TestEnumEvents(t *testing.T) {
	var eventList []Event
	var err error

	// Look for preset events.
	eventList, err = EnumEvents(PRESET_MASK, ENUM_EVENTS)
	if err != nil {
		t.Fatal(err)
	}
	numPresets := len(eventList)
	if numPresets == 0 {
		t.Fatal("List of all preset events is empty")
	}
	eventList, err = EnumEvents(PRESET_MASK, PRESET_ENUM_AVAIL)
	if err != nil {
		t.Fatal(err)
	}
	if len(eventList) == 0 {
		t.Fatal("List of available preset events is empty")
	}
	if len(eventList) > numPresets {
		t.Fatal("More preset events are available than exist in toto")
	}

	// Look for native events.
	eventList, err = EnumEvents(NATIVE_MASK|ComponentMask(0), ENUM_EVENTS)
	if err != nil {
		t.Fatal(err)
	}
	var info ComponentInfo
	if info, err = GetComponentInfo(0); err != nil {
		t.Fatal(err)
	}
	if len(eventList) != info.NumNativeEvents {
		t.Fatalf("Expected to see %d native events but saw %d native events",
			info.NumNativeEvents, len(eventList))
	}
}
