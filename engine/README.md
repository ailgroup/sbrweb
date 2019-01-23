# Engine

This provides the low level connection to Sabre. Everything in here should be
defined well enough to not require modifications beyond what the transmission
package should provide.

The API for this package should hardly ever change.

## Structure
This is a large project with organized subprojects. To get a sense of the number of lines of code, tests, and other files you can `wc -l 'find sbrweb/hotelws -type f'`. 

1. `itin` (itinerary)
  * passenger name record (PNR) details
  * read PNR
  * cancel segment in PNR
1. `hotelws` (hotel web service)
  * Availability, Property and Rate descriptions, various hotel search services.
  * common struct/xml building and parsing
  * common logic for dealing with data formats like timestamps and cancellation policies
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
  * set of primitives to extend specific sabre errors (e.g., itin package `ErrFormat()`)
  * custom handling of sabre web services errors
  * includes handling for:
    * http network errors
    * SOAP faults
    * Sabre Web Service warnings, errors
    * XML (de|en)coding errors