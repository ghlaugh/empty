package controllers

import (
	"github.com/unrolled/render"
)

// 定义一个用于json或者xml等各种渲染的公共渲染模块
var Render *render.Render = render.New()
