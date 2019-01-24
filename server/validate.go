package server

import (
	"github.com/gin-gonic/gin"
	validator "gopkg.in/go-playground/validator.v9"

	terr "github.com/tinklabs/golibs/error"
	"github.com/tinklabs/golibs/log"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type Common struct {
	MsgType   string `json:"msgType" validate:"required"`
	Timestamp int64  `json:"timestamp" validate:"required,gte=0"`
}

type PageInfo struct {
	PageIndex int `json:"pageIndex" validate:"required,gte=1"`
	PageSize  int `json:"pageSize" validate:"required,gte=1,lte=999"`
}

type Request struct {
	*Common `json:"common" validate:"required"`
	Param   map[string]interface{} `json:"param" validate:"required"`
}

type Response struct {
	*Common   `json:"common"`
	*PageInfo `json:"pageInfo,omitempty"`

	Total     int         `json:"total,omitempty"`
	ErrorCode int         `json:"errorCode"`
	ErrorMsg  string      `json:"errorMsg"`
	Data      interface{} `json:"data,omitempty"`
}

func GetPageInfo(c *gin.Context) (*PageInfo, *terr.TError) {
	pi := &PageInfo{
		PageIndex: c.GetInt("pageIndex"),
		PageSize:  c.GetInt("pageSize"),
	}

	if err := validate.Struct(pi); err != nil {
		log.Warn(err)
		return nil, terr.ErrRequest
	}

	return pi, nil
}

func Check() gin.HandlerFunc {

	return func(c *gin.Context) {

		if c.ContentType() == "multipart/form-data" {
			c.Next()
			return
		}

		r := Request{}
		if err := c.ShouldBindJSON(&r); err != nil {
			log.Warn(err)
			Abort(c, terr.ErrRequest)
			return
		}

		if err := validate.Struct(&r); err != nil {
			log.Warn(err)
			Abort(c, terr.ErrRequest)
			return
		}

		if v, isExist := r.Param["pageIndex"]; isExist {
			if v, ok := v.(float64); !ok {
				log.Error("pageIndex type is wrong")
				Abort(c, terr.ErrRequest)
				return
			} else {
				c.Set("pageIndex", int(v))
			}
		}

		if v, isExist := r.Param["pageSize"]; isExist {
			if v, ok := v.(float64); !ok {
				log.Error("pagSize type is wrong")
				Abort(c, terr.ErrRequest)
				return
			} else {
				c.Set("pageSize", int(v))
			}
		}

		c.Set("param", r.Param)

		c.Next()
	}
}
