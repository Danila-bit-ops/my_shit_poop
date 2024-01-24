package pgx

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"danilamukhin/serv_go/internal/model"
	"danilamukhin/serv_go/internal/pgx/filter"

	"github.com/BurntSushi/toml"
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
	from hour_params where timestamp>=$1 and timestamp<=$2 %s limit 100`
	var (
		offsetQ string
	)

	if hourParam.Offset > 0 {
		offsetQ = fmt.Sprintf("offset %d", hourParam.Offset)
	}

	rows, err := r.pool.Query(ctx, fmt.Sprintf(q, offsetQ), hourParam.DateFrom, hourParam.DateTo)
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
	from hour_params %s %s limit 100`

	var (
		addQ    string
		offsetQ string
		params  []any
	)
	if filter.ID > 0 {
		params = append(params, filter.ID)
		addQ = addWhere(addQ, fmt.Sprintf("id = $%d", len(params)))
	}
	if filter.ParamID != 0 {
		params = append(params, filter.ParamID)
		addQ = addWhere(addQ, fmt.Sprintf("param_id = $%d", len(params)))
	}
	if filter.Offset > 0 {
		offsetQ = fmt.Sprintf("offset %d", filter.Offset)
	}

	rows, err := r.pool.Query(ctx, fmt.Sprintf(q, addQ, offsetQ), params...)
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

// TEST
// Конфиг
func ConfigData(url string) (conf model.TableConfig, err error) {
	if _, err := toml.DecodeFile(url, &conf); err != nil {
		fmt.Println("Error decoding config:", err)
		return conf, err
	}
	return conf, nil
}

// Создание схемы и таблицы
func (r Repo) CreateSchemaAndTable(ctx context.Context, conf model.TableConfig, year string, quarter string) (err error) {

	// SQL-запрос на создание схемы Archive,если её нет
	_, err = r.pool.Exec(ctx, `CREATE SCHEMA IF NOT EXISTS Archive;`)
	if err != nil {
		return err
	}

	// Формирование SQL-запроса CREATE TABLE
	var q = "SET search_path TO Archive; CREATE TABLE IF NOT EXISTS "
	var addColumns string
	var addName string

	for _, table := range conf.Tables {
		addName = table.TableName + "_" + year + "_" + quarter
		for _, column := range table.Columns {
			// Добавление столбцов
			addColumns = addColumns + column.Name + " " + column.DataType + ","
		}
	}

	q = q + addName + " (" + addColumns

	// Удаление последней запятой
	qLen := len(q)
	q = q[:qLen-1]

	// Добавление закрывающей скобки и возвращение стандартной схемы поиска
	q = q + "); RESET search_path;"
	fmt.Println(q)
	_, err = r.pool.Exec(ctx, q)
	if err != nil {
		return err
	}
	return nil
}

func (r Repo) MinTimestamp(ctx context.Context, conf model.TableConfig) (time.Time, error) {
	var minTimestamp time.Time
	var addName string
	// Берём имя таблицы из Config.toml
	for _, table := range conf.Tables {
		addName = table.TableName
	}
	// SQL
	_, err := r.pool.Exec(ctx, "SET search_path TO public;")
	if err != nil {
		fmt.Println(err)
	}

	var q = "SELECT MIN(timestamp) FROM "
	q = q + addName + ";"

	_, err = r.pool.Exec(ctx, "RESET search_path;")
	if err != nil {
		fmt.Println(err)
	}

	row := r.pool.QueryRow(ctx, q)
	err = row.Scan(&minTimestamp)

	if err != nil {
		return time.Time{}, err
	}
	return minTimestamp, nil
}

func (r Repo) MoveQuarter(ctx context.Context, tableName string, year string, quarter string) (err error) {
	var q = `
	INSERT INTO %s
	SELECT *
	FROM %s
	WHERE %s
	`
	var where string
	var addYear int

	newTable := "archive." + tableName + "_" + year + "_" + quarter
	tableName = "public." + tableName

	switch quarter {
	case "1":
		where = `
		timestamp >= '%s-01-01'::timestamp
		AND timestamp < '%s-04-01'::timestamp;
		`
		where = fmt.Sprintf(where, year, year)
	case "2":
		where = `
		timestamp >= '%s-04-01'::timestamp
		AND timestamp < '%s-07-01'::timestamp;
		`
		where = fmt.Sprintf(where, year, year)
	case "3":
		where = `
		timestamp >= '%s-07-01'::timestamp
		AND timestamp < '%s-10-01'::timestamp;
		`
		where = fmt.Sprintf(where, year, year)
	case "4":
		where = `
		timestamp >= '%s-10-01'::timestamp
		AND timestamp < '%s-01-01'::timestamp;
		`

		addYear, err = strconv.Atoi(year)
		if err != nil {
			return err
		}

		addYear = addYear + 1
		nextYear := strconv.Itoa(addYear)

		where = fmt.Sprintf(where, year, nextYear)
	default:
		where = `
		timestamp >= '%s-01-01'::timestamp
		AND timestamp < '%s-04-01'::timestamp;
		`
		where = fmt.Sprintf(where, year, year)
	}

	q = fmt.Sprintf(q, newTable, tableName, where)
	fmt.Println(q)

	_, err = r.pool.Exec(ctx, q)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
