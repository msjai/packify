package service

// PackStore defines the storage interface for retrieving pack sizes.
// Follows interface segregation — minimal interface, each service depends only on what it needs.
type PackStore interface {
	GetPackSizes() ([]int, error)
}

// GetService handles retrieving current pack sizes from storage.
type GetService struct {
	name      string
	packStore PackStore
}

// NewGetService creates a new GetService with the given storage.
func NewGetService(packStore PackStore) *GetService {
	return &GetService{
		name:      "get service",
		packStore: packStore,
	}
}

// GetPackSizes returns available pack sizes from storage.
func (s *GetService) GetPackSizes() ([]int, error) {
	return s.packStore.GetPackSizes()
}