# sbrweb

## TODO
 multiple currency

Connects to Sabre Web services endpoints, both SOAP and REST. It is built around three offerings:

1. BookingEngine (BEN): `engine`
1. BookingClient (bClient): `client`
1. BookingCloud (bCloud): `cloud`

## Engine
This is a large project with organized subprojects. To get a sense of the number of lines of code, tests, and other files you can `wc -l 'find sbrweb/hotelws -type f'`. 

1. `srvc` (service)
  * Basic SOAP
    * envelope
    * message header
    * security
    * soap fault
  * Sessions
    * create, close, validate
    * buffered queue as session pool
1. `hotelws` (hotel web service)
  * Availability, Property and Rate descriptions, Reservation services.
  * common struct/xml building and parsing
  * common logic for dealing with data formats like timestamps and cancellation policies
1. `itin` (itinerary)
  * passenger details
  * cancel
1. `hotelrest` (hotel rest)


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