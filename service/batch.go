package service

import (
	"api_orion/model"
	"api_orion/repo"
)

type BatchService interface {
	CreateBatch(batch *model.Batch) error
	GetAllBatch() ([]model.Batch, error)
	GetBatchByID(id int) (*model.Batch, error)
	GetActiveBatch() (*model.Batch, error)
	Update(id int, batch *model.Batch) error
	UpdateActiveStatus(id int, isActive bool) error
	Delete(id int) error
}

type BatchServ struct {
	batchRepo repo.BatchRepository
}

func NewBatchService(batchRepo repo.BatchRepository) *BatchServ {
	return &BatchServ{batchRepo: batchRepo}
}

func (s *BatchServ) CreateBatch(batch *model.Batch) error {
	return s.batchRepo.CreateBatch(batch)
}

func (s *BatchServ) GetAllBatch() ([]model.Batch, error) {
	return s.batchRepo.GetAllBatch()
}

func (s *BatchServ) GetBatchByID(id int) (*model.Batch, error) {
	return s.batchRepo.GetBatchByID(id)
}

func (s *BatchServ) GetActiveBatch() (*model.Batch, error) {
	return s.batchRepo.GetActiveBatch()
}

func (s *BatchServ) Update(id int, batch *model.Batch) error {
	return s.batchRepo.Update(id, batch)
}

func (s *BatchServ) UpdateActiveStatus(id int, isActive bool) error {
	return s.batchRepo.UpdateActiveStatus(id, isActive)
}

func (s *BatchServ) Delete(id int) error {
	return s.batchRepo.Delete(id)
}
