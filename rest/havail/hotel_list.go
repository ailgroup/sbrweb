// Package htlres provides functionality for Sabre Hotel rest endpoints.
package havail

/*
TODO: https://developer.sabre.com/docs/read/rest_apis/hotel/search/get_hotel_list

In Detail: This service allows a series of optional search parameters which determine which hotels are to be returned. The parameters include HotelCode, HotelName, Marketer Codes, Chain Codes, Amenity Codes, and Property Types. Multiple number of these parameters can be passed in the request and the results will be filtered out according to the request. If no hotel is found matching with the Search criteria, nothing will be returned.

duplicates what hotel search criteria already does; so maybe not...?
*/

/*
//accept a type of map[string]string, with values appended to url
func AddParams() {
	v := url.Values{}
	v.Set("name", "Ava")
	v.Add("friend", "Jess")
	v.Add("friend", "Sarah")
	v.Add("friend", "Zoe")
	//v.Encode() == "name=Ava&friend=Jess&friend=Sarah&friend=Zoe"
	fmt.Println(v.Get("name"))
	fmt.Println(v.Get("friend"))
	fmt.Println(v["friend"])
}
*/
