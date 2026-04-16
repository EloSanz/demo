package handler

import (
	"net/http"

	"github.com/elosanz/demo/internal/notification"
	"github.com/elosanz/demo/pkg/web"
)

type NotificationHandler struct {
	svc notification.NotificationService
}

func NewNotificationHandler(svc notification.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

type PublishRequest struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func (h *NotificationHandler) Publish(w http.ResponseWriter, r *http.Request) error {
	var req PublishRequest
	if err := web.DecodeJSON(r, &req); err != nil {
		return err
	}

	n := notification.Notification{
		Type:    req.Type,
		Content: req.Content,
	}

	if err := h.svc.Publish(r.Context(), n); err != nil {
		return err
	}

	return web.EncodeJSON(w, map[string]string{"status": "queued"}, http.StatusAccepted)
}
