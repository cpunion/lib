//go:build windows

// Copyright (c) 2026 The GoPlus Authors. Licensed under the Apache License 2.0.

package pthread

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
)

const (
	errInvalid      = c.Int(87)
	errNotSupported = c.Int(50)
)

//llgo:type C
type RoutineFunc func(c.Pointer) c.Pointer

//llgo:type C
type KeyDestructor func(c.Pointer)

// Thread is an opaque handle to a joinable Windows thread sidecar.
type Thread c.Pointer

//go:linkname create C.llgo_pthread_create
func create(thread *Thread, stackSize uintptr, routine RoutineFunc, arg c.Pointer) c.Int

//go:linkname join C.llgo_pthread_join
func join(thread Thread, retval *c.Pointer) c.Int

//go:linkname threadExit C.llgo_pthread_exit
func threadExit(retval c.Pointer)

func Create(thread *Thread, attr *Attr, routine RoutineFunc, arg c.Pointer) c.Int {
	if thread == nil || routine == nil {
		return errInvalid
	}
	var stackSize uintptr
	if attr != nil {
		if attr.detached != 0 || attr.stackAddr != nil {
			return errNotSupported
		}
		stackSize = attr.stackSize
	}
	return create(thread, stackSize, routine, arg)
}

func Join(thread Thread, retval *c.Pointer) c.Int {
	return join(thread, retval)
}

func Exit(retval c.Pointer) {
	threadExit(retval)
}

// Windows has no safe asynchronous pthread cancellation equivalent.
func Cancel(Thread) c.Int { return errNotSupported }

type Attr struct {
	detached  byte
	stackAddr c.Pointer
	stackSize uintptr
}

func (attr *Attr) Init() c.Int {
	*attr = Attr{}
	return 0
}

func (attr *Attr) Destroy() c.Int {
	*attr = Attr{}
	return 0
}

func (attr *Attr) GetDetached(detached *c.Int) c.Int {
	if detached == nil {
		return errInvalid
	}
	*detached = c.Int(attr.detached)
	return 0
}

func (attr *Attr) SetDetached(detached c.Int) c.Int {
	if detached != 0 {
		return errNotSupported
	}
	attr.detached = byte(detached)
	return 0
}

func (attr *Attr) GetStackSize(stackSize *uintptr) c.Int {
	if stackSize == nil {
		return errInvalid
	}
	*stackSize = attr.stackSize
	return 0
}

func (attr *Attr) SetStackSize(stackSize uintptr) c.Int {
	attr.stackSize = stackSize
	return 0
}

func (attr *Attr) GetStackAddr(stackAddr *c.Pointer) c.Int {
	if stackAddr == nil {
		return errInvalid
	}
	*stackAddr = attr.stackAddr
	return 0
}

func (attr *Attr) SetStackAddr(stackAddr c.Pointer) c.Int {
	if stackAddr != nil {
		return errNotSupported
	}
	attr.stackAddr = stackAddr
	return 0
}

type Key struct {
	index      c.Uint
	destructor KeyDestructor
}

const invalidKey = ^c.Uint(0)

//go:linkname keyCreate C.llgo_pthread_key_create
func keyCreate(index *c.Uint) c.Int

//go:linkname keyDelete C.llgo_pthread_key_delete
func keyDelete(index c.Uint) c.Int

//go:linkname keyGet C.llgo_pthread_getspecific
func keyGet(index c.Uint) c.Pointer

//go:linkname keySet C.llgo_pthread_setspecific
func keySet(index c.Uint, destructor KeyDestructor, value c.Pointer) c.Int

func (key *Key) Create(destructor func(c.Pointer)) c.Int {
	if key == nil {
		return errInvalid
	}
	key.index = invalidKey
	if result := keyCreate(&key.index); result != 0 {
		return result
	}
	key.destructor = KeyDestructor(destructor)
	return 0
}

func (key Key) Delete() c.Int {
	if key.index == invalidKey {
		return 0
	}
	return keyDelete(key.index)
}

func (key Key) Get() c.Pointer {
	if key.index == invalidKey {
		return nil
	}
	return keyGet(key.index)
}

func (key Key) Set(value c.Pointer) c.Int {
	if key.index == invalidKey {
		return errInvalid
	}
	return keySet(key.index, key.destructor, value)
}
