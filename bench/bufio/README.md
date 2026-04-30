# bufio benchmarks

Run the benchmark:

```text
make bench name=bufio
```

Go 1.26.1:

```text
goos: darwin
goarch: arm64
pkg: solod.dev/bench/bufio
cpu: Apple M1
Benchmark_ReaderBuf-8       380958    3089 ns/op      32848 B/op    7 allocs/op
Benchmark_ReaderUnbuf-8     967016    1269 ns/op      16448 B/op    3 allocs/op
Benchmark_WriterBuf-8       400078    2988 ns/op      32848 B/op    7 allocs/op
Benchmark_WriterUnbuf-8     235978    4928 ns/op      65040 B/op    8 allocs/op
Benchmark_Scanner-8        2721681     443.0 ns/op     4264 B/op    6 allocs/op
```

So:

```text
Benchmark_ReaderBuf        1073540    1073 ns/op      32768 B/op    4 allocs/op
Benchmark_ReaderUnbuf      2876020     411.6 ns/op    16384 B/op    1 allocs/op
Benchmark_WriterBuf        1000000    1038 ns/op      32768 B/op    4 allocs/op
Benchmark_WriterUnbuf       725557    1537 ns/op      65024 B/op    7 allocs/op
Benchmark_Scanner         10875376     112.0 ns/op     4096 B/op    1 allocs/op
```
