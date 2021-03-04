package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/clydotron/go-app-test/api/clusterpb"
	"github.com/clydotron/go-app-test/client"
)

func main() {

	cc := client.NewClusterClient()
	defer cc.Close()
	// cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalln("Failed to dial:", err)
	// }
	// defer cc.Close()

	// c := clusterpb.NewClusterServiceClient(cc)

	doHealthCheck(cc.CSC)
}

func doHealthCheck(csc clusterpb.ClusterServiceClient) {

	controlPlanes := []string{"control plane 1", "control plane 2"}
	workerNodes := []string{"node 1", "node 2", "node 3", "node 4", "node 5"}

	req := &clusterpb.HealthCheckRequest{
		ClusterInfo: &clusterpb.ClusterInfo{
			ControlPlaneNodes: controlPlanes,
			WorkerNodes:       workerNodes,
		},
		//WaitTimeout: ,
	}

	stream, err := csc.HealthCheck(context.Background(), req)
	if err != nil {
		log.Fatalln("HealthCheck RPC error:", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// reached end of stream
			break
		}

		if err != nil {
			log.Fatalln("stream Recv error:", err)
		}
		// handle the actual result:
		d := msg.GetMetadata()
		x := msg.GetMessage()
		fmt.Println("msg:", d, x)
	}
}
