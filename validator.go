package adele

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/asaskevich/govalidator"
	"github.com/fatih/camelcase"
)

type Validation struct {
	Data   url.Values // Is this a struct from inertia js? Can we move it?
	Errors map[string]string
}

func (a *Adele) Validator(data url.Values) *Validation {
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
		fieldName := formatFieldName(key)
		v.Errors[key] = strings.ReplaceAll(message, ":attribute", fieldName)
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

// Format field names to a lower case string with spaces
func formatFieldName(field string) string {
	s := camelcase.Split(field)
	var n string
	if len(s) > 0 {
		for i := range s {
			n = n + " " + (s[i])
		}
	}
	return strings.TrimSpace(strings.ToLower(n))
}

// Check for required fields in JSON
func (v *Validation) RequiredJSON(json interface{}, fields ...string) {
	reflectedType := reflect.TypeOf(json)
	reflectedKind := reflectedType.Kind()
	// If we do not have reflection pointer,
	// just add all required fields to the
	// validator.
	if reflectedKind != reflect.Ptr {
		for _, field := range fields {
			v.AddError(field, "The :attribute field is required")
		}
	}

	// Using reflection, search for the required
	// fields, passing errors to the validator
	// as we work through each required field.
	vp := reflect.ValueOf(json)
	vs := reflect.Indirect(vp)
	for _, field := range fields {
		var ok bool
		for i := 0; i < vs.NumField(); i++ {
			name := vs.Type().Field(i).Name
			if strings.EqualFold(field, name) {
				ok = true
			}
		}
		if !ok {
			v.AddError(field, "The :attribute field is required")
		}
	}
}

// The fields under validation must be a valid JSON key and contain a value
func (v *Validation) HasJSON(json interface{}, fields ...string) {
	reflectedType := reflect.TypeOf(json)
	reflectedKind := reflectedType.Kind()
	// If we do not have reflection pointer,
	// just add all required fields to the
	// validator.
	if reflectedKind != reflect.Ptr {
		for _, field := range fields {
			v.AddError(field, "The :attribute field is required")
		}
	}

	// Using reflection, search for the required
	// fields, passing errors to the validator
	// as we work through each required field.
	vp := reflect.ValueOf(json)
	vs := reflect.Indirect(vp)
	for _, field := range fields {
		var ok bool
		for i := 0; i < vs.NumField(); i++ {
			name := vs.Type().Field(i).Name

			if strings.EqualFold(field, name) {
				value := vs.Field(i).Interface()
				if value != "" {
					ok = true
				}
			}
		}
		if !ok {
			v.AddError(field, "The :attribute field is required")
		}
	}
}

// Givin a condition, set a error message in the validation map
func (v *Validation) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Is valid email address
func (v *Validation) IsEmail(field, value string) {
	if !govalidator.IsEmail(value) {
		v.AddError(field, "Invalid email address")
	}
}

// The string under validation is a email address and in the public domain.
func (v *Validation) IsEmailInPublicDomain(field, value string) {
	if !govalidator.IsExistingEmail(value) {
		v.AddError(field, "Invalid email address")
	}
}

// The field under validation contains a valid integer value
func (v *Validation) IsInt(field, value string) {
	_, err := strconv.Atoi(value)
	if err != nil {
		v.AddError(field, "This field must be a integer")
	}
}

// The field under validation contains a valid float value.
func (v *Validation) IsFloat(field, value string) {
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		v.AddError(field, "This field must be a floating point number")
	}
}

// The field under validation contains a valid Date value.
func (v *Validation) isDateISO(field, value string) {
	_, err := time.Parse("2000-01-02", value)
	if err != nil {
		v.AddError(field, "This field must be a date in the form of YYYY-MM-DD")
	}
}

// The field under validation does not contain any spaces
func (v *Validation) NoSpaces(field, value string) {
	if govalidator.HasWhitespace(value) {
		v.AddError(field, "Spaces are not allowed")
	}
}

// The field under validation must not be empty.
func (v *Validation) NotEmpty(field, value string) {
	if strings.TrimSpace(value) == "" {
		v.AddError(field, "This field must contain a value")
	}
}

// Password meets default password rules
func (v *Validation) Password(field string, value string, length ...int) {
	minLength := 12
	if len(length) > 0 {
		minLength = length[0]
	}
	if len(value) < minLength {
		message := fmt.Sprintf("The field does not meet the minimum length of %d characters", minLength)
		v.AddError(field, message)
	}

	// mixed case?
	hasUpper := false
	hasLower := false
	for _, char := range value {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsLower(char) {
			hasLower = true
		}
	}
	if hasUpper == false {
		v.AddError(field, "The field must contain a uppercase character")
	}
	if hasLower == false {
		v.AddError(field, "The field must contain a lowercase character")
	}
}

// Password does not appears in data breach using a search by range
func (v *Validation) PasswordUncompromised(field string, value string, threshold ...int) {
	thresholdVerifier := 1
	if len(threshold) > 0 {
		thresholdVerifier = threshold[0]
	}

	hasher := sha1.New()
	hasher.Write([]byte(value))
	hash := strings.ToUpper(hex.EncodeToString([]byte(hasher.Sum(nil))))
	hashPrefix := hash[0:5]
	hasSuffix := hash[5:]

	uri := fmt.Sprintf("https://%s/%s", "api.pwnedpasswords.com/range", hashPrefix)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Printf(fmt.Sprintf("Error creating new request: %s\n", err))
		return
	}
	req.Header.Set("Add-Padding", "true")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf(fmt.Sprintf("Error making request: %s\n", err))
		return
	}

	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf(fmt.Sprintf("Error reading response from API: %s\n", err))
			return
		}

		s := string(body)
		hashSuffixes := strings.Split(s, "\n")

		isPwned := false
		isThresholdExceeded := false

		for _, suffix := range hashSuffixes {

			pwnedHashCount := strings.Split(suffix, ":")
			pwnedCount, _ := strconv.Atoi(strings.TrimSpace(pwnedHashCount[1]))

			if pwnedCount > 0 {
				pwnedHash := pwnedHashCount[0]

				if pwnedHash == hasSuffix {
					isPwned = true
				}

				if pwnedCount >= thresholdVerifier {

					isThresholdExceeded = true
				}
			}
		}

		if isPwned && isThresholdExceeded {
			v.AddError(field, "The password provided was discovered in a recent data leak; please select another password.")
		}
	} else {
		log.Printf(fmt.Sprintf("API returned a %s status code. Unable to verify password validity.\n", resp.StatusCode))
	}
}
