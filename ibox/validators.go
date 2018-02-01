package ibox

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/inhies/go-bytesize"
	"math"
	"regexp"
	"strconv"
	"strings"
)

const (
	unit_size     int = 512
	pool_min_size int = 1000000000000
)

func parse_size_to_bits(size string) int {

	parsed_size, err := bytesize.Parse(size)
	if err != nil {
		fmt.Printf("[ERROR] parsing %v", err)
		return 0
	}
	parsed_size_string := parsed_size.Format("%.0f", "bytes", false)
	parsed_size_int, _ := strconv.Atoi(parsed_size_string[:len(parsed_size_string)-1])

	return parsed_size_int * 8
}

func checkDivisibleBy(num int, divisor int) bool {
	if math.Mod(float64(num), float64(divisor)) == 0 {
		return true
	}
	return false
}

func round(num int, unit int) int {
	d := float64(num / unit)
	c := int(math.Ceil(d))
	return c * unit
}

func VerifyCapacity(num int, unit int) error {
	if checkDivisibleBy(num, unit) {
		return nil
	} else {
		return fmt.Errorf("[ERROR] Size: %v is not aligned with integral units of %v, the value can be rounded to: %v", num, unit, round(num, unit))
	}
}

func validateMinPoolSize(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < pool_min_size {
		errors = append(errors, fmt.Errorf(
			"%q minimal pool size is: %v bytes", k, pool_min_size))
	}
	return
}

func validateIqn(v interface{}, k string) (ws []string, errors []error) {
	var iscsiInitiatorIqnRegex = regexp.MustCompile(`iqn\.\d{4}-\d{2}\.([[:alnum:]-.]+)(:[^,;*&$|\s]+)$`)
	if !iscsiInitiatorIqnRegex.MatchString(v.(string)) {
		errors = append(errors, fmt.Errorf(
			"%q IQN format is wrong", k))
	}
	return
}

func myvalidateIqn(v interface{}) error {
	var iscsiInitiatorIqnRegex = regexp.MustCompile(`iqn\.\d{4}-\d{2}\.([[:alnum:]-.]+)(:[^,;*&$|\s]+)$`)
	if !iscsiInitiatorIqnRegex.MatchString(v.(string)) {
		return fmt.Errorf("IQN format is wrong")
	}
	return nil
}

func validateFcWWN(v interface{}, k string) (ws []string, errors []error) {
	var fcWWNRegex = regexp.MustCompile(`([0-9aA-fF]{16}`)
	if !fcWWNRegex.MatchString(v.(string)) {
		errors = append(errors, fmt.Errorf(
			"%q FC WWN format is wrong", k))
	}
	return
}

func validateStringMatchesPattern(pattern string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		compiledRegex, err := regexp.Compile(pattern)
		if err != nil {
			errors = append(errors, fmt.Errorf(
				"%q regex does not compile", pattern))
			return
		}

		value := v.(string)
		if !compiledRegex.MatchString(value) {
			errors = append(errors, fmt.Errorf(
				"%q doesn't match the pattern (%q): %q",
				k, pattern, value))
		}

		return
	}
}

func validateIntegerGeqThan(threshold int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if value < threshold {
			errors = append(errors, fmt.Errorf(
				"%q cannot be lower than %q", k, threshold))
		}
		return
	}
}

func validateIntegerInRange(min, max int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if value < min {
			errors = append(errors, fmt.Errorf(
				"%q cannot be lower than %d: %d", k, min, value))
		}
		if value > max {
			errors = append(errors, fmt.Errorf(
				"%q cannot be higher than %d: %d", k, max, value))
		}
		return
	}
}

func validateStringLenghtInRange(min, max int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		if len(value) < min {
			errors = append(errors, fmt.Errorf(
				"%q cannot be lower than %d characters", k, min))
		}
		if len(value) > max {
			errors = append(errors, fmt.Errorf(
				"%q cannot be higher than %d characters", k, max))
		}
		return
	}
}

func validateStringInList(list []string, caseSensitive bool) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (s []string, errors []error) {
		value, ok := v.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("key %s has a wrong type, it must be string", k))
			return
		}
		for _, str := range list {
			if (caseSensitive && strings.ToLower(value) == strings.ToLower(str)) || value == str {
				return
			}
		}
		errors = append(errors, fmt.Errorf("%s value is invalid: %v, valid choices are: %v", k, value, strings.Join(list, ", ")))
		return
	}
}
