package routers

import (
	"github.com/gorilla/mux"
	"strings"
	"fmt"
	"time"
	"sort"
	"net/http"
	"github.com/urfave/negroni"
	"github.com/15125505/zlog/log"
)

// 子路由必须实现的接口
type SubController interface {
	Handle(m *mux.Router, tpl string)
}

// 用于存储用户参数的结构
type handle struct {
	sub SubController
	tpl string
}

// 用户路由信息表
var handles []handle

// 添加控制器
func AddController(sub SubController, tpl string) {
	handles = append(handles, handle{sub, tpl})
}

// 设置路由
func CreateHandle(m *mux.Router) {
	for _, v := range handles {
		v.sub.Handle(m, v.tpl)
	}
}

// 预处理（解析参数和日期打印）
func PreProcess(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	// 在页面被处理之前，你可以做一些工作
	start := time.Now()

	// 这句代码放在next之前，可以避免每次获取参数之前都进行ParseForm操作
	r.ParseForm()

	// 继续后续的处理
	next(rw, r)

	// 为了避免阿里云的SLB日志过多，不打印HEAD请求，
	if r.Method == "HEAD" {
		return
	}

	// 获取http状态码
	res := rw.(negroni.ResponseWriter)
	code := res.Status()
	var color string
	switch {
	case code >= 200 && code < 300:
		color = "\033[01;42;34m" // 绿色
	case code >= 300 && code < 400:
		color = "\033[01;47;34m" // 白色
	case code >= 400 && code < 500:
		color = "\033[01;43;34m" // 黄色
	default:
		color = "\033[01;41;33m" // 红色
	}

	// 获取参数信息
	var keys []string
	for k := range r.Form {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var param string
	for _, k := range keys {
		param += " " + k + ":"
		param += fmt.Sprint(r.Form[k])
	}
	param = strings.TrimLeft(param, " ")

	// 获取请求者IP地址
	ip := "127.0.0.1"
	if forwards := r.Header.Get("X-Forwarded-For"); forwards != "" {
		ips := strings.Split(forwards, ",")
		if len(ips) > 0 {
			ip = ips[0]
		}
	} else {
		ip = r.RemoteAddr
	}
	ips := strings.Split(ip, ":")
	if len(ips) > 0 {
		ip = ips[0]
	}

	// 显示请求详情
	tmpUA := []byte(r.UserAgent())
	if len(tmpUA) > 40 {
		tmpUA = tmpUA[:40]
	}

	log.Info(fmt.Sprintf("%v %v \033[0m\033[37m|%12v\033[32m|%15s\033[01;37m|%5v\033[0m\033[33m|%40v\033[32m|%v\033[33m|%v",
		color,
		res.Status(),
		time.Since(start),
		ip,
		r.Method,
		string(tmpUA),
		r.URL.Path,
		param,
	))
}
