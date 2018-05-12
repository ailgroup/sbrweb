#srbweb
Connects to sabre web services endpoints, both SOAP and REST.

## Tests and Benchmarks
Examples for running coverage reports as well as the tests themselves.

### Coverage
Test coverage is important. Examples of basic coverage stats along with more detailed reporting using `coverprofile`, `test`, and `tool`.

```sh
# Basic test coverage stats on main package and sub-packages
[sbrweb] go test -cover ./...
ok  	github.com/ailgroup/sbrweb/srvc	        0.074s	coverage: 71.5% of statements
ok  	github.com/ailgroup/sbrweb/hotelws	0.026s	coverage: 81.2% of statements
```

```sh
# generate new coverage file inside sbrweb/srvc directory
[srvc] go test -coverprofile=test_data/coverage.out
# coverage to be broken down by function
[srvc] go tool cover -func test_data/coverage.out

# generate new coverage file inside sbrweb/hotelws directory
[hotelws] go test -coverprofile=test_data/coverage.out
# coverage to be broken down by function
[hotelws] go tool cover -func test_data/coverage.out
```

### Running
Run all from root of `sbrweb`.

```sh
go test -v ./...
```

Run specific internal the package directory (e.g., `sbrweb/srvc`, `sbrweb/hotelws`, etc...):

```sh
go test -run TestMessageHeaderBaseUnmarshal
go test -bench=BenchmarkEnvelopeMarshal
go test -bench=.
```