# sbrweb
Connects to Sabre Web Services SOAP endpoints abd provides eaiser to use APIs. Project is built around three core projects:

1. BM Engine: `engine`
  * core set of packages handling all low-level sabre services interaction.
1. BM Transmission: `transmission`
  * higher-level APIs to faciliate interacting with the engine packages.
1. BM clients: `client`
  * client specific integrations leveraging the transmission package.

See the respective projects' README for more details.


## TODO and Currenlty

* multiple currency
* pnr read
* segment cancel


## Tests and Benchmarks
Examples for running coverage reports as well as the tests themselves.

### Coverage
Test coverage is important. Examples of basic coverage stats along with more detailed reporting using `coverprofile`, `test`, and `tool`.

```sh
# Basic test coverage stats on main package and sub-packages
[sbrweb] go test -cover ./...
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

### Tests
We aspire to full test coverage. In some cases 100% test coverage is either not practical nor useful. Additionally, 100% test coverage is often misleading, for this reason we seek to provide testing around multiple scenarios for the same entry points.

Run all from root of `sbrweb`, or internal to packages (e.g., `sbrweb/itin`).

```sh
# run verbose tests for all packages at once
go test -v ./...
# run verbose tests for specific package
go test -v ./engine/itin
# run verbose test on specific test function
go test -v ./engine/hotelws -run TestRateDescCall
```


### Benchmarks
For core functions that are heavily used we provide benchmarks.

```sh
# run all benchmarks
go test ./... -bench=.
go test ./engine/srvc -bench=BenchmarkMessageHeaderBaseMarshal
```