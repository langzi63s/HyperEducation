package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"fmt"
	"encoding/json"
	"bytes"
	"time"
)


type Education struct {
	ObjectType	string	`json:"docType"`
	Name	string	`json:"Name"`		// 姓名
	Gender	string	`json:"Gender"`		// 性别
	Nation	string	`json:"Nation"`		// 民族
	EntityID	string	`json:"EntityID"`		// 身份证号
	Place	string	`json:"Place"`		// 籍贯
	BirthDay	string	`json:"BirthDay"`		// 出生日期

	EnrollDate	string	`json:"EnrollDate"`		// 入学日期
	GraduationDate	string	`json:"GraduationDate"`	// 毕（结）业日期
	SchoolName	string	`json:"SchoolName"`	// 学校名称
	Major	string	`json:"Major"`	// 专业
	QuaType	string	`json:"QuaType"`	// 学历类别
	Length	string	`json:"Length"`	// 学制
	Mode	string	`json:"Mode"`	// 学习形式
	Level	string	`json:"Level"`	// 层次
	Graduation	string	`json:"Graduation"`	// 毕（结）业
	CertNo	string	`json:"CertNo"`	// 证书编号
	Photo	string	`json:"Photo"`	// 照片
	TimeStamp string `json:"TimeStamp"` //交易时间
	PhotoHashCode string `json:"PhotoHashCode"`
	TxID	string  `json:TxID`
	Historys	[]HistoryItem	// 当前edu的历史记录
}
type CertificateObj struct {
	ObjectType	string	`json:"docType"`
	Name	string	`json:"Name"`		// 姓名
	Gender	string	`json:"Gender"`		// 性别
	EntityID	string  `json:"EntityID"`
	TimeStamp   string  `json:"TimeStamp"`
	Level	    string	`json:"Level"`	// 等级 必须是四 或 六
	CertNo	    string	`json:"CertNo"`	// 证书编号
	TestTime    string  `json:"TestTime"` //考试时间 xxxx年xx月
	TestNo      string  `json:"TestNo"` //准考证号
	Score       string  `json:"Score"` //考试分数
	TxID	string  `json:TxID`
}
type HistoryItem struct {
	TxId	string
	Education	Education
}
type CetHistoryItem struct {
	Certificate	CertificateObj
}

type EducationChaincode struct {

}

func (t *EducationChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response{
	fmt.Println(" ==== Init ====")
	return shim.Success(nil)
}

func (t *EducationChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response{
	// 获取用户意图
	fun, args := stub.GetFunctionAndParameters()

	if fun == "addEdu"{
		return t.addEdu(stub, args)		// 添加信息
	}else if fun == "queryEduByCertNoAndName" {
		return t.queryEduByCertNoAndName(stub, args)		// 根据证书编号及姓名查询信息
	}else if fun == "queryEduInfoByEntityID" {
		return t.queryEduInfoByEntityID(stub, args)	// 根据身份证号码及姓名查询详情
	}else if fun == "updateEdu" {
		return t.updateEdu(stub, args)		// 根据证书编号更新信息
	}else if fun == "delEdu"{
		return t.delEdu(stub, args)	// 根据证书编号删除信息
	}else if fun == "addCet"{
		return t.addCet(stub, args)	
	}else if fun == "queryCetByCertNoOrTestNo"{
		return t.queryCetByCertNoOrTestNo(stub, args)
	}else if fun == "queryCetInfoByEntityID"{
		return t.queryCetInfoByEntityID(stub, args)
	}
	return shim.Error("指定的函数名称错误")

}


const DOC_TYPE = "eduObj"
const CET_DOC_TYPE = "cetObj"
const INDEX_NAME = "CET~id~level~test_time"
func GetCetIndexKey(stub shim.ChaincodeStubInterface, cet *CertificateObj) (string,bool){
	indexKey,err := stub.CreateCompositeKey(INDEX_NAME,[]string{cet.EntityID, cet.Level, cet.TestTime})
	if err != nil{
		return "",false
	}
	return indexKey,true
}
func PutEdu(stub shim.ChaincodeStubInterface, edu Education) ([]byte, bool) {

	edu.ObjectType = DOC_TYPE

	b, err := json.Marshal(edu)
	if err != nil {
		return nil, false
	}
	err = stub.PutState(edu.EntityID, b)
	if err != nil {
		return nil, false
	}

	return b, true
}
func PutCet(stub shim.ChaincodeStubInterface, cet CertificateObj) ([]byte, bool){
	cet.ObjectType = CET_DOC_TYPE
	b,err := json.Marshal(cet)
	if err != nil {
		return nil, false
	}
	indexKey,_ := GetCetIndexKey(stub, &cet)
	err = stub.PutState(indexKey, b)
	if err != nil {
		return nil, false
	}

	return b, true
}
func GetEduInfo(stub shim.ChaincodeStubInterface, entityID string) (Education, bool)  {
	var edu Education
	b, err := stub.GetState(entityID)
	if err != nil {
		return edu, false
	}

	if b == nil {
		return edu, false
	}

	err = json.Unmarshal(b, &edu)
	if err != nil {
		return edu, false
	}

	// 返回结果
	return edu, true
}
func GetCetInfo(stub shim.ChaincodeStubInterface, entityID string, level string, test_time string) (CertificateObj, bool)  {
	var cet CertificateObj
	cet.EntityID = entityID
	cet.Level = level
	cet.TestTime = test_time
	indexKey,_ := GetCetIndexKey(stub, &cet)
	b, err := stub.GetState(indexKey)
	if err != nil {
		return cet, false
	}

	if b == nil {
		return cet, false
	}

	err = json.Unmarshal(b, &cet)
	if err != nil {
		return cet, false
	}
	return cet, true
}

// 根据指定的查询字符串实现富查询
func getDataByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil

}

// 添加信息
// args: educationObject
// 身份证号为 key, Education 为 value
func (t *EducationChaincode) addEdu(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2{
		return shim.Error("给定的参数个数不符合要求")
	}

	var edu Education
	err := json.Unmarshal([]byte(args[0]), &edu)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}

	// 查重: 身份证号码必须唯一
	_, exist := GetEduInfo(stub, edu.EntityID)
	if exist {
		return shim.Error("要添加的身份证号码已存在")
	}
	tm,_ := stub.GetTxTimestamp()
	id := stub.GetTxID()
	edu.TxID = id
	edu.TimeStamp = time.Unix(tm.Seconds+ 3600 * 8, int64(tm.Nanos)).Format("2006-01-02 15:04:05")
	_, bl := PutEdu(stub, edu)
	if !bl {
		return shim.Error("保存信息时发生错误")
	}
	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息添加成功"))
}
func (t *EducationChaincode) addCet(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2{
		return shim.Error("给定的参数个数不符合要求")
	}

	var cet CertificateObj
	err := json.Unmarshal([]byte(args[0]), &cet)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}

	_, exist := GetCetInfo(stub, cet.EntityID, cet.Level, cet.TestTime)
	if exist {
		return shim.Error("信息已存在")
	}
	tm,_ := stub.GetTxTimestamp()
	id := stub.GetTxID()
	cet.TxID = id
	cet.TimeStamp = time.Unix(tm.Seconds+ 3600 * 8, int64(tm.Nanos)).Format("2006-01-02 15:04:05")
	_, bl := PutCet(stub, cet)
	if !bl {
		return shim.Error("保存信息时发生错误")
	}
	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("信息添加成功"))
}
// 根据证书编号及姓名查询信息
// args: CertNo, name
func (t *EducationChaincode) queryEduByCertNoAndName(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	CertNo := args[0]
	name := args[1]

	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串)
	// queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"eduObj\", \"CertNo\":\"%s\"}}", CertNo)
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\", \"CertNo\":\"%s\", \"Name\":\"%s\"}}", DOC_TYPE, CertNo, name)

	// 查询数据
	result, err := getDataByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据证书编号及姓名查询信息时发生错误")
	}
	if result == nil {
		return shim.Error("根据指定的证书编号及姓名没有查询到相关的信息")
	}
	return shim.Success(result)
}
func (t *EducationChaincode) queryCetByCertNoOrTestNo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	cetNo := args[0]
	No := args[1]
	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串)
	// queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"eduObj\", \"CertNo\":\"%s\"}}", CertNo)
	var queryString string
	if No == "1"{
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"%s\", \"CertNo\":\"%s\"}}", CET_DOC_TYPE, cetNo)
	}else if No == "2"{
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"%s\", \"TestNo\":\"%s\"}}", CET_DOC_TYPE, cetNo)
	}else{
		return shim.Error("给定的参数args[1]不符合要求")
	}
	// 查询数据
	result, err := getDataByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据证书编号或准考证号查询信息时发生错误")
	}
	if result == nil {
		return shim.Error("根据指定的证书编号或准考证号没有查询到相关的信息")
	}
	return shim.Success(result)
}
// 根据身份证号码查询详情（溯源）
// args: entityID
func (t *EducationChaincode) queryEduInfoByEntityID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}

	b, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("根据身份证号码查询信息失败")
	}

	if b == nil {
		return shim.Error("根据身份证号码没有查询到相关的信息")
	}

	var edu Education
	err = json.Unmarshal(b, &edu)
	if err != nil {
		return  shim.Error("反序列化edu信息失败")
	}

	// 获取历史变更数据
	iterator, err := stub.GetHistoryForKey(edu.EntityID)
	if err != nil {
		return shim.Error("根据指定的身份证号码查询对应的历史变更数据失败")
	}
	defer iterator.Close()

	// 迭代处理
	var historys []HistoryItem
	var hisEdu Education
	for iterator.HasNext() {
		hisData, err := iterator.Next()

		if err != nil {
			return shim.Error("获取edu的历史变更数据失败")
		}
		
		var historyItem HistoryItem
		historyItem.TxId = hisData.TxId
		json.Unmarshal(hisData.Value, &hisEdu)

		if hisData.Value == nil {
			var empty Education
			historyItem.Education = empty
		}else {
			historyItem.Education = hisEdu
		}

		historys = append(historys, historyItem)

	}

	edu.Historys = historys

	result, err := json.Marshal(edu)
	if err != nil {
		return shim.Error("序列化edu信息时发生错误")
	}
	return shim.Success(result)
}
func (t *EducationChaincode) queryCetInfoByEntityID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1{
		return shim.Error("给定的参数个数不符合要求")
	}
	resultsIterator, err := stub.GetStateByPartialCompositeKey(INDEX_NAME,[]string{args[0]})
	if err != nil {
		return shim.Error("根据身份证号码查询信息失败")
	}
	defer resultsIterator.Close()
	// 迭代处理
	var hisCet CertificateObj
	var CetHistorys	[]CetHistoryItem
	for resultsIterator.HasNext() {
		hisData, err := resultsIterator.Next()

		if err != nil {
			return shim.Error("获取该身份证CET数据失败")
		}
		
		//_, compositeKeyParts,_ := stub.SplitCompositeKey(hisData.Key)
		var historyItem CetHistoryItem		
		json.Unmarshal(hisData.Value, &hisCet)
		if hisData.Value == nil {
			var empty CertificateObj
			historyItem.Certificate = empty
		}else {
			historyItem.Certificate = hisCet
		}
		CetHistorys = append(CetHistorys, historyItem)
	}


	// 返回
	result, err := json.Marshal(CetHistorys)
	if err != nil {
		return shim.Error("序列化edu信息时发生错误")
	}
	return shim.Success(result)
}
// 根据身份证号更新信息
// args: educationObject
func (t *EducationChaincode) updateEdu(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2{
		return shim.Error("给定的参数个数不符合要求")
	}

	var info Education
	err := json.Unmarshal([]byte(args[0]), &info)
	if err != nil {
		return  shim.Error("反序列化edu信息失败")
	}

	// 根据身份证号码查询信息
	result, bl := GetEduInfo(stub, info.EntityID)
	if !bl{
		return shim.Error("根据身份证号码查询信息时发生错误")
	}
	
	result.Name = info.Name
	result.BirthDay = info.BirthDay 
	result.Nation = info.Nation 
	result.Gender = info.Gender 
	result.Place = info.Place 
	result.EntityID = info.EntityID
	result.Photo = info.Photo 


	result.EnrollDate = info.EnrollDate
	result.GraduationDate = info.GraduationDate
	result.SchoolName = info.SchoolName
	result.Major = info.Major 
	result.QuaType = info.QuaType 
	result.Length = info.Length 
	result.Mode = info.Mode 
	result.Level = info.Level 
	result.Graduation = info.Graduation
	result.CertNo = info.CertNo
	tm,_ := stub.GetTxTimestamp()
	result.TimeStamp = time.Unix(tm.Seconds + 8 * 3600, int64(tm.Nanos)).Format("2006-01-02 15:04:05")
	result.TxID = stub.GetTxID()
	_, bl = PutEdu(stub, result)
	if !bl {
		return shim.Error("保存信息信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息更新成功"))
}

// 根据身份证号删除信息（暂不提供）
// args: entityID
func (t *EducationChaincode) delEdu(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2{
		return shim.Error("给定的参数个数不符合要求")
	}

	/*var edu Education
	result, bl := GetEduInfo(stub, info.EntityID)
	err := json.Unmarshal(result, &edu)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}*/

	err := stub.DelState(args[0])
	if err != nil {
		return shim.Error("删除信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息删除成功"))
}

func main(){
	fmt.Printf(">>启动EducationChaincode中 ")
	err := shim.Start(new(EducationChaincode))
	if err != nil{
		fmt.Printf(">>启动EducationChaincode时发生错误: %s", err)
	}
}

