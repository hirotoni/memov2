package repos

import (
	"time"

	"github.com/hirotoni/memov2/models"
)

// MockRepo is a mock implementation of the fileRepo interface for testing
type MockRepo struct {
	MockTodosTemplate       func(date time.Time) (models.TodoFileInterface, error)
	MockFindTodosFileByDate func(date time.Time) (models.TodoFileInterface, error)
	MockSave                func(file models.MemoFileInterface, truncate bool) error
	MockMemoEntries         func() ([]models.MemoFileInterface, error)
	MockTodoEntries         func() ([]models.TodoFileInterface, error)
	MockMetadata            func(file models.MemoFileInterface) (map[string]interface{}, error)
}

// TodosTemplate calls the mock implementation
func (m *MockRepo) TodosTemplate(date time.Time) (models.TodoFileInterface, error) {
	if m.MockTodosTemplate != nil {
		return m.MockTodosTemplate(date)
	}
	return nil, nil
}

// FindTodosFileByDate calls the mock implementation
func (m *MockRepo) FindTodosFileByDate(date time.Time) (models.TodoFileInterface, error) {
	if m.MockFindTodosFileByDate != nil {
		return m.MockFindTodosFileByDate(date)
	}
	return nil, nil
}

// Save calls the mock implementation
func (m *MockRepo) Save(file models.MemoFileInterface, truncate bool) error {
	if m.MockSave != nil {
		return m.MockSave(file, truncate)
	}
	return nil
}

func (m *MockRepo) TodoEntries() ([]models.TodoFileInterface, error) {
	if m.MockTodoEntries != nil {
		return m.MockTodoEntries()
	}
	return nil, nil
}

// MemoEntires calls the mock implementation
func (m *MockRepo) MemoEntries() ([]models.MemoFileInterface, error) {
	if m.MockMemoEntries != nil {
		return m.MockMemoEntries()
	}
	return nil, nil
}

// Metadata calls the mock implementation
func (m *MockRepo) Metadata(file models.MemoFileInterface) (map[string]interface{}, error) {
	if m.MockMetadata != nil {
		return m.MockMetadata(file)
	}
	return nil, nil
}

// NewMockFileRepo creates a new instance of MockFileRepo
func NewMockFileRepo() *MockRepo {
	return &MockRepo{}
}
