package adel

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
)

type Validation struct {
	Data   url.Values
	Errors map[string]string
}

func (a *Adel) Validator(data url.Values) *Validation {
	return &Validation{
		Errors: make(map[string]string),
		Data:   data,
	}
}

// Test if any errors exist
func (v *Validation) Valid() bool {
	return len(v.Errors) == 0
}

// Add and error to the validation map
func (v *Validation) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Value is in form post
func (v *Validation) Has(field string, r *http.Request) bool {
	isInRequest := r.Form.Get(field)
	if isInRequest == "" {
		return false
	}
	return true
}

// Check for required fields
func (v *Validation) Required(r *http.Request, fields ...string) {
	for _, field := range fields {
		value := r.Form.Get(field)
		if strings.TrimSpace(value) == "" {
			v.AddError(field, "This field is required")
		}
	}
}

// Givin a condition, set a error message in the validation map
func (v *Validation) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Has valid email address shape
func (v *Validation) IsEmail(field, value string) {
	if !govalidator.IsEmail(value) {
		v.AddError(field, "invalid email address")
	}
}

// Has a valid integer
func (v *Validation) IsInt(field, value string) {
	_, err := strconv.Atoi(value)
	if err != nil {
		v.AddError(field, "this field must be a integer")
	}
}

// Has a valid float
func (v *Validation) IsFloat(field, value string) {
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		v.AddError(field, "this field must be a floating point number")
	}
}

// Has a valid float
func (v *Validation) isDateISO(field, value string) {
	_, err := time.Parse("2000-01-02", value)
	if err != nil {
		v.AddError(field, "this field must be a date in the form of YYYY-MM-DD")
	}
}

// Has a no spaces
func (v *Validation) NoSpaces(field, value string) {
	if govalidator.HasWhitespace(value) {
		v.AddError(field, "spaces are not allowed")
	}
}
