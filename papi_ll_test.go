// Copyright (C) 2011, Los Alamos National Security, LLC.
// Use of this source code is governed by a BSD-style license.

// This file tests the low-level features of the PAPI library.

package papi

import "testing"
import "time"
import "os"


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


// Test various events in an EventSet's lifetime.  This test is
// derived from examples/PAPI_add_remove_events.c in the PAPI
// distribution.
func TestEventSet(t *testing.T) {
	const flops = 1000

	// Start counting a few events.
	var err os.Error
	var events EventSet
	if events, err = CreateEventSet(); err != nil {
		t.Fatal(err)
	}
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
