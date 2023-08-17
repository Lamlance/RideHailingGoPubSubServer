package sse

import (
	"bufio"
	"goserver/libs"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func DriverLoc(c *fiber.Ctx) error {
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
	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {

		ch, close, ok := libs.Subscribe("DriverLoc")
		if !ok {
			return
		}
		defer close()

		w.Write([]byte("id: " + strconv.Itoa(0) + "\n"))
		w.Write([]byte("event: ping \n"))
		w.Write([]byte("data: \n"))
		w.Write([]byte("\n"))

		if w.Flush() != nil {
			return
		}

		for i := 1;; i++ {
			data, ok := <-ch
			if !ok {
				return
			}

			w.Write([]byte("id: " + strconv.Itoa(i) + "\n"))
			w.Write([]byte("event: message \n"))
			w.Write([]byte("data: " + data + "\n"))
			w.Write([]byte("\n"))

			err := w.Flush()
			if err != nil {
				log.Print("Error while flushing:", err)
				break
			}
			time.Sleep(2 * time.Second)
		}
	}))

	return nil
}
