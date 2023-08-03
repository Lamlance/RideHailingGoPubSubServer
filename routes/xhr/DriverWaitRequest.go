package xhr

import (
	"goserver/libs"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func DriverWaitRequest(c *fiber.Ctx) error {
	geo_hash := c.Params("geo_hash")
	if len(geo_hash) < 4 {
		return c.SendStatus(404)
	}



	geo_key := geo_hash[0:4]

	ch, close,ok := libs.Subscribe(geo_key)

	if !ok {
		return c.SendStatus(500)
	}

	defer close()

	time.AfterFunc(10*time.Minute, func() {
		log.Println("Driver wait timeout")
		close()
	})

	
	msg, ok := <-ch	

	if(ok){
		
		c.SendStatus(200)
		return c.SendString(msg)
	}

	return c.SendStatus(408)
}
