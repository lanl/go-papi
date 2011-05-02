// Copyright (C) 2011, Los Alamos National Security, LLC.
// Use of this source code is governed by a BSD-style license.

// This file tests miscellaneous features of the PAPI library.

package papi

import "testing"


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
