package controller

import (
	"net/http"
	"encoding/json"
	"education/service"
	"fmt"
	"time"
)

var data = &struct {
	CurrentUser User
	Flag bool
	Login bool
	Edu service.Education
	CurPicHashCode string
	Msg string
	DataOk bool
	Cets []service.CetHistoryItem
	CurCet service.CertificateObj
}{
	CurrentUser:User{},
	Flag:false,
	Login:true,
	Edu:service.Education{},
	CurPicHashCode:"",
	Msg:"",
	DataOk:true,
}
func dataReset(){
	data.Flag = false
	data.Edu = service.Education{}
	data.CurCet = service.CertificateObj{}
	data.CurPicHashCode = ""
	data.Msg = ""
	data.DataOk = true
	data.Login = true
}
func userCheck(w http.ResponseWriter, r *http.Request){
	if data.CurrentUser == (User{}){
		defer dataReset()
		data.Login = false
		ShowView(w, r, "login.html", data)
	}
}

func (app *Application) LoginView(w http.ResponseWriter, r *http.Request)  {
	ShowView(w, r, "login.html", data)
}
func (app *Application) AdminLoginView(w http.ResponseWriter, r *http.Request)  {
	ShowView(w, r, "login-2.html", data)
}
func (app *Application) Index(w http.ResponseWriter, r *http.Request)  {
	userCheck()
	ShowView(w, r, "index.html", data)
}
func (app *Application) Help(w http.ResponseWriter, r *http.Request)  {

	ShowView(w, r, "help.html", data)
}

// 用户登录
func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	loginName := r.FormValue("loginName")
	password := r.FormValue("password")
	for _, user := range users {
		if user.LoginName == loginName && user.Password == password {
			data.CurrentUser = user
			ShowView(w, r, "index.html", data)
			return
		}
	}
	defer dataReset()
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
	defer dataReset()
	data.Flag = true
	ShowView(w, r, "login-2.html", data)
}

// 用户登出
func (app *Application) LoginOut(w http.ResponseWriter, r *http.Request)  {
	defer dataReset()
	data.CurrentUser = User{}
	ShowView(w, r, "login.html", nil)
}

// 显示添加信息页面
func (app *Application) AddEduShow(w http.ResponseWriter, r *http.Request)  {
	userCheck(w,r)
	ShowView(w, r, "addEdu.html", data)
}
func (app *Application) AddCetShow(w http.ResponseWriter, r *http.Request)  {
	userCheck(w,r)
	ShowView(w, r, "addCet.html", data)
}
// 添加信息
func (app *Application) AddEdu(w http.ResponseWriter, r *http.Request) {
	userCheck(w,r)
	//获取数据
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
	defer dataReset()
	//检查身份证号
	result,_ := app.Setup.FindEduInfoByEntityID(edu.EntityID)
	savedEdu := service.Education{}
	json.Unmarshal(result, &savedEdu)
	if edu.EntityID == savedEdu.EntityID{
		fmt.Printf("身份证号已存在:%s\n",edu.EntityID)

		data.DataOk = false
		data.Msg = "身份证号已存在，请重新填写！"
		ShowView(w, r, "addEdu.html", data)
	}
	//检查证书编号
	result,_ = app.Setup.FindEduByCertNoAndName(edu.CertNo,edu.Name)
	savedEdu = service.Education{}
	json.Unmarshal(result, &savedEdu)
	if edu.CertNo == savedEdu.CertNo{
		fmt.Printf("证书编号已存在:%s\n",edu.CertNo)
		data.DataOk = false
		data.Msg = "证书编号已存在，请重新填写！"
		ShowView(w, r, "addEdu.html", data)
	}
	//提交到待认证列表
	tStr := time.Now().Format("2006-01-02 15:04:05")
	EduWaitingToApproveObj :=EduWaitingToApproveStruct{
		proposer:Proposer {
			Userinfo:data.CurrentUser,
			ProTime:tStr,
		},
		eduItem:edu,
	}

	EduWaitingToApproveList = append(EduWaitingToApproveList,EduWaitingToApproveObj)
	data.Flag = true
	data.Msg = "添加成功！敬请等待认证"
	ShowView(w, r, "addEdu.html", data)
	/*
	r.Form.Set("certNo", edu.CertNo)
	r.Form.Set("name", edu.Name)
	app.FindCertByNoAndName(w, r)
	*/
}
func (app *Application) AddCet(w http.ResponseWriter, r *http.Request) {
	userCheck(w,r)
	cet := service.CertificateObj{
		Name:r.FormValue("name"),
		Gender:r.FormValue("gender"),
		Score:r.FormValue("score"),
		EntityID:r.FormValue("entityID"),
		TestNo:r.FormValue("testNo"),
		TestTime:r.FormValue("testTime"),
		Level:r.FormValue("level"),
		CertNo:r.FormValue("certNo"),
	}
	defer dataReset()
	//检查证书编号准考证号
	result,_ := app.Setup.FindCetByCertNoOrTestNo(cet.CertNo,"1")
	savedCet := service.CertificateObj{}
	json.Unmarshal(result, &savedCet)
	if cet.CertNo == savedCet.CertNo{
		fmt.Printf("证书编号已存在:%s\n",cet.CertNo)
		data.DataOk = false
		data.Msg = "证书编号已存在，请重新填写！"
		ShowView(w, r, "addCet.html", data)
	}
	result,_ = app.Setup.FindCetByCertNoOrTestNo(cet.TestNo,"2")
	savedCet = service.CertificateObj{}
	json.Unmarshal(result, &savedCet)
	if cet.TestNo == savedCet.TestNo{
		fmt.Printf("准考证号已存在:%s\n",cet.TestNo)
		data.DataOk = false
		data.Msg = "准考证号已存在，请重新填写！"
		ShowView(w, r, "addCet.html", data)
	}
	//提交到待认证列表
	tStr := time.Now().Format("2006-01-02 15:04:05")
	CetWaitingToApproveObj :=CetWaitingToApproveStruct{
		proposer:Proposer{
			Userinfo:data.CurrentUser,
			ProTime:tStr,
		},
		cetItem:cet,
	}
	CetWaitingToApproveList = append(CetWaitingToApproveList,CetWaitingToApproveObj)
	data.Msg = "添加成功！敬请等待认证"
	data.Flag = true
	ShowView(w, r, "addCet.html", data)
}

func (app *Application) QueryPage(w http.ResponseWriter, r *http.Request)  {
	userCheck(w,r)
	ShowView(w, r, "query.html", data)
}

// 根据证书编号与姓名查询信息
func (app *Application) FindCertByNoAndName(w http.ResponseWriter, r *http.Request)  {
	userCheck(w,r)
	defer dataReset()
	if r.Method == "GET" {
		ShowView(w, r, "query.html", data)
	}
	certNo := r.FormValue("certNo")
	name := r.FormValue("name")
	result,_ := app.Setup.FindEduByCertNoAndName(certNo, name)
	var edu = service.Education{}
	json.Unmarshal(result, &edu)
	if edu.Name != ""{
		fmt.Println("根据证书编号与姓名查询信息成功：")
		fmt.Println(edu)
		data.Edu = edu
		data.CurPicHashCode = service.GetPicSha256(edu.Photo)
		ShowView(w, r, "queryResult.html", data)
	}else{
		data.Flag = true
		ShowView(w, r, "query.html", data)
	}
}

func (app *Application) QueryPage2(w http.ResponseWriter, r *http.Request)  {
	userCheck(w,r)
	ShowView(w, r, "query2.html", data)
}

// 根据身份证号码查询信息
func (app *Application) FindByID(w http.ResponseWriter, r *http.Request)  {
	userCheck(w,r)
	defer dataReset()
	if r.Method == "GET" {
		ShowView(w, r, "query2.html", data)
	}
	entityID := r.FormValue("entityID")
	result,_ := app.Setup.FindEduInfoByEntityID(entityID)
	var edu = service.Education{}
	json.Unmarshal(result, &edu)
	if edu.Name != ""{
		fmt.Println("根据身份证号码查询信息成功：")
		fmt.Println(edu)
		data.Edu = edu
		data.CurPicHashCode = service.GetPicSha256(edu.Photo)
		ShowView(w, r, "queryResult.html", data)
	}else{
		data.Flag = true
		ShowView(w, r, "query2.html", data)
	}
} 
func (app *Application) QueryPage3(w http.ResponseWriter, r *http.Request)  {
	userCheck(w,r)
	ShowView(w, r, "query3.html", data)
}

// 根据身份证号码查询信息
func (app *Application) FindCetByID(w http.ResponseWriter, r *http.Request)  {
	userCheck(w,r)
	defer dataReset()
	if r.Method == "GET" {
		ShowView(w, r, "query3.html", data)
	}
	entityID := r.FormValue("entityID")
	result,_ := app.Setup.FindCetInfoByEntityID(entityID)
	var Certificates []service.CetHistoryItem
	json.Unmarshal(result, &Certificates)
	if len(Certificates) > 0{
		fmt.Println("根据身份证号码查询信息成功：")
		fmt.Println(Certificates)
		data.Cets = Certificates
		ShowView(w, r, "queryResult2.html", data)
	}else{
		data.Flag = true
		ShowView(w, r, "query3.html", data)
	}
}
func (app *Application) FindCetByCertNoOrTestNoShow(w http.ResponseWriter, r *http.Request)  {
	userCheck(w,r)
	defer dataReset()
	if r.Method == "GET" {
		ShowView(w, r, "query4.html", data)
	}
	cetNo := r.FormValue("cetNo")
	No := r.FormValue("No")
	result,_ := app.Setup.FindCetByCertNoOrTestNo (cetNo,No)
	var Certificate service.CertificateObj
	json.Unmarshal(result, &Certificate)
	if Certificate.CertNo != ""{
		fmt.Println("根据证书编号或准考证号查询信息成功：")
		fmt.Println(Certificate)
		data.CurCet = Certificate
		data.Flag = true
		ShowView(w, r, "query4.html", data)
	}else{
		ShowView(w, r, "query4.html", data)
	}
}
func (app *Application) HistoryShow(w http.ResponseWriter, r *http.Request)  {
	userCheck(w,r)
	defer dataReset()
	entityID := r.FormValue("entityID")
	result,_ := app.Setup.FindEduInfoByEntityID(entityID)

	var edu = service.Education{}
	json.Unmarshal(result, &edu)

	data.Edu = edu
	ShowView(w, r, "history.html", data)
}
// 修改/添加新信息
func (app *Application) ModifyShow(w http.ResponseWriter, r *http.Request)  {
	userCheck(w,r)
	defer dataReset()
	// 根据证书编号与姓名查询信息
	certNo := r.FormValue("certNo")
	name := r.FormValue("name")
	result,_ := app.Setup.FindEduByCertNoAndName(certNo, name)

	var edu = service.Education{}
	json.Unmarshal(result, &edu)

	data.Edu = edu

	ShowView(w, r, "modify.html", data)
}

// 修改/添加新信息
func (app *Application) Modify(w http.ResponseWriter, r *http.Request) {
	userCheck(w,r)
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

	app.Setup.ModifyEdu(edu)

	r.Form.Set("certNo", edu.CertNo)
	r.Form.Set("name", edu.Name)
	app.FindCertByNoAndName(w, r)
}
