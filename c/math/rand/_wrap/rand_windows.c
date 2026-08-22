/* Copyright (c) 2026 The GoPlus Authors. Licensed under Apache-2.0. */

#define _CRT_RAND_S
#include <stdlib.h>

int llgo_rand_r(unsigned int *seed)
{
    *seed = *seed * 214013U + 2531011U;
    return (int)((*seed >> 16) & RAND_MAX);
}

void llgo_sranddev(void)
{
    unsigned int seed;
    if (rand_s(&seed) != 0) seed = 1;
    srand(seed);
}
