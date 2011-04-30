// Copyright (C) 2011, Los Alamos National Security, LLC.
// Use of this source code is governed by a BSD-style license.

package papi

/*
This file was generated semi-automatically from papi.h.
*/

// #include <papi.h>
import "C"

// Define possible attributes of a single level of the memory hierarchy.
const (
	// Cache type -- test with CacheType()
	MH_TYPE_EMPTY   MHAttrs = C.PAPI_MH_TYPE_EMPTY   // Not a cache
	MH_TYPE_INST    MHAttrs = C.PAPI_MH_TYPE_INST    // Instruction cache
	MH_TYPE_DATA    MHAttrs = C.PAPI_MH_TYPE_DATA    // Data cache
	MH_TYPE_VECTOR  MHAttrs = C.PAPI_MH_TYPE_VECTOR  // Vector cache
	MH_TYPE_TRACE   MHAttrs = C.PAPI_MH_TYPE_TRACE   // Trace cache
	MH_TYPE_UNIFIED MHAttrs = C.PAPI_MH_TYPE_UNIFIED // Unified instruction+data cache

	// Write policy -- test with CacheWritePolicy()
	MH_TYPE_WT MHAttrs = C.PAPI_MH_TYPE_WT // Write-through cache
	MH_TYPE_WB MHAttrs = C.PAPI_MH_TYPE_WB // Write-back cache

	// Replacement policy -- test with CacheReplacementPolicy()
	MH_TYPE_UNKNOWN    MHAttrs = C.PAPI_MH_TYPE_UNKNOWN    // Unknown replacement policy
	MH_TYPE_LRU        MHAttrs = C.PAPI_MH_TYPE_LRU        // LRU replacement policy
	MH_TYPE_PSEUDO_LRU MHAttrs = C.PAPI_MH_TYPE_PSEUDO_LRU // Pseudo-LRU replacement policy

	// TLB, prefetch buffer, or cache -- test with CacheUsage()
	MH_TYPE_TLB  MHAttrs = C.PAPI_MH_TYPE_TLB  // TLB, not memory cache
	MH_TYPE_PREF MHAttrs = C.PAPI_MH_TYPE_PREF // Prefetch buffer
)
