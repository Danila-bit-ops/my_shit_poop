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

// func (r Repo) AddNewStringToDB(c *gin.Context) error {
// 	ctx := c.Request.Context()

// 	timest, err := time.Parse("2006-01-02T15:04:05", c.PostForm("AddNewTimestamp"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
// 		return err
// 	}

// 	change_by := c.PostForm("AddNewChange")

// 	comment := c.PostForm("AddNewComment")

// 	val, err := strconv.ParseFloat(c.PostForm("AddNewVal"), 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Val"})
// 		return err
// 	}

// 	paramID, err := strconv.ParseInt(c.PostForm("AddNewParamID"), 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Param ID"})
// 		return err
// 	}

// 	xmlcr, err := strconv.ParseBool(c.PostForm("AddNewXml"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid boolean format"})
// 		return err
// 	}

// 	manual, err := strconv.ParseBool(c.PostForm("AddNewManual"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid boolean format"})
// 		return err
// 	}

// 	newRecord := model.HourParam{
// 		Timestamp: timest,
// 		ChangeBy:  change_by,
// 		Comment:   comment,
// 		Val:       val,
// 		ParamID:   paramID,
// 		XMLCreate: xmlcr,
// 		Manual:    manual,
// 	}

// 	err = r.AddHourParam(ctx, newRecord)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return err
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Record added successfully"})
// }

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
