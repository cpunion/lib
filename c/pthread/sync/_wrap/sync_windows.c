/* Copyright (c) 2026 The GoPlus Authors. Licensed under Apache-2.0. */

/*
 * LLGo's Windows runtime already owns the INIT_ONCE, SRWLOCK, and
 * CONDITION_VARIABLE wrappers used by sync_windows.go. Keep only the
 * semaphore adapter here so programs do not carry a second copy of the
 * runtime synchronization backend.
 */
#define WIN32_LEAN_AND_MEAN
#include <limits.h>
#include <windows.h>

enum {
    LLGO_AGAIN = 11,
    LLGO_INVALID = 22,
};

typedef struct {
    HANDLE handle;
    volatile LONG count;
} llgo_libc_sem;

int llgo_libc_sem_init(llgo_libc_sem *sem, int shared, unsigned int value)
{
    if (sem == NULL || shared != 0 || value > LONG_MAX)
        return LLGO_INVALID;
    sem->handle = CreateSemaphoreW(NULL, (LONG)value, LONG_MAX, NULL);
    if (sem->handle == NULL)
        return (int)GetLastError();
    sem->count = (LONG)value;
    return 0;
}

int llgo_libc_sem_destroy(llgo_libc_sem *sem)
{
    if (sem == NULL || sem->handle == NULL)
        return LLGO_INVALID;
    if (!CloseHandle(sem->handle))
        return (int)GetLastError();
    sem->handle = NULL;
    sem->count = 0;
    return 0;
}

int llgo_libc_sem_post(llgo_libc_sem *sem)
{
    if (sem == NULL || sem->handle == NULL)
        return LLGO_INVALID;
    if (!ReleaseSemaphore(sem->handle, 1, NULL))
        return (int)GetLastError();
    InterlockedIncrement(&sem->count);
    return 0;
}

int llgo_libc_sem_wait(llgo_libc_sem *sem, int try_only)
{
    DWORD result;
    if (sem == NULL || sem->handle == NULL)
        return LLGO_INVALID;
    result = WaitForSingleObject(sem->handle, try_only ? 0 : INFINITE);
    if (result == WAIT_OBJECT_0) {
        InterlockedDecrement(&sem->count);
        return 0;
    }
    if (result == WAIT_TIMEOUT)
        return LLGO_AGAIN;
    return (int)GetLastError();
}

int llgo_libc_sem_value(llgo_libc_sem *sem, int *value)
{
    if (sem == NULL || sem->handle == NULL || value == NULL)
        return LLGO_INVALID;
    *value = (int)InterlockedCompareExchange(&sem->count, 0, 0);
    return 0;
}
