package main

import (
	"context"
	pb "shippy/vessel-service/proto/vessel"

	"gopkg.in/mgo.v2"
)

type handler struct {
	session *mgo.Session
}

func (h *handler) GetRepo() Repository {
	return &VesselRepository{h.session.Clone()}
}

// 实现微服务
func (h *handler) FindAvailable(ctx context.Context, spec *pb.Specification, resp *pb.Response) error {
	defer h.GetRepo().Close()
	// 调用内部方法查找
	v, err := h.GetRepo().FindAvailable(spec)
	if err != nil {
		return err
	}
	resp.Vessel = v
	return nil
}

// 实现微服务的服务端
func (h *handler) Create(ctx context.Context, req *pb.Vessel, resp *pb.Response) error {
	defer h.GetRepo().Close()
	if err := h.GetRepo().Create(req); err != nil {
		return err
	}
	resp.Vessel = req
	resp.Created = true
	return nil
}
