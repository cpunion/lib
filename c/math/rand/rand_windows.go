//go:build windows

// Copyright (c) 2026 The GoPlus Authors. Licensed under the Apache License 2.0.

package rand

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
)

const (
	LLGoPackage = "link"
	LLGoFiles   = "_wrap/rand_windows.c"
)

//go:linkname RandR C.llgo_rand_r
func RandR(*c.Uint) c.Int

//go:linkname Sranddev C.llgo_sranddev
func Sranddev()

// MSVC has rand/srand rather than the BSD random family. Windows long and int
// have the same width, so these retain the public API's result type.
//
//go:linkname Random C.rand
func Random() c.Long

//go:linkname Srandom C.srand
func Srandom(c.Uint)

//go:linkname Srandomdev C.llgo_sranddev
func Srandomdev()
