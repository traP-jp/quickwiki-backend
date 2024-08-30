package handler

import (
	"context"
	"log"
	"net/http"
	"os"
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
		AccessToken:          os.Getenv("TRAQ_BOT_TOKEN_BN256"),
		DisableAutoReconnect: true,
	})
	if err != nil {
		log.Println("traQ BOT access :", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized WebSocket access.")
	}

	botChannelId := "01913f8b-8c05-76a2-b51f-bb83e9e93615"
	_, _, err = bot.API().
		MessageApi.
		PostMessage(context.Background(), botChannelId).
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
