//go:build windows

/*
 * Copyright (c) 2026 The GoPlus Authors (goplus.org). All rights reserved.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package os

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
)

const (
	LLGoPackage = "link"
	LLGoFiles   = "_os/os_windows.c"
	PATH_MAX    = 260
)

type (
	ModeT c.Int
	UidT  c.Uint
	GidT  c.Uint
	OffT  int64
	DevT  c.Uint
)

const (
	F_GETFL = 3
	F_SETFL = 4

	O_RDONLY   = 0x0000
	O_WRONLY   = 0x0001
	O_RDWR     = 0x0002
	O_ACCMODE  = 0x0003
	O_APPEND   = 0x0008
	O_CREAT    = 0x0100
	O_TRUNC    = 0x0200
	O_EXCL     = 0x0400
	O_TEXT     = 0x4000
	O_BINARY   = 0x8000
	O_NONBLOCK = 0

	EAGAIN = 11
)

//go:linkname Clearenv C.llgo_os_clearenv
func Clearenv() c.Int

//go:linkname Getcwd C._getcwd
func Getcwd(buffer c.Pointer, size uintptr) *c.Char

//go:linkname Open C.llgo_os_open
func Open(path *c.Char, flags c.Int, __llgo_va_list ...any) c.Int

// Fcntl reports ENOSYS because the Universal CRT has no fcntl equivalent.
//
//go:linkname Fcntl C.llgo_os_fcntl
func Fcntl(fd c.Int, cmd c.Int, __llgo_va_list ...any) c.Int

//go:linkname Close C._close
func Close(fd c.Int) c.Int

//go:linkname Read C.llgo_os_read
func Read(fd c.Int, buf c.Pointer, count uintptr) int

//go:linkname Write C.llgo_os_write
func Write(fd c.Int, buf c.Pointer, count uintptr) int

//go:linkname Execlp C._execlp
func Execlp(file *c.Char, arg0 *c.Char, __llgo_va_list ...any) c.Int
