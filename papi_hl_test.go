// Copyright (C) 2011, Los Alamos National Security, LLC.
// Use of this source code is governed by a BSD-style license.

// This file tests the high-level features of the PAPI library.

package papi

import "testing"
import "time"


// Do most of the work for TestFlipFlops() by testing any of Flips(),
// Flops(), or Ipc() high-level PAPI functions in a consistent manner.
func testFlipFlopsHelper(t *testing.T, funcName string, hlFunc func() (float32, float32, int64, float32, Errno)) {
	const sleep_usecs = 10000
	const flops = 100
	var someValue float64 = 123.456
	counterValues := make([]int64, 3)

	// Test Flips().
	rtime1, ptime1, other1, _, err := hlFunc()
	if err != OK {
		t.Fatal(err)
	}
	for i := 0; i < flops; i++ {
		// This should do one floating-point instruction and
		// one floating-point operation per iteration and be
		// unlikely for the compiler to optimize away.
		someValue = someValue * float64(i%7)
	}
	time.Sleep(sleep_usecs * 1000)
	rtime2, ptime2, other2, _, err := hlFunc()
	if err != OK {
		t.Fatal(err)
	}
	if rtime2-rtime1 < sleep_usecs/1.0e6 {
		t.Fatalf("%s() real time increases too slowly: %f vs. %f",
			funcName, rtime2-rtime1, sleep_usecs/1.0e6)
	}
	if ptime1 > ptime2 {
		t.Fatalf("%s() process time decreases: %f >= %f",
			funcName, ptime1, ptime2)
	}
	if (other2 - other1) < flops {
		t.Fatalf("%s() observed too few counts: %d >= %d",
			other2-other1, flops)
	}
	if err := StopCounters(counterValues); err != OK {
		t.Fatal(err)
	}
}


// Ensure that the time, floating-point, and instruction counters
// return believable values.
func TestFlipFlops(t *testing.T) {
	testFlipFlopsHelper(t, "Flips", Flips)
	testFlipFlopsHelper(t, "Flops", Flops)
	testFlipFlopsHelper(t, "Ipc", Ipc)
}


// Ensure that the high-level counters actually count something.
func TestHLCounters(t *testing.T) {
	// Start counting a few events (but not more than NumCounters).
	eventList := []Event{LD_INS, SR_INS, TOT_CYC, TOT_INS}
	var usedEvents []Event
	if len(eventList) <= NumCounters {
		usedEvents = eventList
	} else {
		usedEvents = eventList[0 : NumCounters-1]
	}
	if err := StartCounters(usedEvents); err != OK {
		t.Fatal(err)
	}

	// Read the counters right away.
	counterValues := make([]int64, len(usedEvents))
	if err := ReadCounters(counterValues); err != OK {
		t.Fatal(err)
	}

	// Burn some cycles.
	eventNames := make(map[Event]string)
	for _, ecode := range eventList {
		eventNames[ecode] = ecode.String()
	}

	// Ensure that at least one of our counters increased.  All
	// should increase, but perhaps some don't on certain hardware.
	if err := AccumCounters(counterValues); err != OK {
		t.Fatal(err)
	}
	anyChanged := false
	for _, value := range counterValues {
		if value != 0 {
			anyChanged = true
			break
		}
	}
	if !anyChanged {
		t.Fatalf("None of %v appear to count", usedEvents)
	}

	// Burn some more cycles.
	for _, ecode := range eventList {
		eventNames[ecode] = eventNames[ecode] + " event"
	}

	// Ensure that at least one of our counters increased.  All
	// should increase, but perhaps some don't on certain hardware.
	counterValues = make([]int64, len(counterValues)) // Clear all counters
	if err := StopCounters(counterValues); err != OK {
		t.Fatal(err)
	}
	anyChanged = false
	for _, value := range counterValues {
		if value != 0 {
			anyChanged = true
			break
		}
	}
	if !anyChanged {
		t.Fatalf("None of %v appear to count", usedEvents)
	}
}
