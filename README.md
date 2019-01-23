# sbrweb
Connects to Sabre Web Services SOAP endpoints abd provides eaiser to use APIs. Project is built around three core projects:

1. BM Engine: `engine`
    * core set of packages handling all low-level sabre services interaction.
1. BM Transmission: `transmission`
    * higher-level APIs to faciliate interacting with the engine packages.
1. BM clients: `client`
    * client specific integrations leveraging the transmission package.

See the respective projects' README for more details.


## TODO and Currently

* Pnr read
* Segment cancel
* GeoServices (city and state lookups independent of hotel web services)
* Multiple currency
* Sabre rest endpoints for hotel content


## Tests and Benchmarks
Examples for running coverage reports as well as the tests themselves.
We aspire to full test coverage. In some cases 100% test coverage is neither practical nor useful. Additionally, 100% test coverage is often misleading. In many cases you want to test one function multiple times. We seek to provide testing for all executable code _and_ multiple scenarios for the same entry points.

### Coverage
Test coverage is important. Examples of basic coverage stats along with more detailed reporting using `coverprofile`, `test`, and `tool`.

```sh
#coverage stats for whole project
[sbrweb] go test -cover ./...
#coverage stats for specific package
go test -cover ./engine/hotelws

#generate new coverage file for project
go test -coverprofile=coverage.out ./...
#generate new coverage file for specific package
go test ./engine/srvc -coverprofile=engine/srvc/test_data/coverage.out

#cover tool shows coverage by function
go tool cover -func coverage.out
```

### Tests
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