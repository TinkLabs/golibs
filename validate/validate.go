package validate

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"

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
	RequestID string `json:"requestID" validate:"required"`
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
		r := Request{}

		if err := c.ShouldBindJSON(&r); err != nil {
			log.Warn(err)
			abort(c, terr.ErrRequest)
			return
		}

		if err := validate.Struct(&r); err != nil {
			log.Warn(err)
			abort(c, terr.ErrRequest)
			return
		}

		c.Set("requestID", r.Common.RequestID)

		if v, isExist := r.Param["pageIndex"]; isExist {
			if v, ok := v.(float64); !ok {
				log.Error("pageIndex type is wrong")
				abort(c, terr.ErrRequest)
				return
			} else {
				c.Set("pageIndex", int(v))
			}
		}

		if v, isExist := r.Param["pageSize"]; isExist {
			if v, ok := v.(float64); !ok {
				log.Error("pagSize type is wrong")
				abort(c, terr.ErrRequest)
				return
			} else {
				c.Set("pageSize", int(v))
			}
		}

		c.Set("param", r.Param)

		c.Next()
	}
}

func abort(c *gin.Context, err *terr.TError) {
	ts := time.Now().UnixNano() / 1000000

	c.AbortWithStatusJSON(http.StatusOK, &Response{
		Common: &Common{
			MsgType:   "response",
			Timestamp: ts,
			RequestID: c.GetString("requestID"),
		},
		ErrorCode: int(err.Code),
		ErrorMsg:  err.Desc,
	})
}
