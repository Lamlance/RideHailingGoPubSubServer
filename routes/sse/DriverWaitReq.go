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

func DriverWaitReq(c *fiber.Ctx) error {
	geo_hash := c.Params("geo_hash")
	if len(geo_hash) < 4 {
		return c.SendStatus(400)
	}

	geo_key := geo_hash[0:4]

	c.Response().Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	c.SendStatus(200)
	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		ch, close, _ := libs.Subscribe(geo_key)
		for i:=0;;i++ {
			defer close()
			time.AfterFunc(10*time.Minute, func() {
				log.Println("Driver wait timeout")
				close()
			})
			msg, ok := <-ch
			if !ok {
				break
			}
			w.Write([]byte("id: " + strconv.Itoa(i) + "\n"))
			w.Write([]byte("event: " + "message" + "\n"))
			w.Write([]byte("data: " + msg + "\n"))
			w.Write([]byte("\n"))

			err := w.Flush()
			if err != nil {
				log.Println(err)
				break
			}
		}
	}))

	return nil
}
