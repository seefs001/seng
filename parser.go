package seng

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

// BodyParser parser struct application/json
func (c *Context) BodyParser(out interface{}) (err error) {
	contentType := strings.ToLower(c.GetContentType())
	// json parser
	if strings.HasPrefix(contentType, MINEApplicationJSON) {
		// read data from request
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, out)
		if err != nil {
			return err
		}
	}
	return
}
