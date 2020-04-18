package yangkai

import (
	"errors"
	"io"
	"net/http"
	"reflect"
	"time"
)

var REQUEST_JSON_DEFAULT interface{} //全局的json的默认返回配置

//控制器模版
type ControllerTemplate func(Response, Request) string

//全局map类型配置
type GlobalMap map[string]interface{}

func (this *GlobalMap) Get(key string) string {
	if d, ok := (*this)[key]; ok {
		return d.(string)
	} else {
		return ""
	}
}

func (this *GlobalMap) Set(key string, value string) {
	if len(*this) == 0 {
		*this = make(map[string]interface{})
	}
	(*this)[key] = value
}

type GroupConfig struct {
	Path string
	Middleware []func(response Response,request Request) bool
}


const (
	ROUTER_HTTP_GET        = "GET:"
	ROUTER_HTTP_POST       = "POST:"
	ROUTER_CONFIG_RESPONSE = "response"
	ROUTER_CONFIG_REQUEST  = "request"
)

type Router struct {
	c          map[string]ControllerTemplate
	config     map[string]map[string]interface{}
	groupC     map[string]int
	groupPath  string
	groupSlice []GroupConfig
	httpServer *http.Server
}

func (this *Router) SetResponseConfig(config map[string]interface{}) {
	this.config[ROUTER_CONFIG_RESPONSE] = config
}

func (this *Router) SetRequestConfig(config map[string]interface{}) {
	this.config[ROUTER_CONFIG_REQUEST] = config
}

func (this *Router) GET(path string, function ControllerTemplate) {
	path = this.PathModify(path)
	this.c[ROUTER_HTTP_GET+this.groupPath+path] = function
}

func (this *Router) POST(path string, function ControllerTemplate) {
	path = this.PathModify(path)
	this.c[ROUTER_HTTP_POST+this.groupPath+path] = function
}

func (this *Router) ALL(path string, function ControllerTemplate) {
	path = this.PathModify(path)
	this.c[this.groupPath+path] = function
}

func (this *Router) Group(config GroupConfig, function func()) {
	var c = map[string]string{}
	for k,_ := range this.c {
		c[k] = ""
	}

	config.Path = this.PathModify(config.Path)
	this.groupSlice = append(this.groupSlice, config)
	this.groupPath = config.Path
	function()
	this.groupPath = ""

	for k,_ := range this.c {
		if _,ok := c[k]; !ok {
			this.groupC[k] = len(this.groupSlice)-1
		}
	}
}

func (this *Router) Start(Addr string) {
	this.httpServer = &http.Server{
		Addr:           Addr,
		Handler:        this,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//配置初始化
	this.c = make(map[string]ControllerTemplate)
	this.config = make(map[string]map[string]interface{})
	this.groupC = make(map[string]int)
}

//全局配置&&监听http
func (this *Router) Run() error {

	//对预注入的结构体验证
	var configSetStruct = this.config[ROUTER_CONFIG_RESPONSE][RESPONSE_CONFIG_SET_STRUCT]
	for key, value := range this.config[ROUTER_CONFIG_RESPONSE] {
		if err := new(Response).ConfigVerify(key, value); err != nil {
			return err
		}
	}

	if configSetStruct != nil {
		var typeOf = reflect.TypeOf(configSetStruct)
		var typeOfElem = typeOf.Elem()
		if typeOfElem.PkgPath() == "" {
			return errors.New("You need to define a data type")
		}
		for i := 0; i < typeOfElem.NumField(); i++ {
			if typeOfElem.Field(i).Tag.Get(TAG_JSON) == "" {
				return errors.New(typeOfElem.Field(i).Name + ":not tag json")
			}
		}
	}

	return this.httpServer.ListenAndServe()
}

//路径的开头必须是/ 结尾不能是/
func (this *Router) PathModify(path string) string {
	if path != "" {
		if path[0] != '/' {
			path = "/"+path
		}
		if path[len(path)-1] == '/' {
			path = path[0:len(path)-1]
		}
	}
	return path
}

func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var path = this.PathModify(r.URL.Path)

	//首先验证全局的路由
	if _, ok := this.c[path]; !ok {
		switch r.Method {
		case http.MethodGet:
			path = ROUTER_HTTP_GET + r.URL.Path
		case http.MethodPost:
			path = ROUTER_HTTP_POST + r.URL.Path
		}
	}

	//用这个handler实现路由转发，相应的路由调用相应func
	if method, ok := this.c[path]; ok {

		//用户访问配置
		response := Response{err: nil, httpResponse: w}
		response.ConfigNotVerify(this.config[ROUTER_CONFIG_RESPONSE])

		//用户请求http
		request := Request{request: r}
		request.ConfigNotVerify(this.config[ROUTER_CONFIG_REQUEST])
		request.New()

		//组的中间件验证
		if number, ok := this.groupC[path]; ok {
			for _,groupMiddleware := range this.groupSlice[number].Middleware {
				if groupMiddleware(response,request) == false {
					_ = r.Close
					return
				}
			}
		}

		var returnData = method(response, request)

		_, _ = io.WriteString(w, returnData)
		_ = r.Close

	} else {
		w.WriteHeader(404)
		_, _ = io.WriteString(w, "router not create URL:"+r.URL.String())
		_ = r.Close
	}

}
