package service

import (
	"context"
	"time"

	"danilamukhin/serv_go/internal/model"
	"danilamukhin/serv_go/internal/pgx/filter"
)

func (s Service) CreateHourParam(ctx context.Context, hourParam model.HourParam) (err error) {
	return s.repo.InsertHourParam(ctx, hourParam)
}

func (s Service) UpdateHourParam(ctx context.Context, hourParam model.HourParam) (err error) {
	return s.repo.UpdateHourParam(ctx, hourParam)
}

func (s Service) RangeHourParam(ctx context.Context, hourParam filter.HourParam) (_ model.HourParamList, err error) {
	return s.repo.RangeHourParam(ctx, hourParam)
}

func (s Service) DeleteHourParam(ctx context.Context, hourParam model.HourParam) (err error) {
	return s.repo.DelHourParam(ctx, hourParam)
}

func (s Service) GetHourParamList(ctx context.Context, filter filter.HourParam) (_ model.HourParamList, err error) {
	return s.repo.GetHourParamList(ctx, filter)
}

//TEST

func (s Service) CreateSchemaAndTable(ctx context.Context, conf model.TableConfig, year string, quarter string) (err error) {
	return s.repo.CreateSchemaAndTable(ctx, conf, year, quarter)
}

func (s Service) MinTimestamp(ctx context.Context, conf model.TableConfig) (minTimestamp time.Time, err error) {
	return s.repo.MinTimestamp(ctx, conf)
}

func (s Service) MoveQuarter(ctx context.Context, tableName string, year string, quarter string) (err error) {
	return s.repo.MoveQuarter(ctx, tableName, year, quarter)
}
