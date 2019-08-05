package sharedlib_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ilovelili/dongfeng-shared-lib"
)

func TestVerifyLogin(t *testing.T) {
	appid := "xxx"
	secretkey := "yyy"
	idtoken := "zzz"
	client := sharedlib.NewAuthClient(appid, secretkey)
	status, err := client.VerifyLogin(idtoken)
	assert.True(t, status)
	assert.Nil(t, err)
}
