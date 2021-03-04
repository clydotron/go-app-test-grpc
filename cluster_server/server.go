package main

import (
	"fmt"
	"log"
	"net"

	"github.com/clydotron/go-app-test/api/clusterpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type clusterServer struct{}

func (cs *clusterServer) HealthCheck(req *clusterpb.HealthCheckRequest, stream clusterpb.ClusterService_HealthCheckServer) error {

	fmt.Println("HealthCheck:", req)

	//@todo this includes a timeout/duration, implement a timer and pass

	// send info for each control plane
	cplanes := req.ClusterInfo.GetControlPlaneNodes()
	for _, cplane := range cplanes {

		resp := &clusterpb.HealthCheckProgress{
			Metadata: &clusterpb.Metadata{
				Hostname: cplane,
				Error:    "none",
			},
			Message: "active",
		}
		stream.Send(resp)
	}

	// send info for each worker node
	//wnodes := req.ClusterInfo.GetWorkerNodes()
	for _, node := range req.ClusterInfo.GetWorkerNodes() {
		resp := &clusterpb.HealthCheckProgress{
			Metadata: &clusterpb.Metadata{
				Hostname: node,
				Error:    "none",
			},
			Message: "active",
		}
		stream.Send(resp)
	}

	return nil
}

func (cs *clusterServer) GreetMany(req *clusterpb.GreetManyRequest, srv clusterpb.ClusterService_GreetManyServer) error {
	return status.Errorf(codes.Unimplemented, "method GreetMany not implemented")
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	s := grpc.NewServer()

	clusterpb.RegisterClusterServiceServer(s, &clusterServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
