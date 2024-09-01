package handler

import (
	"github.com/labstack/echo"
	"net/http"
	"os"
	"quickwiki-backend/model"
)

func (h *Handler) GetUserInfo(c echo.Context) (model.Me_Response, error) {
	if os.Getenv("DEV_MODE") == "true" {
		return model.Me_Response{
			TraqID:      "kavos",
			DisplayName: "kavos",
			IconUri:     "https://q.trap.jp/api/v3/public/icon/kavos",
		}, nil
	}
	userTraqID := c.Request().Header.Get("X-Forwarded-User")
	if username != "" {
		return model.Me_Response{}, echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	res, err := h.scraper.GetUserDetail(userTraqID[0])
	if err != nil {
		return model.Me_Response{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return res, nil
}
