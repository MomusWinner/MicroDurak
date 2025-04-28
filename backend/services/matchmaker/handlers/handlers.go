package handlers

import (
	"context"
	"encoding/json"
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

type FindMatchResponse struct {
	// MatchStatus string
	Status    string `json:"status"`
	GameId    string `json:"game_id,omitempty"`
	GroupSize int    `json:"group_size,omitzero"`
}

type Handler struct {
	Queue         chan<- types.MatchChan
	Cancel        chan<- types.MatchCancel
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

		metrics.SearchDuration.WithLabelValues(h.Config.PodName, h.Config.Namespace).Observe(duration)

		metrics.HTTPRequestsTotal.WithLabelValues(
			c.Request().Method,
			c.Path(),
			strconv.Itoa(statusCode),
			h.Config.PodName,
			h.Config.Namespace,
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
		metrics.WebsocketUpgradeErrors.WithLabelValues(h.Config.PodName, h.Config.Namespace).Inc()
		return err
	}
	defer ws.Close()

	doneChan := make(chan bool, 1)
	closeHandler := ws.CloseHandler()
	ws.SetCloseHandler(func(code int, text string) error {
		select {
		case doneChan <- true:
		default:
		}
		err := closeHandler(code, text)
		return err
	})

	go func() {
		defer func() {
			select {
			case doneChan <- true:
			default:
			}
		}()
		for {
			if _, _, err := ws.ReadMessage(); err != nil {
				break
			}
		}
	}()

	metrics.PlayersSearching.WithLabelValues(h.Config.PodName, h.Config.Namespace).Inc()
	defer metrics.PlayersSearching.WithLabelValues(h.Config.PodName, h.Config.Namespace).Dec()

	returnChan := make(chan types.MatchResponse)
	h.Queue <- types.MatchChan{
		PlayerId:   playerId,
		Rating:     int(player.Rating),
		SentTime:   time.Now(),
		ReturnChan: returnChan,
	}

	for {
		select {
		case matchReturn := <-returnChan:
			switch matchReturn.Status {
			case types.MatchCreated:
				roomId := matchReturn.RoomId
				resp := FindMatchResponse{
					Status: types.MatchCreated.String(),
					GameId: roomId,
				}

				respString, _ := json.Marshal(resp)
				ws.WriteMessage(websocket.TextMessage, respString)

				closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Connection closed")
				err = ws.WriteControl(websocket.CloseMessage, closeMsg, time.Now().Add(time.Second))
				if err != nil {
					return err
				}

				return nil
			case types.MatchError:
				status := FindMatchResponse{
					Status: matchReturn.Status.String(),
				}

				stautsString, _ := json.Marshal(status)

				ws.WriteMessage(websocket.TextMessage, stautsString)
				if err != nil {
					return err
				}

				return matchReturn.Error
			default:
				status := FindMatchResponse{
					Status:    matchReturn.Status.String(),
					GroupSize: matchReturn.GroupSize,
				}

				stautsString, _ := json.Marshal(status)

				err := ws.WriteMessage(websocket.TextMessage, stautsString)
				if err != nil {
					return err
				}
			}
		case <-doneChan:
			h.Cancel <- types.MatchCancel{PlayerId: playerId}
			return nil
		}
	}
}
