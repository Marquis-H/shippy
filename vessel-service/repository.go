package main

import (
	pb "shippy/vessel-service/proto/vessel"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DB_NAME        = "shippy"
	CON_COLLECTION = "vessels"
)

type Repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
	Create(*pb.Vessel) error
	Close()
}

type VesselRepository struct {
	session *mgo.Session
}

func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	// 选择最近一条容量、载重都符合的货轮
	var vessel *pb.Vessel
	err := repo.collection().Find(
		bson.M{
			"capacity":  bson.M{"$gte": spec.Capacity},
			"maxweight": bson.M{"$gte": spec.MaxWeight},
		},
	).One(&vessel)
	if err != nil {
		return nil, err
	}

	return vessel, nil
}

func (repo *VesselRepository) collection() *mgo.Collection {
	return repo.session.DB(DB_NAME).C(CON_COLLECTION)
}

func (repo *VesselRepository) Create(v *pb.Vessel) error {
	return repo.collection().Insert(v)
}

// 关闭连接
func (repo *VesselRepository) Close() {
	repo.session.Close()
}
