package booking

type Service interface {
	CreateDraft(req CreateDraftRequest) (string, error)
	GetDraft(id string) (*BookingDraft, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateDraft(req CreateDraftRequest) (string, error) {
	// 這裡可以加入業務邏輯驗證
	return s.repo.Save(req.StudioID, req.DraftData)
}

func (s *service) GetDraft(id string) (*BookingDraft, error) {
	return s.repo.FindByID(id)
}
