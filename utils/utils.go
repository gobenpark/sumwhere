package utils

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo/v4"
	"net/http"
	"sumwhere/factory"
)

type ApiResult struct {
	Result  interface{} `json:"result" extensions:"x-nullable,x-abc=def"`
	Success bool        `json:"success" example:"true"`
	Error   ApiError    `json:"error"`
}

type ApiError struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

type ArrayResult struct {
	Items      interface{} `json:"items"`
	TotalCount int64       `json:"totalCount"`
}

var (
	// System Error
	ApiErrorSystem             = ApiError{Code: 10001, Message: "System Error"}
	ApiErrorServiceUnavailable = ApiError{Code: 10002, Message: "Service unavailable"}
	ApiErrorRemoteService      = ApiError{Code: 10003, Message: "Remote service error"}
	ApiErrorIPLimit            = ApiError{Code: 10004, Message: "IP limit"}
	ApiErrorPermissionDenied   = ApiError{Code: 10005, Message: "Permission denied"}
	ApiErrorIllegalRequest     = ApiError{Code: 10006, Message: "Illegal request"}
	ApiErrorHTTPMethod         = ApiError{Code: 10007, Message: "HTTP method is not suported for this request"}
	ApiErrorParameter          = ApiError{Code: 10008, Message: "Parameter error"}
	ApiErrorMissParameter      = ApiError{Code: 10009, Message: "Miss required parameter"}
	ApiErrorDB                 = ApiError{Code: 10010, Message: "DB error, please contact the administator"}
	ApiErrorTokenInvaild       = ApiError{Code: 10011, Message: "Token invaild"}
	ApiErrorMissToken          = ApiError{Code: 10012, Message: "Miss token"}
	ApiErrorVersion            = ApiError{Code: 10013, Message: "API version %s invalid"}
	ApiErrorNotFound           = ApiError{Code: 10014, Message: "Resource not found"}
	ApiErrorWebSocket          = ApiError{Code: 10015, Message: "Websocket error"}
	ApiErrorRedis              = ApiError{Code: 10016, Message: "Redis error, please contact the administator"}
	ApiErrorUpdate             = ApiError{Code: 10017, Message: "Update error"}
	// Business Error
	ApiErrorUserNotExists  = ApiError{Code: 20001, Message: "User does not exists"}
	ApiErrorPassword       = ApiError{Code: 20002, Message: "Password error"}
	ApiErrorKakaoAuth      = ApiError{Code: 30001, Message: "Get Kakao Information does not exists"}
	ApiErrorFacebookAuth   = ApiError{Code: 30002, Message: "Get Facebook Information does not exists"}
	ApiErrorNotEnoughPoint = ApiError{Code: 30003, Message: "Not Enough Point"}
	ApiErrorFirebase       = ApiError{Code: 30004, Message: "Firebase Error"}
)

func ReturnApiFail(c echo.Context, status int, apiError ApiError, err error, v ...interface{}) error {
	str := ""
	if err != nil {
		str = err.Error()
	}
	return c.JSON(status, ApiResult{
		Success: false,
		Error: ApiError{
			Code:    apiError.Code,
			Message: fmt.Sprintf(apiError.Message, v...),
			Details: str,
		},
	})
}

func ReturnApiSucc(c echo.Context, status int, result interface{}) error {
	req := c.Request()
	if req.Method == "POST" || req.Method == "PUT" || req.Method == "DELETE" {
		if session, ok := factory.DB(req.Context()).(*xorm.Session); ok {
			err := session.Commit()
			if err != nil {
				return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
			}
		}
	}
	return c.JSON(status, ApiResult{
		Success: true,
		Result:  result,
	})
}
