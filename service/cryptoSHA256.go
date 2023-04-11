package service

import (
	"crypto/sha256"
	"os"
	"encoding/hex"
)

func Sha256(src string) string{
	m := sha256.New()
	m.Write([]byte(src))
	res := hex.EncodeToString(m.Sum(nil))
	
	return res
}
func GetPicSha256(photoPath string) string{
	if(photoPath != ""){
		content,err := os.ReadFile("/root/go/src/education/web"+ photoPath)
		if err != nil{
			return ""
		}
		return Sha256(string(content))
	}
	return ""
}