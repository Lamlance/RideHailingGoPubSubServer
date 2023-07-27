package xhr

import (
	"goserver/libs"
	"goserver/routes"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func DriverWaitRequest(c *fiber.Ctx) error {
	ch, close := libs.Subscribe(routes.PubSub, "Driver")
	defer close()

	time.AfterFunc(10*time.Minute, func() {
		log.Println("Driver wait timeout")
		close()
	})

	
	msg, ok := <-ch
	
	if(ok){
		return c.SendString(msg)
	}

	return c.SendStatus(408)
}
