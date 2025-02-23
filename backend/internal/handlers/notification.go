package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
    "todo-app/internal/models"
)

type NotificationHandler struct{}

func NewNotificationHandler() *NotificationHandler {
    return &NotificationHandler{}
}

func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
    userID := getUserIDFromToken(r)
    notifications, err := models.GetUserNotifications(userID)
    if err != nil {
        http.Error(w, "Could not get notifications", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(notifications)
}

func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    notificationID, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid notification ID", http.StatusBadRequest)
        return
    }

    err = models.MarkNotificationAsRead(uint(notificationID))
    if err != nil {
        http.Error(w, "Could not mark notification as read", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *NotificationHandler) CheckDueTasks(w http.ResponseWriter, r *http.Request) {
    err := models.CheckDueTasks()
    if err != nil {
        http.Error(w, "Could not check due tasks", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
} 