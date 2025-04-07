package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
)

type ChatHandler struct {
	chatDal dal.PubSubDAL
}

func NewChatHandler(chatDal dal.PubSubDAL) *ChatHandler {
	return &ChatHandler{
		chatDal: chatDal,
	}
}

func (h *ChatHandler) Subscribe(ctx *gin.Context) {
	// Subscribe to the chat topic
	messageChannel, err := h.chatDal.Subscribe(context.Background(), "chat")
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to subscribe to chat topic", err)
		return
	}

	// Listen for messages in a separate goroutine
	go func() {
		for msgStr := range messageChannel {
			var message models.Message
			if err := json.Unmarshal([]byte(msgStr), &message); err != nil {
				continue // Skip invalid messages
			}
			models.ResponseSuccess(ctx, http.StatusOK, "Received message", message)
		}
	}()

	models.ResponseSuccess(ctx, http.StatusOK, "Subscribed to chat topic", nil)
}

func (h *ChatHandler) Unsubscribe(ctx *gin.Context) {
	// Unsubscribe from the chat topic
	err := h.chatDal.Unsubscribe(context.Background(), "chat")
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to unsubscribe from chat topic", err)
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Unsubscribed successfully", nil)
}

func (h *ChatHandler) GetAllMessages(ctx *gin.Context) {
	// Get all messages from the chat topic
	messageStrings, err := h.chatDal.GetAllMessages(context.Background(), "chat")
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to retrieve messages", err)
		return
	}

	// Convert string messages to Message structs
	messages := make([]models.Message, 0, len(messageStrings))
	for _, msgStr := range messageStrings {
		var message models.Message
		if err := json.Unmarshal([]byte(msgStr), &message); err != nil {
			continue // Skip invalid messages
		}
		messages = append(messages, message)
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Messages retrieved successfully", messages)
}

func (h *ChatHandler) PublishMessage(ctx *gin.Context) {
	// Get the message from the request body
	var message models.Message
	if err := ctx.ShouldBindJSON(&message); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid message format", err)
		return
	}

	// Set the creation timestamp if not provided
	if message.CreatedAt == "" {
		message.CreatedAt = time.Now().Format(time.RFC3339)
	}

	// Convert message to JSON string
	messageJSON, err := json.Marshal(message)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to serialize message", err)
		return
	}

	// Publish the message to the chat topic
	err = h.chatDal.Publish(context.Background(), "chat", string(messageJSON))
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to publish message", err)
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Message published successfully", message)
}

func (h *ChatHandler) ReceiveMessage(ctx *gin.Context) {
	// Create a channel to receive messages
	messageChannel, err := h.chatDal.Subscribe(context.Background(), "chat")
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to subscribe to chat topic", err)
		return
	}

	// Listen for messages in a separate goroutine
	go func() {
		for msgStr := range messageChannel {
			var message models.Message
			if err := json.Unmarshal([]byte(msgStr), &message); err != nil {
				continue // Skip invalid messages
			}
			models.ResponseSuccess(ctx, http.StatusOK, "Message received", message)
		}
	}()

	models.ResponseSuccess(ctx, http.StatusOK, "Subscribed to chat topic", nil)
}
