package controller

import (
	"net/http"
	"encoding/json"
	"education/service"
	"fmt"
	"strconv"
)

var data = &struct {
	CurrentUser User
	Flag bool
	Login bool
	Edu service.Education
	CurPicHashCode string
	Msg string
	DataOk bool
	CetListEmpty bool
	EduListEmpty bool
	Index int
	PersonalSpace PersonalSpace
	Cets []service.CetHistoryItem
	CurCet service.CertificateObj
	CetWTBAList []*CetWaitingToApproveStruct
	EduWTBAList []*EduWaitingToApproveStruct
}{
	CurrentUser:User{},
	Flag:false,
	Login:true,
	Edu:service.Education{},
	CurPicHashCode:"",
	Msg:"",
	DataOk:true,
	CetListEmpty:true,
	EduListEmpty:true,
}
func dataReset(){
	data.Flag = false
	data.CurPicHashCode = ""
	data.Msg = ""
	data.DataOk = true
	data.Login = true
	data.CetListEmpty=true
	data.EduListEmpty=true
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
func (app *Application) Index(w http.ResponseWriter, r *http.Request) {
	userCheck(w,r)
	defer dataReset()
	obj := r.FormValue("obj")
	indexStr := r.FormValue("index")
	index,_ := strconv.Atoi(indexStr) 
	data.PersonalSpace = PersonalSpaceMap[data.CurrentUser.LoginName]
	if(len(data.PersonalSpace.CetPtrList) > 0){
		data.CetListEmpty = false
	}
	if(len(data.PersonalSpace.EduPtrList) > 0){
		data.EduListEmpty = false
	}
	if data.CetListEmpty && data.EduListEmpty{
		data.Flag = true
	}
	if(obj == "edu"){
		data.PersonalSpace.EduPtrList[index].UpdateStatusCode(-1,"")
	}else if(obj == "cet"){
		data.PersonalSpace.CetPtrList[index].UpdateStatusCode(-1,"")
	}
	ShowView(w, r, "index.html", data)
}
func (app *Application) Help(w http.ResponseWriter, r *http.Request)  {

	ShowView(w, r, "help.html", data)
}

// 用户登录
func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	loginName := r.FormValue("loginName")
	password := r.FormValue("password")
	res,_ := MySqlLoginCheck(loginName,password)
	if res{
		data.CurrentUser = user
		app.Index(w,r)
		return
	}
	defer dataReset()
	data.Flag = true
	ShowView(w, r, "login.html", data)
}

// 用户登出
func (app *Application) LoginOut(w http.ResponseWriter, r *http.Request)  {
	data.CurrentUser = User{}
	ShowView(w, r, "login.html", data)
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
		return
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
		return
	}
	//提交到待认证列表
	AddEduProposal(&edu,data.CurrentUser.LoginName)
	data.Flag = true
	data.Msg = "添加成功！敬请等待认证"
	ShowView(w, r, "addEdu.html", data)
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
		return
	}
	result,_ = app.Setup.FindCetByCertNoOrTestNo(cet.TestNo,"2")
	savedCet = service.CertificateObj{}
	json.Unmarshal(result, &savedCet)
	if cet.TestNo == savedCet.TestNo{
		fmt.Printf("准考证号已存在:%s\n",cet.TestNo)
		data.DataOk = false
		data.Msg = "准考证号已存在，请重新填写！"
		ShowView(w, r, "addCet.html", data)
		return
	}
	//提交到待认证列表
	AddCetProposal(&cet,data.CurrentUser.LoginName)
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

func (app *Application) CetConfirmShow(w http.ResponseWriter, r *http.Request){
	userCheck(w,r)
	defer dataReset()
	indexStr := r.FormValue("index")
	event := r.FormValue("event")
	var index int
	if indexStr != ""{
		index,_ = strconv.Atoi(indexStr)
	}
	if event == "withdraw"{
		CetWaitingToApproveList[index].UpdateStatusCode(-1,"")
	}else if event == "detail"{
		data.CurCet = CetWaitingToApproveList[index].CetItem
		data.Index = index
		ShowView(w, r, "confirmResult2.html", data)
		return
	}else{
		if(len(data.CetWTBAList) < len(CetWaitingToApproveList)){
			for i := len(data.CetWTBAList);i < len(CetWaitingToApproveList);i++{
				data.CetWTBAList = append(data.CetWTBAList,&CetWaitingToApproveList[i])
			}
		}
	}
	ShowView(w, r, "cetconfirm.html", data)
}
func (app *Application) EduConfirmShow(w http.ResponseWriter, r *http.Request){
	userCheck(w,r)
	defer dataReset()
	indexStr := r.FormValue("index")
	event := r.FormValue("event")
	var index int
	if indexStr != ""{
		index,_ = strconv.Atoi(indexStr)
	}
	if event == "withdraw"{
		EduWaitingToApproveList[index].UpdateStatusCode(-1,"")
	}else if event == "detail"{
		data.Edu = EduWaitingToApproveList[index].EduItem
		data.Index = index
		ShowView(w, r, "confirmResult.html", data)
	}else{
		if(len(data.EduWTBAList) < len(EduWaitingToApproveList)){
			for i := len(data.EduWTBAList);i < len(EduWaitingToApproveList);i++{
				data.EduWTBAList = append(data.EduWTBAList,&EduWaitingToApproveList[i])
			}
		}
	}
	ShowView(w, r, "educonfirm.html", data)
}
func (app *Application) EduConfirm(w http.ResponseWriter, r *http.Request){
	defer dataReset()
	indexStr := r.FormValue("index")
	index,_ := strconv.Atoi(indexStr)
	txId,err := app.Setup.SaveEdu(EduWaitingToApproveList[index].EduItem)
	if err != nil{
		data.Flag = true
		data.Msg = "添加时发生错误！"
		ShowView(w, r, "educonfirm.html", data)
	}
	data.Flag = true
	data.Msg = "交易成功！交易编号：" + txId
	EduWaitingToApproveList[index].UpdateStatusCode(1,data.CurrentUser.LoginName)
	ShowView(w, r, "confirmResult.html", data)
}
func (app *Application) CetConfirm(w http.ResponseWriter, r *http.Request){
	defer dataReset()
	indexStr := r.FormValue("index")
	index,_ := strconv.Atoi(indexStr)
	txId,err := app.Setup.SaveCet(CetWaitingToApproveList[index].CetItem)
	if err != nil{
		data.Flag = true
		data.Msg = "添加时发生错误！"
		ShowView(w, r, "cetconfirm.html", data)
	}
	data.Flag = true
	data.Msg = "交易成功！交易编号：" + txId
	CetWaitingToApproveList[index].UpdateStatusCode(1,data.CurrentUser.LoginName)
	ShowView(w, r, "confirmResult2.html", data)
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
