package main

import (
	"context"
	"flag"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	gateway "github.com/rauljordan/minimal-grpc-gateway"
	pb "github.com/rauljordan/minimal-grpc-gateway/example/proto/api/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	allowedOriginsFlag = flag.String(
		"gateway-allowed-origins",
		"*",
		"comma-separated, allowed origins for gateway cors",
	)
	grpcServerAddressFlag = flag.String(
		"grpc-server-address",
		"127.0.0.1:4000",
		"host:port address for the grpc server",
	)
	grpcGatewayAddressFlag = flag.String(
		"grpc-gateway-address",
		"127.0.0.1:8080",
		"host:port address for the grpc-JSON gateway server",
	)
)

// Example gRPC server implementation.
type server struct{}

// SignupUser implements the protobuf service defined by our API.
func (s *server) SignupUser(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	return &pb.SignupResponse{
		JwtKey: []byte("hello-world"),
	}, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	grpcGatewayAddress := *grpcGatewayAddressFlag
	grpcServerAddress := *grpcServerAddressFlag
	lis, err := net.Listen("tcp", grpcGatewayAddress)
	if err != nil {
		logrus.Errorf("Could not listen to port in Start() %s: %v", grpcServerAddress, err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAPIServer(grpcServer, &server{})

	go func() {
		logrus.WithField("address", grpcGatewayAddress).Info("gRPC server listening on port")
		if err := grpcServer.Serve(lis); err != nil {
			logrus.Errorf("Could not serve gRPC: %v", err)
		}
	}()

	allowedOrigins := []string{*allowedOriginsFlag}
	if strings.Contains(*allowedOriginsFlag, ",") {
		allowedOrigins = strings.Split(*allowedOriginsFlag, ",")
	}
	gatewaySrv := gateway.New(ctx, &gateway.Config{
		GatewayAddress:      grpcGatewayAddress,
		RemoteAddress:       grpcServerAddress,
		Mux:                 nil, /*optional http mux */
		AllowedOrigins:      allowedOrigins,
		EndpointsToRegister: []gateway.RegistrationFunc{pb.RegisterAPIHandlerFromEndpoint},
	})
	gatewaySrv.Start()

	// Listen for any process interrupts.
	stop := make(chan struct{})
	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigc)
		<-sigc
		logrus.Info("Got interrupt, shutting down...")
		grpcServer.GracefulStop()
		stop <- struct{}{}
	}()

	// Wait for stop channel to be closed.
	<-stop
}
