package handler

import (
	"context"
	"log"
	"net/http"
	"quickwiki-backend/model"

	"github.com/labstack/echo"
	traq "github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
)

func (h *Handler) PostMessageToTraQ(c echo.Context) error {

	var message model.MessageToTraQ_POST
	err := c.Bind(&message)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
	}

	bot, err := traqwsbot.NewBot(&traqwsbot.Options{
		//AccessToken: os.Getenv("TRAQ_BOT_TOKEN"),
		AccessToken:          "fX6iHTxwkZR7zle4vXzlxQIZbSWXFbnbj5GA",
		DisableAutoReconnect: true,
	})
	if err != nil {
		log.Println("traQ BOT access :", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized WebSocket access.")
	}

	_, _, err = bot.API().
		MessageApi.
		PostMessage(context.Background(), "1aec50b2-0cdf-46d2-8877-5f0eebd11bd4").
		PostMessageRequest(traq.PostMessageRequest{
			Content: message.Content,
		}).
		Execute()

	if err != nil {
		log.Println("post message err : ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
	}

	if err := bot.Start(); err != nil {
		log.Println("BOT start err : ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
	}

	return c.JSON(http.StatusOK, "ok")
}
