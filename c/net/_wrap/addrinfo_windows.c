/* Copyright (c) 2026 The GoPlus Authors. Licensed under Apache-2.0. */

#define WIN32_LEAN_AND_MEAN

#include <winsock2.h>
#include <ws2tcpip.h>

static INIT_ONCE llgo_winsock_once = INIT_ONCE_STATIC_INIT;
static int llgo_winsock_error;

static BOOL CALLBACK llgo_init_winsock(PINIT_ONCE once, PVOID parameter, PVOID *context)
{
    (void)once;
    (void)parameter;
    (void)context;
    WSADATA data;
    llgo_winsock_error = WSAStartup(MAKEWORD(2, 2), &data);
    return TRUE;
}

int llgo_net_getaddrinfo(
    const char *host, const char *service, const ADDRINFOA *hints,
    PADDRINFOA *result)
{
    InitOnceExecuteOnce(&llgo_winsock_once, llgo_init_winsock, NULL, NULL);
    if (llgo_winsock_error != 0) return llgo_winsock_error;
    return getaddrinfo(host, service, hints, result);
}

void llgo_net_freeaddrinfo(PADDRINFOA value)
{
    freeaddrinfo(value);
}
