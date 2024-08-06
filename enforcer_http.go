package aibalance

import (
	"encoding/json"
	"fmt"
	"github.com/belief428/aigw-balance/model"
	"github.com/belief428/aigw-balance/persist"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type Context struct {
	*http.Request
	http.ResponseWriter
}

type Page struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (this *Response) Marshal() []byte {
	_bytes, _ := json.Marshal(this)
	return _bytes
}

func (this *Response) Write() {

}

// getParams 获取参数信息
func getParams(enforcer *Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, &Response{
			Message: "ok",
			Data:    enforcer.GetParams(),
		})
	}
}

// setParams 设置参数信息
func setParams(enforcer *Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := &Response{Code: 0, Message: "ok"}

		_params := make(map[string]interface{}, 0)

		_bytes, err := io.ReadAll(c.Request.Body)

		if err != nil {
			resp.Code = -1
			resp.Message = err.Error()
			goto LOOP
		}
		json.Unmarshal(_bytes, &_params)

		if err = enforcer.SetParams(_params); err != nil {
			resp.Code = -1
			resp.Message = err.Error()
		}
	LOOP:
		if enforcer.watcher != nil && enforcer.watcher.GetParamsCallbackFunc() != nil {
			enforcer.watcher.GetParamsCallbackFunc()(_params)
		}
		c.JSON(http.StatusOK, resp)
	}
}

// getArchive 获取档案信息
func getArchive(enforcer *Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := &Response{Code: 0, Message: "ok"}

		_params := &struct {
			GatewayCode string `json:"gateway_code" form:"gateway_code"`
			Kind        int    `json:"kind" form:"kind"`
		}{}
		if err := c.ShouldBindQuery(_params); err != nil {
			resp.Code = -1
			resp.Message = err.Error()
			c.JSON(http.StatusOK, resp)
			return
		}
		if enforcer.watcher == nil || enforcer.watcher.GetArchiveFunc() == nil {
			c.JSON(http.StatusOK, resp)
			return
		}
		archives := enforcer.watcher.GetArchiveFunc()(&persist.WatcherArchiveParams{
			Code: _params.GatewayCode,
			Kind: _params.Kind,
		})
		for _, v := range archives {
			attribute := EnforcerArchive(enforcer.archives).filter(_params.GatewayCode, v.GetCode())
			v.SetRegulate(attribute.Regulate > 0)
			v.SetWeight(attribute.Weight)
		}
		resp.Data = archives
		c.JSON(http.StatusOK, resp)
	}
}

// setArchive 设置档案信息
func setArchive(enforcer *Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := &Response{Code: 0, Message: "ok"}

		_params := &struct {
			GatewayCode string   `json:"gateway_code"`
			Codes       []string `json:"codes"`
			model.ArchiveAttribute
		}{}
		_bytes, err := io.ReadAll(c.Request.Body)

		if err != nil {
			resp.Code = -1
			resp.Message = err.Error()
			c.JSON(http.StatusOK, resp)
			return
		}
		json.Unmarshal(_bytes, &_params)

		if _params.GatewayCode == "" || len(_params.Codes) <= 0 {
			resp.Code = -1
			resp.Message = "网关编号/档案编号不能为空"
			c.JSON(http.StatusOK, resp)
			return
		}
		archives := make([]*model.Archive, 0)
		mArchive := new(model.Archive)
		// 查询数据库是否存在
		err = enforcer.engine.Table((&model.Archive{}).TableName()).Where("gateway_code = ?", _params.GatewayCode).
			Where("code IN (?)", _params.Codes).Find(&archives).Error

		if err != nil {
			resp.Code = -1
			resp.Message = err.Error()
			c.JSON(http.StatusOK, resp)
			return
		}
		now := time.Now()

		tx := enforcer.engine.Begin()

		_archives := make(map[string]*model.Archive, 0)

		for _, v := range archives {
			_archives[v.Code] = v
		}
		list := make([]*model.Archive, 0)

		for _, v := range _params.Codes {
			if data, has := _archives[v]; has {
				// 更新数据库
				data.Attribute = _params.ArchiveAttribute
				data.UpdatedAt = now

				if err = tx.Table(mArchive.TableName()).Where("id = ?", data.ID).Updates(data).Error; err != nil {
					tx.Rollback()
					resp.Code = -1
					resp.Message = err.Error()
					c.JSON(http.StatusOK, resp)
					return
				}
				continue
			}
			list = append(list, &model.Archive{GatewayCode: _params.GatewayCode, Code: v,
				Attribute: _params.ArchiveAttribute, CreatedAt: now, UpdatedAt: now,
			})
		}
		if len(list) > 0 {
			if err = tx.Table(mArchive.TableName()).CreateInBatches(list, 50).Error; err != nil {
				tx.Rollback()
				resp.Code = -1
				resp.Message = err.Error()
				c.JSON(http.StatusOK, resp)
				return
			}
		}
		tx.Commit()

		for _, v := range _params.Codes {
			_, has := enforcer.archives[_params.GatewayCode]

			if !has {
				enforcer.archives[_params.GatewayCode] = map[string]model.ArchiveAttribute{v: _params.ArchiveAttribute}
				continue
			}
			enforcer.archives[_params.GatewayCode][v] = _params.ArchiveAttribute
		}
		c.JSON(http.StatusOK, resp)
	}
}

func setArchiveDeg(enforcer *Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := &Response{Code: 0, Message: "ok"}

		_params := &struct {
			GatewayCode string   `json:"gateway_code"`
			Kind        int      `json:"kind" form:"kind"`
			Codes       []string `json:"codes"`
			Deg         uint8    `json:"deg"`
		}{}
		_bytes, err := io.ReadAll(c.Request.Body)

		if err != nil {
			resp.Code = -1
			resp.Message = err.Error()
			c.JSON(http.StatusOK, resp)
			return
		}
		json.Unmarshal(_bytes, &_params)

		if _params.GatewayCode == "" || len(_params.Codes) <= 0 {
			resp.Code = -1
			resp.Message = "网关编号/档案编号不能为空"
			c.JSON(http.StatusOK, resp)
			return
		}
		if enforcer.watcher != nil && enforcer.watcher.GetRegulateCallbackFunc() != nil {
			for _, v := range _params.Codes {
				_ = enforcer.watcher.GetRegulateCallbackFunc()(&persist.WatcherRegulateParams{
					Code:        _params.GatewayCode,
					ArchiveCode: v,
					Kind:        _params.Kind,
					Value:       _params.Deg,
				})
			}
		}
		c.JSON(http.StatusOK, resp)
	}
}

// getRegulate 获取调控历史信息
func getRegulate(enforcer *Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := &Response{Code: -1, Message: "ok", Data: struct {
			Data  interface{} `json:"data"`
			Count int64       `json:"count"`
		}{}}
		_params := &struct {
			Kind int `json:"kind" form:"kind"`
			Page
		}{}
		err := c.ShouldBindQuery(_params)

		if err != nil {
			resp.Code = -1
			resp.Message = err.Error()
			c.JSON(http.StatusOK, resp)
			return
		}
		var iModel persist.IModel

		var out interface{}

		if _params.Kind == 1 {
			iModel = new(model.RegulateHouse)
			out = make([]*model.RegulateHouse, 0)
		} else {
			iModel = new(model.RegulateBuild)
			out = make([]*model.RegulateBuild, 0)
		}
		query := enforcer.engine.Table(iModel.TableName())

		var count int64

		if err = query.Count(&count).Error; err != nil {
			resp.Message = err.Error()
			return
		}
		if err = query.Offset((_params.Page.Page - 1) * _params.Limit).Limit(_params.Limit).Find(&out).Error; err != nil {
			resp.Message = err.Error()
			c.JSON(http.StatusOK, resp)
			return
		}
		resp.Code = 0
		resp.Data = struct {
			Data  interface{} `json:"data"`
			Count int64       `json:"count"`
		}{Data: out, Count: count}
		c.JSON(http.StatusOK, resp)
	}
}

func (this *Enforcer) http() {
	defer func() {
		if err := recover(); err != nil {
			this.errorf("Aigw-balance http recover error：%v", err)
		}
		go this.http()
	}()
	gin.SetMode("release")
	app := gin.Default()
	app.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers, application/octet-stream, text/event-stream")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "HEAD, POST, GET, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
	_catalogue := "dist"

	static.Serve("/", static.LocalFile(path.Join(_catalogue, "index.html"), true))
	app.StaticFile("/favicon.ico", path.Join(_catalogue, "favicon.ico"))
	app.StaticFS("/static", http.Dir(path.Join(_catalogue, "static")))

	app.NoRoute(func(c *gin.Context) {
		accept := c.Request.Header.Get("Accept")
		flag := strings.Contains(accept, "text/html")

		if flag {
			content, err := os.ReadFile("dist/index.html")
			if err != nil {
				c.Writer.WriteHeader(404)
				c.Writer.WriteString("Not Found")
				return
			}
			c.Writer.WriteHeader(200)
			c.Writer.Header().Add("Accept", "text/html")
			c.Writer.Write(content)
			c.Writer.Flush()
		}
	})
	// 注册API
	v1 := app.Group("/api/v1")
	{
		v1.GET("/params", getParams(this))
		v1.POST("/params/set", setParams(this))
		v1.GET("/archive", getArchive(this))
		v1.POST("/archive/set", setArchive(this))
		v1.POST("/archive/set_deg", setArchiveDeg(this))
		v1.GET("/regulate", getRegulate(this))
	}
	serve := &http.Server{
		Addr:           fmt.Sprintf(":%d", this.port),
		Handler:        app,
		MaxHeaderBytes: 1 << 20,
	}
	_ = serve.ListenAndServe()
}
