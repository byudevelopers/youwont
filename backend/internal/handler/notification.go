package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"youwont.api/internal/middleware"
	"youwont.api/internal/service"
)

type NotificationHandler struct {
	svc *service.NotificationService
}

func NewNotificationHandler(svc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

// List handles GET /notifications?page=0&limit=20.
// @Summary      List notifications
// @Description  Returns the authenticated user's notification feed, newest first. Paginated.
// @Tags         notifications
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Page number (0-indexed)" default(0)
// @Param        limit query int false "Items per page (max 50)" default(20)
// @Success      200 {object} NotificationListResponse
// @Failure      401 {object} ErrorResponse
// @Router       /notifications [get]
func (h *NotificationHandler) List(c *echo.Context) error {
	user := middleware.UserFromContext(c)

	page, _ := strconv.ParseInt(c.QueryParam("page"), 10, 64)
	limit, _ := strconv.ParseInt(c.QueryParam("limit"), 10, 64)
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	notifs, hasMore, err := h.svc.List(c.Request().Context(), user.ID, page, limit)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"notifications": notifs,
		"has_more":      hasMore,
	})
}

// UnreadCount handles GET /notifications/unread-count.
// @Summary      Get unread notification count
// @Description  Returns the number of unread notifications. Used for the bell badge.
// @Tags         notifications
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} UnreadCountResponse
// @Failure      401 {object} ErrorResponse
// @Router       /notifications/unread-count [get]
func (h *NotificationHandler) UnreadCount(c *echo.Context) error {
	user := middleware.UserFromContext(c)

	count, err := h.svc.UnreadCount(c.Request().Context(), user.ID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"count": count,
	})
}

// MarkRead handles POST /notifications/:id/read.
// @Summary      Mark notification as read
// @Description  Marks a single notification as read.
// @Tags         notifications
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Notification ID"
// @Success      200 {object} MarkReadResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Router       /notifications/{id}/read [post]
func (h *NotificationHandler) MarkRead(c *echo.Context) error {
	user := middleware.UserFromContext(c)

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return badRequest(c, "invalid notification id")
	}

	if err := h.svc.MarkRead(c.Request().Context(), id, user.ID); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":   id.Hex(),
		"read": true,
	})
}

// MarkAllRead handles POST /notifications/read-all.
// @Summary      Mark all notifications as read
// @Description  Marks all unread notifications as read for the authenticated user.
// @Tags         notifications
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} MarkAllReadResponse
// @Failure      401 {object} ErrorResponse
// @Router       /notifications/read-all [post]
func (h *NotificationHandler) MarkAllRead(c *echo.Context) error {
	user := middleware.UserFromContext(c)

	count, err := h.svc.MarkAllRead(c.Request().Context(), user.ID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"updated": count,
	})
}
