package pgx

import (
	"context"
	"fmt"

	"danilamukhin/serv_go/internal/model"
	"danilamukhin/serv_go/internal/pgx/filter"
)

func addWhere(addQ, q string) string {
	if len(addQ) > 0 {
		return addQ + " and " + q
	}
	return " where " + q
}

func (r Repo) InsertHourParam(ctx context.Context, hourParam model.HourParam) (err error) {
	const q = `insert into hour_params
	(val, param_id, timestamp, change_by, xml_create, manual, comment)
	values ($1, $2, $3, $4, $5, $6, $7)`

	_, err = r.pool.Exec(ctx, q, hourParam.Val, hourParam.ParamID, hourParam.Timestamp,
		hourParam.ChangeBy, hourParam.XMLCreate, hourParam.Manual, hourParam.Comment)
	if err != nil {
		return err
	}

	return nil
}

func (r Repo) RangeHourParam(ctx context.Context, hourParam filter.HourParam) (_ model.HourParamList, err error) {
	const q = `select id, val, param_id, timestamp, change_by, xml_create, manual, comment
	from hour_params where timestamp>$1 and timestamp<$2`
	rows, err := r.pool.Query(ctx, q, hourParam.DateFrom, hourParam.DateTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make(model.HourParamList, 0)
	for rows.Next() {
		var row model.HourParam
		if err = rows.Scan(
			&row.ID,
			&row.Val,
			&row.ParamID,
			&row.Timestamp,
			&row.ChangeBy,
			&row.XMLCreate,
			&row.Manual,
			&row.Comment,
		); err != nil {
			return nil, err
		}

		list = append(list, row)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return list, nil
}

func (r Repo) UpdateHourParam(ctx context.Context, hourParam model.HourParam) (err error) {
	const q = `update hour_params
	set val=$2, param_id=$3, timestamp=$4, change_by=$5, xml_create=$6, manual=$7, comment=$8
	where id = $1`
	_, err = r.pool.Exec(ctx, q, hourParam.ID, hourParam.Val, hourParam.ParamID, hourParam.Timestamp,
		hourParam.ChangeBy, hourParam.XMLCreate, hourParam.Manual, hourParam.Comment)
	if err != nil {
		return err
	}
	return nil
}

func (r Repo) DelHourParam(ctx context.Context, hourParam model.HourParam) (err error) {
	const q = `delete from hour_params where id = $1`

	_, err = r.pool.Exec(ctx, q, hourParam.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r Repo) GetHourParamList(ctx context.Context, filter filter.HourParam) (_ model.HourParamList, err error) {
	const q = `select id, val, param_id, timestamp, change_by, xml_create, manual, comment
	from hour_params %s %s`

	var (
		addQ   string
		limitQ string
		params []any
	)
	if filter.ID > 0 {
		params = append(params, filter.ID)
		addQ = addWhere(addQ, fmt.Sprintf("id = $%d", len(params)))
	}
	if filter.ParamID != 0 {
		params = append(params, filter.ParamID)
		addQ = addWhere(addQ, fmt.Sprintf("param_id = $%d", len(params)))
	}
	if filter.Limit > 0 {
		limitQ = fmt.Sprintf("limit %d", filter.Limit)
	}

	rows, err := r.pool.Query(ctx, fmt.Sprintf(q, addQ, limitQ), params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make(model.HourParamList, 0)
	for rows.Next() {
		var row model.HourParam
		if err = rows.Scan(
			&row.ID,
			&row.Val,
			&row.ParamID,
			&row.Timestamp,
			&row.ChangeBy,
			&row.XMLCreate,
			&row.Manual,
			&row.Comment,
		); err != nil {
			return nil, err
		}

		list = append(list, row)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return list, nil
}
