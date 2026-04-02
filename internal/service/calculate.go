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

// Calculate returns optimal pack combination for a given order.
// TODO: replace mock with DP algorithm.
func (s *CalculateService) Calculate(order int) (map[int]int, error) {
	return map[int]int{250: 2, 500: 1}, nil
}