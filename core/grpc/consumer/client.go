package consumer

import (
	"context"
	"dss/common/log"
	"dss/core/config"
	"dss/core/discover"
	pb "dss/core/grpc/proto"
	"dss/core/host"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"time"
)

func Startup(ctx context.Context) {
	var (
		err         error
		conn        *grpc.ClientConn
		client      pb.StreamService_RecordClient
		serviceName = config.CoreConf.ServiceName
	)
	dis := discover.NewServiceDiscovery(ctx, serviceName)
	go func() {
		for {
			select {
			case sig := <-dis.Wait():
				if sig != discover.Ok {
					log.WarnF("receive a ng signal,quit...")
					return
				}
				conn, err = getClientConn(dis.Get())
				if err != nil {
					log.Fatal(err.Error())
				}
				client, err = pb.NewStreamServiceClient(conn).Record(ctx)
				if err != nil {
					log.Fatal(err.Error())
				}

				req := pb.StreamRequest{Pt: &pb.StreamPoint{
					Host:            host.LocalIP(),
					Name:            host.Name.Load().(string),
					Platform:        host.Platform.Load().(string),
					PlatformVersion: host.PlatformVersion.Load().(string),
				}}
				if err = client.Send(&req); err != nil {
					log.Errorf(err.Error())
				}
				func() {
					defer func(conn *grpc.ClientConn) {
						if err = conn.Close(); err != nil {
							log.Errorf(err.Error())
						}
					}(conn)
					defer func(client pb.StreamService_RecordClient) {
						if err = client.CloseSend(); err != nil {
							log.Errorf(err.Error())
						}
					}(client)
				}()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func getClientConn(srv string) (*grpc.ClientConn, error) {
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
	gRpcServerAddress, _, _ := net.SplitHostPort(srv)
	conn, err = grpc.DialContext(context.Background(), fmt.Sprintf("%s:%d", gRpcServerAddress, config.CoreConf.GrpcPort), options...)
	if err != nil {
		return nil, fmt.Errorf("grpc.dial err: %s", err.Error())
	}
	return conn, nil
}
