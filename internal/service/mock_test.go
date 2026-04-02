package service

// mockPackStore implements PackStore and PackSubmitter for testing.
type mockPackStore struct {
	sizes []int
	err   error
}

func (m *mockPackStore) GetPackSizes() ([]int, error) {
	return m.sizes, m.err
}

func (m *mockPackStore) UpdatePackSizes(sizes []int) error {
	if m.err != nil {
		return m.err
	}
	m.sizes = sizes
	return nil
}
