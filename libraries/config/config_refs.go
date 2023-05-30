package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"reflect"
	"regexp"
	"strings"
)

var refRegExp *regexp.Regexp
var refGroupIndex = -1

func init() {
	refRegExp = regexp.MustCompile(`\${(?P<ref>[\w\d\.]+)}`)

	for i, name := range refRegExp.SubexpNames() {
		if name == "ref" {
			refGroupIndex = i
		}
	}

	if refGroupIndex == -1 {
		panic("Reference group must exist")
	}
}

func resolveReferencesInStringFields(values MapStr) mapstructure.DecodeHookFunc {
	return func(f reflect.Kind, t reflect.Kind, data interface{}) (interface{}, error) {
		if f != reflect.String {
			return data, nil
		}

		return resolveStringReferences(values, data)
	}
}

func resolveStringReferences(values MapStr, data interface{}) (interface{}, error) {
	source, isSourceString := data.(string)

	if !isSourceString {
		return data, nil
	}

	matches := refRegExp.FindAllStringSubmatch(source, -1)

	if len(matches) == 0 {
		return data, nil
	}

	target := source

	for _, match := range matches {
		if refGroupIndex < len(match) {
			fullRef := match[0]
			ref := match[refGroupIndex]

			if len(ref) > 0 {
				path := strings.Replace(ref, ".", "/", -1)
				val := values.FindPath(path)

				if val == nil {
					return nil, errors.New(fmt.Sprintf("Reference '%s' could not be found", ref))
				}

				val = castValueToImpliedType(val)

				val, err := resolveStringReferences(values, val)

				if err != nil {
					return nil, errors.WithMessage(err, "Failed to resolve references")
				}

				if fullRef == target {
					return val, nil
				}

				strVal := fmt.Sprintf("%v", val)
				target = strings.ReplaceAll(target, fullRef, strVal)
			}
		}
	}

	return target, nil
}

func castValueToImpliedType(val interface{}) interface{} {
	strVal, isStr := val.(string)
	if !isStr {
		return val
	}

	switch {
	case strings.EqualFold(strVal, "true"):
		return true
	case strings.EqualFold(strVal, "false"):
		return false
	}

	return val
}

func resolveReferencesInUntypedMapFields(values MapStr) mapstructure.DecodeHookFunc {
	return func(f reflect.Kind, t reflect.Kind, data interface{}) (interface{}, error) {
		if f != reflect.Map || t != reflect.Interface {
			return data, nil
		}

		return resolveInterfaceReferences(values, data)
	}
}

func resolveArrayReferences(values MapStr, in []interface{}) ([]interface{}, error) {
	result := make([]interface{}, len(in))
	for i, v := range in {
		res, err := resolveInterfaceReferences(values, v)
		if err != nil {
			return nil, err
		}
		result[i] = res
	}
	return result, nil
}

func resolveMapReferences(values MapStr, in MapStr) (MapStr, error) {
	result := make(MapStr)
	for k, v := range in {
		res, err := resolveInterfaceReferences(values, v)
		if err != nil {
			return nil, err
		}
		result[k] = res
	}
	return result, nil
}

func resolveInterfaceReferences(values MapStr, v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case []interface{}:
		return resolveArrayReferences(values, v)
	case MapStr:
		return resolveMapReferences(values, v)
	case string:
		return resolveStringReferences(values, v)
	default:
		return v, nil
	}
}
