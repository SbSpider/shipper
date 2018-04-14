package main

import (
    "fmt"
    "log"

    // Import the generated protobuf code
    pb "github.com/sbspider/shipper/consignment-service/proto/consignment"
    vesselProto "github.com/sbspider/shipper/vessel-service/proto/vessel"
    micro "github.com/micro/go-micro"
    "golang.org/x/net/context"
)

const (
    port = ":50051"
)

type IRepository interface {
    Create(*pb.Consignment) (*pb.Consignment, error)
    GetAll() []*pb.Consignment
}

// Repository - Dummy repository, this simulates the use of a datastore
// of some kind. WE'll replace this with a real implementation later on
type Repository struct {
    consignments []*pb.Consignment
}

func (repo * Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
    updated := append(repo.consignments, consignment)
    repo.consignments = updated
    return consignment, nil
}

func (repo * Repository) GetAll() []*pb.Consignment {
    return repo.consignments 
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. you can check the interface
// in the generated code itself of rthe exact method signatures etc
// to give you a better idea
type service struct {
    repo IRepository
    vesselClient vesselProto.VesselServiceClient
}

// CreateConsignment - we created just one method on our service,
// which is  acreate method, which takes a context and a request as an 
// argument, these are handled by the gRpc server
func (s * service) CreateConsignment(ctx context.Context, req *pb.Consignment, res * pb.Response) error {

    // here we call a client instance of our vessel service wih our consignment weight,
    // and the amount of containers as the capacity value
    vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
        MaxWeight: req.Weight,
        Capacity: int32(len(req.Containers)),
    })
    log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
    if err != nil {
        return err
    }

    req.VesselId = vesselResponse.Vessel.Id

    // Save our consignment
    consignment, err := s.repo.Create(req)
    if err != nil {
        return err
    }

    // Return matching the 'Response` message we created in our
    // protobuf definition
    res.Created = true
    res.Consignment = consignment
    return nil
}

func (s * service) GetConsignments(ctx context.Context, req *pb.GetRequest, res * pb.Response) error {
    consignments := s.repo.GetAll()
    res.Consignments = consignments
    return nil
}

func main() {

    repo := &Repository{}

    // Create a new service. Optionally, include some options here
    srv := micro.NewService(
        // This name must match the package name given in your protobuf definition
        micro.Name("go.micro.srv.consignment"),
        micro.Version("latest"),
    )

    vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

    // Init will parse the command line flags.
    srv.Init()

    // Register handler
    pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})

    // Run the server
    if err := srv.Run(); err != nil {
        fmt.Println(err)
    }
}
