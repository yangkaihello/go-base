# go-base
* go 的一个 http 快速创建的框架，内部支持，本仓库中提供go的一些组件支持
* https://github.com/yangkaihello/go-sql-orm #sql 的连贯操作版本的更新请查看库中的描述
* https://github.com/yangkaihello/go-validate #验证库，可以自定义验证规则


### 定义开始

```
//实例路由
Router := new(yangkai.Router)
Router.Start(":1230")

//全局配置response config
Router.SetResponseConfig(map[string]interface{}{
    yangkai.RESPONSE_CONFIG_SET_STRUCT: &struct {
        Status bool `json:"status"`
        Data interface{} `json:"data"`
    }{true,""},
    yangkai.RESPONSE_CONFIG_SET_ACCESS_ALL: true,
})
//全局配置request
Router.SetRequestConfig(map[string]interface{}{
    yangkai.REQUEST_CONFIG_PARSE_MULTIPART_FORM: 10,
})

//定义路由的url
Router.ALL("/",new(controller.Index).Index)
Router.ALL("/start",new(controller.Index).Start)
Router.ALL("/server/add",new(controller.Server).Add)
Router.ALL("/server/index",new(controller.Server).Index)

//http项目开始运行
Router.Run()
```
