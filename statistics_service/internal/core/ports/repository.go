package ports

import "statistics_service/internal/core/domain"

type StatisticsRepository interface {
	SaveOrderStat(stat domain.OrderStat) error
	SaveInventoryStat(stat domain.InventoryStat) error
	GetUserOrderStats(userID string) ([]domain.OrderStat, error)
}
