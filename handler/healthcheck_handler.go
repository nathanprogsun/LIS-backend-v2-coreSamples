package handler

import (
	"context"
	"coresamples/common"
	pb "coresamples/proto"
	"log"
	"net/http"
)

type HealthcheckHandler struct{}

func (h *HealthcheckHandler) Check(ctx context.Context, request *pb.HealthCheckRequest, response *pb.HealthCheckResponse) error {
	switch request.GetService() {
	case "readiness", "liveness":
		if common.Env.ServiceDeregistered {
			response.Status = pb.HealthCheckResponse_NOT_SERVING
		} else {
			response.Status = pb.HealthCheckResponse_SERVING
		}
	default:
		response.Status = pb.HealthCheckResponse_UNKNOWN
	}
	//logger.Info(request.GetService() + " " + response.Status.String())
	return nil
}

func HTTPCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

func HandleHTTPHealthCheck() {
	http.HandleFunc("/healthcheck", HTTPCheck)
	go func() {
		if err := http.ListenAndServe(":8083", nil); err != nil {
			log.Fatal(err)
		}
	}()
}
