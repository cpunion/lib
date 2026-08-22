/* Copyright (c) 2026 The GoPlus Authors. Licensed under Apache-2.0. */

#include <stddef.h>

typedef unsigned long llgo_dword;
typedef unsigned int llgo_uint;
typedef int llgo_bool;
typedef void *llgo_handle;

#if defined(_WIN64)
#define LLGO_WINAPI
#else
#define LLGO_WINAPI __attribute__((stdcall))
#endif

typedef llgo_dword(LLGO_WINAPI *llgo_thread_start)(void *arg);
typedef void(LLGO_WINAPI *llgo_fls_callback)(void *value);

__declspec(dllimport) llgo_handle LLGO_WINAPI CreateThread(
    void *, size_t, llgo_thread_start, void *, llgo_dword, llgo_dword *);
__declspec(dllimport) void LLGO_WINAPI ExitThread(llgo_dword);
__declspec(dllimport) llgo_bool LLGO_WINAPI CloseHandle(llgo_handle);
__declspec(dllimport) llgo_dword LLGO_WINAPI GetLastError(void);
__declspec(dllimport) llgo_dword LLGO_WINAPI WaitForSingleObject(llgo_handle, llgo_dword);
__declspec(dllimport) llgo_handle LLGO_WINAPI GetProcessHeap(void);
__declspec(dllimport) void *LLGO_WINAPI HeapAlloc(llgo_handle, llgo_dword, size_t);
__declspec(dllimport) llgo_bool LLGO_WINAPI HeapFree(llgo_handle, llgo_dword, void *);
__declspec(dllimport) llgo_dword LLGO_WINAPI FlsAlloc(llgo_fls_callback);
__declspec(dllimport) llgo_bool LLGO_WINAPI FlsFree(llgo_dword);
__declspec(dllimport) void *LLGO_WINAPI FlsGetValue(llgo_dword);
__declspec(dllimport) llgo_bool LLGO_WINAPI FlsSetValue(llgo_dword, void *);

#if defined(LLGO_USE_BDWGC)
llgo_handle LLGO_WINAPI GC_CreateThread(
    void *, size_t, llgo_thread_start, void *, llgo_dword, llgo_dword *);
void LLGO_WINAPI GC_ExitThread(llgo_dword);
#endif

enum {
    llgo_error_not_enough_memory = 8,
    llgo_error_invalid_parameter = 87,
    llgo_wait_object_0 = 0,
};
#define LLGO_INFINITE ((llgo_dword)0xffffffffUL)
#define LLGO_FLS_OUT_OF_INDEXES ((llgo_dword)0xffffffffUL)

typedef void *(*llgo_routine)(void *);

typedef struct llgo_thread {
    llgo_handle handle;
    llgo_routine routine;
    void *arg;
    void *result;
} llgo_thread;

static __declspec(thread) llgo_thread *llgo_current_thread;

static llgo_dword LLGO_WINAPI llgo_thread_entry(void *raw)
{
    llgo_thread *thread = (llgo_thread *)raw;
    llgo_current_thread = thread;
    thread->result = thread->routine(thread->arg);
    return 0;
}

int llgo_pthread_create(void **out, size_t stack_size, llgo_routine routine, void *arg)
{
    if (out == 0 || routine == 0) return llgo_error_invalid_parameter;
    llgo_thread *thread = (llgo_thread *)HeapAlloc(GetProcessHeap(), 0, sizeof(*thread));
    if (thread == 0) return llgo_error_not_enough_memory;
    thread->routine = routine;
    thread->arg = arg;
    thread->result = 0;
#if defined(LLGO_USE_BDWGC)
    thread->handle = GC_CreateThread(0, stack_size, llgo_thread_entry, thread, 0, 0);
#else
    thread->handle = CreateThread(0, stack_size, llgo_thread_entry, thread, 0, 0);
#endif
    if (thread->handle == 0) {
        int error = (int)GetLastError();
        HeapFree(GetProcessHeap(), 0, thread);
        return error;
    }
    *out = thread;
    return 0;
}

int llgo_pthread_join(void *raw, void **result)
{
    llgo_thread *thread = (llgo_thread *)raw;
    if (thread == 0) return llgo_error_invalid_parameter;
    llgo_dword wait = WaitForSingleObject(thread->handle, LLGO_INFINITE);
    if (wait != llgo_wait_object_0) return (int)GetLastError();
    if (result != 0) *result = thread->result;
    CloseHandle(thread->handle);
    HeapFree(GetProcessHeap(), 0, thread);
    return 0;
}

void llgo_pthread_exit(void *result)
{
    if (llgo_current_thread != 0) llgo_current_thread->result = result;
#if defined(LLGO_USE_BDWGC)
    GC_ExitThread(0);
#else
    ExitThread(0);
#endif
}

typedef void (*llgo_destructor)(void *);
typedef struct llgo_fls_value {
    llgo_destructor destructor;
    void *value;
} llgo_fls_value;

static void LLGO_WINAPI llgo_fls_destructor(void *raw)
{
    llgo_fls_value *slot = (llgo_fls_value *)raw;
    if (slot == 0) return;
    if (slot->destructor != 0 && slot->value != 0) slot->destructor(slot->value);
    HeapFree(GetProcessHeap(), 0, slot);
}

int llgo_pthread_key_create(llgo_uint *index)
{
    if (index == 0) return llgo_error_invalid_parameter;
    llgo_dword value = FlsAlloc(llgo_fls_destructor);
    if (value == LLGO_FLS_OUT_OF_INDEXES) return (int)GetLastError();
    *index = value;
    return 0;
}

int llgo_pthread_key_delete(llgo_uint index)
{
    return FlsFree(index) ? 0 : (int)GetLastError();
}

void *llgo_pthread_getspecific(llgo_uint index)
{
    llgo_fls_value *slot = (llgo_fls_value *)FlsGetValue(index);
    return slot == 0 ? 0 : slot->value;
}

int llgo_pthread_setspecific(llgo_uint index, llgo_destructor destructor, void *value)
{
    llgo_fls_value *slot = (llgo_fls_value *)FlsGetValue(index);
    if (value == 0) {
        if (slot == 0) return 0;
        if (!FlsSetValue(index, 0)) return (int)GetLastError();
        HeapFree(GetProcessHeap(), 0, slot);
        return 0;
    }
    if (slot == 0) {
        slot = (llgo_fls_value *)HeapAlloc(GetProcessHeap(), 0, sizeof(*slot));
        if (slot == 0) return llgo_error_not_enough_memory;
        slot->destructor = destructor;
        slot->value = value;
        if (!FlsSetValue(index, slot)) {
            int error = (int)GetLastError();
            HeapFree(GetProcessHeap(), 0, slot);
            return error;
        }
        return 0;
    }
    slot->destructor = destructor;
    slot->value = value;
    return 0;
}
