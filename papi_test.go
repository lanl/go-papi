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


// Ensure that the time and floating-point functions return believable
// values.
func TestFlipFlops(t *testing.T) {
	const sleep_usecs = 10000
	const flops = 100
	var someValue float64 = 123.456

	// Test Flips().
	rtime1, ptime1, flpins1, _, err := Flips()
	if err != OK {
		t.Fatal(err)
	}
	for i := 0; i < flops/2; i++ {
		someValue = someValue*float64(i) + float64(flops-i)
	}
	time.Sleep(sleep_usecs * 1000)
	rtime2, ptime2, flpins2, _, err := Flips()
	if err != OK {
		t.Fatal(err)
	}
	if rtime2-rtime1 < sleep_usecs/1.0e6 {
		t.Fatalf("Flips() real time increases too slowly: %f vs. %f",
			rtime2-rtime1, sleep_usecs/1.0e6)
	}
	if ptime1 > ptime2 {
		t.Fatalf("Flips() process time decreases: %f >= %f",
			ptime1, ptime2)
	}
	if (flpins2 - flpins1) < flops {
		t.Fatalf("Flips() observed too few flips: %d >= %d",
			flpins2-flpins1, flops)
	}

	// Test Flops().
	/*
		rtime1, ptime1, flpops1, _, err := Flops()
		if err != OK {
			t.Fatal(err)
		}
		for i := 0; i<flops/2; i++ {
			someValue = someValue * float64(i) + float64(flops-i)
		}
		time.Sleep(sleep_usecs * 1000)
		rtime2, ptime2, flpops2, _, err := Flops()
		if err != OK {
			t.Fatal(err)
		}
		if rtime2 - rtime1 < sleep_usecs/1.0e6 {
			t.Fatalf("Flops() real time increases too slowly: %f vs. %f",
				rtime2-rtime1, sleep_usecs/1.0e6)
		}
		if ptime1 > ptime2 {
			t.Fatalf("Flops() process time decreases: %f >= %f",
				ptime1, ptime2)
		}
		if (flpops2 - flpops1) < flops {
			t.Fatalf("Flops() observed too few flops: %d >= %d",
				flpops2-flpops1, flops)
		}
	*/
}
