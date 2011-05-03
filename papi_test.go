// Copyright (C) 2011, Los Alamos National Security, LLC.
// Use of this source code is governed by a BSD-style license.

// This file tests miscellaneous features of the PAPI library.

package papi

import "testing"

// Prevent SomeValue from being optimized away by exporting it.
var SomeValue float64 = 123.456


// Peform a given number of floating-point operations, just to burn cycle.
func performWork(numFlops int) {
	for i := 0; i < numFlops; i++ {
		// This should do one floating-point instruction and
		// one floating-point operation per iteration and be
		// unlikely for the compiler to optimize away.
		SomeValue = SomeValue * float64(i%7)
	}
}


// Ensure that we can map back-and-forth between event names and event codes.
func TestEventNames(t *testing.T) {
	eventCodes := []Event{
		BR_INS,
		CA_INV,
		CSR_SUC,
		FMA_INS,
		HW_INT,
		L1_DCA,
		RES_STL,
		TLB_TL,
	}
	for _, ecode1 := range eventCodes {
		ename := ecode1.String()
		ecode2, err := StringToEvent(ename)
		if err != nil {
			t.Fatal(err)
		}
		if ecode1 != ecode2 {
			t.Fatalf("Event code got mangled: %d --> %s --> %d",
				int(ecode1), ename, int(ecode2))
		}
	}
}


// Ensure that we can map event modifiers to strings.
func TestEventModifiers(t *testing.T) {
	// papi.go makes some assumptions about a couple of PAPI's
	// event-modifier values.  Verify that these assumptions are
	// correct.
	if int(ENUM_EVENTS) != 0 {
		t.Fatalf("Expected ENUM_EVENTS to be 0, but it's actually %d", int(ENUM_EVENTS))
	}
	if int(ENUM_FIRST) != 1 {
		t.Fatalf("Expected ENUM_FIRST to be 1, but it's actually %d", int(ENUM_FIRST))
	}
	if int(PRESET_ENUM_AVAIL) != 2 {
		t.Fatalf("Expected PRESET_ENUM_AVAIL to be 2, but it's actually %d", int(PRESET_ENUM_AVAIL))
	}

	// Try mapping a few names to strings and verifying that we
	// get what we expected to get.
	expectedToActual := map[EventModifier]string{
		ENUM_EVENTS:                     "PAPI_ENUM_EVENTS",
		ENUM_FIRST:                      "PAPI_ENUM_FIRST",
		PRESET_ENUM_AVAIL:               "PAPI_PRESET_ENUM_AVAIL",
		PRESET_BIT_CACH | PRESET_BIT_L3: "PAPI_PRESET_BIT_CACH|PAPI_PRESET_BIT_L3"}
	for emod, str := range expectedToActual {
		if emod.String() != str {
			t.Fatalf("Expected to map %d to \"%s\" but instead got \"%s\"",
				int32(emod), str, emod.String())
		}
	}
}
