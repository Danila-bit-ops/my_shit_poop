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

	}
}

func (a *api) TableHourParam(c *gin.Context) {
	ctx := c.Request.Context()
	// inputData := c.Query("FilterID")
	// fmt.Println("Полученные данные:", inputData)
	// ID, err := strconv.ParseInt(inputData, 10, 64)
	// if err != nil {
	// 	log.Err(err).Msg("ParseInt")
	// 	c.Status(http.StatusBadRequest)
	// 	return
	// }
	list, err := a.srv.GetHourParamList(ctx, filter.HourParam{Limit: 200})

	if err != nil {
		log.Err(err).Msg("GetHourParamList")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"list": list})
}

func (a *api) LoadIndexHTML(c *gin.Context) {
	// type Data struct {
	// 	Test string
	// }

	// ctx := c.Request.Context()

	// list, err := a.srv.GetHourParamList(ctx, filter.HourParam{Limit: 100})
	// if err != nil {
	// 	log.Err(err).Msg("GetHourParamList")
	// 	c.Status(http.StatusInternalServerError)
	// 	return
	// }

	c.HTML(http.StatusOK, "index.html", nil)
	// c.JSON(http.StatusOK, gin.H{"list": list})
}

// Поиск по Param_ID
func (a *api) GetByParamID(c *gin.Context) {
	ctx := c.Request.Context()
	sParamID := c.Query("FilterParamID")
	fmt.Println(sParamID)
	paramID, err := strconv.ParseInt(sParamID, 10, 64)
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
		list, err := a.srv.GetHourParamList(ctx, filter.HourParam{ParamID: paramID, Limit: Limit})
		if err != nil {
			log.Err(err).Msg("GetHourParamList")
			c.Status(http.StatusInternalServerError)
			return
		}
		// fmt.Println(list)

		c.JSON(http.StatusOK, gin.H{"list": list})
	} else {
		list, err := a.srv.GetHourParamList(ctx, filter.HourParam{ParamID: paramID})
		if err != nil {
			log.Err(err).Msg("GetHourParamList")
			c.Status(http.StatusInternalServerError)
			return
		}
		// fmt.Println(list)

		c.JSON(http.StatusOK, gin.H{"list": list})
	}

}

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

// Удалить запись
func (a *api) RemoveRecord(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseInt(c.PostForm("DelID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	err = a.srv.DeleteHourParam(ctx, model.HourParam{
		ID: id,
	})
	if err != nil {
		log.Err(err).Msg("DeleteHourParam")
		c.Status(http.StatusInternalServerError)
		return
	}
	list, err := a.srv.GetHourParamList(ctx, filter.HourParam{Limit: 50})
	if err != nil {
		log.Err(err).Msg("GetHourParamList")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{"list": list})
}

// Редактировать запись
func (a *api) UpdateRecord(c *gin.Context) {
	ctx := c.Request.Context()

	Id, err := strconv.ParseInt(c.PostForm("UpdID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Param ID"})
		return
	}

	timestamp, err := time.Parse("2006-01-02T15:04", c.PostForm("UpdTimestamp"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
		return
	}

	changeBy := c.PostForm("UpdChange")

	comment := c.PostForm("UpdComment")

	val, err := strconv.ParseFloat(c.PostForm("UpdVal"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Val"})
		return
	}

	paramID, err := strconv.ParseInt(c.PostForm("UpdParamID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Param ID"})
		return
	}
	xmlcr, err := strconv.ParseBool(c.PostForm("UpdXml"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid boolean format"})
		return
	}

	manual, err := strconv.ParseBool(c.PostForm("UpdManual"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid boolean format"})
		return
	}

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

	list, err := a.srv.GetHourParamList(ctx, filter.HourParam{Limit: 50})
	if err != nil {
		log.Err(err).Msg("GetHourParamList")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{"list": list})
}

// Найти записи в определённом временном интервале
func (a *api) FindRecordByRange(c *gin.Context) {
	ctx := c.Request.Context()
	// fmt.Println(c.Query("RngStart"))
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

	lim := c.Query("Limit")
	if lim != "" {
		Limit, err := strconv.ParseInt(lim, 10, 64)
		if err != nil {
			log.Err(err).Msg("ParseInt")
			c.Status(http.StatusBadRequest)
			return
		}
		list, err := a.srv.GetHourParamList(ctx, filter.HourParam{DateFrom: timestampStart, DateTo: timestampEnd, Limit: Limit})
		if err != nil {
			log.Err(err).Msg("GetHourParamList")
			c.Status(http.StatusInternalServerError)
			return
		}
		// fmt.Println(list)

		c.JSON(http.StatusOK, gin.H{"list": list})
	} else {
		list, err := a.srv.GetHourParamList(ctx, filter.HourParam{DateFrom: timestampStart, DateTo: timestampEnd})
		if err != nil {
			log.Err(err).Msg("GetHourParamList")
			c.Status(http.StatusInternalServerError)
			return
		}
		// fmt.Println(list)

		c.JSON(http.StatusOK, gin.H{"list": list})
	}
}

// Добавить запись
func (a *api) AddNewRecord(c *gin.Context) {
	ctx := c.Request.Context()
	timestamp, err := time.Parse("2006-01-02T15:04", c.PostForm("AddNewTimestamp"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
		return
	}

	changeBy := c.PostForm("AddNewChange")

	comment := c.PostForm("AddNewComment")

	val, err := strconv.ParseFloat(c.PostForm("AddNewVal"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Val"})
		return
	}

	paramID, err := strconv.ParseInt(c.PostForm("AddNewParamID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Param ID"})
		return
	}
	xmlcr, err := strconv.ParseBool(c.PostForm("AddNewXml"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid boolean format"})
		return
	}

	manual, err := strconv.ParseBool(c.PostForm("AddNewManual"))
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

	list, err := a.srv.GetHourParamList(ctx, filter.HourParam{Limit: 50})
	if err != nil {
		log.Err(err).Msg("GetHourParamList")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{"list": list})
}
