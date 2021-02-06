package util

import (
	"errors"
)

func CodeToIdentity(code int32) (string, error) {
	switch code {
	case 1:
		return "学生", nil
	case 2:
		return "老师", nil
	case 3:
		return "管理员", nil
	default:
		return "", errors.New("the code is invalid")

	}
}

func IdentityToCode(identity string) (int32, error) {
	switch identity {
	case "学生":
		return 1, nil
	case "老师":
		return 2, nil
	case "管理员":
		return 3, nil
	default:
		return -1, errors.New("the identity is invalid")
	}
}
