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
