package external

import (
	"context"
	"coresamples/common"
	pb "coresamples/proto"
	cgrpc "github.com/asim/go-micro/plugins/client/grpc/v4"
	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"google.golang.org/grpc"
	"strconv"
)

type CustomerService struct {
	endpoint    string
	serviceName string
	service     micro.Service
}

var customerSvc *CustomerService

func InitCustomerService(consulToken string, consulURL string) {
	if customerSvc == nil {
		customerSvc = &CustomerService{
			endpoint:    "192.168.60.6:30276",
			serviceName: "lis-core-grpc",
		}
		if common.Env.RunEnv == common.AksProductionEnv || common.Env.RunEnv == common.AksStagingEnv {
			var consulRegistry registry.Registry
			consulRegistry = consul.NewRegistry(
				consul.Config(&api.Config{
					Address: consulURL,
					Token:   consulToken,
				}))

			service := micro.NewService(
				micro.Client(cgrpc.NewClient()),
				micro.Registry(consulRegistry),
			)
			service.Init()
			customerSvc.service = service
		}
	}
}

func GetCustomerService() *CustomerService {
	return customerSvc
}

func (cs *CustomerService) GetCustomerClinics(customerID int64) ([]int64, error) {
	if common.Env.RunEnv == common.AksProductionEnv || common.Env.RunEnv == common.AksStagingEnv {
		return cs.getCustomerClinicsConsul(strconv.FormatInt(customerID, 10))
	}
	return cs.getCustomerClinicsDial(strconv.FormatInt(customerID, 10))
}

func (cs *CustomerService) getCustomerClinicsConsul(customerIDs ...string) ([]int64, error) {
	cli := pb.NewCustomerService(cs.serviceName, cs.service.Client())
	resp, err := cli.ListCustomerAllClinics(context.Background(), &pb.ListCustomerAllClinicsRequest{
		CustomerIds: customerIDs,
	})
	if err != nil {
		common.Error(err)
		return nil, err
	}
	var ids []int64
	for _, clinics := range resp.CustomerClinics {
		for _, clinic := range clinics.Clinics {
			ids = append(ids, int64(clinic.ClinicId))
		}
	}
	return ids, nil
}

func (cs *CustomerService) getCustomerClinicsDial(customerIDs ...string) ([]int64, error) {
	trackingId := uuid.NewString()
	conn, err := grpc.DialContext(context.Background(), cs.endpoint, grpc.WithInsecure(), grpc.WithUnaryInterceptor(common.UnaryClientInterceptor(trackingId)))
	if err != nil {
		common.Error(err)
		return nil, err
	}
	defer conn.Close()
	cli := pb.NewCustomerServiceClient(conn)
	resp, err := cli.ListCustomerAllClinics(context.Background(), &pb.ListCustomerAllClinicsRequest{
		CustomerIds: customerIDs,
	})
	if err != nil {
		common.Error(err)
		return nil, err
	}
	var ids []int64
	for _, clinics := range resp.CustomerClinics {
		for _, clinic := range clinics.Clinics {
			ids = append(ids, int64(clinic.ClinicId))
		}
	}
	return ids, nil
}
