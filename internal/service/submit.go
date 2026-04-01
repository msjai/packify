package service

// PackSubmitter defines the minimal storage interface for updating pack sizes.
// Kept to a single method, follows Go's interface segregation principle.
type PackSubmitter interface {
	UpdatePackSizes(sizes []int) error
}

// SubmitService handles updating pack sizes in storage.
type SubmitService struct {
	name          string
	packSubmitter PackSubmitter
}

// NewSubmitService creates a new SubmitService with the given storage.
func NewSubmitService(packSubmitter PackSubmitter) *SubmitService {
	return &SubmitService{
		name:          "submit service",
		packSubmitter: packSubmitter,
	}
}

// UpdatePackSizes replaces all pack sizes in storage.
func (s *SubmitService) UpdatePackSizes(sizes []int) error {
	return s.packSubmitter.UpdatePackSizes(sizes)
}