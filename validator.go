package seng

import (
	"errors"
	"reflect"
	"strings"
)

type Validator struct {
}

const (
	ValidateParams   = "validate"
	SplitMultiParams = "|"
	SplitTagMessage  = "#"
)

func (c *Context) Validate(data interface{}) error {
	validator := c.engine.validatorPool.Get().(*Validator)
	defer c.engine.validatorPool.Put(validator)

	return validator.Validate(data)
}

func (s *Validator) Validate(data interface{}) error {
	d := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	fieldNum := d.NumField()
	for i := 0; i < fieldNum; i++ {
		field := d.Field(i)
		tags := field.Tag
		validateParam, ok := tags.Lookup(ValidateParams)
		if !ok {
			continue
		}
		success, msg := validate(v.Field(i).String(), validateParam)
		if !success {
			return errors.New(msg)
		}
	}
	return nil
}

func splitParams(tags string) []string {
	return strings.Split(tags, SplitMultiParams)
}

func validate(data, tags string) (bool, string) {
	params := splitParams(tags)
	for _, param := range params {
		tagAndMsg := strings.Split(param, SplitTagMessage)
		// required
		if tagAndMsg[0] == "required" {
			if data == "" {
				return false, tagAndMsg[1]
			}
		}
	}
	return true, ""
}
