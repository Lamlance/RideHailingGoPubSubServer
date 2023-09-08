package middlewares

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

func GetUserDetailInfo(user_id string) (*ClientDetail, error) {
	res, err := http.Get("http://localhost:3000/api/users/" + user_id)
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
	random_string := uuid.New().String()

	rideInfo := &RideReqInfo{
		SLon: c.QueryFloat("slon"),
		SLat: c.QueryFloat("slat"),
		SAdr: c.Query("sadr"),

		ELon: c.QueryFloat("elon"),
		ELat: c.QueryFloat("elat"),
		EAdr: c.Query("eadr"),

		User_id: c.Query("user_id"),
		Trip_id: random_string,

		Price:      c.QueryFloat("price"),
		User_Name:  c.Query("user_name", ""),
		User_Phone: c.Query("user_phone", ""),
	}

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
	GlobalRoomMap.Data[random_string] = room

	c.Locals("trip_id", random_string)
	c.Locals("room", room)
	c.Locals("ride_info", rideInfo)
	GlobalRoomMap.Lock.Unlock()

	return c.Next()
}
