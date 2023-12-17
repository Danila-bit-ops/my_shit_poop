package service

import (
	"context"

	"danilamukhin/serv_go/internal/model"
	"danilamukhin/serv_go/internal/pgx/filter"
)

func (s Service) CreateHourParam(ctx context.Context, hourParam model.HourParam) (err error) {
	return s.repo.InsertHourParam(ctx, hourParam)
}

func (s Service) GetHourParamList(ctx context.Context, filter filter.HourParam) (_ model.HourParamList, err error) {
	return s.repo.GetHourParamList(ctx, filter)
}
