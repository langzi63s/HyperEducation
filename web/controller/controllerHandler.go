package controller

import (
	"net/http"
	"encoding/json"
	"education/service"
	"fmt"
)

var data = &struct {
	CurrentUser User
	Flag bool
	Edu service.Education
	History bool
}{
	CurrentUser:User{},
	Flag:false,
	Edu:service.Education{},
	History:false,
}
func (app *Application) LoginView(w http.ResponseWriter, r *http.Request)  {
	ShowView(w, r, "login.html", nil)
}
func (app *Application) AdminLoginView(w http.ResponseWriter, r *http.Request)  {
	ShowView(w, r, "login-2.html", nil)
}
func (app *Application) Index(w http.ResponseWriter, r *http.Request)  {
	ShowView(w, r, "index.html", data)
}

func (app *Application) Help(w http.ResponseWriter, r *http.Request)  {
	ShowView(w, r, "help.html", data)
}

// 用户登录
func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	data.Flag = false
	loginName := r.FormValue("loginName")
	password := r.FormValue("password")

	for _, user := range users {
		if user.LoginName == loginName && user.Password == password {
			data.CurrentUser = user
			ShowView(w, r, "index.html", data)
			return
		}
	}
	data.Flag = true
	ShowView(w, r, "login.html", data)
}
func (app *Application) Login_2(w http.ResponseWriter, r *http.Request) {
	loginName := r.FormValue("loginName")
	password := r.FormValue("password")

	for _, user := range users {
		if user.LoginName == loginName && user.Password == password {
			data.CurrentUser = user
			ShowView(w, r, "index.html", data)
			return
		}
	}
	data.Flag = true
	ShowView(w, r, "login-2.html", data)
}

// 用户登出
func (app *Application) LoginOut(w http.ResponseWriter, r *http.Request)  {
	data.CurrentUser = User{}
	data.Flag = false
	ShowView(w, r, "login.html", nil)
}

// 显示添加信息页面
func (app *Application) AddEduShow(w http.ResponseWriter, r *http.Request)  {
	ShowView(w, r, "addEdu.html", data)
}

// 添加信息
func (app *Application) AddEdu(w http.ResponseWriter, r *http.Request) {
	edu := service.Education{
		Name:r.FormValue("name"),
		Gender:r.FormValue("gender"),
		Nation:r.FormValue("nation"),
		EntityID:r.FormValue("entityID"),
		Place:r.FormValue("place"),
		BirthDay:r.FormValue("birthDay"),
		EnrollDate:r.FormValue("enrollDate"),
		GraduationDate:r.FormValue("graduationDate"),
		SchoolName:r.FormValue("schoolName"),
		Major:r.FormValue("major"),
		QuaType:r.FormValue("quaType"),
		Length:r.FormValue("length"),
		Mode:r.FormValue("mode"),
		Level:r.FormValue("level"),
		Graduation:r.FormValue("graduation"),
		CertNo:r.FormValue("certNo"),
		Photo:r.FormValue("photo"),
	}

	app.Setup.SaveEdu(edu)
	/*transactionID, err := app.Setup.SaveEdu(edu)

	if err != nil {
		data.Msg = err.Error()
	}else{
		data.Msg = "信息添加成功:" + transactionID
	}*/

	//ShowView(w, r, "addEdu.html", data)
	r.Form.Set("certNo", edu.CertNo)
	r.Form.Set("name", edu.Name)
	app.FindCertByNoAndName(w, r)
}

func (app *Application) QueryPage(w http.ResponseWriter, r *http.Request)  {
	data.Flag = false
	ShowView(w, r, "query.html", data)
}

// 根据证书编号与姓名查询信息
func (app *Application) FindCertByNoAndName(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "GET" {
		data.Flag = false
		ShowView(w, r, "query.html", data)
	}
	certNo := r.FormValue("certNo")
	name := r.FormValue("name")
	result, err := app.Setup.FindEduByCertNoAndName(certNo, name)
	var edu = service.Education{}
	json.Unmarshal(result, &edu)
	if edu.Name != ""{
		fmt.Println("根据证书编号与姓名查询信息成功：")
		fmt.Println(edu)
		data.Edu = edu
		data.History = false
		if err != nil {
			data.Flag = true
		}
		ShowView(w, r, "queryResult.html", data)
	}else{
		data.Flag = true
		ShowView(w, r, "query.html", data)
	}
}

func (app *Application) QueryPage2(w http.ResponseWriter, r *http.Request)  {
	data.Flag = false
	ShowView(w, r, "query2.html", data)
}

// 根据身份证号码查询信息
func (app *Application) FindByID(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "GET" {
		data.Flag = false
		ShowView(w, r, "query2.html", data)
	}
	entityID := r.FormValue("entityID")
	result, err := app.Setup.FindEduInfoByEntityID(entityID)
	var edu = service.Education{}
	json.Unmarshal(result, &edu)
	if edu.Name != ""{
		fmt.Println("根据身份证号码查询信息成功：")
		fmt.Println(edu)
		data.Edu = edu
		data.History = true
		if err != nil {
			data.Flag = true
		}
		ShowView(w, r, "queryResult.html", data)
	}else{
		data.Flag = true
		ShowView(w, r, "query2.html", data)
	}
}

// 修改/添加新信息
func (app *Application) ModifyShow(w http.ResponseWriter, r *http.Request)  {
	// 根据证书编号与姓名查询信息
	certNo := r.FormValue("certNo")
	name := r.FormValue("name")
	result, err := app.Setup.FindEduByCertNoAndName(certNo, name)

	var edu = service.Education{}
	json.Unmarshal(result, &edu)

	data.Edu = edu
	data.Flag=false
	data.History = true
	if err != nil {
		data.Flag = true
	}
	ShowView(w, r, "modify.html", data)
}

// 修改/添加新信息
func (app *Application) Modify(w http.ResponseWriter, r *http.Request) {
	edu := service.Education{
		Name:r.FormValue("name"),
		Gender:r.FormValue("gender"),
		Nation:r.FormValue("nation"),
		EntityID:r.FormValue("entityID"),
		Place:r.FormValue("place"),
		BirthDay:r.FormValue("birthDay"),
		EnrollDate:r.FormValue("enrollDate"),
		GraduationDate:r.FormValue("graduationDate"),
		SchoolName:r.FormValue("schoolName"),
		Major:r.FormValue("major"),
		QuaType:r.FormValue("quaType"),
		Length:r.FormValue("length"),
		Mode:r.FormValue("mode"),
		Level:r.FormValue("level"),
		Graduation:r.FormValue("graduation"),
		CertNo:r.FormValue("certNo"),
		Photo:r.FormValue("photo"),
	}

	//transactionID, err := app.Setup.ModifyEdu(edu)
	app.Setup.ModifyEdu(edu)

	/*data := &struct {
		Edu service.Education
		CurrentUser User
		Msg string
		Flag bool
	}{
		CurrentUser:cuser,
		Flag:true,
		Msg:"",
	}

	if err != nil {
		data.Msg = err.Error()
	}else{
		data.Msg = "新信息添加成功:" + transactionID
	}

	ShowView(w, r, "modify.html", data)
	*/
	r.Form.Set("certNo", edu.CertNo)
	r.Form.Set("name", edu.Name)
	app.FindCertByNoAndName(w, r)
}
