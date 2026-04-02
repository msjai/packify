package service

// mockPackStore implements PackStore and PackSubmitter for testing.
type mockPackStore struct {
	sizes     []int
	getErr    error
	updateErr error
}

func (m *mockPackStore) GetPackSizes() ([]int, error) {
	return m.sizes, m.getErr
}

func (m *mockPackStore) UpdatePackSizes(sizes []int) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.sizes = sizes
	return nil
}
