package repo

import (
	"api_orion/model"

	"gorm.io/gorm"
)

type MemberRepository interface {
	CreateMember(member *model.NewMember) error
	GetAllMember(limit int, page int, query string, batchID int) ([]model.NewMember, error)
	GetMemberByID(id int) (*model.NewMember, error)
	GetMemberByNim(nim string) (*model.NewMember, error)
	Update(id int, member *model.NewMember) error
	UpdateStatus(id int, status string) error
	GetRegistrationTrend(batchID int) ([]model.RegistrationTrend, error)
	Delete(id int) error
}

type MemberRepo struct {
	db *gorm.DB
}

func NewMemberRepo(db *gorm.DB) *MemberRepo {
	return &MemberRepo{db: db}
}

func (r *MemberRepo) CreateMember(member *model.NewMember) error {
	return r.db.Create(member).Error
}

func (r *MemberRepo) GetAllMember(limit int, page int, query string, batchID int) ([]model.NewMember, error) {
	var members []model.NewMember
	offset := (page - 1) * limit

	db := r.db

	if batchID > 0 {
		db = db.Where("batch_id = ?", batchID)
	}

	if query != "" {
		db = db.Where("full_name ILIKE ?", "%"+query+"%")
	}

	err := db.Limit(limit).Offset(offset).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (r *MemberRepo) GetMemberByID(id int) (*model.NewMember, error) {
	var member model.NewMember
	err := r.db.First(&member, id).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *MemberRepo) GetMemberByNim(nim string) (*model.NewMember, error) {
	var member model.NewMember
	err := r.db.Where("nim = ?", nim).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *MemberRepo) Update(id int, member *model.NewMember) error {
	return r.db.Where("id = ?", id).Updates(member).Error
}

func (r *MemberRepo) UpdateStatus(id int, status string) error {
	return r.db.Model(&model.NewMember{}).Where("id = ?", id).Update("status", status).Error
}

func (r *MemberRepo) Delete(id int) error {
	return r.db.Delete(&model.NewMember{}, id).Error
}

func (r *MemberRepo) GetRegistrationTrend(batchID int) ([]model.RegistrationTrend, error) {
	var trends []model.RegistrationTrend

	query := `
		SELECT 
			TO_CHAR(created_at, 'Dy') as day,
			COALESCE(SUM(CASE WHEN devision = 'PROGRAMMING' THEN 1 ELSE 0 END), 0) as programming,
			COALESCE(SUM(CASE WHEN devision = 'ELECTRONICS' THEN 1 ELSE 0 END), 0) as electronic,
			COALESCE(SUM(CASE WHEN devision = 'MECHANICAL' THEN 1 ELSE 0 END), 0) as mechanic
		FROM new_members
		WHERE batch_id = ?
		GROUP BY TO_CHAR(created_at, 'Dy'), EXTRACT(ISODOW FROM created_at)
		ORDER BY EXTRACT(ISODOW FROM created_at)
	`

	err := r.db.Raw(query, batchID).Scan(&trends).Error
	if err != nil {
		return nil, err
	}
	return trends, nil
}
