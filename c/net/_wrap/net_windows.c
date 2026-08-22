/* Copyright (c) 2026 The GoPlus Authors. Licensed under Apache-2.0. */

#define WIN32_LEAN_AND_MEAN
#define _WINSOCK_DEPRECATED_NO_WARNINGS

#include <limits.h>
#include <stdint.h>
#include <winsock2.h>
#include <ws2tcpip.h>

static INIT_ONCE llgo_winsock_once = INIT_ONCE_STATIC_INIT;
static int llgo_winsock_error;

static BOOL CALLBACK llgo_init_winsock(
    PINIT_ONCE once, PVOID parameter, PVOID *context)
{
    WSADATA data;
    (void)once;
    (void)parameter;
    (void)context;
    llgo_winsock_error = WSAStartup(MAKEWORD(2, 2), &data);
    return TRUE;
}

static int llgo_ensure_winsock(void)
{
    if (!InitOnceExecuteOnce(
            &llgo_winsock_once, llgo_init_winsock, NULL, NULL))
        return (int)GetLastError();
    return llgo_winsock_error;
}

int llgo_net_getaddrinfo(
    const char *host, const char *service, const ADDRINFOA *hints,
    PADDRINFOA *result)
{
    int error = llgo_ensure_winsock();
    if (error != 0)
        return error;
    return getaddrinfo(host, service, hints, result);
}

void llgo_net_freeaddrinfo(PADDRINFOA value)
{
    freeaddrinfo(value);
}

uintptr_t llgo_net_socket(int domain, int type, int protocol)
{
    int error = llgo_ensure_winsock();
    if (error != 0) {
        WSASetLastError(error);
        return (uintptr_t)INVALID_SOCKET;
    }
    return (uintptr_t)socket(domain, type, protocol);
}

int llgo_net_bind(uintptr_t socket_value, const struct sockaddr_in *address,
                  int address_length)
{
    return bind((SOCKET)socket_value, (const struct sockaddr *)address,
                address_length);
}

int llgo_net_connect(uintptr_t socket_value, const struct sockaddr *address,
                     int address_length)
{
    return connect((SOCKET)socket_value, address, address_length);
}

int llgo_net_listen(uintptr_t socket_value, int backlog)
{
    return listen((SOCKET)socket_value, backlog);
}

uintptr_t llgo_net_accept(uintptr_t socket_value, struct sockaddr_in *address,
                          int *address_length)
{
    return (uintptr_t)accept((SOCKET)socket_value,
                             (struct sockaddr *)address, address_length);
}

int llgo_net_close(uintptr_t socket_value)
{
    return closesocket((SOCKET)socket_value);
}

struct hostent *llgo_net_gethostbyname(const char *name)
{
    int error = llgo_ensure_winsock();
    if (error != 0) {
        WSASetLastError(error);
        return NULL;
    }
    return gethostbyname(name);
}

const char *llgo_net_inet_ntop(int family, const void *source,
                               char *destination, size_t size)
{
    int error = llgo_ensure_winsock();
    if (error != 0) {
        WSASetLastError(error);
        return NULL;
    }
    return inet_ntop(family, source, destination, size);
}

unsigned long llgo_net_inet_addr(const char *value)
{
    int error = llgo_ensure_winsock();
    if (error != 0) {
        WSASetLastError(error);
        return INADDR_NONE;
    }
    return inet_addr(value);
}

int llgo_net_send(uintptr_t socket_value, const void *buffer, size_t length,
                  int flags)
{
    if (length > INT_MAX) {
        WSASetLastError(WSAEMSGSIZE);
        return SOCKET_ERROR;
    }
    return send((SOCKET)socket_value, (const char *)buffer, (int)length, flags);
}

int llgo_net_recv(uintptr_t socket_value, void *buffer, size_t length,
                  int flags)
{
    if (length > INT_MAX) {
        WSASetLastError(WSAEMSGSIZE);
        return SOCKET_ERROR;
    }
    return recv((SOCKET)socket_value, (char *)buffer, (int)length, flags);
}

int llgo_net_setsockopt(uintptr_t socket_value, int level, int option_name,
                        const void *option_value, int option_length)
{
    return setsockopt((SOCKET)socket_value, level, option_name,
                      (const char *)option_value, option_length);
}
