package api

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"danilamukhin/serv_go/internal/model"
	"danilamukhin/serv_go/internal/pgx/filter"
	"danilamukhin/serv_go/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type HourParamRequest struct {
	Timestamp string `json:"AddNewTimestamp"`
	ChangeBy  string `json:"AddNewChange"`
	Comment   string `json:"AddNewComment"`
	Val       string `json:"AddNewVal"`
	ParamID   string `json:"AddNewParamID"`
	XMLCreate string `json:"AddNewXml"`
	Manual    string `json:"AddNewManual"`
}
type HourParamUpd struct {
	ID        string `json:"UpdID"`
	Timestamp string `json:"UpdTimestamp"`
	ChangeBy  string `json:"UpdChange"`
	Comment   string `json:"UpdComment"`
	Val       string `json:"UpdVal"`
	ParamID   string `json:"UpdParamID"`
	XMLCreate string `json:"UpdXml"`
	Manual    string `json:"UpdManual"`
}

func InitApi(srv *service.Service) *api {
	return &api{
		srv: srv,
	}
}

type api struct {
	srv *service.Service
}

func (a *api) InitRouter() *gin.Engine {
	router := gin.Default()
	a.initHandlers(router)
	return router
}

func (a *api) initHandlers(r *gin.Engine) {
	cmdDir, _ := filepath.Abs(filepath.Dir("./cmd"))
	assetsDir := filepath.Join(cmdDir, ".", "assets")
	r.LoadHTMLGlob(filepath.Join(assetsDir, "public/*.html"))
	r.GET("/index", a.LoadIndexHTML)
	r.Static("/assets/src", filepath.Join(assetsDir, "src"))

	// api
	api := r.Group("/api")
	{
		api.GET("/get-by-param-id", a.GetByParamID)
		api.GET("/get-by-id", a.GetByID)
		api.POST("/add-new-record", a.AddNewRecord)
		api.POST("/del-by-id", a.RemoveRecord)
		api.POST("/update-by-id", a.UpdateRecord)
		api.GET("/get-range", a.FindRecordByRange)
		api.GET("/table-hour-params", a.TableHourParam)
		api.GET("/lazy-loading", a.LazyLoading)
		// api.GET("/test", a.CreateArchiveTable)
	}
}

// Lazy-Loading
func (a *api) LazyLoading(c *gin.Context) {
	ctx := c.Request.Context()

	sOffset := c.Query("Offset")
	Offset, err := strconv.ParseInt(sOffset, 10, 64)
	if err != nil {
		log.Err(err).Msg("ParseInt")
		c.Status(http.StatusBadRequest)
		return
	}

	sFilter := c.Query("FilterParamID")

	if c.Query("RngStart") != "" {
		timestampStart, err := time.Parse("2006-01-02T15:04", c.Query("RngStart"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
			return
		}

		timestampEnd, err := time.Parse("2006-01-02T15:04", c.Query("RngEnd"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
			return
		}

		list, err := a.srv.RangeHourParam(ctx, filter.HourParam{DateFrom: timestampStart, DateTo: timestampEnd, Offset: Offset})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"list": list})

	} else if sFilter != "" {
		Param_ID, err := strconv.ParseInt(sFilter, 10, 64)

		if err != nil {
			log.Err(err).Msg("ParseInt")
			c.Status(http.StatusBadRequest)
			return
		}

		list, err := a.srv.GetHourParamList(ctx, filter.HourParam{Offset: Offset, ParamID: Param_ID})
		if err != nil {
			log.Err(err).Msg("GetHourParamList")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"list": list})
	} else {
		list, err := a.srv.GetHourParamList(ctx, filter.HourParam{Offset: Offset})
		if err != nil {
			log.Err(err).Msg("GetHourParamList")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"list": list})
	}
}

//End

func (a *api) TableHourParam(c *gin.Context) {
	ctx := c.Request.Context()
	list, err := a.srv.GetHourParamList(ctx, filter.HourParam{Limit: 200})

	if err != nil {
		log.Err(err).Msg("GetHourParamList")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"list": list})
}

func (a *api) LoadIndexHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// Поиск по Param_ID
func (a *api) GetByParamID(c *gin.Context) {
	ctx := c.Request.Context()
	sParamID := c.Query("FilterParamID")
	paramID, err := strconv.ParseInt(sParamID, 10, 64)
	if err != nil {
		log.Err(err).Msg("ParseInt")
		c.Status(http.StatusBadRequest)
		return
	}
	list, err := a.srv.GetHourParamList(ctx, filter.HourParam{ParamID: paramID})
	if err != nil {
		log.Err(err).Msg("GetHourParamList")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": list})
}

// Конец

// Поиск по ID
func (a *api) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	sID := c.Query("FilterID")
	fmt.Println("sID =  " + sID)
	ID, err := strconv.ParseInt(sID, 10, 64)
	if err != nil {
		log.Err(err).Msg("ParseInt")
		c.Status(http.StatusBadRequest)
		return
	}

	lim := c.Query("Limit")
	if lim != "" {
		Limit, err := strconv.ParseInt(lim, 10, 64)
		if err != nil {
			log.Err(err).Msg("ParseInt")
			c.Status(http.StatusBadRequest)
			return
		}
		list, err := a.srv.GetHourParamList(ctx, filter.HourParam{ID: ID, Limit: Limit})
		if err != nil {
			log.Err(err).Msg("GetHourParamList")
			c.Status(http.StatusInternalServerError)
			return
		}
		// fmt.Println(list)

		c.JSON(http.StatusOK, gin.H{"list": list})
	} else {
		list, err := a.srv.GetHourParamList(ctx, filter.HourParam{ID: ID})
		if err != nil {
			log.Err(err).Msg("GetHourParamList")
			c.Status(http.StatusInternalServerError)
			return
		}
		// fmt.Println(list)

		c.JSON(http.StatusOK, gin.H{"list": list})
	}
}

// Конец

// Удалить запись
func (a *api) RemoveRecord(c *gin.Context) {
	ctx := c.Request.Context()

	var data struct {
		ID int64 `json:"ID"`
	}
	// Парсим JSON-данные из тела запроса
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	list, err := a.srv.GetHourParamList(ctx, filter.HourParam{ID: data.ID})
	if err != nil {
		log.Err(err).Msg("GetHourParamList")
		c.Status(http.StatusInternalServerError)
		return
	}

	fmt.Println(len(list))
	if len(list) != 0 {
		err := a.srv.DeleteHourParam(ctx, model.HourParam{
			ID: data.ID,
		})
		if err != nil {
			log.Err(err).Msg("DeleteHourParam")
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Данные успешно удалены из базы данных"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Нет такого ID"})
	}
}

// Конец

// Редактировать запись
func (a *api) UpdateRecord(c *gin.Context) {
	ctx := c.Request.Context()
	var data HourParamUpd
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	Id, err := strconv.ParseInt(data.ID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Param ID"})
		return
	}

	timestamp, err := time.Parse("2006-01-02T15:04", data.Timestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
		return
	}

	changeBy := data.ChangeBy

	comment := data.Comment

	val, err := strconv.ParseFloat(data.Val, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Val"})
		return
	}

	paramID, err := strconv.ParseInt(data.ParamID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Param ID"})
		return
	}
	xmlcr, err := strconv.ParseBool(data.XMLCreate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid boolean format"})
		return
	}

	manual, err := strconv.ParseBool(data.Manual)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid boolean format"})
		return
	}

	list, err := a.srv.GetHourParamList(ctx, filter.HourParam{ID: Id})
	if err != nil {
		log.Err(err).Msg("GetHourParamList")
		c.Status(http.StatusInternalServerError)
		return
	}
	if len(list) != 0 {
		err = a.srv.UpdateHourParam(ctx, model.HourParam{
			ID:        Id,
			Timestamp: timestamp,
			ChangeBy:  changeBy,
			Comment:   comment,
			Val:       val,
			ParamID:   paramID,
			XMLCreate: xmlcr,
			Manual:    manual,
		})

		if err != nil {
			log.Err(err).Msg("UpdateHourParam")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Данные успешно изменены"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Нет такого ID"})
	}
}

// Конец

// Найти записи в определённом временном интервале
func (a *api) FindRecordByRange(c *gin.Context) {
	ctx := c.Request.Context()

	timestampStart, err := time.Parse("2006-01-02T15:04", c.Query("RngStart"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
		return
	}

	timestampEnd, err := time.Parse("2006-01-02T15:04", c.Query("RngEnd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
		return
	}

	list, err := a.srv.RangeHourParam(ctx, filter.HourParam{DateFrom: timestampStart, DateTo: timestampEnd})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": list})
}

// Конец

// Добавить запись
func (a *api) AddNewRecord(c *gin.Context) {
	ctx := c.Request.Context()
	var data HourParamRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	timestamp, err := time.Parse("2006-01-02T15:04", data.Timestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
		return
	}
	changeBy := data.ChangeBy
	comment := data.Comment
	val, err := strconv.ParseFloat(data.Val, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Val"})
		return
	}
	paramID, err := strconv.ParseInt(data.ParamID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Param ID"})
		return
	}
	xmlcr, err := strconv.ParseBool(data.XMLCreate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid boolean format"})
		return
	}
	manual, err := strconv.ParseBool(data.Manual)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid boolean format"})
		return
	}

	err = a.srv.CreateHourParam(ctx, model.HourParam{
		Timestamp: timestamp,
		ChangeBy:  changeBy,
		Comment:   comment,
		Val:       val,
		ParamID:   paramID,
		XMLCreate: xmlcr,
		Manual:    manual,
	})
	if err != nil {
		log.Err(err).Msg("CreateHourParam")
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Данные успешно добавлены в базу данных"})
}

// Конец
