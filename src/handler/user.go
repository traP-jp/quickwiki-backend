package handler

import (
	"github.com/labstack/echo"
	"net/http"
	"quickwiki-backend/model"
)

func (h *Handler) GetUserInfo(c echo.Context) (model.Me_Response, error) {
	userTraqID, ok := c.Request().Header["X-Forwarded-User"]
	if !ok {
		return model.Me_Response{}, echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	res, err := h.scraper.GetUserDetail(userTraqID[0])
	if err != nil {
		return model.Me_Response{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return res, nil
}