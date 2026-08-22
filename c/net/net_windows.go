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

package net

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
)

const (
	LLGoPackage = "link: -lws2_32"
	LLGoFiles   = "_wrap/net_windows.c"
)

type (
	SocketT  = uintptr
	SocklenT = c.Int
)

const InvalidSocket SocketT = ^SocketT(0)

const (
	AF_UNSPEC    = 0
	AF_UNIX      = 1
	AF_LOCAL     = AF_UNIX
	AF_INET      = 2
	AF_IPX       = 6
	AF_APPLETALK = 16
	AF_NETBIOS   = 17
	AF_INET6     = 23
	AF_IRDA      = 26
	AF_BTH       = 32
	AF_MAX       = 34
)

const (
	SOCK_STREAM    = 1
	SOCK_DGRAM     = 2
	SOCK_RAW       = 3
	SOCK_RDM       = 4
	SOCK_SEQPACKET = 5
)

const (
	EAI_AGAIN    = 11002
	EAI_BADFLAGS = 10022
	EAI_FAIL     = 11003
	EAI_FAMILY   = 10047
	EAI_MEMORY   = 8
	EAI_NONAME   = 11001
	EAI_SERVICE  = 10109
	EAI_SOCKTYPE = 10044
)

const INET_ADDRSTRLEN = 16

// SockaddrIn matches the Winsock SOCKADDR_IN layout. Unlike BSD sockaddr_in,
// it has no leading length byte and its address family is 16 bits.
type SockaddrIn struct {
	Family uint16
	Port   uint16
	Addr   InAddr
	Zero   [8]c.Char
}

type SockaddrIn6 struct {
	Family   uint16
	Port     uint16
	Flowinfo c.Uint
	Addr     In6Addr
	ScopeId  c.Uint
}

type SockaddrStorage struct {
	Family uint16
	pad1   [6]c.Char
	align  c.LongLong
	pad2   [112]c.Char
}

type InAddr struct {
	Addr c.Uint
}

type In6Addr struct {
	U6Addr [16]uint8
}

type SockAddr struct {
	Family uint16
	Data   [14]c.Char
}

type Hostent struct {
	Name     *c.Char
	Aliases  **c.Char
	AddrType int16
	Length   int16
	AddrList **c.Char
}

//go:linkname Socket C.llgo_net_socket
func Socket(domain c.Int, typ c.Int, protocol c.Int) SocketT

//go:linkname Bind C.llgo_net_bind
func Bind(sockfd SocketT, addr *SockaddrIn, addrlen SocklenT) c.Int

//go:linkname Connect C.llgo_net_connect
func Connect(sockfd SocketT, addr *SockAddr, addrlen SocklenT) c.Int

//go:linkname Listen C.llgo_net_listen
func Listen(sockfd SocketT, backlog c.Int) c.Int

//go:linkname Accept C.llgo_net_accept
func Accept(sockfd SocketT, addr *SockaddrIn, addrlen *SocklenT) SocketT

//go:linkname Close C.llgo_net_close
func Close(sockfd SocketT) c.Int

//go:linkname GetHostByName C.llgo_net_gethostbyname
func GetHostByName(name *c.Char) *Hostent

//go:linkname InetNtop C.llgo_net_inet_ntop
func InetNtop(af c.Int, src c.Pointer, dst *c.Char, size uintptr) *c.Char

//go:linkname InetAddr C.llgo_net_inet_addr
func InetAddr(value *c.Char) c.Uint

//go:linkname Send C.llgo_net_send
func Send(sockfd SocketT, buffer c.Pointer, length uintptr, flags c.Int) c.Long

//go:linkname Recv C.llgo_net_recv
func Recv(sockfd SocketT, buffer c.Pointer, length uintptr, flags c.Int) c.Long

//go:linkname SetSockOpt C.llgo_net_setsockopt
func SetSockOpt(sockfd SocketT, level c.Int, optionName c.Int, optionValue c.Pointer, optionLength SocklenT) c.Int

//go:linkname Ntohs C.ntohs
func Ntohs(value uint16) uint16

//go:linkname Htons C.htons
func Htons(value uint16) uint16

//go:linkname Ntohl C.ntohl
func Ntohl(value c.Uint) c.Uint

//go:linkname Htonl C.htonl
func Htonl(value c.Uint) c.Uint

//go:linkname LastError C.WSAGetLastError
func LastError() c.Int
