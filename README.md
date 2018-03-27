#srbweb
Connects to sabre web services endpoints, both SOAP and REST.

## Tests and Benchmarks

Run all: 

```sh
go test -bench=.
go test -v
```

Run specific:

```sh
go test -run TestMessageHeaderBaseUnmarshal
go test -bench=BenchmarkEnvelopeMarshal
```