package aibalance

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/belief428/aigw-balance/model"
	"github.com/belief428/aigw-balance/persist"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	Page  int `json:"page"`
	Limit int `json:"limit"`
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
func getArchive(enforcer *Enforcer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := &Response{Code: 0, Message: "ok"}

		_params := &struct {
			GatewayCode string `json:"gateway_code"`
			Kind        int    `json:"kind"`
		}{}
		_bytes, err := io.ReadAll(r.Body)

		if err != nil {
			resp.Code = -1
			resp.Message = err.Error()
			w.Write(resp.Marshal())
			return
		}
		json.Unmarshal(_bytes, &_params)

		if enforcer.watcher == nil || enforcer.watcher.GetArchiveFunc() == nil {
			w.Write(resp.Marshal())
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
		w.Write(resp.Marshal())
	}
}

func setArchive(enforcer *Enforcer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := &Response{Code: 0, Message: "ok"}

		_params := &struct {
			GatewayCode string `json:"gateway_code"`
			Code        string `json:"code"`
			model.ArchiveAttribute
		}{}
		_bytes, err := io.ReadAll(r.Body)

		if err != nil {
			resp.Code = -1
			resp.Message = err.Error()
			w.Write(resp.Marshal())
			return
		}
		json.Unmarshal(_bytes, &_params)

		if _params.GatewayCode == "" || _params.Code == "" {
			resp.Code = -1
			resp.Message = "网关编号/档案编号不能为空"
			w.Write(resp.Marshal())
			return
		}
		mArchive := new(model.Archive)
		// 查询数据库是否存在
		err = enforcer.engine.Table(mArchive.TableName()).Where("gateway_code = ?", _params.GatewayCode).Where("code = ?", _params.Code).First(mArchive).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Code = -1
			resp.Message = err.Error()
			w.Write(resp.Marshal())
			return
		}
		now := time.Now()

		if mArchive.ID > 0 {
			mArchive.Attribute = _params.ArchiveAttribute
			mArchive.UpdatedAt = now

			if err = enforcer.engine.Table(mArchive.TableName()).Save(mArchive).Error; err != nil {
				resp.Code = -1
				resp.Message = err.Error()
				w.Write(resp.Marshal())
				return
			}
		} else {
			mArchive.GatewayCode = _params.GatewayCode
			mArchive.Code = _params.Code
			mArchive.Attribute = _params.ArchiveAttribute
			mArchive.CreatedAt = now
			mArchive.UpdatedAt = now

			if err = enforcer.engine.Table(mArchive.TableName()).Create(mArchive).Error; err != nil {
				resp.Code = -1
				resp.Message = err.Error()
				w.Write(resp.Marshal())
				return
			}
		}
		_, has := enforcer.archives[_params.GatewayCode]

		if !has {
			enforcer.archives[_params.GatewayCode] = map[string]model.ArchiveAttribute{
				_params.GatewayCode: _params.ArchiveAttribute,
			}
		} else {
			enforcer.archives[_params.GatewayCode][_params.Code] = _params.ArchiveAttribute
		}
		w.Write(resp.Marshal())
	}
}

// getHorizontalHistory 获取水平调控历史信息
func getHorizontalHistory(enforcer *Enforcer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := &Response{Code: -1, Message: "ok", Data: make([]string, 0)}

		pagination := &Page{Page: 1, Limit: 10}

		_bytes, err := io.ReadAll(r.Body)

		if err != nil {
			resp.Message = err.Error()
			return
		}
		json.Unmarshal(_bytes, pagination)

		if pagination.Page <= 0 {
			pagination.Page = 1
		}
		query := enforcer.engine.Table((&model.RegulateBuild{}).TableName())

		var count int64

		if err = query.Count(&count).Error; err != nil {
			resp.Message = err.Error()
			return
		}
		out := make([]*model.RegulateBuild, 0)

		if err = query.Offset((pagination.Page - 1) * pagination.Limit).Limit(pagination.Limit).Find(&out).Error; err != nil {
			resp.Message = err.Error()
			w.Write(resp.Marshal())
			return
		}
		resp.Code = 0
		resp.Data = struct {
			Data  interface{} `json:"data"`
			Count int64       `json:"count"`
		}{Data: out, Count: count}
		w.Write(resp.Marshal())
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
	}

	serve := &http.Server{
		Addr:           fmt.Sprintf(":%d", this.port),
		Handler:        app,
		MaxHeaderBytes: 1 << 20,
	}
	_ = serve.ListenAndServe()
}
