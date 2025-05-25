package repos

import (
	"time"

	"github.com/hirotoni/memov2/models"
)

// MockFileRepo is a mock implementation of the fileRepo interface for testing
type MockFileRepo struct {
	MockTodosTemplate       func(date time.Time) (models.FileInterface, error)
	MockFindTodosFileByDate func(date time.Time) (models.FileInterface, error)
	MockSave                func(file models.FileInterface, truncate bool) error
	MockMemoEntries         func() ([]models.MemoFileInterface, error)
	MockMetadata            func(file models.MemoFileInterface) (map[string]interface{}, error)
}

// Ensure MockFileRepo implements fileRepo interface
var _ fileRepo = (*MockFileRepo)(nil)

// TodosTemplate calls the mock implementation
func (m *MockFileRepo) TodosTemplate(date time.Time) (models.FileInterface, error) {
	if m.MockTodosTemplate != nil {
		return m.MockTodosTemplate(date)
	}
	return nil, nil
}

// FindTodosFileByDate calls the mock implementation
func (m *MockFileRepo) FindTodosFileByDate(date time.Time) (models.FileInterface, error) {
	if m.MockFindTodosFileByDate != nil {
		return m.MockFindTodosFileByDate(date)
	}
	return nil, nil
}

// Save calls the mock implementation
func (m *MockFileRepo) Save(file models.FileInterface, truncate bool) error {
	if m.MockSave != nil {
		return m.MockSave(file, truncate)
	}
	return nil
}

// MemoEntires calls the mock implementation
func (m *MockFileRepo) MemoEntries() ([]models.MemoFileInterface, error) {
	if m.MockMemoEntries != nil {
		return m.MockMemoEntries()
	}
	return nil, nil
}

// Metadata calls the mock implementation
func (m *MockFileRepo) Metadata(file models.MemoFileInterface) (map[string]interface{}, error) {
	if m.MockMetadata != nil {
		return m.MockMetadata(file)
	}
	return nil, nil
}

// NewMockFileRepo creates a new instance of MockFileRepo
func NewMockFileRepo() *MockFileRepo {
	return &MockFileRepo{}
}
