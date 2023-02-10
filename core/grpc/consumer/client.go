package consumer

import (
	"context"
	"dss/common/log"
	"dss/core/config"
	"dss/core/global"
	pb "dss/core/grpc/proto"
	"dss/core/host"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"strings"
	"time"
)

func Startup() {
	var (
		err    error
		conn   *grpc.ClientConn
		client pb.StreamService_RecordClient
	)
	conn, err = getClientConn()
	if err != nil {
		log.Fatal(err.Error())
	}
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			client, err = pb.NewStreamServiceClient(conn).Record(context.Background())
			if err != nil {
				log.Errorf(err.Error())
			}
			req := pb.StreamRequest{Pt: &pb.StreamPoint{
				Host:            strings.Join(host.PrivateIPv4.Load().([]string), ","),
				Name:            host.Name.Load().(string),
				Platform:        host.Platform.Load().(string),
				PlatformVersion: host.PlatformVersion.Load().(string),
			}}
			if err = client.Send(&req); err != nil {
				log.Errorf(err.Error())
			}
		case <-global.Ctx.Done():
			if _, err = client.CloseAndRecv(); err != nil {
				log.Errorf(err.Error())
			}
			return
		}
	}
}

func getClientConn() (*grpc.ClientConn, error) {
	var (
		err  error
		conn *grpc.ClientConn
	)
	options := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithReturnConnectionError(),
		grpc.FailOnNonTempDialError(true),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(16 * 1024 * 1024)),
		grpc.WithConnectParams(grpc.ConnectParams{Backoff: backoff.Config{MaxDelay: time.Second * 2}, MinConnectTimeout: time.Second * 2}),
	}
	conn, err = grpc.DialContext(context.Background(), fmt.Sprintf(":%d", config.CoreConf.Consumer.GrpcPort), options...)
	if err != nil {
		return nil, fmt.Errorf("grpc.dial err: %s", err.Error())
	}
	return conn, nil
}
