//go:build windows

/*
 * Copyright (c) 2026 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package time

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
)

const (
	LLGoPackage = "link"
	LLGoFiles   = "_wrap/time_windows.c"
)

type TimeT int64

// Tm matches the Universal CRT struct tm. The Unix-only tm_gmtoff and tm_zone
// extensions are intentionally absent on Windows.
type Tm struct {
	Sec   c.Int
	Min   c.Int
	Hour  c.Int
	Mday  c.Int
	Mon   c.Int
	Year  c.Int
	Wday  c.Int
	Yday  c.Int
	Isdst c.Int
}

//go:linkname Time C.llgo_time
func Time(timer *TimeT) TimeT

//go:linkname Mktime C.llgo_mktime
func Mktime(timer *Tm) TimeT

//go:linkname Ctime C.llgo_ctime
func Ctime(timer *TimeT) string

//go:linkname Difftime C.llgo_difftime
func Difftime(end, start TimeT) float64

//go:linkname Gmtime C.llgo_gmtime
func Gmtime(timer *TimeT) *Tm

//go:linkname Localtime C.llgo_localtime
func Localtime(timer *TimeT) *Tm

//go:linkname Strftime C.llgo_strftime
func Strftime(buf *c.Char, bufSize uintptr, format *c.Char, timeptr *Tm) uintptr

type ClockT c.Long

//go:linkname Clock C.llgo_clock
func Clock() ClockT

type ClockidT c.Int

const (
	CLOCK_REALTIME           ClockidT = 0
	CLOCK_MONOTONIC          ClockidT = 1
	CLOCK_PROCESS_CPUTIME_ID ClockidT = 2
	CLOCK_THREAD_CPUTIME_ID  ClockidT = 3
	CLOCK_MONOTONIC_RAW      ClockidT = 4
)

type Timespec struct {
	Sec  TimeT
	Nsec c.Long
}

//go:linkname ClockGettime C.llgo_clock_gettime
func ClockGettime(clkID ClockidT, tp *Timespec) c.Int

//go:linkname ClockSettime C.llgo_clock_settime
func ClockSettime(clkID ClockidT, tp *Timespec) c.Int

//go:linkname ClockGetres C.llgo_clock_getres
func ClockGetres(clkID ClockidT, res *Timespec) c.Int
