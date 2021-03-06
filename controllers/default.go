package controllers

import (
	"github.com/astaxie/beego"
	"jethome/models"
)

type MainController struct {
	beego.Controller
	Info string
}

var main = `
# 周报管理平台

--------------------------
#### REST接口
>项目查询 Get请求 /api/p/\<project\> 获取任务详情
`

var about = `
### 出现背景
主要是因为写周报是merge太麻烦，wiki又不给力，所以才出现的这个周报整合平台。

该平台采用go语言的beego框架编写。
虽说go语言的特性很好，不过开发还是挺艰辛的。

贴几张设计图，留个纪念。

![page1](/static/img/draft/draft1.jpg)
![page2](/static/img/draft/draft2.jpg)

### 更新历史
**2013-6-17**
号晚上开始构思。

**2013-6-18**
出现了第一个demo版。

### 致谢
感谢astaxie提供了这么好用的beego框架

感谢bootstrap对前端页面的大力支持。

author: ssx205@gmail.com
`

// return data for html use
func NameList(cur string) []map[string]interface{} {
	ps := models.Names() // projects
	beego.Info("count projs:", len(ps))
	pm := make([]map[string]interface{}, 0, len(ps))
	for _, projname := range ps {
		m := make(map[string]interface{}, 2)
		m["Name"] = projname
		if cur == projname {
			m["Cru"] = true
		} else {
			m["Cru"] = false
		}
		pm = append(pm, m)
	}
	return pm
}

func (this *MainController) Get() {
	name, pname := this.Ctx.Params[":name"], this.Ctx.Params[":pname"]
	cruproj, _ := models.GetJob(pname, 0, -1)
	this.Data["CruProj"] = cruproj

	pm := NameList(pname)
	this.Data["ProjList"] = pm

	content := main
	if name != "" {
		beego.Debug("name:", name)
		switch name {
		case "about":
			content = about
		}
	}

	beego.Debug("proj name:", pname)
	if pname != "" {
		job, _ := models.GetJob(pname, 0, -1)
		content = job.Description
		this.Data["Project"] = job.Name
		this.Data["QA"] = job.QA
		this.Data["RD"] = job.RD
		this.Data["Type"] = job.Type
	}

	this.Data["Content"] = content
	this.TplNames = "index.tpl"
}

func (this *MainController) Post() {

}
