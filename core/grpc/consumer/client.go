package consumer

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"dss/common/log"
	"dss/common/utils"
	"dss/core/config"
	"dss/core/discover"
	pb "dss/core/grpc/proto"
	"dss/core/host"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"net"
	"os"
	"path/filepath"
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
		root = filepath.Join(utils.WorkingDirectory(), "common/cert")
	)
	cred, err := credential(filepath.Join(root, "client.pem"), filepath.Join(root, "client.key"), filepath.Join(root, "ca.pem"))
	if err != nil {
		return nil, fmt.Errorf("get credential error:%s", err.Error())
	}
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(cred),
		grpc.WithBlock(),
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

func credential(crtFile, keyFile, caFile string) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("load X509 error:%s", err.Error())
	}
	caBytes, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("read ca file error:%s", err.Error())
	}
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caBytes); !ok {
		log.Errorf("append cert error")
		return nil, fmt.Errorf("append cert error")
	}
	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   config.CoreConf.ServiceName,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		RootCAs:      certPool,
	}), nil
}
