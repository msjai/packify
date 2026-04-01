package service

// CalculateService handles optimal pack calculation for orders.
type CalculateService struct {
	name      string
	packStore PackStore
}

// NewCalculateService creates a new CalculateService with the given storage.
func NewCalculateService(packStore PackStore) *CalculateService {
	return &CalculateService{
		name:      "calculate service",
		packStore: packStore,
	}
}