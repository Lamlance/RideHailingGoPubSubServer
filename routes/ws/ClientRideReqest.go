package ws

import (
	"goserver/libs"
	"goserver/routes"
	"strconv"

	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func publish_ride_request_loop(geo_key string,
	rideInfo *routes.RideReqInfo,
	room *CommunicationRoom,
	res []redis.GeoLocation) {
	for i := 0; i < 3; i++ {
		res = append(res, redis.GeoLocation{
			Name: "Driver #" + strconv.Itoa(i),
		})
	}
	room.lock.Lock()
	req_chan := room.Ride_requst_channel
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
			// b, err := json.Marshal(rideInfo)
			driver_found = true
			// if err != nil {
			// 	driver_found = false
			// 	break
			// }
			// // room.client_msg.lock.Lock()
			// room.client_msg.data = libs.Enque(room.client_msg.data, DriverFound+string(b))
			// room.client_msg.lock.Unlock()

			// room.driver_msg.lock.Lock()
			// room.driver_msg.data = libs.Enque(room.driver_msg.data, DriverFound+string(b))
			// room.driver_msg.lock.Unlock()

		case DriverDecline:
			log.Println("Driver has declined")

			if !timer.Stop() {
				<-timer.C
			}
			rideInfo.Driver_id = pos.Name
			routes.GlobalRideReqToPubChannel <- &routes.RideReqToPub{
				Channel:     geo_key,
				RideReqInfo: *rideInfo,
			}
			timer.Reset(10 * time.Second)
		case DriverTimeOut:
			log.Println("Driver has timeout")

			rideInfo.Driver_id = pos.Name
			routes.GlobalRideReqToPubChannel <- &routes.RideReqToPub{
				Channel:     geo_key,
				RideReqInfo: *rideInfo,
			}
			if i < len(res)-1 {
				timer.Reset(10 * time.Second)
			}
		}
	}
	log.Println("Driver Req loop done")

	if !driver_found {
		log.Println("No driver found")

		room.driver_msg.lock.Lock()
		room.driver_msg.data = libs.Enque(room.driver_msg.data, NoDriver)
		room.driver_msg.lock.Unlock()

	}

}

func ClientCheckMiddleware(c *fiber.Ctx) error {
	random_string := "Secrete_trip_id" //uuid.New().String()

	GlobalRoomMap.Lock.Lock()

	room := MakeEmptyCommunicationRoom()
	GlobalRoomMap.Data[random_string] = room

	c.Locals("trip_id", random_string)
	c.Locals("room", room)

	GlobalRoomMap.Lock.Unlock()

	return c.Next()
}

func ClientRideRequest(c *fiber.Ctx) error {
	rideInfo := &routes.RideReqInfo{
		SLon: c.QueryFloat("slon"),
		SLat: c.QueryFloat("slat"),
		SAdr: c.Query("sadr"),

		ELon: c.QueryFloat("elon"),
		ELat: c.QueryFloat("elat"),
		EAdr: c.Query("eadr"),

		User_id: c.Query("user_id"),
		Trip_id: c.Locals("trip_id").(string),
	}

	geo_hash := c.Params("geo_hash")
	room, ok_room := c.Locals("room").(*CommunicationRoom)

	if len(geo_hash) < 4 {
		c.SendStatus(400)
		return c.SendString("Invalid geo hash")
	}

	if !ok_room {
		log.Println("Server can't find communication room")
		return c.SendStatus(500)
	}

	room.RideInfo = rideInfo

	lon := rideInfo.SLon
	lat := rideInfo.SLat
	user_id := rideInfo.User_id
	geo_key := geo_hash[0:4]

	log.Println(lon, lat, user_id, geo_key)
	routes.GlobalRedis.Mutex.Lock()
	res, err := routes.GlobalRedis.Client.GeoRadius(routes.GlobalRedis.Context, geo_key, lon, lat, &redis.GeoRadiusQuery{
		Radius: 1,
		Unit:   "km",
		Count:  10,
		Sort:   "ASC",
	}).Result()
	routes.GlobalRedis.Mutex.Unlock()

	if err != nil {
		log.Println("Ride request error: ", err)
		return c.SendStatus(500)
	}
	if len(res) <= 0 {
		c.SendString("Server can't find any driver")
		return c.SendStatus(500)
	}
	go publish_ride_request_loop(geo_key, rideInfo, room, res)
	return c.Next()
}
