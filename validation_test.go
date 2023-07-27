package belajargovalidation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
)

// validasi struct
func TestValidationStruct(t *testing.T) {
	var validate = validator.New()
	if validate == nil {
		t.Error("validate is nil")
	}
}

func TestValidationVariable(t *testing.T) {
	var validate = validator.New()
	user := ""

	err := validate.Var(user, "required")

	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestValidationDuaVariable(t *testing.T) {
	var validate = validator.New()
	variable1 := ""
	variable2 := "beda"

	err := validate.VarWithValue(variable1, variable2, "eqfield")

	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestBakedInVariable(t *testing.T) {
	// check Baked-in Validation
	// https://pkg.go.dev/github.com/go-playground/validator/v10#section-readme
}

func TestMultipleTagValidation(t *testing.T) {
	validate := validator.New()
	userName := "dicki123"

	err := validate.Var(userName, "required,numeric")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestTagParameterValidation(t *testing.T) {
	validate := validator.New()
	userName := "99"

	err := validate.Var(userName, "required,numeric,min=5,max=100")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestStructValidation(t *testing.T) {
	type LoginRequest struct {
		Username string `validate:"required,email"`
		Password string `validate:"required,min=5"`
	}

	validate := validator.New()
	login := LoginRequest{
		Username: "dicki",
		Password: "dicki",
	}

	err := validate.Struct(login)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestValidationErrors(t *testing.T) {
	type LoginRequest struct {
		Username string `validate:"required,email"`
		Password string `validate:"required,min=5"`
	}

	validate := validator.New()
	login := LoginRequest{
		Username: "dicki",
		Password: "ki",
	}

	err := validate.Struct(login)
	if err != nil {
		// konversi error ke validation error bawaan package validator
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			fmt.Println("error", fieldError.Field(), "on tag", fieldError.Tag(), "with error", fieldError.Error())
		}
	}
}

func TestValidationCrossField(t *testing.T) {
	type RegisterUser struct {
		Username        string `validate:"required,email"`
		Password        string `validate:"required,min=5"`
		ConfirmPassword string `validate:"required,min=5,eqfield=Password"`
	}

	validate := validator.New()
	register := RegisterUser{
		Username:        "dicki@mail.com",
		Password:        "dicki",
		ConfirmPassword: "ki",
	}

	err := validate.Struct(register)
	if err != nil {
		// konversi error ke validation error bawaan package validator
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			fmt.Println("error", fieldError.Field(), "on tag", fieldError.Tag(), "with error", fieldError.Error())
		}
	}
}

func TestValidationNestedStruct(t *testing.T) {
	type Address struct {
		City    string `validate:"required"`
		Country string `validate:"required"`
	}
	type User struct {
		Id      string  `validate:"required"`
		Name    string  `validate:"required"`
		Address Address `validate:"required"`
	}

	validate := validator.New()
	request := User{
		Id:   "dicki@mail.com",
		Name: "dicki",
		Address: Address{
			City:    "pontianak",
			Country: "",
		},
	}

	err := validate.Struct(request)
	if err != nil {
		// konversi error ke validation error bawaan package validator
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			fmt.Println("error", fieldError.Field(), "on tag", fieldError.Tag(), "with error", fieldError.Error())
		}
	}
}

func TestValidationCollection(t *testing.T) {
	type Address struct {
		City    string `validate:"required"`
		Country string `validate:"required"`
	}
	type User struct {
		Id   string `validate:"required"`
		Name string `validate:"required"`
		// gunakan tag dive jika ingin isi slice di validasi
		Address []Address `validate:"required,dive"`
	}

	validate := validator.New()
	request := User{
		Id:   "dicki@mail.com",
		Name: "dicki",
		Address: []Address{
			{
				City:    "",
				Country: "",
			},
			{
				City:    "",
				Country: "",
			},
		},
	}

	err := validate.Struct(request)
	if err != nil {
		// konversi error ke validation error bawaan package validator
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			fmt.Println("error", fieldError.Field(), "on tag", fieldError.Tag(), "with error", fieldError.Error())
		}
	}
}

func TestValidationBasicCollection(t *testing.T) {
	type Address struct {
		City    string `validate:"required"`
		Country string `validate:"required"`
	}
	type User struct {
		Id   string `validate:"required"`
		Name string `validate:"required"`
		// gunakan tag dive jika ingin isi slice di validasi
		Address []Address `validate:"required,dive"`
		Hobbies []string  `validate:"dive,required,min=3"`
	}

	validate := validator.New()
	request := User{
		Id:   "dicki@mail.com",
		Name: "dicki",
		Address: []Address{
			{
				City:    "x",
				Country: "x",
			},
			{
				City:    "x",
				Country: "x",
			},
		},
		Hobbies: []string{"Berenang", "Gaming", "X", ""},
	}

	err := validate.Struct(request)
	if err != nil {
		// konversi error ke validation error bawaan package validator
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			fmt.Println("error", fieldError.Field(), "on tag", fieldError.Tag(), "with error", fieldError.Error())
		}
	}
}

func TestValidationMap(t *testing.T) {
	type Address struct {
		City    string `validate:"required"`
		Country string `validate:"required"`
	}
	type School struct {
		Name string `validate:"required"`
	}
	type User struct {
		Id   string `validate:"required"`
		Name string `validate:"required"`
		// gunakan tag dive jika ingin isi slice di validasi
		Address []Address `validate:"required,dive"`
		Hobbies []string  `validate:"dive,required,min=3"`
		// untuk map perlu tag keys & endkeys untuk validasi key value di map
		// karna mapnya berupa struct maka perlu tambahkan dive lagi di endkeys
		Schools map[string]School `validate:"dive,keys,required,endkeys,dive"`
	}

	validate := validator.New()
	request := User{
		Id:   "dicki@mail.com",
		Name: "dicki",
		Address: []Address{
			{
				City:    "x",
				Country: "x",
			},
			{
				City:    "x",
				Country: "x",
			},
		},
		Hobbies: []string{"Berenang", "Gaming", "XYZ"},
		Schools: map[string]School{
			"SD": {
				Name: "SD X",
			},
			"SMP": {
				Name: "SMP X",
			},
			"": {
				Name: "",
			},
		},
	}

	err := validate.Struct(request)
	if err != nil {
		// konversi error ke validation error bawaan package validator
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			fmt.Println("error", fieldError.Field(), "on tag", fieldError.Tag(), "with error", fieldError.Error())
		}
	}
}

func TestValidationBasicMap(t *testing.T) {
	type Address struct {
		City    string `validate:"required"`
		Country string `validate:"required"`
	}
	type School struct {
		Name string `validate:"required"`
	}
	type User struct {
		Id   string `validate:"required"`
		Name string `validate:"required"`
		// gunakan tag dive jika ingin isi slice di validasi
		Address []Address         `validate:"required,dive"`
		Hobbies []string          `validate:"dive,required,min=3"`
		Schools map[string]School `validate:"dive,keys,required,endkeys,dive"`
		Wallets map[string]int    `validate:"dive,keys,required,endkeys,required,gt=500"`
	}

	validate := validator.New()
	request := User{
		Id:   "dicki@mail.com",
		Name: "dicki",
		Address: []Address{
			{
				City:    "x",
				Country: "x",
			},
			{
				City:    "x",
				Country: "x",
			},
		},
		Hobbies: []string{"Berenang", "Gaming", "XYZ"},
		Schools: map[string]School{
			"SD": {
				Name: "SD X",
			},
			"SMP": {
				Name: "SMP X",
			},
		},
		Wallets: map[string]int{
			"BCA":     1000,
			"Mandiri": 0,
			"":        200,
		},
	}

	err := validate.Struct(request)
	if err != nil {
		// konversi error ke validation error bawaan package validator
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			fmt.Println("error", fieldError.Field(), "on tag", fieldError.Tag(), "with error", fieldError.Error())
		}
	}
}

func TestAliasTag(t *testing.T) {
	validate := validator.New()
	validate.RegisterAlias("varchar", "required,max=255")

	type Seller struct {
		Id   string `validate:"varchar,min=3"`
		Name string `validate:"varchar"`
	}

	seller := Seller{
		Id:   "12",
		Name: "x",
	}

	err := validate.Struct(seller)
	if err != nil {
		// konversi error ke validation error bawaan package validator
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			fmt.Println("error", fieldError.Field(), "on tag", fieldError.Tag(), "with error", fieldError.Error())
		}
	}
}

// custom validation func
func MustValidUsername(field validator.FieldLevel) bool {
	value, ok := field.Field().Interface().(string)
	if ok {
		if value != strings.ToUpper(value) {
			return false
		}

		if len(value) < 5 {
			return false
		}
	}

	return true
}

func TestCustomValidationFunc(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("username", MustValidUsername)

	type LoginRequest struct {
		Username string `validate:"required,username"`
		Password string `validate:"required"`
	}

	request := LoginRequest{
		Username: "x",
		Password: "x",
	}

	err := validate.Struct(request)
	if err != nil {
		// konversi error ke validation error bawaan package validator
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			fmt.Println("error", fieldError.Field(), "on tag", fieldError.Tag(), "with error", fieldError.Error())
		}
	}
}

// custom parameter validation
var regexNumber = regexp.MustCompile("^[0-9]+$")

func MustValidPin(field validator.FieldLevel) bool {
	length, err := strconv.Atoi(field.Param())
	if err != nil {
		panic(err)
	}

	value := field.Field().String()
	if !regexNumber.MatchString(value) {
		return false
	}

	return len(value) == length
}

func TestCustomParameterValidation(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("pin", MustValidPin)

	type Login struct {
		Phone string `validate:"required,number"`
		Pin   string `validate:"required,pin=4"`
	}

	request := Login{
		Phone: "0800",
		Pin:   "123",
	}

	err := validate.Struct(request)
	if err != nil {
		// konversi error ke validation error bawaan package validator
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			fmt.Println("error", fieldError.Field(), "on tag", fieldError.Tag(), "with error", fieldError.Error())
		}
	}
}

func TestOrRule(t *testing.T) {
	validate := validator.New()

	type Login struct {
		// gunakan tanda | (pipe) untuk or jika koma artinya and
		Username string `validate:"required,email|numeric"`
		Password string `validate:"required"`
	}

	request := Login{
		Username: "dicki@mail.com",
		Password: "123",
	}

	err := validate.Struct(request)
	if err != nil {
		// konversi error ke validation error bawaan package validator
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			fmt.Println("error", fieldError.Field(), "on tag", fieldError.Tag(), "with error", fieldError.Error())
		}
	}
}

// cross field validation
func MustEqualsIgnoreCase(field validator.FieldLevel) bool {
	// reflect.Value, reflect.Kind, bool, bool
	value, _, _, ok := field.GetStructFieldOK2()
	if !ok {
		panic("field not ok")
	}

	firstValue := strings.ToUpper(field.Field().String())
	secondValue := strings.ToUpper(value.String())

	return firstValue == secondValue
}

func TestCrossFieldValidation(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("field_equals_ignore_case", MustEqualsIgnoreCase)

	type User struct {
		Username string `validate:"required,field_equals_ignore_case=Email|field_equals_ignore_case=Phone"`
		Email    string `validate:"required,email"`
		Phone    string `validate:"required,numeric,max=12"`
		Name     string `validate:"required"`
	}

	user := User{
		Username: "beda",
		Email:    "dicki@mail.com",
		Phone:    "0800",
		Name:     "Dicki",
	}

	err := validate.Struct(user)
	if err != nil {
		// konversi error ke validation error bawaan package validator
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			fmt.Println("error", fieldError.Field(), "on tag", fieldError.Tag(), "with error", fieldError.Error())
		}
	}
}

// struct level validation
type RegisterRequest struct {
	Username string `validate:"required"`
	Email    string `validate:"required,email"`
	Phone    string `validate:"required,numeric"`
	Password string `validate:"required"`
}

func MustValidRegisterSuccess(level validator.StructLevel) {
	registerRequest := level.Current().Interface().(RegisterRequest)

	if registerRequest.Username == registerRequest.Email || registerRequest.Username == registerRequest.Phone {
		// sukses
	} else {
		// gagal
		level.ReportError(registerRequest.Username, "Username", "Username", "username", "")
	}
}

// digunakan jika validasi kompleks
func TestStructLevelValidation(t *testing.T) {
	validate := validator.New()
	validate.RegisterStructValidation(MustValidRegisterSuccess, RegisterRequest{})

	request := RegisterRequest{
		Username: "dicki",
		Email:    "dicki@mail.com",
		Phone:    "0800",
		Password: "12345",
	}

	err := validate.Struct(request)
	if err != nil {
		fmt.Println(err.Error())
	}
}
