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
#include <fcntl.h>
#include <io.h>
#include <limits.h>
#include <stdarg.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

int llgoErrno(void) { return *_errno(); }

int llgo_os_clearenv(void)
{
    char **env = *__p__environ();
    while (env != NULL && env[0] != NULL) {
        const char *entry = env[0];
        const char *separator = strchr(entry, '=');
        if (separator == NULL || separator == entry) {
            ++env;
            continue;
        }
        size_t length = (size_t)(separator - entry);
        char *name = (char *)malloc(length + 1);
        if (name == NULL) return -1;
        memcpy(name, entry, length);
        name[length] = '\0';
        int result = _putenv_s(name, "");
        free(name);
        if (result != 0) return -1;
        env = *__p__environ();
    }
    return 0;
}

int llgo_os_open(const char *path, int flags, ...)
{
    if ((flags & _O_CREAT) != 0) {
        va_list args;
        va_start(args, flags);
        int mode = va_arg(args, int);
        va_end(args);
        return _open(path, flags, mode);
    }
    return _open(path, flags);
}

int llgo_os_fcntl(int fd, int cmd, ...)
{
    (void)fd;
    (void)cmd;
    errno = ENOSYS;
    return -1;
}

int llgo_os_read(int fd, void *buffer, size_t count)
{
    if (count > UINT_MAX) count = UINT_MAX;
    return _read(fd, buffer, (unsigned int)count);
}

int llgo_os_write(int fd, const void *buffer, size_t count)
{
    if (count > UINT_MAX) count = UINT_MAX;
    return _write(fd, buffer, (unsigned int)count);
}
