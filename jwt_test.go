package sharedlib_test

import (
	"fmt"
	"testing"

	sharedlib "github.com/ilovelili/dongfeng-shared-lib"
	"github.com/stretchr/testify/assert"
)

const (
	clientid  = "xxx"
	secretkey = "yyy"
)

func TestVerifyLogin(t *testing.T) {
	idtoken := "zzz"
	client := sharedlib.NewAuthClient(clientid, secretkey)
	status, err := client.VerifyLogin(idtoken)
	assert.True(t, status)
	assert.Nil(t, err)
}

func TestParseUserInfo(t *testing.T) {
	userid := "5d3d973d75115d234ffa762b"
	client := sharedlib.NewAuthClient(clientid, secretkey)
	userinfo, err := client.ParseUserInfo(userid)

	fmt.Printf("user info is %s\n", userinfo)
	assert.Nil(t, err)
}
