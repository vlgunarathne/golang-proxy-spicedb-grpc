package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	spicedb_pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	pb "github.com/vlgunarathne/golang-proxy-spicedb/pkg/spicedb"
)

type server struct {
	pb.UnimplementedProxyServiceSpiceDBServer
	spiceDbClient *authzed.Client
}

func (s *server) SayHelloProxy(ctx context.Context, request *pb.HelloProxyRequest) (*pb.HelloProxyReply, error) {
	req := &spicedb_pb.CheckPermissionRequest{
		Permission: "allowed",
		Resource: &spicedb_pb.ObjectReference{
			ObjectType: "permissions",
			ObjectId:   request.Resource.ObjectId,
		},
		Subject: &spicedb_pb.SubjectReference{
			Object: &spicedb_pb.ObjectReference{
				ObjectType: "projectgroups",
				ObjectId:   request.Subject.Object.ObjectId,
			},
		},
	}

	resp, err := s.spiceDbClient.CheckPermission(ctx, req)
	if err != nil {
		log.Printf("failed to check permission: %s", err)
	}

	response := &pb.HelloProxyReply{
		Permissionship: "Proxy " + resp.GetPermissionship().String(),
	}
	return response, nil
}

func main() {

	serverAddr := ":9090"

	client, err := authzed.NewClient(
		"spicedb:50051",
		grpcutil.WithInsecureBearerToken("foobar"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// systemCerts,
	)

	if err != nil {
		log.Fatalf("unable to initialize client: %s", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterProxyServiceSpiceDBServer(grpcServer, &server{spiceDbClient: client})

	lis, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Server listening on %s", serverAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
