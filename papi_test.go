// See if the PAPI package works.
// By Scott Pakin <pakin@lanl.gov>.

package papi

import (
	"testing"
	"time"
)


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
		if err != OK {
			t.Fatal(err)
		}
		if ecode1 != ecode2 {
			t.Fatalf("Event code got mangled: %d --> %s --> %d",
				int(ecode1), ename, int(ecode2))
		}
	}
}


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
