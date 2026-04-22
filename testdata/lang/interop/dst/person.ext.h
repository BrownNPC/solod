#include <stdio.h>
#include <stdint.h>
#include "so/builtin/builtin.h"

typedef struct Account Account;

typedef void (*write_func_t)(Account* a, const char* fmt, ...);

struct Account {
    so_String name;
    int64_t balance;
    so_Slice flags;
    write_func_t write;
};

int64_t account_inc_balance(Account* a, int64_t amount);

void account_set_name(Account* a, so_String name);

void write_acc(Account* a, const char* fmt, ...);
