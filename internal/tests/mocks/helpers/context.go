package helpers

import (
	"bytes"
	"encoding/json"
	"net/http"

	userpb "github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/gin-gonic/gin"
)

func CreateMockContext(requestData map[string]interface{}, url string, currentUserRoleName string) *gin.Context {
	data, _ := json.Marshal(requestData)
	buffer := bytes.NewBuffer(data)
	request, _ := http.NewRequest("POST", url, buffer)
	request.Header.Add("Content-Type", "application/json")

	c := &gin.Context{Request: request}

	if len(currentUserRoleName) > 0 {
		c.Set("_user", &userpb.User{
			FirstName: "Fname",
			LastName:  "Lname",
			Email:     "usr@example.com",
			RoleName:  currentUserRoleName,
			GroupId:   1,
			UID:       "c8e1a5b7-7457-4fc4-af10-2aeafc9bf9f9sdvsadv", // random
		})
	}

	return c
}
