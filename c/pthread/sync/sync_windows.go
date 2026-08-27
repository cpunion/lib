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

package sync

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
	ctime "github.com/goplus/lib/c/time"
)

const (
	LLGoPackage = "link"
	LLGoFiles   = "_wrap/sync_windows.c"
	errInvalid  = c.Int(22)
)

// Once uses the zero-initialized Windows INIT_ONCE representation.
type Once struct{ state uintptr }

var OnceInit Once

//go:linkname winOnce C.llgo_win_once
func winOnce(once *Once, f OnceFunc) c.Int

func (o *Once) Do(f OnceFunc) c.Int { return winOnce(o, f) }

type MutexType c.Int

const (
	MUTEX_NORMAL MutexType = iota
	MUTEX_ERRORCHECK
	MUTEX_RECURSIVE
	MUTEX_DEFAULT = MUTEX_NORMAL
)

type MutexAttr struct{ typ MutexType }

func (a *MutexAttr) Init(_ *MutexAttr) c.Int { a.typ = MUTEX_DEFAULT; return 0 }
func (a *MutexAttr) Destroy()                {}
func (a *MutexAttr) SetType(typ MutexType) c.Int {
	if typ != MUTEX_NORMAL && typ != MUTEX_DEFAULT {
		return errInvalid
	}
	a.typ = typ
	return 0
}

// Mutex uses the zero-initialized Windows SRWLOCK representation. Windows
// SRWLOCK does not provide recursive or error-checking mutex modes.
type Mutex struct{ state uintptr }

//go:linkname winMutexLock C.llgo_win_mutex_lock
func winMutexLock(*Mutex)

//go:linkname winMutexUnlock C.llgo_win_mutex_unlock
func winMutexUnlock(*Mutex)

//go:linkname winMutexTryLock C.llgo_win_mutex_trylock
func winMutexTryLock(*Mutex) c.Int

func (m *Mutex) Init(attr *MutexAttr) c.Int {
	if attr != nil && attr.typ != MUTEX_NORMAL && attr.typ != MUTEX_DEFAULT {
		return errInvalid
	}
	m.state = 0
	return 0
}
func (m *Mutex) Destroy()       {}
func (m *Mutex) TryLock() c.Int { return winMutexTryLock(m) }
func (m *Mutex) Lock()          { winMutexLock(m) }
func (m *Mutex) Unlock()        { winMutexUnlock(m) }

type RWLockAttr struct{}

func (*RWLockAttr) Init(*RWLockAttr) c.Int { return 0 }
func (*RWLockAttr) Destroy()               {}
func (*RWLockAttr) SetPShared(value c.Int) c.Int {
	if value != 0 {
		return errInvalid
	}
	return 0
}
func (*RWLockAttr) GetPShared(value *c.Int) c.Int {
	if value == nil {
		return errInvalid
	}
	*value = 0
	return 0
}

// RWLock uses the zero-initialized Windows SRWLOCK representation.
type RWLock struct{ state uintptr }

//go:linkname winRWLockRLock C.llgo_win_rwlock_rlock
func winRWLockRLock(*RWLock)

//go:linkname winRWLockTryRLock C.llgo_win_rwlock_tryrlock
func winRWLockTryRLock(*RWLock) c.Int

//go:linkname winRWLockRUnlock C.llgo_win_rwlock_runlock
func winRWLockRUnlock(*RWLock)

//go:linkname winRWLockLock C.llgo_win_rwlock_lock
func winRWLockLock(*RWLock)

//go:linkname winRWLockTryLock C.llgo_win_rwlock_trylock
func winRWLockTryLock(*RWLock) c.Int

//go:linkname winRWLockUnlock C.llgo_win_rwlock_unlock
func winRWLockUnlock(*RWLock)

func (rw *RWLock) Init(*RWLockAttr) c.Int { rw.state = 0; return 0 }
func (*RWLock) Destroy()                  {}
func (rw *RWLock) RLock()                 { winRWLockRLock(rw) }
func (rw *RWLock) TryRLock() c.Int        { return winRWLockTryRLock(rw) }
func (rw *RWLock) RUnlock()               { winRWLockRUnlock(rw) }
func (rw *RWLock) Lock()                  { winRWLockLock(rw) }
func (rw *RWLock) TryLock() c.Int         { return winRWLockTryLock(rw) }
func (rw *RWLock) Unlock()                { winRWLockUnlock(rw) }

type CondAttr struct{ clock ctime.ClockidT }

func (a *CondAttr) Init(*CondAttr) c.Int { a.clock = ctime.CLOCK_REALTIME; return 0 }
func (*CondAttr) Destroy()               {}
func (a *CondAttr) SetClock(clock ctime.ClockidT) c.Int {
	if clock != ctime.CLOCK_REALTIME {
		return errInvalid
	}
	a.clock = clock
	return 0
}
func (a *CondAttr) GetClock(clock *ctime.ClockidT) c.Int {
	if clock == nil {
		return errInvalid
	}
	*clock = a.clock
	return 0
}

// Cond uses the zero-initialized Windows CONDITION_VARIABLE representation.
type Cond struct{ state uintptr }

//go:linkname winCondSignal C.llgo_win_cond_signal
func winCondSignal(*Cond) c.Int

//go:linkname winCondBroadcast C.llgo_win_cond_broadcast
func winCondBroadcast(*Cond) c.Int

//go:linkname winCondWait C.llgo_win_cond_wait
func winCondWait(*Cond, *Mutex) c.Int

//go:linkname winCondTimedWait C.llgo_win_cond_timedwait
func winCondTimedWait(*Cond, *Mutex, *ctime.Timespec) c.Int

func (cond *Cond) Init(*CondAttr) c.Int { cond.state = 0; return 0 }
func (*Cond) Destroy()                  {}
func (cond *Cond) Signal() c.Int        { return winCondSignal(cond) }
func (cond *Cond) Broadcast() c.Int     { return winCondBroadcast(cond) }
func (cond *Cond) Wait(m *Mutex) c.Int  { return winCondWait(cond, m) }
func (cond *Cond) TimedWait(m *Mutex, deadline *ctime.Timespec) c.Int {
	return winCondTimedWait(cond, m, deadline)
}

// Sem stores the Windows semaphore handle and its observable count. The handle
// is closed by Destroy; unlike the native SRW types, Sem must be initialized.
type Sem struct {
	handle uintptr
	count  c.Long
}

//go:linkname winSemInit C.llgo_libc_sem_init
func winSemInit(*Sem, c.Int, c.Uint) c.Int

//go:linkname winSemDestroy C.llgo_libc_sem_destroy
func winSemDestroy(*Sem) c.Int

//go:linkname winSemPost C.llgo_libc_sem_post
func winSemPost(*Sem) c.Int

//go:linkname winSemWait C.llgo_libc_sem_wait
func winSemWait(*Sem, c.Int) c.Int

//go:linkname winSemValue C.llgo_libc_sem_value
func winSemValue(*Sem, *c.Int) c.Int

func (s *Sem) Init(shared c.Int, value c.Uint) c.Int { return winSemInit(s, shared, value) }
func (s *Sem) Destroy() c.Int                        { return winSemDestroy(s) }
func (s *Sem) Post() c.Int                           { return winSemPost(s) }
func (s *Sem) Wait() c.Int                           { return winSemWait(s, 0) }
func (s *Sem) TryWait() c.Int                        { return winSemWait(s, 1) }
func (s *Sem) GetValue(value *c.Int) c.Int           { return winSemValue(s, value) }
