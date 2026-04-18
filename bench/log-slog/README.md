# log/slog benchmarks

Run the benchmark:

```text
make bench name=log-slog
```

Go 1.26.1:

```text
goos: darwin
goarch: arm64
pkg: solod.dev/bench/log-slog
cpu: Apple M1
Benchmark_NoAttr-8    6130304    165.6 ns/op      0 B/op    0 allocs/op
Benchmark_Attr-8      4617334    259.1 ns/op    144 B/op    3 allocs/op
```

So:

```text
Benchmark_NoAttr     30484750     38.62 ns/op     0 B/op    0 allocs/op
Benchmark_Attr       31404568     38.15 ns/op     0 B/op    0 allocs/op
```
