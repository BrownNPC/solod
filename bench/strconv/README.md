# strconv benchmarks

Run the benchmark:

```text
make bench name=strconv
```

## Parsing

Go 1.26.1:

```text
goos: darwin
goarch: arm64
pkg: solod.dev/bench/strconv
cpu: Apple M1
Benchmark_Atof64_Decimal-8         53919762    20.76 ns/op
Benchmark_Atof64_Float-8           47406704    24.28 ns/op
Benchmark_Atof64_Exp-8             46986270    25.44 ns/op
Benchmark_Atof64_Big-8             31725643    37.66 ns/op
Benchmark_ParseInt_7bit-8         123124688     9.747 ns/op
Benchmark_ParseInt_26bit-8         83328512    14.47 ns/op
Benchmark_ParseInt_31bit-8         73286917    16.36 ns/op
Benchmark_ParseInt_56bit-8         50235303    23.86 ns/op
Benchmark_ParseInt_62bit-8         45842499    26.42 ns/op
```

So:

```text
Benchmark_Atof64_Decimal           98232759    12.07 ns/op
Benchmark_Atof64_Float             78673047    14.92 ns/op
Benchmark_Atof64_Exp               59550393    20.74 ns/op
Benchmark_Atof64_Big               48915700    25.00 ns/op
Benchmark_ParseInt_7bit           298545337     3.853 ns/op
Benchmark_ParseInt_26bit          165213476     7.258 ns/op
Benchmark_ParseInt_31bit          139547053     8.593 ns/op
Benchmark_ParseInt_56bit           82124281    14.59 ns/op
Benchmark_ParseInt_62bit           70671378    16.91 ns/op
```

## Formatting

Go 1.26.1:

```text
goos: darwin
goarch: arm64
pkg: solod.dev/bench/strconv
cpu: Apple M1
Benchmark_FormatFloat_Decimal-8    40641112    29.88 ns/op
Benchmark_FormatFloat_Float-8      28667326    42.79 ns/op
Benchmark_FormatFloat_Exp-8        34375335    34.75 ns/op
Benchmark_FormatFloat_Big-8        30486578    39.15 ns/op
Benchmark_FormatInt_7bit-8         82486035    14.45 ns/op
Benchmark_FormatInt_26bit-8        73198256    16.50 ns/op
Benchmark_FormatInt_31bit-8        61575829    19.55 ns/op
Benchmark_FormatInt_56bit-8        51147170    23.65 ns/op
Benchmark_FormatInt_62bit-8        45434968    26.46 ns/op
```

So:

```text
Benchmark_FormatFloat_Decimal      44751072    26.92 ns/op
Benchmark_FormatFloat_Float        36485253    33.67 ns/op
Benchmark_FormatFloat_Exp          43064775    29.64 ns/op
Benchmark_FormatFloat_Big          36714088    33.26 ns/op
Benchmark_FormatInt_7bit          252425917     4.753 ns/op
Benchmark_FormatInt_26bit         167144420     7.178 ns/op
Benchmark_FormatInt_31bit         144028574     8.334 ns/op
Benchmark_FormatInt_56bit         100000000    11.76 ns/op
Benchmark_FormatInt_62bit          91240875    13.12 ns/op
```
