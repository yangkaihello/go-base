package yangkai

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
)

var TAG_JSON = "json"

const (
	RESPONSE_CONFIG_SET_ACCESS_ALL = "set-access-all"
	RESPONSE_CONFIG_SET_STRUCT     = "set-struct"
)

var configResponseTemp = map[string]interface{}{
	RESPONSE_CONFIG_SET_ACCESS_ALL: reflect.Bool,
	RESPONSE_CONFIG_SET_STRUCT:     reflect.Ptr,
}

type Response struct {
	err          error
	body         string
	jsonStruct   interface{}
	httpResponse http.ResponseWriter
}

func (this *Response) configLoad(key, a interface{}) {
	switch key {
	case RESPONSE_CONFIG_SET_ACCESS_ALL:
		this.SetAccessAll(a.(bool))
	case RESPONSE_CONFIG_SET_STRUCT:
		this.SetStruct(a)
	}
}

func (this *Response) ConfigVerify(key string, value interface{}) error {
	var err error = nil
	if configResponseTemp[key] != reflect.TypeOf(value).Kind() {
		err = errors.New(key + ": typeof error")
	}
	return err
}

func (this *Response) ConfigNotVerify(config map[string]interface{}) {
	for key, value := range config {
		this.configLoad(key, value)
	}
}

func (this *Response) Config(config map[string]interface{}) {
	for key, value := range config {
		if err := this.ConfigVerify(key, value); err == nil {
			this.configLoad(key, value)
		}
	}
}

func (this *Response) SetAccessAll(b bool) *Response {
	if b == true {
		this.httpResponse.Header().Set("Access-Control-Allow-Origin", "*")
		this.httpResponse.Header().Set("Access-Control-Allow-Methods", "*")
		this.httpResponse.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	} else {
		this.httpResponse.Header().Del("Access-Control-Allow-Origin")
		this.httpResponse.Header().Del("Access-Control-Allow-Methods")
		this.httpResponse.Header().Del("Access-Control-Allow-Headers")
	}

	return this
}

func (this *Response) SetStatus(code int) *Response {
	this.httpResponse.WriteHeader(code)
	return this
}

func (this *Response) SetStructValue(key string, value interface{}) *Response {
	var typeOfElem = reflect.TypeOf(this.jsonStruct).Elem()
	var valueOfElem = reflect.ValueOf(this.jsonStruct).Elem()

	if typeOfElem.Kind() == reflect.Struct && typeOfElem.NumField() != 0 {
		field, _ := typeOfElem.FieldByName(key)
		if field.Name != "" {
			valueOfElem.Field(field.Index[0]).Set(reflect.ValueOf(value))
		}
	}
	this.Json(this.jsonStruct)
	return this
}

func (this *Response) SetStruct(ptr interface{}) *Response {
	var typeOf = reflect.TypeOf(ptr)
	var valueOf = reflect.ValueOf(ptr)

	if typeOf.Kind() == reflect.Ptr && typeOf.Elem().PkgPath() != "" {
		var typeOfElem = typeOf.Elem()
		var valueOfElem = valueOf.Elem()
		newPtr := reflect.New(reflect.TypeOf(ptr).Elem()).Interface()
		this.jsonStruct = newPtr

		for i := 0; i < typeOfElem.NumField(); i++ {
			this.SetStructValue(typeOfElem.Field(i).Name, valueOfElem.Field(i).Interface())
		}
	} else {
		this.err = errors.New("struct not ptr AND not type")
	}
	return this
}

func (this *Response) Json(a interface{}) *Response {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)

	this.err = jsonEncoder.Encode(a)
	this.body = bf.String()
	return this
}

func (this *Response) Data(data string) *Response {
	this.body = data
	return this
}

func (this *Response) Send() string {
	return this.body
}

func (this *Response) GetErr() error {
	return this.err
}

func (this *Response) GetHttpResponse() http.ResponseWriter {
	return this.httpResponse
}
