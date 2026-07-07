//go:build cgo && amx_c

/* The bundled AMX source retains its original license header in camx/amx.c. */
#define LINUX 1
#define AMX_NODYNALOAD 1
#define HAVE_STDINT_H 1
#define HAVE_ALLOCA_H 1
#include "camx/amx.c"
#include "c_bridge.h"

struct pawntest_amx {
  AMX amx;
  unsigned char *program;
};

pawntest_amx *pawntest_amx_load(const void *data, size_t size, int *error) {
  if (size < sizeof(AMX_HEADER)) {
    *error = AMX_ERR_FORMAT;
    return NULL;
  }
  AMX_HEADER header;
  memcpy(&header, data, sizeof(header));
  amx_Align16(&header.magic);
  amx_Align32((uint32_t *)&header.size);
  amx_Align32((uint32_t *)&header.stp);
  if (header.magic != AMX_MAGIC || header.size < 0 || (size_t)header.size > size || header.stp < header.size) {
    *error = AMX_ERR_FORMAT;
    return NULL;
  }
  pawntest_amx *vm = calloc(1, sizeof(*vm));
  if (vm == NULL) {
    *error = AMX_ERR_MEMORY;
    return NULL;
  }
  vm->program = calloc(1, (size_t)header.stp);
  if (vm->program == NULL) {
    free(vm);
    *error = AMX_ERR_MEMORY;
    return NULL;
  }
  memcpy(vm->program, data, (size_t)header.size);
  *error = amx_Init(&vm->amx, vm->program);
  if (*error != AMX_ERR_NONE) {
    free(vm->program);
    free(vm);
    return NULL;
  }
  return vm;
}

void pawntest_amx_free(pawntest_amx *vm) {
  if (vm != NULL) {
    amx_Cleanup(&vm->amx);
    free(vm->program);
    free(vm);
  }
}

int pawntest_amx_num_publics(pawntest_amx *vm) {
  int count = 0;
  return amx_NumPublics(&vm->amx, &count) == AMX_ERR_NONE ? count : -1;
}

int pawntest_amx_num_natives(pawntest_amx *vm) {
  int count = 0;
  return amx_NumNatives(&vm->amx, &count) == AMX_ERR_NONE ? count : -1;
}

int pawntest_amx_public_name(pawntest_amx *vm, int index, char *name, size_t size) {
  (void)size;
  return amx_GetPublic(&vm->amx, index, name);
}

int pawntest_amx_native_name(pawntest_amx *vm, int index, char *name, size_t size) {
  (void)size;
  return amx_GetNative(&vm->amx, index, name);
}

int pawntest_amx_exec(pawntest_amx *vm, int index, const int32_t *args, int count, int32_t *result) {
  for (int i = count - 1; i >= 0; --i) {
    int error = amx_Push(&vm->amx, (cell)args[i]);
    if (error != AMX_ERR_NONE) return error;
  }
  cell value = 0;
  int error = amx_Exec(&vm->amx, &value, index);
  *result = (int32_t)value;
  return error;
}

const char *pawntest_amx_error(int error) {
  switch (error) {
    case AMX_ERR_NONE: return "none";
    case AMX_ERR_BOUNDS: return "array index out of bounds";
    case AMX_ERR_MEMACCESS: return "invalid memory access";
    case AMX_ERR_INVINSTR: return "invalid instruction";
    case AMX_ERR_DIVIDE: return "divide by zero";
    case AMX_ERR_FORMAT: return "invalid AMX format";
    case AMX_ERR_INDEX: return "public index out of range";
    default: return "AMX execution error";
  }
}
