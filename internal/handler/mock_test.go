package handler

// mockSubmitter implements Submitter for testing.
type mockSubmitter struct {
	sizes []int
	err   error
}

func (m *mockSubmitter) UpdatePackSizes(sizes []int) error {
	if m.err != nil {
		return m.err
	}
	m.sizes = sizes
	return nil
}

// mockCalculator implements Calculator for testing.
type mockCalculator struct {
	gotOrder int
	result   map[int]int
	err      error
}

func (m *mockCalculator) Calculate(order int) (map[int]int, error) {
	m.gotOrder = order
	return m.result, m.err
}

// mockGetter implements Getter for testing.
type mockGetter struct {
	sizes []int
	err   error
}

func (m *mockGetter) GetPackSizes() ([]int, error) {
	return m.sizes, m.err
}
