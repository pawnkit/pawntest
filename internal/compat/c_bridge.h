#ifndef PAWNTEST_C_BRIDGE_H
#define PAWNTEST_C_BRIDGE_H

#include <stddef.h>
#include <stdint.h>

typedef struct pawntest_amx pawntest_amx;

pawntest_amx *pawntest_amx_load(const void *data, size_t size, int *error);
void pawntest_amx_free(pawntest_amx *vm);
int pawntest_amx_num_publics(pawntest_amx *vm);
int pawntest_amx_num_natives(pawntest_amx *vm);
int pawntest_amx_public_name(pawntest_amx *vm, int index, char *name, size_t size);
int pawntest_amx_native_name(pawntest_amx *vm, int index, char *name, size_t size);
int pawntest_amx_exec(pawntest_amx *vm, int index, const int32_t *args, int count, int32_t *result);
const char *pawntest_amx_error(int error);

#endif
