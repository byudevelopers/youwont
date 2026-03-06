package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"youwont.api/internal/middleware"
	"youwont.api/internal/service"
)

type InviteHandler struct {
	svc *service.InviteService
}

func NewInviteHandler(svc *service.InviteService) *InviteHandler {
	return &InviteHandler{svc: svc}
}

// Send handles POST /groups/:id/invites.
// @Summary      Send group invite
// @Description  Invite a user to a group. Sender must be a group member. Creates a notification for the invitee.
// @Tags         invites
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Group ID"
// @Param        body body SendInviteRequest true "User to invite"
// @Success      201 {object} model.Invite
// @Failure      400 {object} ErrorResponse
// @Failure      403 {object} ErrorResponse
// @Failure      409 {object} ErrorResponse
// @Router       /groups/{id}/invites [post]
func (h *InviteHandler) Send(c *echo.Context) error {
	user := middleware.UserFromContext(c)

	groupID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return badRequest(c, "invalid group id")
	}

	var body struct {
		UserID string `json:"user_id"`
	}
	if err := c.Bind(&body); err != nil {
		return badRequest(c, "invalid request body")
	}

	inviteeID, err := primitive.ObjectIDFromHex(body.UserID)
	if err != nil {
		return badRequest(c, "invalid user_id")
	}

	invite, err := h.svc.Send(c.Request().Context(), user, groupID, inviteeID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, invite)
}

// ListMine handles GET /invites.
// @Summary      List my pending invites
// @Description  Returns all pending group invites for the authenticated user.
// @Tags         invites
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} InviteListResponse
// @Failure      401 {object} ErrorResponse
// @Router       /invites [get]
func (h *InviteHandler) ListMine(c *echo.Context) error {
	user := middleware.UserFromContext(c)

	invites, err := h.svc.ListMine(c.Request().Context(), user.ID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"invites": invites,
	})
}

// Accept handles POST /invites/:id/accept.
// @Summary      Accept an invite
// @Description  Accept a pending group invite. Adds the user to the group and notifies the inviter.
// @Tags         invites
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Invite ID"
// @Success      200 {object} InviteAcceptResponse
// @Failure      400 {object} ErrorResponse
// @Failure      403 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Router       /invites/{id}/accept [post]
func (h *InviteHandler) Accept(c *echo.Context) error {
	user := middleware.UserFromContext(c)

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return badRequest(c, "invalid invite id")
	}

	invite, err := h.svc.Accept(c.Request().Context(), user, id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":       invite.ID,
		"status":   invite.Status,
		"group_id": invite.GroupID,
	})
}

// Decline handles POST /invites/:id/decline.
// @Summary      Decline an invite
// @Description  Decline a pending group invite. The user can be re-invited later.
// @Tags         invites
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Invite ID"
// @Success      200 {object} InviteDeclineResponse
// @Failure      400 {object} ErrorResponse
// @Failure      403 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Router       /invites/{id}/decline [post]
func (h *InviteHandler) Decline(c *echo.Context) error {
	user := middleware.UserFromContext(c)

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return badRequest(c, "invalid invite id")
	}

	invite, err := h.svc.Decline(c.Request().Context(), user, id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":     invite.ID,
		"status": invite.Status,
	})
}
