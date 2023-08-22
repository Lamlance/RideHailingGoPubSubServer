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
		defer close()
		time.AfterFunc(10*time.Minute, func() {
			log.Println("Driver wait timeout")
			close()
		})

		w.Write([]byte("id: " + strconv.Itoa(0) + "\n"))
		w.Write([]byte("event: " + "ping" + "\n"))
		w.Write([]byte("data: \n"))
		w.Write([]byte("\n"))
		if err:=w.Flush(); err != nil {
			return
		}

		for i := 1; ; i++ {
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
