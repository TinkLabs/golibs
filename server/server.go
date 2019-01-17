package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tinklabs/golibs/cmd"
	"github.com/tinklabs/golibs/config"
	"github.com/tinklabs/golibs/consul"
	"github.com/tinklabs/golibs/db"
	terr "github.com/tinklabs/golibs/error"
	"github.com/tinklabs/golibs/log"
	"github.com/tinklabs/golibs/utils"
)

var (
	Server *http.Server
	router *gin.Engine
)

func Init() {
	cmd.Init()
	log.Init()
	consul.Init()
	config.Init()
	db.Init()

	if !cmd.IsDebug() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(log.Logger())
	r.Use(Check())

	router = r
}

func GetRouter() *gin.Engine {
	return router
}

func Start() {
	cc := consul.GetConsulClient()
	cc.Register()

	cf := cmd.GetCmdFlag()
	Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", cf.ServerPort),
		Handler: router,
	}

	log.Info("Server is listening on ", cf.ServerPort)
	if err := Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		cc.Deregister()
		log.Fatal(err)
	}
}

func Stop() {
	cc := consul.GetConsulClient()
	cc.Deregister()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := Server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Info("Server exiting")
}

func Register(version, method, source string, callback func(*gin.Context)) {
	cf := cmd.GetCmdFlag()
	url := fmt.Sprintf("/api/%s/%s%s", cf.ServerName, version, source)

	switch method {
	case "GET":
		router.GET(url, callback)
	case "POST":
		router.POST(url, callback)
	case "PUT":
		router.PUT(url, callback)
	case "PATCH":
		router.PATCH(url, callback)
	case "DELETE":
		router.DELETE(url, callback)
	default:
		panic(fmt.Sprintf("unsupported method: %s", method))
	}
}

func OK(c *gin.Context, data interface{}) {
	var pi *PageInfo

	i, s := c.GetInt("pageIndex"), c.GetInt("pageSize")
	if i > 0 && s > 0 {
		pi = &PageInfo{
			PageIndex: i,
			PageSize:  s,
		}
	}
	c.JSON(http.StatusOK, &Response{
		Common: &Common{
			MsgType:   "response",
			Timestamp: utils.GetNowTs(),
			RequestID: c.GetString("requestID"),
		},
		PageInfo:  pi,
		Total:     c.GetInt("total"),
		ErrorCode: 0,
		ErrorMsg:  "",
		Data:      data,
	})
}

func Fail(c *gin.Context, err *terr.TError) {
	c.JSON(http.StatusOK, &Response{
		Common: &Common{
			MsgType:   "response",
			Timestamp: utils.GetNowTs(),
			RequestID: c.GetString("requestID"),
		},
		ErrorCode: int(err.Code),
		ErrorMsg:  err.Desc,
	})
}

func Abort(c *gin.Context, err *terr.TError) {
	c.AbortWithStatusJSON(http.StatusOK, &Response{
		Common: &Common{
			MsgType:   "response",
			Timestamp: utils.GetNowTs(),
			RequestID: c.GetString("requestID"),
		},
		ErrorCode: int(err.Code),
		ErrorMsg:  err.Desc,
	})
}
