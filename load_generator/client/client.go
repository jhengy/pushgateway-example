package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	pb "github.com/grpc-ecosystem/go-grpc-prometheus/examples/grpc-server-with-prometheus/protobuf"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// instrument grpc_client with default grpc-client-metrics
	// Create a metrics registry.
	reg := prometheus.NewRegistry()
	// Create some standard client metrics.
	grpcMetrics := grpc_prometheus.NewClientMetrics()
	// Register client metrics to registry.
	reg.MustRegister(grpcMetrics)
	// Create a insecure gRPC channel to communicate with the server.
	conn, err := grpc.Dial(
		fmt.Sprintf("grpcserver:%v", 9093),
		grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(grpcMetrics.StreamClientInterceptor()),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	// Create a HTTP server for prometheus.
	// to export collected metrics to prometheus through a http server
	httpServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: fmt.Sprintf("0.0.0.0:%d", 9094)}
	// Start your http server for prometheus.
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		} else {
			log.Println("successfully started the http server to handle prometheus receiver")
		}
	}()

	// Create a gRPC server client.
	client := pb.NewDemoServiceClient(conn)
	fmt.Println("Start to call the method called SayHello every 3 seconds")
	for {
		// Call “SayHello” method and wait for response from gRPC Server.
		_, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "Test"})
		if err != nil {
			log.Printf("Calling the SayHello method unsuccessfully. ErrorInfo: %+v", err)
			log.Printf("You should to stop the process")
			return
		} else {
			log.Println("client.SayHello success, sleep for 3s and then repeat")
		}
		time.Sleep(3 * time.Second)
	}

	// commented out as when run in docker container the process will exit automatically for some reason
	//scanner := bufio.NewScanner(os.Stdin)
	//fmt.Println("You can press n or N to stop the process of client")
	//for scanner.Scan() {
	//	if strings.ToLower(scanner.Text()) == "n" {
	//		os.Exit(0)
	//	}
	//}
}
