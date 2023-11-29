package entity

import (
	"context"
	"encoding/json"
	"errors"
	chiM "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

var InternalServerErrorResponse, _ = json.Marshal(LogicError{ResponseMessage: "internal server error"})
var NotFoundError = &LogicError{ResponseMessage: "not found", Code: 404}
var URLNotEnteredError = &LogicError{ResponseMessage: "url not entered", Code: 404}
var ShortURLNotEnteredError = &LogicError{ResponseMessage: "shortURL not entered", Code: 404}

type LogicError struct {
	ResponseMessage string `json:"error"`
	Err             error  `json:"-"`
	Code            int    `json:"-"`
}

func NewLogicError(err error, responseMessage string, code int) *LogicError {
	return &LogicError{
		ResponseMessage: responseMessage,
		Err:             err,
		Code:            code,
	}
}

func (e *LogicError) Error() string {
	if e == nil || e.Err == nil {
		return ""
	}

	return e.Err.Error()
}

func (e *LogicError) JsonMarshal() []byte {
	if e == nil || len(e.ResponseMessage) == 0 {
		return nil
	}

	b, _ := json.Marshal(e)
	return b
}

func HandleError(ctx context.Context, logger *zap.Logger, err error) ([]byte, int) {
	var logicErr *LogicError
	code := 500
	if errors.As(err, &logicErr) {
		code = logicErr.Code
		logger.Error(logicErr.Error(), zap.String("RequestId", chiM.GetReqID(ctx)), zap.Int("ResponseCode", code))
		return logicErr.JsonMarshal(), logicErr.Code
	}

	logger.Error(err.Error(), zap.String("RequestId", chiM.GetReqID(ctx)), zap.Int("ResponseCode", code))

	return InternalServerErrorResponse, code
}
