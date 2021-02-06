package en_decryption

import (
	"encoding/base64"
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	"github.com/Chandler-WQ/experiment/common/pb"
)

func TestAesEncrypt(t *testing.T) {
	//key的长度必须是16、24或者32字节，分别用于选择AES-128, AES-192, or AES-256
	key := "12345678abcdefgh"
	session := pb.Session{
		Name:         "zhangsan",
		SessionId:    23123127312963912,
		UserId:       12312312731291,
		PassportName: "zhangsansdasa",
		College:      "计算机院",
		UserType:     1,
	}
	sessiobyte, err := proto.Marshal(&session)
	assert.Nil(t, err)
	str, err := AesEncrypt(sessiobyte, []byte(key))
	assert.Nil(t, err)
	pass64 := base64.StdEncoding.EncodeToString(str)
	log.Infof("the str is %s", pass64)

	bytesPass, err := base64.StdEncoding.DecodeString(pass64)
	if err != nil {
		log.Infof("the err is %s", err)
		return
	}

	tpass, err := AesDecrypt(bytesPass, []byte(key))
	if err != nil {
		log.Infof("the str is %s", tpass)
		return
	}
	fmt.Printf("解密后:%s\n", tpass)
	sessionTem := pb.Session{}
	err = proto.Unmarshal(tpass, &sessionTem)
	if err != nil {
		log.Infof("the err is %s", err)
		return
	}
	log.Infof("the sessionTem is %v", sessionTem)

}
