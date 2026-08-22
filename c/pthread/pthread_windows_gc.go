//go:build windows && !nogc

// Copyright (c) 2026 The GoPlus Authors. Licensed under the Apache License 2.0.

package pthread

const (
	LLGoPackage = "link"
	LLGoFiles   = "_wrap/pthread_windows_gc.c"
)
