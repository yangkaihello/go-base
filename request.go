package yangkai

import (
	"errors"
	"net/http"
	"reflect"
)

const (
	REQUEST_CONFIG_PARSE_MULTIPART_FORM = "set-parse-multipart-form"
)

var configRequestTemp = map[string]interface{}{
	REQUEST_CONFIG_PARSE_MULTIPART_FORM: reflect.Int,
}

type Request struct {
	request     *http.Request
	formData    GlobalMap
	formAll     GlobalMap
	formPostAll GlobalMap
}

func (this *Request) New() *Request {
	this.request.ParseForm()

	for key, value := range this.request.Form {
		this.formAll.Set(key, value[0])
	}

	for key, value := range this.request.PostForm {
		this.formPostAll.Set(key, value[0])
	}

	if this.request.MultipartForm != nil {
		for key, value := range this.request.MultipartForm.Value {
			this.formData.Set(key, value[0])
		}
	}

	return this
}

func (this *Request) configLoad(key, a interface{}) {
	switch key {
	case REQUEST_CONFIG_PARSE_MULTIPART_FORM:
		this.SetParseMultipartForm(a.(int))
	}
}

func (this *Request) ConfigVerify(key string, value interface{}) error {
	var err error = nil
	if configRequestTemp[key] != reflect.TypeOf(value).Kind() {
		err = errors.New(key + ": typeof error")
	}
	return err
}

func (this *Request) ConfigNotVerify(config map[string]interface{}) {
	for key, value := range config {
		this.configLoad(key, value)
	}
}

func (this *Request) Config(config map[string]interface{}) {
	for key, value := range config {
		if err := this.ConfigVerify(key, value); err == nil {
			this.configLoad(key, value)
		}
	}
}

func (this *Request) SetParseMultipartForm(i int) *Request {
	var maxMemory int64
	if i == 0 {
		maxMemory = 1024 * 1024 * 10
	} else {
		maxMemory = 1024 * 1024 * int64(i)
	}
	this.request.ParseMultipartForm(maxMemory)
	return this
}

func (this *Request) SetMultiForm(key string, value string) {
	this.formData.Set(key, value)
}

func (this *Request) SetAllForm(key string, value string) {
	this.formAll.Set(key, value)
}

func (this *Request) SetPostForm(key string, value string) {
	this.formPostAll.Set(key, value)
}

func (this *Request) GetKeyMultiForm(key string, def string) string {
	if this.formData.Get(key) != "" {
		return this.formData.Get(key)
	} else {
		return def
	}
}

func (this *Request) GetKeyAllForm(key string, def string) string {
	if this.formAll.Get(key) != "" {
		return this.formAll.Get(key)
	} else {
		return def
	}
}

func (this *Request) GetKeyPostForm(key string, def string) string {
	if this.formPostAll.Get(key) != "" {
		return this.formPostAll.Get(key)
	} else {
		return def
	}
}

func (this *Request) GetMultiForm() GlobalMap {
	return this.formData
}

func (this *Request) GetAllForm() GlobalMap {
	return this.formAll
}

func (this *Request) GetPostForm() GlobalMap {
	return this.formPostAll
}

/*func (this *Request) ParseFormData()  {
	bodyBytes,_ := ioutil.ReadAll(this.request.Body)
	this.request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	var multi *multipart.Reader
	var p *multipart.Part
	var err error
	var data []byte

	if multi,err = this.request.MultipartReader(); err == nil {
		for  {
			if p,err = multi.NextPart(); err != nil {
				break
			}

			if data,err = ioutil.ReadAll(p); err != nil {
				break
			}
			this.formData.Set(p.FormName(),string(data))
		}
	}
}*/

func (this *Request) GetHttpRequest() *http.Request {
	return this.request
}
