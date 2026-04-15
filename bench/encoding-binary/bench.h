#include "so/builtin/builtin.h"

typedef struct {
    uint8_t buf[8];
    so_Slice res;
} Sink;

volatile Sink sink = {0};
