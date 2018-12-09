# Workflows
Notes on specific sabre soap services to achieve an end-goal.

## Itineraries

https://developer.sabre.com/docs/read/soap_apis/management/itinerary


## Modify Passenger Name Record
The following workflow allows you to modify an existing passenger name record (PNR). This workflow requires a passenger name record to be created in advance.
Steps

1. Retrieve the passenger name record using the Retrieve Itinerary API (GetReservationRQ).
2. Call the Modify Itinerary API (TravelItineraryModifyInfoLLSRQ) with the updated passenger name record in the request.
3. End the transaction of the passenger name record using the End Transaction API (EndTransactionLLSRQ).

NOTE you must use a session token (with a session pool) to consume these APIs. You may need to call other APIs in addition to the Modify Itinerary API to complete your transaction.

## Book Hotel Reservation
The following workflow allows you to search and book a hotel room.
Steps

1. Retrieve hotel availability using OTA_HotelAvailLLSRQ.
2. Retrieve hotel rates using HotelPropertyDescriptionLLSRQ.
3. Retrieve hotel rules and policies using HotelRateDescriptionLLSRQ.\*
4. Add any additional (required) information to create the passenger name record (PNR) using PassengerDetailsRQ.\*\*
5. Book a room for the selected hotel using OTA_HotelResLLSRQ.
6. End the transaction of the passenger name record using EndTransactionLLSRQ.
Note

\* Mandatory only if selected option in response of HotelPropertyDescriptionLLSRQ contains HRD_RequiredForSell="true".

\*\* Ensure Agency address is added within call to PassengerDetails, so as the OTA_HotelResLLSRQ call is not rejected.



## Post Booking Transaction (cancel booking)
The following workflow demonstrates how to take action (cancel) over an existing passenger name record. This workflow requires a passenger name record to be created in advance.
Steps

1. Retrieve the passenger name record using GetReservationRQ.
2. Cancel the existing itinerary using OTA_CancelLLSRQ.
3. End the transaction of the passenger name record using EndTransactionLLSRQ.


## Car Reservation
The following workflow allows you to search and book a rental car.
Steps

1. Retrieve rental car availability using OTA_VehAvailRateLLSRQ .
2. Add passenger information using PassengerDetailsRQ.
3. Book desired rental car using OTA_VehResLLSRQ.
4. End the transaction of the passenger name record using EndTransactionLLSRQ.
