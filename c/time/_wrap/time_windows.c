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

#define _CRT_SECURE_NO_WARNINGS

#include <errno.h>
#include <stdint.h>
#include <time.h>
#include <windows.h>

typedef struct llgo_timespec {
    __time64_t sec;
    long nsec;
} llgo_timespec;

__time64_t llgo_time(__time64_t *timer) { return _time64(timer); }
__time64_t llgo_mktime(struct tm *value) { return _mktime64(value); }
char *llgo_ctime(const __time64_t *timer) { return _ctime64(timer); }
double llgo_difftime(__time64_t end, __time64_t start) { return _difftime64(end, start); }
struct tm *llgo_gmtime(const __time64_t *timer) { return _gmtime64(timer); }
struct tm *llgo_localtime(const __time64_t *timer) { return _localtime64(timer); }
size_t llgo_strftime(char *buf, size_t size, const char *format, const struct tm *value)
{
    return strftime(buf, size, format, value);
}
clock_t llgo_clock(void) { return clock(); }

static void llgo_from_100ns(uint64_t ticks, llgo_timespec *value)
{
    value->sec = (__time64_t)(ticks / 10000000ULL);
    value->nsec = (long)((ticks % 10000000ULL) * 100ULL);
}

int llgo_clock_gettime(int clock_id, llgo_timespec *value)
{
    if (value == NULL) {
        errno = EINVAL;
        return -1;
    }
    if (clock_id == 0) {
        struct timespec now;
        if (timespec_get(&now, TIME_UTC) != TIME_UTC) {
            errno = EINVAL;
            return -1;
        }
        value->sec = now.tv_sec;
        value->nsec = now.tv_nsec;
        return 0;
    }
    if (clock_id == 1 || clock_id == 4) {
        LARGE_INTEGER counter;
        LARGE_INTEGER frequency;
        if (!QueryPerformanceCounter(&counter) || !QueryPerformanceFrequency(&frequency)) {
            errno = EINVAL;
            return -1;
        }
        value->sec = counter.QuadPart / frequency.QuadPart;
        value->nsec = (long)(((counter.QuadPart % frequency.QuadPart) * 1000000000ULL) /
            frequency.QuadPart);
        return 0;
    }

    FILETIME created, exited, kernel, user;
    BOOL ok = clock_id == 2
        ? GetProcessTimes(GetCurrentProcess(), &created, &exited, &kernel, &user)
        : clock_id == 3
            ? GetThreadTimes(GetCurrentThread(), &created, &exited, &kernel, &user)
            : FALSE;
    if (!ok) {
        errno = EINVAL;
        return -1;
    }
    ULARGE_INTEGER k, u;
    k.LowPart = kernel.dwLowDateTime;
    k.HighPart = kernel.dwHighDateTime;
    u.LowPart = user.dwLowDateTime;
    u.HighPart = user.dwHighDateTime;
    llgo_from_100ns(k.QuadPart + u.QuadPart, value);
    return 0;
}

int llgo_clock_settime(int clock_id, const llgo_timespec *value)
{
    (void)clock_id;
    (void)value;
    errno = ENOTSUP;
    return -1;
}

int llgo_clock_getres(int clock_id, llgo_timespec *value)
{
    if (value == NULL) {
        errno = EINVAL;
        return -1;
    }
    value->sec = 0;
    if (clock_id == 1 || clock_id == 4) {
        LARGE_INTEGER frequency;
        if (!QueryPerformanceFrequency(&frequency)) {
            errno = EINVAL;
            return -1;
        }
        value->nsec = (long)(1000000000ULL / frequency.QuadPart);
        if (value->nsec == 0) value->nsec = 1;
        return 0;
    }
    if (clock_id >= 0 && clock_id <= 3) {
        value->nsec = 100;
        return 0;
    }
    errno = EINVAL;
    return -1;
}
