/**
  @Author : hanxiaodong
*/

package web

import (
	"net/http"
	"fmt"
	"education/web/controller"
)


// 启动Web服务并指定路由信息
func WebStart(app controller.Application)  {

	fs:= http.FileServer(http.Dir("web/static"))
	fs2:= http.FileServer(http.Dir("web/assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs2))
	// 指定路由信息(匹配请求)
	http.HandleFunc("/", app.LoginView)
	http.HandleFunc("/adminlogin", app.AdminLoginView)
	http.HandleFunc("/login", app.Login)
	http.HandleFunc("/login-2", app.Login_2)
	http.HandleFunc("/loginout", app.LoginOut)

	http.HandleFunc("/index", app.Index)
	http.HandleFunc("/help", app.Help)

	http.HandleFunc("/addEduInfo", app.AddEduShow)	// 显示添加信息页面
	http.HandleFunc("/addEdu", app.AddEdu)	// 提交信息请求

	http.HandleFunc("/addCetInfo", app.AddCetShow)	
	http.HandleFunc("/addCet", app.AddCet)	
	
	http.HandleFunc("/queryPage", app.QueryPage)	// 转至根据证书编号与姓名查询EDU
	http.HandleFunc("/query", app.FindCertByNoAndName)	// 根据证书编号与姓名查询信息

	http.HandleFunc("/queryPage2", app.QueryPage2)	// 转至根据身份证号码查询EDU
	http.HandleFunc("/query2", app.FindByID)	

	http.HandleFunc("/queryPage3", app.QueryPage3)	// 转至根据身份证号码查询信息页面
	http.HandleFunc("/query3", app.FindCetByID)	
	
	http.HandleFunc("/history",app.HistoryShow)

	http.HandleFunc("/modifyPage", app.ModifyShow)	// 修改信息页面
	http.HandleFunc("/modify", app.Modify)	//  修改信息

	http.HandleFunc("/upload", app.UploadFile)

	fmt.Println("启动Web服务, 请访问127.0.0.1:9000")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Printf("Web服务启动失败: %v", err)
	}

}



