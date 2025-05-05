package usecase

import (
	"statistics_service/internal/core/domain"
	"statistics_service/internal/core/ports"
)

type StatisticsUseCase struct {
	repo ports.StatisticsRepository
}

func New(repo ports.StatisticsRepository) *StatisticsUseCase {
	return &StatisticsUseCase{repo: repo}
}

func (uc *StatisticsUseCase) SaveOrderStat(stat domain.OrderStat) error {
	return uc.repo.SaveOrderStat(stat)
}

func (uc *StatisticsUseCase) SaveInventoryStat(stat domain.InventoryStat) error {
	return uc.repo.SaveInventoryStat(stat)
}

func (uc *StatisticsUseCase) GetUserOrderStats(userID string) ([]domain.OrderStat, error) {
	return uc.repo.GetUserOrderStats(userID)
}
