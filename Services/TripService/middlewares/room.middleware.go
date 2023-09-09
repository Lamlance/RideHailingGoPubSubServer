package middlewares

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type CommunicationMsg struct {
	data []string
	lock *sync.Mutex
}

type CommunicationRoom struct {
	client_msg *CommunicationMsg
	driver_msg *CommunicationMsg

	RideInfo *RideReqInfo

	lock                *sync.Mutex
	Ride_requst_channel chan int
}

type GlobalCommunicationMsg struct {
	Data map[string]*CommunicationRoom
	Lock *sync.Mutex
}

type ResponDriver struct {
	Lon       float64 `json:"Lon"`
	Lat       float64 `json:"Lat"`
	Dist      float64 `json:"dist"`
	Driver_id string  `json:"driver_id"`
}

func GetDrivers(lon float64, lat float64, geo_hash string) ([]redis.GeoLocation, error) {

	url := fmt.Sprintf("http://localhost:3083/ridehail/geo/find/drivers?lon=%f&lat=%f&g=%s", lon, lat, geo_hash)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println("Body drivers", string(resBody))

	resp_drivers := make([]ResponDriver, 0)
	err = json.Unmarshal(resBody, &resp_drivers)
	if err != nil {
		return nil, err
	}
	drivers := make([]redis.GeoLocation, 0)
	for _, d := range resp_drivers {
		drivers = append(drivers, redis.GeoLocation{
			Longitude: d.Lon,
			Latitude:  d.Lat,
			Dist:      d.Dist,
			Name:      d.Driver_id,
		})
	}
	return drivers, nil
}

func publish_ride_request_loop(room *CommunicationRoom, res []redis.GeoLocation) {
	for i := 0; len(res) < 5; i++ {
		res = append(res, redis.GeoLocation{
			Name: "Driver #" + strconv.Itoa(i),
		})
	}
	room.lock.Lock()
	req_chan := room.Ride_requst_channel
	rideInfo := room.RideInfo
	room.lock.Unlock()

	const (
		DriverAccept  int = 0
		DriverDecline int = 1
		DriverTimeOut int = 2
	)

	timer := time.AfterFunc(0, func() {
		req_chan <- 2
	})
	driver_found := false
	for i, pos := range res {
		if driver_found {
			break
		}
		code, ok := <-req_chan
		if !ok {
			break
		}
		switch code {
		case DriverAccept:
			log.Println("Driver has accepted")
			if !timer.Stop() {
				<-timer.C
			}
			driver_found = true

		case DriverDecline:
			log.Println("Driver has declined")

			if !timer.Stop() {
				<-timer.C
			}
			rideInfo.Driver_id = pos.Name
			RideReqToPub <- rideInfo
			timer.Reset(20 * time.Second)
		case DriverTimeOut:
			log.Println("Driver has timeout")

			rideInfo.Driver_id = pos.Name
			RideReqToPub <- rideInfo
			if i < len(res)-1 {
				timer.Reset(20 * time.Second)
			}
		}
	}
	log.Println("Driver Req loop done")

	if !driver_found {
		log.Println("No driver found")

		room.driver_msg.lock.Lock()
		room.driver_msg.data = Enque(room.driver_msg.data, NoDriver)
		room.driver_msg.lock.Unlock()

	}
}

func ClientRideRequest(c *fiber.Ctx) error {
	rideInfo, ok_rideInfo := c.Locals("ride_info").(*RideReqInfo)
	geo_hash := c.Params("geo_hash")
	room, ok_room := c.Locals("room").(*CommunicationRoom)

	if !ok_rideInfo {
		log.Println("Server can't find ride info")
		return c.SendStatus(500)
	}

	if len(geo_hash) < 4 {
		c.SendStatus(400)
		return c.SendString("Invalid geo hash")
	}

	if !ok_room {
		log.Println("Server can't find communication room")
		return c.SendStatus(500)
	}

	room.RideInfo = rideInfo
	drivers, err := GetDrivers(rideInfo.SLon, rideInfo.SLat, geo_hash)
	if err != nil {
		log.Println("Get drivers errors: ", err)
		go publish_ride_request_loop(room, make([]redis.GeoLocation, 0))

	} else {
		log.Println("Find: ", len(drivers), " for trip: ", rideInfo.Trip_id)
		go publish_ride_request_loop(room, drivers)
	}

	return c.Next()
}

func DriverRideRequest(c *fiber.Ctx) error {
	trip_id := c.Params("trip_id")
	driver_id := c.Query("driver_id")
	if trip_id == "" || driver_id == "" {
		return c.SendStatus(400)
	}

	GlobalRoomMap.Lock.Lock()
	room, ok := GlobalRoomMap.Data[trip_id]
	if !ok {
		log.Println("Driver can't find trip Id in com channels")
		GlobalRoomMap.Lock.Unlock()
		return c.SendStatus(404)
	}
	c.Locals("room", room)
	log.Printf("Driver is stopping ride req loop")
	room.Ride_requst_channel <- 0
	log.Printf("Driver has stopping ride req loop")
	GlobalRoomMap.Lock.Unlock()

	return c.Next()
}
