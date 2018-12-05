# Client (BookingClient, bClient)


```sh
curl -v -H "Accept: application/json" -X GET 'http://localhost:8080/rates/hotel/id?hotel_id=4' | jq

curl -v -H "Accept: application/json" -X GET 'http://localhost:8080/rates/hotel/id?hotel_id=4&guest_count=2' | jq

curl -H "Accept: application/json" -X GET 'http://localhost:8080/rates/hotel/id?guest_count=3&arrive=2018-08-10&depart=2018-08-11&hotel_id=1004' | jq

curl -H "Accept: application/json" -X GET 'http://localhost:8080/rates/hotel/id?guest_count=3&arrive=2018-11-10&depart=2018-11-11&hotel_id=1004' | jq


curl -H "Accept: application/json" -X GET 'http://localhost:8080/book/hotel/room?room_meta=YXJ2OjExLTEwfGRwdDoxMS0xMXxnc3Q6M3xoYzpIT0QxMDA0LzEwTk9WLTExTk9WM3xoaWQ6MTAwNHxycGg6MDA2fHJtdDo2WTlDTUVXfFtjdXI6VVNELXJxczpmYWxzZS1hbXQ6MTM3LjE3XQ=='
```

