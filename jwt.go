package sharedlib

import (
	"encoding/json"
	"sync"

	authing "github.com/Authing/authing-go-sdk"
	"github.com/kelvinji2009/graphql"
	"fmt"
)

var (
	once     sync.Once
	instance *Client
)

// Client authing client wrapper
type Client struct {
	AuthingClient *authing.Client
	ClientID      string
}

// NewAuthClient  new authing client
func NewAuthClient(clientID, appSecret string) *Client {
	once.Do(func() {
		client := authing.NewClient(clientID, appSecret, false)
		instance = &Client{
			AuthingClient: client,
			ClientID:      clientID,
		}
	})

	return instance
}

// VerifyLogin verifiy login status by idtoken
func (c *Client) VerifyLogin(idtoken string) (bool, error) {
	p := authing.CheckLoginStatusQueryParameter{
		Token: graphql.String(idtoken),
	}

	q, err := c.AuthingClient.CheckLoginStatus(&p)
	if err != nil {
		return false, err
	}

	return bool(q.CheckLoginStatus.Status), nil
}

// ParseUserInfo parse user info based on user id
func (c *Client) ParseUserInfo(userID string) ([]byte, error) {
	p := authing.UserQueryParameter{
		ID:               graphql.String(userID),
		RegisterInClient: graphql.String(c.ClientID),
	}

	q, err := c.AuthingClient.User(&p)
	if err != nil {
		return []byte{}, err
	}
	user := q.User
	if user.Blocked {
		return []byte{}, fmt.Errorf("This user is blocked")
	}
	if user.IsDeleted {
		return []byte{}, fmt.Errorf("This user has been deleted")
	}
	return json.Marshal(user)
}
