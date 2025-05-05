package nats

import (
	"encoding/json"
	"log"
	"statistics_service/internal/core/domain"
	"statistics_service/internal/usecase"
	"time"

	"github.com/nats-io/nats.go"
)

type NATSHandler struct {
	useCase *usecase.StatisticsUseCase
}

func NewNATSHandler(uc *usecase.StatisticsUseCase) *NATSHandler {
	return &NATSHandler{useCase: uc}
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func (h *NATSHandler) Subscribe(nc *nats.Conn) {
	_, err := nc.Subscribe("order.created", func(msg *nats.Msg) {
		var payload struct {
			UserID    string `json:"user_id"`
			OrderID   string `json:"order_id"`
			CreatedAt string `json:"created_at"`
		}
		if err := json.Unmarshal(msg.Data, &payload); err != nil {
			log.Println("failed to unmarshal order.created:", err)
			return
		}

		err := h.useCase.SaveOrderStat(domain.OrderStat{
			UserID:    payload.UserID,
			OrderID:   payload.OrderID,
			CreatedAt: parseTime(payload.CreatedAt),
		})
		if err != nil {
			log.Println("failed to save order stat:", err)
		}
	})
	if err != nil {
		log.Println("NATS subscribe error:", err)
	}
}
