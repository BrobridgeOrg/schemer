# schemer
Data serialization library

## Benchmark

Here is the benchmark result of the JavaScript VM and the transformer.

```shell
$ go test -bench . --benchmem -cpu 4,8,16
goos: darwin
goarch: amd64
pkg: github.com/BrobridgeOrg/schemer/runtime/goja
cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
BenchmarkJavaScriptVM-4                	 4894088	       243.6 ns/op	     104 B/op	       3 allocs/op
BenchmarkJavaScriptVM-8                	 4795806	       247.2 ns/op	     104 B/op	       3 allocs/op
BenchmarkJavaScriptVM-16               	 5060337	       244.3 ns/op	     104 B/op	       3 allocs/op
BenchmarkTransformer-4                 	   72811	     15771 ns/op	    6984 B/op	     171 allocs/op
BenchmarkTransformer-8                 	   76226	     15830 ns/op	    6984 B/op	     171 allocs/op
BenchmarkTransformer-16                	   78670	     15990 ns/op	    6984 B/op	     171 allocs/op
BenchmarkTransformer_PassThrough-4     	 1329184	       889.8 ns/op	     408 B/op	      10 allocs/op
BenchmarkTransformer_PassThrough-8     	 1359368	       907.0 ns/op	     408 B/op	      10 allocs/op
BenchmarkTransformer_PassThrough-16    	 1348467	       888.5 ns/op	     408 B/op	      10 allocs/op
PASS
ok  	github.com/BrobridgeOrg/schemer/runtime/goja	15.254s
```
