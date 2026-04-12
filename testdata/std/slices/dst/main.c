#include "main.h"

// -- Implementation --

int main(void) {
    {
        // Make a slice.
        so_Slice s = slices_Make(so_int, ((mem_Allocator){0}), (3));
        so_at(so_int, s, 0) = 11;
        so_at(so_int, s, 1) = 22;
        so_at(so_int, s, 2) = 33;
        if (so_len(s) != 3 || so_cap(s) != 3) {
            so_panic("Make failed");
        }
        if (so_at(so_int, s, 0) != 11 || so_at(so_int, s, 1) != 22 || so_at(so_int, s, 2) != 33) {
            so_panic("Make failed");
        }
        slices_Free(so_int, ((mem_Allocator){0}), (s));
    }
    {
        // Append within capacity.
        so_Slice s = slices_MakeCap(so_int, ((mem_Allocator){0}), (0), (8));
        s = slices_Append(so_int, ((mem_Allocator){0}), (s), (10), (20), (30));
        if (so_len(s) != 3 || so_at(so_int, s, 0) != 10 || so_at(so_int, s, 1) != 20 || so_at(so_int, s, 2) != 30) {
            so_panic("Append failed");
        }
        slices_Free(so_int, ((mem_Allocator){0}), (s));
    }
    {
        // Append that triggers growth.
        so_Slice s = slices_MakeCap(so_int, ((mem_Allocator){0}), (0), (2));
        s = slices_Append(so_int, ((mem_Allocator){0}), (s), (1), (2));
        s = slices_Append(so_int, ((mem_Allocator){0}), (s), (3), (4), (5));
        if (so_len(s) != 5 || so_at(so_int, s, 0) != 1 || so_at(so_int, s, 4) != 5) {
            so_panic("Append grow failed");
        }
        slices_Free(so_int, ((mem_Allocator){0}), (s));
    }
    {
        // Extend from another slice.
        so_Slice s = slices_MakeCap(so_int, ((mem_Allocator){0}), (0), (8));
        so_Slice other = (so_Slice){(so_int[3]){100, 200, 300}, 3, 3};
        s = slices_Extend(so_int, ((mem_Allocator){0}), (s), (other));
        if (so_len(s) != 3 || so_at(so_int, s, 0) != 100 || so_at(so_int, s, 2) != 300) {
            so_panic("Extend failed");
        }
        slices_Free(so_int, ((mem_Allocator){0}), (s));
    }
    {
        // Clone a slice.
        so_Slice s1 = (so_Slice){(so_int[3]){11, 22, 33}, 3, 3};
        so_Slice s2 = slices_Clone(so_int, ((mem_Allocator){0}), (s1));
        so_at(so_int, s2, 0) = 99;
        if (so_at(so_int, s1, 0) != 11 || so_at(so_int, s2, 0) != 99) {
            so_panic("Clone failed");
        }
        slices_Free(so_int, ((mem_Allocator){0}), (s2));
    }
    {
        // Equal slices.
        so_Slice s1 = (so_Slice){(so_int[3]){1, 2, 3}, 3, 3};
        so_Slice s2 = (so_Slice){(so_int[3]){1, 2, 3}, 3, 3};
        so_Slice s3 = (so_Slice){(so_int[3]){1, 2, 4}, 3, 3};
        so_Slice s4 = (so_Slice){(so_int[2]){1, 2}, 2, 2};
        so_Slice s5 = (so_Slice){&so_Nil, 0, 0};
        so_Slice s6 = (so_Slice){&so_Nil, 0, 0};
        if (!slices_Equal(so_int, (s1), (s2))) {
            so_panic("want s1 == s2");
        }
        if (slices_Equal(so_int, (s1), (s3))) {
            so_panic("want s1 != s3");
        }
        if (slices_Equal(so_int, (s1), (s4))) {
            so_panic("want s1 != s4");
        }
        if (!slices_Equal(so_int, (s5), (s6))) {
            so_panic("want empty and nil slices equal");
        }
    }
    {
        // Equal string slices.
        so_Slice s1 = (so_Slice){(so_String[3]){so_str("a"), so_str("b"), so_str("c")}, 3, 3};
        so_Slice s2 = (so_Slice){(so_String[3]){so_str("a"), so_str("b"), so_str("c")}, 3, 3};
        so_Slice s3 = (so_Slice){(so_String[3]){so_str("a"), so_str("b"), so_str("d")}, 3, 3};
        if (!slices_Equal(so_String, (s1), (s2))) {
            so_panic("want s1 == s2");
        }
        if (slices_Equal(so_String, (s1), (s3))) {
            so_panic("want s1 != s3");
        }
    }
    {
        // Equal struct slices.
        typedef struct point {
            so_int x;
            so_int y;
        } point;
        so_Slice s1 = (so_Slice){(point[2]){(point){1, 2}, (point){3, 4}}, 2, 2};
        so_Slice s2 = (so_Slice){(point[2]){(point){1, 2}, (point){3, 4}}, 2, 2};
        so_Slice s3 = (so_Slice){(point[2]){(point){1, 2}, (point){3, 5}}, 2, 2};
        if (!slices_Equal(point, (s1), (s2))) {
            so_panic("want s1 == s2");
        }
        if (slices_Equal(point, (s1), (s3))) {
            so_panic("want s1 != s3");
        }
    }
}
