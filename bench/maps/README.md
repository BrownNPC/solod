# maps benchmarks

Requires GCC/Clang and mimalloc (for heap allocations in So). If mimalloc isn't available, the benchmarks will use the default libc allocator, which is much slower.

Run the benchmark:

```text
make bench name=maps
```

Go 1.26.1:

```text
goos: darwin
goarch: arm64
pkg: solod.dev/bench/maps
cpu: Apple M1

Benchmark_IntSet-8       31677    35580 ns/op     74264 B/op    20 allocs/op
Benchmark_IntGet-8      218179     5573 ns/op         0 B/op     0 allocs/op
Benchmark_IntDel-8       50260    23892 ns/op     36944 B/op     5 allocs/op

Benchmark_StrSet-8       24082    48677 ns/op    108760 B/op    20 allocs/op
Benchmark_StrGet-8      134481     8990 ns/op         0 B/op     0 allocs/op
Benchmark_StrDel-8       34094    33878 ns/op     54608 B/op     5 allocs/op
```

So (mimalloc):

```text
Benchmark_IntSet         20305    56696 ns/op     98112 B/op    27 allocs/op
Benchmark_IntGet        734978     1638 ns/op         0 B/op     0 allocs/op
Benchmark_IntDel         30958    38556 ns/op     73728 B/op     6 allocs/op

Benchmark_StrSet         18986    71879 ns/op    130816 B/op    27 allocs/op
Benchmark_StrGet        117218    10206 ns/op         0 B/op     0 allocs/op
Benchmark_StrDel         23670    50111 ns/op     98304 B/op     6 allocs/op
```

So (arena):

```text
Benchmark_IntSet         21157    57661 ns/op     98112 B/op    27 allocs/op
Benchmark_IntGet        752787     1583 ns/op         0 B/op     0 allocs/op
Benchmark_IntDel         30339    38821 ns/op     73728 B/op     6 allocs/op

Benchmark_StrSet         18884    63500 ns/op    130816 B/op    27 allocs/op
Benchmark_StrGet        119041    10083 ns/op         0 B/op     0 allocs/op
Benchmark_StrDel         24212    49507 ns/op     98304 B/op     6 allocs/op
```

So (built-in map):

```text
Benchmark_IntSet        374192     3242 ns/op
Benchmark_IntGet        446527     2733 ns/op
Benchmark_StrSet        170769     6970 ns/op
Benchmark_StrGet        114100    10735 ns/op
```
