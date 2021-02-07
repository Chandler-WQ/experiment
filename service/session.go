package service

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"github.com/Chandler-WQ/experiment/common"
	"github.com/Chandler-WQ/experiment/common/model"
	"github.com/Chandler-WQ/experiment/common/pb"
	"github.com/Chandler-WQ/experiment/db"
	"github.com/Chandler-WQ/experiment/en_decryption"
)

func GetSession(token string) (*pb.Session, error) {
	if len(token) < 6 {
		return nil, errors.New("the token is error")
	}
	md5StrTar := token[len(token)-5:]
	md5StrNow := md5.Sum([]byte(token[:len(token)-5]))
	if md5StrTar != fmt.Sprintf("%x", md5StrNow)[:5] {
		return nil, errors.New("the token is error")
	}
	tokenBytes, err := base64.StdEncoding.DecodeString(token[:len(token)-5])
	if err != nil {
		log.Errorf("the err is %s", err)
		return nil, err

	}

	tpass, err := en_decryption.AesDecrypt(tokenBytes, en_decryption.AesKey)
	if err != nil {
		log.Errorf("the err is %s", err)
		return nil, errors.New("the token is error")
	}

	session := pb.Session{}
	err = proto.Unmarshal(tpass, &session)

	if err != nil {
		log.Errorf("the err is %s", err)
		return nil, errors.New("the Unmarshal is error")
	}
	if session.ExpireTime < time.Now().Unix() {
		return nil, errors.New("session到期失效")
	}
	return &session, nil
}

func CreateSession(session *pb.Session) (string, error) {
	sessioByte, err := proto.Marshal(session)
	if err != nil {
		log.Errorf("[CreateSession]the err is %s", err)
		return "", errors.New("the Marshal is error")
	}

	str, err := en_decryption.AesEncrypt(sessioByte, en_decryption.AesKey)
	if err != nil {
		log.Errorf("[CreateSession]the err is %s", err)
		return "", errors.New("the AesEncrypt is error")
	}

	token := base64.StdEncoding.EncodeToString(str)

	has := md5.Sum([]byte(token))
	md5Str := fmt.Sprintf("%x", has)
	tokenMd5 := token + md5Str[:5]
	err = db.Db.CreateSession(&model.Session{
		SessionId:  session.SessionId,
		UserId:     session.UserId,
		ExpireTime: time.Now().Unix() + common.SessionAge,
		Data:       tokenMd5,
	})
	if err != nil {
		log.Errorf("[CreateSession]the db err is %s", err)
		return "", errors.New("the db save session error")
	}

	return tokenMd5, nil
}
