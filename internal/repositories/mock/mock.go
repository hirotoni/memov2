package mock

import (
	"time"

	"github.com/hirotoni/memov2/internal/domain"
)

// Mock is a mock implementation of the fileRepo interface for testing
type Mock struct {
	MockTodosTemplate       func(date time.Time) (domain.TodoFileInterface, error)
	MockFindTodosFileByDate func(date time.Time) (domain.TodoFileInterface, error)
	MockSave                func(file domain.MemoFileInterface, truncate bool) error
	MockMemoEntries         func() ([]domain.MemoFileInterface, error)
	MockTodoEntries         func() ([]domain.TodoFileInterface, error)
	MockMetadata            func(file domain.MemoFileInterface) (map[string]interface{}, error)
}

// TodosTemplate calls the mock implementation
func (m *Mock) TodosTemplate(date time.Time) (domain.TodoFileInterface, error) {
	if m.MockTodosTemplate != nil {
		return m.MockTodosTemplate(date)
	}
	return nil, nil
}

// FindTodosFileByDate calls the mock implementation
func (m *Mock) FindTodosFileByDate(date time.Time) (domain.TodoFileInterface, error) {
	if m.MockFindTodosFileByDate != nil {
		return m.MockFindTodosFileByDate(date)
	}
	return nil, nil
}

// Save calls the mock implementation
func (m *Mock) Save(file domain.MemoFileInterface, truncate bool) error {
	if m.MockSave != nil {
		return m.MockSave(file, truncate)
	}
	return nil
}

func (m *Mock) TodoEntries() ([]domain.TodoFileInterface, error) {
	if m.MockTodoEntries != nil {
		return m.MockTodoEntries()
	}
	return nil, nil
}

// MemoEntires calls the mock implementation
func (m *Mock) MemoEntries() ([]domain.MemoFileInterface, error) {
	if m.MockMemoEntries != nil {
		return m.MockMemoEntries()
	}
	return nil, nil
}

// Metadata calls the mock implementation
func (m *Mock) Metadata(file domain.MemoFileInterface) (map[string]interface{}, error) {
	if m.MockMetadata != nil {
		return m.MockMetadata(file)
	}
	return nil, nil
}

// NewMock creates a new instance of MockFileRepo
func NewMock() *Mock {
	return &Mock{}
}

// MockEditor is a mock implementation of the Editor interface for testing
type MockEditor struct {
	OpenFunc func(basedir, path string) error
	Calls    []EditorCall // Track all calls made to the editor
}

// EditorCall tracks a call to the Open method
type EditorCall struct {
	Basedir string
	Path    string
}

// Open implements the Editor interface
func (m *MockEditor) Open(basedir, path string) error {
	// Record the call
	m.Calls = append(m.Calls, EditorCall{
		Basedir: basedir,
		Path:    path,
	})

	// Call the custom function if provided
	if m.OpenFunc != nil {
		return m.OpenFunc(basedir, path)
	}

	// Default behavior: do nothing (successful no-op)
	return nil
}

// NewMockEditor creates a new mock editor instance
func NewMockEditor() *MockEditor {
	return &MockEditor{
		Calls: []EditorCall{},
	}
}
