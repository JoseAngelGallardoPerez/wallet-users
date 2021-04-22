package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Confialink/wallet-users/internal/config"
	pb "github.com/Confialink/wallet-users/rpc/proto/users"
)

var requestUID = "c8e1a5b7-7457-4fc4-af10-2aeafc9bf9f9"

func main() {
	// Retrieve config options.
	conf := config.GetConf()

	addr := fmt.Sprintf(":%s", conf.RPC.GetUsersServerPort())

	client := pb.NewUserHandlerProtobufClient(addr, &http.Client{})

	var (
		res *pb.Response
		err error
	)

	res, err = client.GetByUID(context.Background(), &pb.Request{UID: requestUID})
	if err != nil {
		fmt.Printf("oh no: %v", err)
		os.Exit(1)
	}

	fmt.Printf("%+v", res)
}
