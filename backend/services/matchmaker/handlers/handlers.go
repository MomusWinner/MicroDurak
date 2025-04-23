package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MommusWinner/MicroDurak/internal/players/v1"
	"github.com/MommusWinner/MicroDurak/services/matchmaker/config"
	"github.com/MommusWinner/MicroDurak/services/matchmaker/metrics"
	"github.com/MommusWinner/MicroDurak/services/matchmaker/types"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	upgrader = websocket.Upgrader{}
)

type Handler struct {
	Queue         chan<- types.MatchChan
	Config        *config.Config
	PlayersClient players.PlayersClient
}

func (h *Handler) FindMatch(c echo.Context) error {
	start := time.Now()
	var err error

	defer func() {
		duration := time.Since(start).Seconds()
		var statusCode int

		// Determine status code from error if available
		if httpErr, ok := err.(*echo.HTTPError); ok {
			statusCode = httpErr.Code
		} else if err != nil {
			statusCode = http.StatusInternalServerError
		}

		metrics.SearchDuration.Observe(duration)

		metrics.HTTPRequestsTotal.WithLabelValues(
			c.Request().Method,
			c.Path(),
			strconv.Itoa(statusCode),
		).Inc()
	}()

	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	playerId, ok := c.Get("playerId").(string)
	if !ok {
		panic("Missing jwt middleware")
	}

	player, err := h.PlayersClient.GetPlayer(ctx, &players.GetPlayerRequest{Id: playerId})
	if err != nil {
		s := status.Convert(err)
		if s.Code() == codes.NotFound {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unknown Player")
		}
		c.Logger().Error(err)
		return err
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		metrics.WebsocketUpgradeErrors.Inc()
		return err
	}
	defer ws.Close()

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := ws.WriteMessage(websocket.PingMessage, []byte{})
				if err != nil {
					metrics.WebsocketWriteErrors.Inc()
					c.Logger().Error(err)
				}

				_, _, err = ws.ReadMessage()
				if err != nil {
					c.Logger().Error(err)
				}
			}
		}
	}()

	metrics.PlayersSearching.Inc()
	defer metrics.PlayersSearching.Dec()

	returnChan := make(chan types.MatchResponse)
	h.Queue <- types.MatchChan{
		PlayerId:   playerId,
		Rating:     int(player.Rating),
		SentTime:   time.Now(),
		ReturnChan: returnChan,
	}

	roomId := <-returnChan

	type FindMatchResponse struct {
		RoomId string `json:"room_id"`
	}
	ws.WriteMessage(websocket.TextMessage, []byte(roomId.RoomId))

	closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Connection closed")
	err = ws.WriteControl(websocket.CloseMessage, closeMsg, time.Now().Add(time.Second))
	if err != nil {
		return err
	}

	return nil
}
