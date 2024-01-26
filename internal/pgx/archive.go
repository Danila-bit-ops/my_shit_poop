package pgx

import (
	"context"
	"danilamukhin/serv_go/internal/model"
	"fmt"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
)

// Архивирование
// Конфиг
func ConfigData(url string) (conf model.TableConfig, err error) {
	if _, err := toml.DecodeFile(url, &conf); err != nil {
		fmt.Println("Error decoding config:", err)
		return conf, err
	}
	return conf, nil
}

// Нахождение минимального Timestamp
func (r Repo) MinTimestamp(ctx context.Context, conf model.TableConfig) ([]string, []string, error) {
	var minTimestamp time.Time
	var addName, year, quarter []string

	// Берём имена таблиц из Config.toml
	for _, table := range conf.Tables {
		addName = append(addName, table.TableName)
	}

	// Проходимся по всем таблицам
	for i := range conf.Tables {

		// SQL
		_, err := r.pool.Exec(ctx, "SET search_path TO public;")
		if err != nil {
			fmt.Println(err)
		}

		var q = "SELECT MIN(timestamp) FROM "
		q = q + addName[i] + ";"

		_, err = r.pool.Exec(ctx, "RESET search_path;")
		if err != nil {
			fmt.Println(err)
		}

		row := r.pool.QueryRow(ctx, q)
		err = row.Scan(&minTimestamp)
		if err != nil {
			return nil, nil, err
		}

		year = append(year, strconv.Itoa(minTimestamp.Year()))
		quarter = append(quarter, strconv.Itoa((int(minTimestamp.Month())-1)/3+1))
	}
	return year, quarter, nil
}

// Создание схемы и таблицы
func (r Repo) CreateSchemaAndTable(ctx context.Context, conf model.TableConfig, year string, quarter string, tableName string, i int) (err error) {

	// SQL-запрос на создание схемы Archive,если её нет
	_, err = r.pool.Exec(ctx, `CREATE SCHEMA IF NOT EXISTS Archive;`)
	if err != nil {
		return err
	}

	// Формирование SQL-запроса CREATE TABLE
	var q = "SET search_path TO Archive; CREATE TABLE IF NOT EXISTS "

	addName := tableName + "_" + year + "_" + quarter

	// table := conf.Tables[i]
	// for _, column := range table.Columns {
	// 	// Добавление столбцов
	// 	addColumns = addColumns + column.Name + " " + column.DataType + ","
	// }

	q = q + addName + " (LIKE public." + tableName + " INCLUDING ALL);"

	// Добавление закрывающей скобки и возвращение стандартной схемы поиска
	q = q + "RESET search_path;"
	fmt.Println(q)
	_, err = r.pool.Exec(ctx, q)
	if err != nil {
		return err
	}
	return nil
}

// Перенести квартал из исходной таблицы в новую
func (r Repo) MoveQuarter(ctx context.Context, conf model.TableConfig, year string, quarter string, tableName string) (err error) {
	var q = `
	INSERT INTO %s
	SELECT *
	FROM %s
	WHERE %s
	`

	var where string
	var addYear int

	// Имя новой таблицы типа имя_год_квартал + путь через схему Archive
	newTable := "archive." + tableName + "_" + year + "_" + quarter
	// Путь к таблице hour_params
	tableName = "public." + tableName
	// В зависимости от значения квартала формируется WHERE
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
	// Полностью формируется строка SQL
	q = fmt.Sprintf(q, newTable, tableName, where)
	fmt.Println(q)
	// Отправление запроса
	_, err = r.pool.Exec(ctx, q)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// Удалить данные из исходной таблицы
func (r Repo) DeleteQuarter(ctx context.Context, conf model.TableConfig, year string, quarter string, tableName string) (err error) {
	var q = `
	DELETE FROM %s 
	WHERE %s
	`
	var where string
	var addYear int

	// Путь к таблице hour_params
	tableName = "public." + tableName
	// В зависимости от значения квартала формируется WHERE
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
	// Полностью формируется строка SQL
	q = fmt.Sprintf(q, tableName, where)
	fmt.Println(q)

	// Отправление запроса SQL
	_, err = r.pool.Exec(ctx, q)
	if err != nil {
		return err
	}

	return nil
}
