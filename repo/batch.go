package repo

import (
	"api_orion/model"

	"gorm.io/gorm"
)

type BatchRepository interface {
	CreateBatch(batch *model.Batch) error
	GetAllBatch() ([]model.Batch, error)
	GetActiveBatch() (*model.Batch, error)
	GetBatchByID(id int) (*model.Batch, error)
	Update(id int, batch *model.Batch) error
	UpdateActiveStatus(id int, isActive bool) error
	Delete(id int) error
}

type BatchRepo struct {
	db *gorm.DB
}

func NewBatchRepo(db *gorm.DB) *BatchRepo {
	return &BatchRepo{db: db}
}

func (r *BatchRepo) CreateBatch(batch *model.Batch) error {
	return r.db.Create(batch).Error
}

func (r *BatchRepo) GetAllBatch() ([]model.Batch, error) {
	var batches []model.Batch
	err := r.db.Find(&batches).Error
	if err != nil {
		return nil, err
	}
	return batches, nil
}

func (r *BatchRepo) GetBatchByID(id int) (*model.Batch, error) {
	var batch model.Batch
	err := r.db.First(&batch, id).Error
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *BatchRepo) GetActiveBatch() (*model.Batch, error) {
	var batch model.Batch
	err := r.db.Where("is_active = ?", true).First(&batch).Error
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *BatchRepo) Update(id int, batch *model.Batch) error {
	return r.db.Where("id = ?", id).Updates(batch).Error
}

func (r *BatchRepo) UpdateActiveStatus(id int, isActive bool) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if isActive {
			if err := tx.Model(&model.Batch{}).Where("id != ?", id).Update("is_active", false).Error; err != nil {
				return err
			}
		}
		return tx.Model(&model.Batch{}).Where("id = ?", id).Update("is_active", isActive).Error
	})
}

func (r *BatchRepo) Delete(id int) error {
	return r.db.Where("id = ?", id).Delete(&model.Batch{}).Error
}
