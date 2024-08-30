package handler

import (
	"log"
	"net/http"
	"quickwiki-backend/model"

	"github.com/labstack/echo"
)

func (h *Handler) PostMessageToTraQ(c echo.Context) error {
	var message model.MessageToTraQ_POST
	err := c.Bind(&message)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
	}

	err = h.scraper.MessageToTraQ(message.Content)

	if err != nil {
		log.Println("post message err : ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
	}

	return c.JSON(http.StatusOK, "Post Message To TraQ is clear")
}
