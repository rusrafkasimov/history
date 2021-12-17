package errs

import (
	"github.com/gin-gonic/gin"
	"github.com/rusrafkasimov/history/internal/types"
	"github.com/rusrafkasimov/history/pkg/dto"
)

func ErrorHandler(c *gin.Context, err types.ErrorWithCode) {

	// TODO Evaluating the level of logging API errors

	//switch err.ErrorCode() {
	//case 400:
	//	log.Infof("Bad request. %s", err)
	//case 404:
	//	log.Infof("Record not found")
	//case 418:
	//	log.Infof("Teapot error. %s", err)
	//case 422:
	//	log.Infof("Unprocessable entity error. %s", err)
	//case 500:
	//	log.Errorf("Internal server error. %s", err)
	//default:
	//	log.Warnf("Unexpected error. %s", err)
	//}

	c.JSON(err.ErrorCode(), ConvertError(err))
}

func ConvertError(errApi types.ErrorWithCode) *dto.Error {
	return &dto.Error{
		Status:  errApi.ErrorCode(),
		Message: errApi.ErrorMessage(),
		Error:   errApi.Error(),
	}
}