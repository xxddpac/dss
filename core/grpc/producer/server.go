package producer

import (
	"dss/common/log"
	"dss/core/config"
	"dss/core/dao"
	"dss/core/global"
	pb "dss/core/grpc/proto"
	"dss/core/models"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

const (
	defaultMinPingTime    = 5 * time.Second
	defaultMaxConnIdle    = 20 * time.Minute
	defaultPingTime       = 10 * time.Minute
	defaultPingAckTimeout = 5 * time.Second
	maxMsgSize            = 1024 * 1024 * 10
)

var (
	gRPC *grpc.Server
)

type StreamService struct{}

func Grpc() {
	NewGrpcServer()
	addr := fmt.Sprintf(":%d", config.CoreConf.GrpcPort)
	lis, err := net.Listen(global.TCP, addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.InfoF("start gRPC server, listen:%s", addr)
	go func() {
		if err = gRPC.Serve(lis); err != nil {
			log.Fatal(err.Error())
		}
	}()
}

func NewGrpcServer() {
	knfP := keepalive.EnforcementPolicy{
		MinTime:             defaultMinPingTime,
		PermitWithoutStream: true,
	}
	sp := keepalive.ServerParameters{
		MaxConnectionIdle: defaultMaxConnIdle,
		Time:              defaultPingTime,
		Timeout:           defaultPingAckTimeout,
	}
	options := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(knfP),
		grpc.KeepaliveParams(sp),
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	}
	server := grpc.NewServer(options...)
	reflection.Register(server)
	pb.RegisterStreamServiceServer(server, &StreamService{})
	gRPC = server
}

func (s *StreamService) Record(stream pb.StreamService_RecordServer) error {
	response, err := stream.Recv()
	if err != nil {
		return err
	}
	if err = dao.Repo(global.GrpcClient).RemoveAll(bson.M{"host": response.Pt.Host}); err != nil {
		log.Errorf("remove gRPC client %v record err:%v", response.Pt.Host, err)
		return err
	}
	if err = dao.Repo(global.GrpcClient).Insert(
		models.ClientInsertFunc(
			models.Client{
				Host:            response.Pt.Host,
				Name:            response.Pt.Name,
				Platform:        response.Pt.Platform,
				PlatformVersion: response.Pt.PlatformVersion,
			})); err != nil {
		log.Errorf("insert gRPC client %v record err:", response.Pt.Host, err)
	}
	return nil
}

func CloseGrpc() {
	go func() {
		gRPC.Stop()
	}()
}
