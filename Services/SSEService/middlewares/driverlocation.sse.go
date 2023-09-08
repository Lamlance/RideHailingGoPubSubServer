package middlewares

import (
	"bufio"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

//	type DriverLocToAdd struct {
//		Lon       float64
//		Lat       float64
//		GeoKey    string
//		Driver_id string
//	}
//
//	type DriverLocationPostBody struct {
//		Lon float64 `json:"lon"`
//		Lat float64 `json:"lat"`
//		G   string  `json:"g"`
//	}

func DriverLocationWatch(c *fiber.Ctx) error {
	driver_id := c.Params("driver_id")
	if driver_id == "" {
		return c.SendStatus(400)
	}

	c.Response().Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	//c.Set("Transfer-Encoding", "chunked")
	c.SendStatus(200)
	log.Println("A client start watching driver: ", driver_id)

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		log.Println("Start subscribe: ", driver_id)
		ch, close, ok := Subscribe("DriverCoord")
		if !ok {
			log.Println("Cant subscribe: ", driver_id)
			return
		}
		log.Println("Had subscribe: ", driver_id)

		defer close()
		w.Write([]byte("id: " + strconv.Itoa(0) + "\n"))
		w.Write([]byte("event: ping \n"))
		w.Write([]byte("data: \n"))
		w.Write([]byte("\n"))

		if w.Flush() != nil {
			return
		}

		for i := 1; ; i++ {
			data, ok := <-ch
			if !ok {
				return
			}
			// log.Println("SSE Client get driver loc: ", data)
			w.Write([]byte("id: " + strconv.Itoa(i) + "\n"))
			w.Write([]byte("event: message \n"))
			w.Write([]byte("data: " + data + "\n"))
			w.Write([]byte("\n"))

			err := w.Flush()
			if err != nil {
				log.Print("Error while flushing:", err)
				break
			}
		}
	}))

	return nil
}
