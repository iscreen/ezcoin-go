package validator

import (
	"reflect"
	"sync"

	"ezcoin.cc/ezcoin-go/server/app/service"
	"ezcoin.cc/ezcoin-go/server/global"
	val "github.com/go-playground/validator/v10"
)

// var MyValidations = map[string]val.Func{
// 	"bookabledate": BookableDate,
// }

type CustomValidator struct {
	Once     sync.Once
	Validate *val.Validate
}

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{}
}

func (v *CustomValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyinit()
		if err := v.Validate.Struct(obj); err != nil {
			return err
		}
	}

	return nil
}

func (v *CustomValidator) Engine() interface{} {
	v.lazyinit()
	return v.Validate
}

func (v *CustomValidator) lazyinit() {
	v.Once.Do(func() {
		v.Validate = val.New()
		v.Validate.SetTagName("binding")
		global.GVA_LOG.Debug("lazyinit")

		v.Validate.RegisterStructValidation(UserStructLevelValidation, service.UpdateUser{})
		v.Validate.RegisterStructValidation(UpdateUserRobotStructLevelValidation, service.UpdateUserRobot{})
		v.Validate.RegisterStructValidation(CreateUserRobotStructLevelValidation, service.CreateUserRobot{})
		// v.Validate.RegisterStructValidation(UserRobotStructLevelValidation, request.CreateUserRobot{})

		// for key, value := range MyValidations {
		// 	global.GVA_LOG.Debug(fmt.Sprintf("----> load customize validator : %s", key))
		// 	// v.Validate.RegisterValidation(key, value)

		// }
	})
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
