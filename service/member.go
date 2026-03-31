package service

import (
	"api_orion/model"
	"api_orion/repo"
)

type MemberService interface {
	CreateMember(member *model.NewMember) error
	GetAllMember(limit int, page int, query string) ([]model.NewMember, error)
	GetMemberByID(id int) (*model.NewMember, error)
	GetMemberByNim(nim string) (*model.NewMember, error)
	Update(id int, member *model.NewMember) error
	UpdateStatus(id int, status string) error
	GetRegistrationTrend(batchID int) ([]model.RegistrationTrend, error)
	Delete(id int) error
}

type memberService struct {
	memberRepo repo.MemberRepository
}

func NewMemberService(memberRepo repo.MemberRepository) MemberService {
	return &memberService{memberRepo}
}

func (s *memberService) CreateMember(member *model.NewMember) error {
	return s.memberRepo.CreateMember(member)
}

func (s *memberService) GetAllMember(limit int, page int, query string) ([]model.NewMember, error) {
	return s.memberRepo.GetAllMember(limit, page, query)
}

func (s *memberService) GetMemberByID(id int) (*model.NewMember, error) {
	return s.memberRepo.GetMemberByID(id)
}

func (s *memberService) GetMemberByNim(nim string) (*model.NewMember, error) {
	return s.memberRepo.GetMemberByNim(nim)
}

func (s *memberService) Update(id int, member *model.NewMember) error {
	return s.memberRepo.Update(id, member)
}

func (s *memberService) UpdateStatus(id int, status string) error {
	return s.memberRepo.UpdateStatus(id, status)
}

func (s *memberService) Delete(id int) error {
	return s.memberRepo.Delete(id)
}

func (s *memberService) GetRegistrationTrend(batchID int) ([]model.RegistrationTrend, error) {
	return s.memberRepo.GetRegistrationTrend(batchID)
}
