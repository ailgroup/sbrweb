# sbrweb
Connects to Sabre® SOAP APIs, formerly known as Sabre® Web Services, and REST endpoints. Sabre has over 100 APIs, a number of them implemented here. For counting lines of code I recommend [scc](https://github.com/boyter/scc), which can be run in the root of the project as `scc .`.

## Structure
Project is built around three core projects:

1. `rest`
    * hotel queries; sabre does not include REST endpoints for hotel reservations, must use `soap` package.
1. `sbrerr`
    * standard formatting for errors.
1. `soap`
    * hotel, itinerary, and sessions for all SOAP endpoints.


These packages provide the low level connections and functionality for Sabre. Everything here should be defined well enough to not require modification beyond what external custom API packages might need.

Sabre Web Services provides both SOAP and REST endpoints for many endpoints. However, they do not provide REST endpoints for Itinerary services, which are required for making reservations and cancellations. Since reservations require using the AAA workspace one needs to load availability requests first; this necessitates providing those endpoints as SOAP even if they exist as REST (for more information see the `workflows.md` document).

The API for this package should rarely change as Sabre endpoints are both versioned, stable, and have long term support.

### rest

### soap

1. `hotel` (hotel web service) implements many SOAP endpoints for the hotel portion of Sabre SOAP.
    * Availability, Property/Rate descriptions, various hotel search services.
    * Common struct/xml building and parsing
    * Common logic for dealing with data formats like timestamps and cancellation policies
1. `itin` (itinerary) deals with Itinerary, PNR, Reservation, Cancelations, and Profiles.
    * Passenger name record (PNR) details
    * Read PNR
    * Copy profile to PNR
    * Cancel segment in PNR
1. `srvc` (service) core set of functionality for common SOAP and session management.
    * Basic SOAP
      * envelope
      * message header
      * security
      * soap fault
    * Sessions
      * create, close, validate
      * buffered queue as session pool

### sbrerr

1. `sbrerr` (sabre errors)
    * Set of primitives to extend specific sabre errors (e.g., itin package `ErrFormat()`)
    * Custom handling of sabre web services errors
    * Includes handling for:
      * http network errors
      * SOAP faults
      * Sabre Web Service warnings, errors
      * XML (de|en)coding errors

## Documentation

```
godoc -http=:6060
```

Go to `http://localhost:6060/pkg/CODE-HOST/COMPANY-USER-NAME/sbrweb/`.

## TODO and Currently

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
go test -cover ./soap/htlsp

#generate new coverage file for project
go test -coverprofile=coverage.out ./...
#generate new coverage file for specific package
go test ./soap/srvc -coverprofile=soap/srvc/test_data/coverage.out

#cover tool shows coverage by function
go tool cover -func coverage.out
```

### Tests
Run all from root of `sbrweb`, or internal to packages (e.g., `sbrweb/itin`).

```sh
# run verbose tests for all packages at once
go test -v ./...
# run verbose tests for specific package
go test -v ./soap/itin
# run verbose test on specific test function
go test -v ./soap/htlsp -run TestRateDescCall
```


### Benchmarks
For core functions that are heavily used we provide benchmarks.

```sh
# run all benchmarks
go test ./... -bench=.
go test ./soap/srvc -bench=BenchmarkMessageHeaderBaseMarshal
```