//go:build !windows

// Copyright (c) 2026 The GoPlus Authors. Licensed under the Apache License 2.0.

package net

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
)

const LLGoPackage = true

type AddrInfo struct {
	Flags     c.Int
	Family    c.Int
	SockType  c.Int
	Protocol  c.Int
	AddrLen   c.Uint
	CanOnName *c.Char
	Addr      *SockAddr
	Next      *AddrInfo
}

//go:linkname Getaddrinfo C.getaddrinfo
func Getaddrinfo(host *c.Char, port *c.Char, addrInfo *AddrInfo, result **AddrInfo) c.Int

//go:linkname Freeaddrinfo C.freeaddrinfo
func Freeaddrinfo(addrInfo *AddrInfo)
