package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/bobfive1/user-management-api/internal/chrono"

	errInt "github.com/bobfive1/user-management-api/internal/error"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	Validate   *validator.Validate
	uni        *ut.UniversalTranslator
	trans      ut.Translator
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

const (
	BirthdateFormat = "2006-01-02"
)

func Init() {

	en := en.New()
	uni = ut.New(en, en)

	trans, _ = uni.GetTranslator("en")

	Validate = validator.New()
	Validate = binding.Validator.Engine().(*validator.Validate)

	Validate.RegisterTagNameFunc(registerTagJson)

	Validate.RegisterValidation("checkyear", ValidateBirthdate)
	Validate.RegisterValidation("email", VailidEmail)

	RegisterTranslation("checkyear", "{0} Must be at least 18 years old")
	RegisterTranslation("email", "{0} Email invalid format")
}

func registerTagJson(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}

func ValidateBirthdate(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(chrono.DateOnly)

	current := time.Now()
	birthdate18 := value.AddDate(18, 0, 0)

	if current.Equal(birthdate18) || current.After(birthdate18) {
		return true
	}

	return false
}

func ShouldBindJSONWithValidate[T any](c *gin.Context, body T) (T, error) {
	if err := c.ShouldBindJSON(&body); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			out := make(map[string]string)

			for _, fe := range ve {
				// fe.Field() คือชื่อฟิลด์, fe.Tag() คือเงื่อนไขที่พัง (เช่น required, gte)
				out[fe.Field()] = fe.Translate(trans) //"Invalid value for " + fe.Tag()
			}
			return body, errInt.NewFieldValidationError(out)
		}
		return body, errInt.NewFieldValidationError(err)
	}
	return body, nil
}

func VailidEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()

	if len(email) < 3 || len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

func RegisterTranslation(tag, message string) {

	registerFn := func(ut ut.Translator) error {
		return ut.Add(tag, message, true) // see universal-translator for details
	}

	translationFn := func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fmt.Sprintf("[%v]", fe.Tag()))

		return t
	}

	Validate.RegisterTranslation(tag, trans, registerFn, translationFn)
}
