# Engine

This provides the low level connection to Sabre. Everything in here should be
defined well enough to not require modifications beyond what the transmission
package should provide.

The API for this package should hardly ever change.

## Structure
This is a large project with organized subprojects. To get a sense of the number of lines of code, tests, and other files you can `wc -l 'find sbrweb/hotelws -type f'`. 

Another tool I recommend for counting lines of code is [scc](https://github.com/boyter/scc), which can be run in the root of the project as `scc .`.

1. `itin` (itinerary)
  * Passenger name record (PNR) details
  * Read PNR
  * Copy profile to PNR
  * Cancel segment in PNR
1. `hotelws` (hotel web service)
  * Availability, Property/Rate descriptions, various hotel search services.
  * Common struct/xml building and parsing
  * Common logic for dealing with data formats like timestamps and cancellation policies
1. `srvc` (service)
  * Basic SOAP
    * envelope
    * message header
    * security
    * soap fault
  * Sessions
    * create, close, validate
    * buffered queue as session pool
1. `sbrerr` (sabre errors)
  * Set of primitives to extend specific sabre errors (e.g., itin package `ErrFormat()`)
  * Custom handling of sabre web services errors
  * Includes handling for:
    * http network errors
    * SOAP faults
    * Sabre Web Service warnings, errors
    * XML (de|en)coding errors