/* Copyright (c) 2026 The GoPlus Authors. Licensed under Apache-2.0. */

#define WIN32_LEAN_AND_MEAN
#include <limits.h>
#include <stdint.h>
#include <windows.h>

enum { LLGO_BUSY = 16, LLGO_INVALID = 22, LLGO_AGAIN = 11, LLGO_TIMEDOUT = 110 };

typedef struct { void *code; void *context; } llgo_go_func;
extern void llgo_win_once_invoke(llgo_go_func *fn);

static BOOL CALLBACK llgo_once_callback(PINIT_ONCE once, PVOID parameter, PVOID *context)
{
    (void)once; (void)context;
    llgo_win_once_invoke((llgo_go_func *)parameter);
    return TRUE;
}

int llgo_win_once(PINIT_ONCE once, llgo_go_func *fn)
{
    if (fn == NULL || fn->code == NULL) return LLGO_INVALID;
    return InitOnceExecuteOnce(once, llgo_once_callback, fn, NULL) ? 0 : (int)GetLastError();
}

void llgo_win_mutex_lock(PSRWLOCK lock) { AcquireSRWLockExclusive(lock); }
void llgo_win_mutex_unlock(PSRWLOCK lock) { ReleaseSRWLockExclusive(lock); }
int llgo_win_mutex_trylock(PSRWLOCK lock) { return TryAcquireSRWLockExclusive(lock) ? 0 : LLGO_BUSY; }
void llgo_win_rwlock_rlock(PSRWLOCK lock) { AcquireSRWLockShared(lock); }
int llgo_win_rwlock_tryrlock(PSRWLOCK lock) { return TryAcquireSRWLockShared(lock) ? 0 : LLGO_BUSY; }
void llgo_win_rwlock_runlock(PSRWLOCK lock) { ReleaseSRWLockShared(lock); }
void llgo_win_rwlock_lock(PSRWLOCK lock) { AcquireSRWLockExclusive(lock); }
int llgo_win_rwlock_trylock(PSRWLOCK lock) { return TryAcquireSRWLockExclusive(lock) ? 0 : LLGO_BUSY; }
void llgo_win_rwlock_unlock(PSRWLOCK lock) { ReleaseSRWLockExclusive(lock); }

int llgo_win_cond_signal(PCONDITION_VARIABLE cond) { WakeConditionVariable(cond); return 0; }
int llgo_win_cond_broadcast(PCONDITION_VARIABLE cond) { WakeAllConditionVariable(cond); return 0; }
int llgo_win_cond_wait(PCONDITION_VARIABLE cond, PSRWLOCK lock)
{
    return SleepConditionVariableSRW(cond, lock, INFINITE, 0) ? 0 : (int)GetLastError();
}

typedef struct { int64_t sec; long nsec; } llgo_timespec;
static uint64_t llgo_unix_time_100ns(void)
{
    FILETIME value;
    ULARGE_INTEGER ticks;
    GetSystemTimeAsFileTime(&value);
    ticks.LowPart = value.dwLowDateTime;
    ticks.HighPart = value.dwHighDateTime;
    return ticks.QuadPart - 116444736000000000ULL;
}

int llgo_win_cond_timedwait(PCONDITION_VARIABLE cond, PSRWLOCK lock, const llgo_timespec *deadline)
{
    if (deadline == NULL || deadline->sec < 0 || deadline->nsec < 0 || deadline->nsec >= 1000000000L)
        return LLGO_INVALID;
    uint64_t target = (uint64_t)deadline->sec * 10000000ULL + (uint64_t)deadline->nsec / 100ULL;
    uint64_t now = llgo_unix_time_100ns();
    if (target <= now) return LLGO_TIMEDOUT;
    uint64_t millis = (target - now + 9999ULL) / 10000ULL;
    DWORD timeout = millis >= INFINITE ? INFINITE - 1 : (DWORD)millis;
    if (SleepConditionVariableSRW(cond, lock, timeout, 0)) return 0;
    DWORD error = GetLastError();
    return error == ERROR_TIMEOUT ? LLGO_TIMEDOUT : (int)error;
}

typedef struct { HANDLE handle; volatile LONG count; } llgo_sem;

int llgo_win_sem_init(llgo_sem **out, int shared, unsigned int value)
{
    if (out == NULL || shared != 0 || value > LONG_MAX) return LLGO_INVALID;
    llgo_sem *sem = (llgo_sem *)HeapAlloc(GetProcessHeap(), 0, sizeof(*sem));
    if (sem == NULL) return ERROR_NOT_ENOUGH_MEMORY;
    sem->handle = CreateSemaphoreW(NULL, (LONG)value, LONG_MAX, NULL);
    if (sem->handle == NULL) {
        int error = (int)GetLastError();
        HeapFree(GetProcessHeap(), 0, sem);
        return error;
    }
    sem->count = (LONG)value;
    *out = sem;
    return 0;
}

int llgo_win_sem_destroy(llgo_sem **slot)
{
    if (slot == NULL || *slot == NULL) return LLGO_INVALID;
    llgo_sem *sem = *slot;
    if (!CloseHandle(sem->handle)) return (int)GetLastError();
    HeapFree(GetProcessHeap(), 0, sem);
    *slot = NULL;
    return 0;
}

int llgo_win_sem_post(llgo_sem **slot)
{
    if (slot == NULL || *slot == NULL) return LLGO_INVALID;
    if (!ReleaseSemaphore((*slot)->handle, 1, NULL)) return (int)GetLastError();
    InterlockedIncrement(&(*slot)->count);
    return 0;
}

int llgo_win_sem_wait(llgo_sem **slot, int try_only)
{
    if (slot == NULL || *slot == NULL) return LLGO_INVALID;
    DWORD result = WaitForSingleObject((*slot)->handle, try_only ? 0 : INFINITE);
    if (result == WAIT_OBJECT_0) { InterlockedDecrement(&(*slot)->count); return 0; }
    if (result == WAIT_TIMEOUT) return LLGO_AGAIN;
    return (int)GetLastError();
}

int llgo_win_sem_value(llgo_sem **slot, int *value)
{
    if (slot == NULL || *slot == NULL || value == NULL) return LLGO_INVALID;
    *value = (int)InterlockedCompareExchange(&(*slot)->count, 0, 0);
    return 0;
}
