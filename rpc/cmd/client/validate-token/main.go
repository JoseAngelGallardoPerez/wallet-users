package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Confialink/wallet-users/internal/config"
	pb "github.com/Confialink/wallet-users/rpc/proto/users"
)

var accessToken = "eyJraWQiOiI4TUpteXliTTR5bEJDUjg0ajlldmticzVia0J4V1wvNlBrUkdKREtmaStSUT0iLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiI2MTBlY2NiYy0xNmM0LTQyMmUtOGJhYi1kMzliZDZkODlmNWQiLCJldmVudF9pZCI6ImU1OWVlYjI0LTliNTYtMTFlOC1hMGYyLTVmZTdkY2IzNWVlOSIsInRva2VuX3VzZSI6ImFjY2VzcyIsInNjb3BlIjoiYXdzLmNvZ25pdG8uc2lnbmluLnVzZXIuYWRtaW4iLCJhdXRoX3RpbWUiOjE1MzM3NjU4MTUsImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC51cy1lYXN0LTEuYW1hem9uYXdzLmNvbVwvdXMtZWFzdC0xX1h5dVhXakdtRCIsImV4cCI6MTUzMzc2OTQxNSwiaWF0IjoxNTMzNzY1ODE1LCJqdGkiOiI0YjlkODFjMS05NDY4LTRlMjEtODg5ZS03ODc4OGI1OGEzOGEiLCJjbGllbnRfaWQiOiI2aWg5ZGVlaXFjcGY2cmN1bGU2c250M2o1NCIsInVzZXJuYW1lIjoiNjEwZWNjYmMtMTZjNC00MjJlLThiYWItZDM5YmQ2ZDg5ZjVkIn0.TTHbyMIL07dkEpI7lOnnFJXOF54669CJLkYoPH6y1oIdGXqO6ckHnbXSl8Alxy5EM0JIp19BtHV7fu7PCR1kvdgTKytPFpElx-RZ0qv_LHZHWnV3AGS_DJ4iSTCHdzYl4akhCMd0Bw8n87V9YpZNM-wXrHrOWciYLTW2eWRR2Z15nTe73OwF7tJvRuMTt0J1w3tEbrxzAUq7zvbA-k0Q8f42zM_WjXVvgcSMmY2V8_n9td-vQGM5IpXWUda7BcK2DpCQlzne2604PwiI6UJPcWYNKNkqoNxtmVak_RXRTQe65kJRrFS36r40EXAILT4lGFPEtX6fhyLIIG1XMkQznw"

func main() {
	// Retrieve config options.
	conf := config.GetConf()

	addr := fmt.Sprintf(":%s", conf.RPC.GetUsersServerPort())

	client := pb.NewUserHandlerProtobufClient(addr, &http.Client{})

	var (
		res *pb.Response
		err error
	)

	res, err = client.ValidateAccessToken(context.Background(), &pb.Request{AccessToken: accessToken})
	if err != nil {
		fmt.Printf("oh no: %v", err)
		os.Exit(1)
	}

	fmt.Printf("%+v", res)
}
