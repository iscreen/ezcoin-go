package app

import (
	"fmt"
	"strings"

	"ezcoin.cc/ezcoin-go/server/app/model/common/request"
	"ezcoin.cc/ezcoin-go/server/global"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	val "github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

type ValidError struct {
	Key     string
	Message string
}
type ValidErrors []*ValidError

func (v *ValidError) Error() string {
	return v.Message
}

func (v ValidErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

func (v ValidErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

// Get error fields message
// key is been converted to camelcase
// value is been removed key in content
func (v ValidErrors) ErrorFields() map[string]interface{} {
	errs := map[string]interface{}{}
	for _, err := range v {
		keys := strings.Split(err.Key, ".")
		key := strcase.ToLowerCamel(keys[len(keys)-1])
		errs[key] = []string{strings.Trim(err.Message, key)}
	}
	return errs
}

func BindTableQuery(c *gin.Context, req *request.TableQuery) error {
	if err := c.ShouldBindQuery(req); err != nil {
		return err
	}
	req.Filter = c.QueryMap("filter")
	return nil
}

func BindAndValid(c *gin.Context, v interface{}) (bool, ValidErrors) {
	var errs ValidErrors
	err := c.ShouldBindJSON(v)
	global.GVA_LOG.Debug(fmt.Sprintf("ShouldBindJSON err: %v", err))
	if err != nil {
		v := c.Value("trans")
		// global.GVA_LOG.Debug(fmt.Sprintf("trans: %v", v))
		trans, _ := v.(ut.Translator)
		verrs, ok := err.(val.ValidationErrors)
		if !ok {
			return false, errs
		}
		for key, value := range verrs.Translate(trans) {
			errs = append(errs, &ValidError{
				Key:     key,
				Message: value,
			})
		}
		return false, errs
	}

	return true, nil
}
