package domain

import (
	"testing"
	"time"
)

func TestNewTodosFile(t *testing.T) {
	t.Run("valid date", func(t *testing.T) {
		date := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

		todoFile, err := NewTodosFile(date)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if todoFile == nil {
			t.Error("expected TodosFile, got nil")
		}

		// Check internal fields
		tf := todoFile.(*TodoFile)
		if !tf.date.Equal(date) {
			t.Errorf("expected date %v, got %v", date, tf.date)
		}

		expectedTitle := date.Format(FileNameDateLayoutTodo)
		if tf.title != expectedTitle {
			t.Errorf("expected title %s, got %s", expectedTitle, tf.title)
		}

		if tf.fileType != FileTypeTodos {
			t.Errorf("expected fileType %v, got %v", FileTypeTodos, tf.fileType)
		}
	})

	t.Run("zero date", func(t *testing.T) {
		var zeroDate time.Time

		todoFile, err := NewTodosFile(zeroDate)

		if err == nil {
			t.Error("expected error for zero date, got nil")
		}

		if todoFile != nil {
			t.Error("expected nil TodosFile for zero date, got non-nil")
		}

		expectedError := "invalid date"
		if err.Error() != expectedError {
			t.Errorf("expected error message '%s', got '%s'", expectedError, err.Error())
		}
	})
}

func TestNewTodoTemplateFile(t *testing.T) {
	todoTemplate, err := NewTodoTemplateFile()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if todoTemplate == nil {
		t.Error("expected TodosFile, got nil")
	}

	// Check internal fields
	tf := todoTemplate.(*TodoFile)

	expectedTitle := "todos_template"
	if tf.title != expectedTitle {
		t.Errorf("expected title %s, got %s", expectedTitle, tf.title)
	}

	if tf.fileType != FileTypeTemplate {
		t.Errorf("expected fileType %v, got %v", FileTypeTemplate, tf.fileType)
	}

	// Check that date is set (should be recent)
	if tf.date.IsZero() {
		t.Error("expected date to be set, got zero date")
	}

	// Note: HeadingBlocks are set via SetHeadingBlocks but there's no getter method
	// We can only verify that the SetHeadingBlocks call doesn't panic
	// The actual verification would need access to the internal state or a getter method
}

func TestTodosFile_FileName(t *testing.T) {
	t.Run("regular todos file", func(t *testing.T) {
		date := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC) // Monday
		todoFile, _ := NewTodosFile(date)

		filename := todoFile.FileName()

		// Expected format: 20240115Mon + Sep + FileTypeTodos.String() + Ext
		expectedDateString := "20240115Mon"
		expectedFilename := expectedDateString + FileSeparator + FileTypeTodos.String() + FileExtension

		if filename != expectedFilename {
			t.Errorf("expected filename %s, got %s", expectedFilename, filename)
		}
	})

	t.Run("template file", func(t *testing.T) {
		todoTemplate, _ := NewTodoTemplateFile()

		filename := todoTemplate.FileName()

		expectedFilename := "todos_template" + FileExtension

		if filename != expectedFilename {
			t.Errorf("expected filename %s, got %s", expectedFilename, filename)
		}
	})

	t.Run("different weekdays", func(t *testing.T) {
		testCases := []struct {
			date     time.Time
			expected string
		}{
			{time.Date(2024, 1, 14, 0, 0, 0, 0, time.UTC), "20240114Sun"}, // Sunday
			{time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), "20240115Mon"}, // Monday
			{time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC), "20240116Tue"}, // Tuesday
			{time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC), "20240117Wed"}, // Wednesday
			{time.Date(2024, 1, 18, 0, 0, 0, 0, time.UTC), "20240118Thu"}, // Thursday
			{time.Date(2024, 1, 19, 0, 0, 0, 0, time.UTC), "20240119Fri"}, // Friday
			{time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC), "20240120Sat"}, // Saturday
		}

		for _, tc := range testCases {
			todoFile, _ := NewTodosFile(tc.date)
			filename := todoFile.FileName()
			expected := tc.expected + FileSeparator + FileTypeTodos.String() + FileExtension

			if filename != expected {
				t.Errorf("for date %v, expected filename %s, got %s", tc.date, expected, filename)
			}
		}
	})
}

func TestFileNameDateLayoutTodo(t *testing.T) {
	// Test that the date layout constant works as expected
	date := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC) // Monday
	formatted := date.Format(FileNameDateLayoutTodo)
	expected := "20240115Mon"

	if formatted != expected {
		t.Errorf("expected formatted date %s, got %s", expected, formatted)
	}
}

func TestTodoFileInterface(t *testing.T) {
	var _ TodoFileInterface = &TodoFile{}
	var _ FileInterface = &TodoFile{}
}

func TestTodoFile_FileName(t *testing.T) {
	tests := []struct {
		name string
		f    *TodoFile
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.FileName(); got != tt.want {
				t.Errorf("TodoFile.FileName() = %v, want %v", got, tt.want)
			}
		})
	}
}
