package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type Handler struct{}

func (h *Handler) FindMatch(c echo.Context) error {
	cookie, err := c.Cookie("Authorization")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		msg := ""
		err = websocket.Message.Receive(ws, &msg)
		if err != nil {
			c.Logger().Error(err)
		}
		
		stopChan := make(chan bool, 1)
		go func() {
			for {
				switch {
				case <-stopChan:
					return
				default:
					err := websocket.Message.Send(ws, "")
					if err != nil {
						c.Logger().Error(err)
					}

					msg := ""
					err = websocket.Message.Receive(ws, &msg)
					if err != nil {
						c.Logger().Error(err)
					}
					time.Sleep(10 * time.Second)
				}
			}
		}()

		

		for {
			if 
		}
		stopChan <- true
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
