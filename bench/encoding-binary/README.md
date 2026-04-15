# encoding/binary benchmarks

Run the benchmark:

```text
make bench name=encoding/binary
```

Go 1.26.1:

```text
goos: darwin
goarch: arm64
pkg: solod.dev/bench/encoding/binary
cpu: Apple M1
Benchmark_BE_PutUint64-8       1000000000    0.6283 ns/op    12732.19 MB/s
Benchmark_BE_AppendUint64-8    695606808     1.772 ns/op      4514.40 MB/s
Benchmark_LE_PutUint64-8       1000000000    0.6321 ns/op    12655.26 MB/s
Benchmark_LE_AppendUint64-8    692428867     1.726 ns/op      4634.84 MB/s
```

So:

```text
Benchmark_BE_PutUint64         1000000000    0.3180 ns/op    25159.13 MB/s
Benchmark_BE_AppendUint64      1000000000    0.9517 ns/op     8405.88 MB/s
Benchmark_LE_PutUint64         1000000000    0.3142 ns/op    25459.79 MB/s
Benchmark_LE_AppendUint64      1000000000    0.9514 ns/op     8408.69 MB/s
```
