package main

import (
	"context"
	"log"
	pb "shippy/consignment-service/proto/consignment"
	vesselPb "shippy/vessel-service/proto/vessel"

	"gopkg.in/mgo.v2"
)

type handler struct {
	session      *mgo.Session
	vesselClient vesselPb.VesselServiceClient
}

func (h *handler) GetRepo() Repository {
	return &ConsignmentRepository{h.session.Clone()}
}

func (h *handler) CreateConsignment(ctx context.Context, req *pb.Consignment, resp *pb.Response) error {
	defer h.GetRepo().Close()

	// 检查是否有适合的货轮
	vReq := &vesselPb.Specification{
		Capacity:  int32(len(req.Containers)),
		MaxWeight: req.Weight,
	}
	vResp, err := h.vesselClient.FindAvailable(context.Background(), vReq)
	if err != nil {
		return err
	}

	// 货物被承运
	log.Printf("found vessel: %s\n", vResp.Vessel.Name)
	req.VesselId = vResp.Vessel.Id
	err = h.GetRepo().Create(req)
	if err != nil {
		return err
	}
	resp.Created = true
	resp.Consignment = req
	return nil
}

func (h *handler) GetConsignments(ctx context.Context, req *pb.GetRequest, resp *pb.Response) error {
	defer h.GetRepo().Close()
	consignments, err := h.GetRepo().GetAll()
	if err != nil {
		return err
	}
	resp.Consignments = consignments
	return nil
}
