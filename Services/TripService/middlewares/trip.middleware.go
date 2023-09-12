package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type RideReqInfo struct {
	SLon float64 `json:"slon"`
	SLat float64 `json:"slat"`
	SAdr string  `json:"sadr"`

	ELon float64 `json:"elon"`
	ELat float64 `json:"elat"`
	EAdr string  `json:"eadr"`

	User_id   string `json:"user_id"`
	Driver_id string `json:"driver_id"`
	Trip_id   string `json:"trip_id"`

	Price float64 `json:"price"`

	User_Name  string `json:"user_name"`
	User_Phone string `json:"user_phone"`
}

type ClientDetail struct {
	Phone string `json:"phone"`
	Name  string `json:"name"`
}

func CreateTrip(info *RideReqInfo) (string, error) {
	b, err := json.Marshal(info)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "http://ride_hailing_webapp:8080/api/rides", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Println("Body", string(body))
	trip_response := &struct {
		Id string `json:"id"`
	}{}
	err = json.Unmarshal(body, trip_response)
	if err != nil {
		return "", err
	}
	return trip_response.Id, nil
}

func GetUserDetailInfo(user_id string) (*ClientDetail, error) {
	res, err := http.Get("http://ride_hailing_webapp:8080/api/users/" + user_id)
	if err != nil {
		return nil, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	client_detail := &ClientDetail{
		Phone: "",
		Name:  "",
	}
	err = json.Unmarshal(resBody, client_detail)
	if err != nil {
		return nil, err
	}

	return client_detail, nil
}

func TripMiddleware(c *fiber.Ctx) error {

	// random_string := uuid.New().String()

	rideInfo := &RideReqInfo{
		SLon: c.QueryFloat("slon"),
		SLat: c.QueryFloat("slat"),
		SAdr: c.Query("sadr"),

		ELon: c.QueryFloat("elon"),
		ELat: c.QueryFloat("elat"),
		EAdr: c.Query("eadr"),

		User_id: c.Query("user_id"),
		Trip_id: "",

		Price:      c.QueryFloat("price"),
		User_Name:  c.Query("user_name", ""),
		User_Phone: c.Query("user_phone", ""),
	}
	trip_id, err := CreateTrip(rideInfo)
	if err != nil {
		log.Println("Trip create error: ", err)
		return c.SendStatus(500)
	}

	rideInfo.Trip_id = trip_id
	log.Println("Created trip: ", trip_id)

	if rideInfo.User_Name == "" || rideInfo.User_Phone == "" {
		res, err := GetUserDetailInfo(rideInfo.User_id)
		if err != nil {
			c.SendStatus(400)
			return c.SendString("Invalid user id")
		}

		rideInfo.User_Name = res.Name
		rideInfo.User_Phone = res.Phone
	}

	GlobalRoomMap.Lock.Lock()

	room := MakeEmptyCommunicationRoom()
	GlobalRoomMap.Data[trip_id] = room

	c.Locals("trip_id", trip_id)
	c.Locals("room", room)
	c.Locals("ride_info", rideInfo)
	GlobalRoomMap.Lock.Unlock()

	return c.Next()
}
