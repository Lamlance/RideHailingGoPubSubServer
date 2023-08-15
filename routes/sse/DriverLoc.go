package sse

import (
	"bufio"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type SSEBody struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func DriverLoc(c *fiber.Ctx) error {
	c.Response().Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	//c.Set("Transfer-Encoding", "chunked")
	c.SendStatus(200)
	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		for i := 0; i < 10; i++ {
			msg := &SSEBody{
				Data:  "Hello #" + strconv.Itoa(i),
				Event: "message",
			}

			w.Write([]byte("id: " + strconv.Itoa(i) + "\n"))
			w.Write([]byte("event: " + msg.Event + "\n"))
			w.Write([]byte("data: " + msg.Data + "\n"))
			w.Write([]byte("\n"))

			err := w.Flush()
			if err != nil {
				log.Print("Error while flushing:", err)
				break
			}
			time.Sleep(2 * time.Second)
		}
	}))

	return nil;
}
