package main

import (
	pb "shippy/consignment-service/proto/consignment"

	"gopkg.in/mgo.v2"
)

const (
	DB_NAME        = "shippy"
	CON_COLLECTION = "consignments"
)

type Repository interface {
	Create(*pb.Consignment) error
	GetAll() ([]*pb.Consignment, error)
	Close()
}

type ConsignmentRepository struct {
	session *mgo.Session
}

// 接口实现
func (repo *ConsignmentRepository) Create(c *pb.Consignment) error {
	return repo.collection().Insert(c)
}

// 获取全部数据
func (repo *ConsignmentRepository) GetAll() ([]*pb.Consignment, error) {
	var cons []*pb.Consignment
	err := repo.collection().Find(nil).All(&cons)
	return cons, err
}

// 关闭连接
func (repo *ConsignmentRepository) Close() {
	repo.session.Close()
}

// 返回所有货物信息
func (repo *ConsignmentRepository) collection() *mgo.Collection {
	return repo.session.DB(DB_NAME).C(CON_COLLECTION)
}
