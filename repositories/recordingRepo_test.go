package repositories

import (
	"context"
	"errors"
	"field_archive/server/entities"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

// MOCK ROW IMPLEMENATTION --------------------
type MockRow struct {
	mockScan func(dest ...any) error
}

func (r *MockRow) Scan(dest ...any) error {
	return r.mockScan(dest)
}

// --------------------------------------------

// MOCK ROWS IMPLEMENTATION -------------------
type MockRows struct {
	mockNext   func() bool
	mockScan   func(dest ...any) error
	mockErr    func() error
	commandTag pgconn.CommandTag
	closed     bool
}

func (m *MockRows) Conn() *pgx.Conn {
	panic("unimplemented")
}

func (m *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	panic("unimplemented")
}

func (m *MockRows) RawValues() [][]byte {
	panic("unimplemented")
}

func (m *MockRows) Values() ([]any, error) {
	panic("unimplemented")
}

func (m *MockRows) Next() bool {
	if m.closed {
		return false
	}
	return m.mockNext()
}

func (m *MockRows) Scan(dest ...any) error {
	if m.closed {
		return errors.New("rows are closed")
	}
	return m.mockScan(dest...)
}

func (m *MockRows) Err() error {
	return m.mockErr()
}

func (m *MockRows) Close() {
	m.closed = true
}

func (m *MockRows) CommandTag() pgconn.CommandTag {
	return m.commandTag
}

// --------------------------------------------

// MOCK DATABASE IMPLEMENTATION ---------------
type MockDatabase struct {
	mockExec     func(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	mockQueryRow func(ctx context.Context, query string, args ...any) pgx.Row
	mockQuery    func(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
}

func (m *MockDatabase) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return m.mockExec(ctx, query, args...)
}
func (m *MockDatabase) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return m.mockQueryRow(ctx, query, args...)
}

func (m *MockDatabase) Query(ctx context.Context, query string, args any) (pgx.Rows, error) {
	return m.mockQuery(ctx, query, args)
}

// --------------------------------------------

// REPOSITORY METHOD TESTS START HERE

func TestInsert(t *testing.T) {
	check := `INSERT INTO recordings` +
		`(title, audio_location, artwork_location, date_uploaded, recording_date, location_id, user_id, ` +
		`duration, format, description, equipment, file_size, channels, license) ` +
		`VALUES ` +
		`(@title, @audio_location, @date_uploaded, @recording_date, @location_id, @user_id, @duration, @format, @description, @equipment, @file_Size, @channels, @license) ` +
		`RETURNING id`

	mockDB := &MockDatabase{
		mockQueryRow: func(ctx context.Context, query string, args ...any) pgx.Row {
			if check == query {
				return &MockRow{mockScan: func(dest ...any) error {
					innerSlice, ok := dest[0].([]any)
					if !ok {
						return errors.New("Can't access inner slice")
					}
					*(innerSlice[0].(*int)) = 3
					return nil
				}}
			}
			return &MockRow{mockScan: func(dest ...any) error {
				return errors.New("row not found")
			},
			}
		},
	}
	repo := &RecordingRepoImplement{conn: mockDB}
	recording := entities.Recording{
		Title:         "Test Title",
		AudioLocation: "test/audio/location.mp3",
		DateUploaded:  nil,
		RecordingDate: time.Now(),
		LocationID:    1,
		UserID:        1,
		Duration:      120,
		Format:        "mp3",
		Description:   "Test description",
		Equipment:     "Test equipment",
		Size:          2048,
		Channels:      "2",
		License:       "Creative Commons",
	}
	id, err := repo.Insert(recording, context.Background())
	assert.NoError(t, err)
	if id != 3 {
		t.Fatalf("INSERT ERROR: returned ids not equal")
	}
}

func TestGetRowByID(t *testing.T) {
	check := `SELECT * FROM recordings WHERE id=@id`
	expectedRecording := entities.Recording{
		ID:              1,
		Title:           "Test Title",
		AudioLocation:   "test/audio/location.mp3",
		ArtworkLocation: nil,
		DateUploaded:    nil,
		RecordingDate:   time.Date(2025, 1, 6, 20, 2, 57, 0, time.UTC),
		LocationID:      1,
		UserID:          1,
		Duration:        120,
		Format:          "mp3",
		Description:     "test description",
		Equipment:       "test equipment",
		Size:            2048,
		Channels:        "2",
		License:         "Creative Commons",
	}
	mockDB := &MockDatabase{
		mockQueryRow: func(ctx context.Context, query string, args ...any) pgx.Row {
			if check == query {
				return &MockRow{mockScan: func(dest ...any) error {
					innerSlice, ok := dest[0].([]any)
					if !ok {
						return errors.New("Can't access inner slice")
					}
					*(innerSlice[0].(*int)) = 1
					*(innerSlice[1].(*string)) = "Test Title"
					*(innerSlice[2].(*string)) = "test/audio/location.mp3"
					*(innerSlice[3].(**string)) = nil
					*(innerSlice[4].(**time.Time)) = nil
					*(innerSlice[5].(*time.Time)) = time.Date(2025, 1, 6, 20, 2, 57, 0, time.UTC)
					*(innerSlice[6].(*int)) = 1
					*(innerSlice[7].(*int)) = 1
					*(innerSlice[8].(*int)) = 120
					*(innerSlice[9].(*string)) = "mp3"
					*(innerSlice[10].(*string)) = "test description"
					*(innerSlice[11].(*string)) = "test equipment"
					*(innerSlice[12].(*float64)) = 2048
					*(innerSlice[13].(*string)) = "2"
					*(innerSlice[14].(*string)) = "Creative Commons"
					return nil
				}}
			}
			return &MockRow{mockScan: func(dest ...any) error {
				return errors.New("row not found")
			}}
		},
	}
	repo := &RecordingRepoImplement{conn: mockDB}
	test, err := repo.GetRowByID(1, context.Background())
	if test != expectedRecording {
		t.Errorf("GetRowByID failed: return = %v, but expected %v", test, expectedRecording)
	}
	if err != nil {
		t.Errorf("getRowByID failed: Error returned from func %v", err)
	}
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	check := `UPDATE recordings ` +
		`SET title = @title, audio_location = @audio_location, ` +
		`artwork_location = @artwork_location, date_uploaded = @date_uploaded, ` +
		`recording_date = @recording_date, location_id = @location_id, user_id = @user_id, duration = @duration, ` +
		`format = @format, description = @description, equipment = @equipment, file_size = @file_size, ` +
		`channels = @channels, license = @license ` +
		`WHERE id = @id`
	mockDB := MockDatabase{mockExec: func(ctx context.Context,
		query string, args ...any) (pgconn.CommandTag, error) {
		if check == query {
			return pgconn.CommandTag{}, nil
		} else {
			return pgconn.CommandTag{}, errors.New("Query mismatch")
		}
	}}
	repo := &RecordingRepoImplement{conn: &mockDB}

	recording := entities.Recording{
		ID:              1,
		Title:           "Test Title",
		AudioLocation:   "test/audio/location.mp3",
		ArtworkLocation: nil,
		DateUploaded:    nil,
		RecordingDate:   time.Date(2025, 1, 6, 20, 2, 57, 0, time.UTC),
		LocationID:      1,
		UserID:          1,
		Duration:        120,
		Format:          "mp3",
		Description:     "test description",
		Equipment:       "test equipment",
		Size:            2048,
		Channels:        "2",
		License:         "Creative Commons",
	}

	err := repo.Update(recording, ctx)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	check := `DELETE FROM recordings WHERE id = @id`
	mockDB := MockDatabase{mockExec: func(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
		if check == query {
			return pgconn.CommandTag{}, nil
		} else {
			return pgconn.CommandTag{}, errors.New("DELETE ERROR: Query did not match check")
		}
	}}
	repo := &RecordingRepoImplement{conn: &mockDB}
	err := repo.Delete(1, ctx)
	assert.NoError(t, err)
}

func TestList(t *testing.T) {
	ctx := context.Background()
	check := `SELECT * FROM recordings LIMIT $1::int`

	expectedRecordings := []entities.Recording{
		{
			ID:              1,
			Title:           "Recording 1",
			AudioLocation:   "audio/location1.mp3",
			ArtworkLocation: nil,
			DateUploaded:    nil,
			RecordingDate:   time.Date(2025, 1, 6, 20, 0, 0, 0, time.UTC),
			LocationID:      1,
			UserID:          1,
			Duration:        180,
			Format:          "mp3",
			Description:     "Description 1",
			Equipment:       "Equipment 1",
			Size:            2048,
			Channels:        "2",
			License:         "Creative Commons",
		},
		{
			ID:              2,
			Title:           "Recording 2",
			AudioLocation:   "audio/location2.mp3",
			ArtworkLocation: nil,
			DateUploaded:    nil,
			RecordingDate:   time.Date(2025, 1, 7, 20, 0, 0, 0, time.UTC),
			LocationID:      2,
			UserID:          1,
			Duration:        200,
			Format:          "wav",
			Description:     "Description 2",
			Equipment:       "Equipment 2",
			Size:            4096,
			Channels:        "2",
			License:         "Creative Commons",
		},
	}

	mockDB := MockDatabase{mockQuery: func(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
		if check == query {
			return &MockRows{
				mockNext: func() bool {
					return len(expectedRecordings) > 0
				},
				mockScan: func(dest ...any) error {
					if len(expectedRecordings) == 0 {
						return errors.New("No recordings left!")
					}
					recording := expectedRecordings[0]
					*(dest[0].(*int)) = recording.ID
					*(dest[1].(*string)) = recording.Title
					*(dest[2].(*string)) = recording.AudioLocation
					*(dest[3].(**string)) = recording.ArtworkLocation
					*(dest[4].(**time.Time)) = recording.DateUploaded
					*(dest[5].(*time.Time)) = recording.RecordingDate
					*(dest[6].(*int)) = recording.LocationID
					*(dest[7].(*int)) = recording.UserID
					*(dest[8].(*int)) = recording.Duration
					*(dest[9].(*string)) = recording.Format
					*(dest[10].(*string)) = recording.Description
					*(dest[11].(*string)) = recording.Equipment
					*(dest[12].(*float64)) = recording.Size
					*(dest[13].(*string)) = recording.Channels
					*(dest[14].(*string)) = recording.License
					expectedRecordings = expectedRecordings[1:]
					return nil
				},
			}, nil
		}
		fmt.Println(query)
		fmt.Println(check)
		return nil, errors.New("Query did not match Check")
	}}
	repo := &RecordingRepoImplement{conn: &mockDB}

	res, err := repo.List(ctx, 2)

	if err != nil {
		t.Errorf("Error listing recordings: %v", err)
	}

	if slices.Equal(res, expectedRecordings) {
		t.Errorf("Error listing recordings!\nreceived: %v\nBut expected: %v", res, expectedRecordings)
	}

}

func TestCount(t *testing.T) {
	ctx := context.Background()
	check := `SELECT COUNT(id) FROM recordings`
	expectedReturn := 1

	mockDB := MockDatabase{mockQueryRow: func(ctx context.Context, query string, args ...any) pgx.Row {
		if check == query {
			return &MockRow{mockScan: func(dest ...any) error {
				innerSlice, ok := dest[0].([]any)
				if !ok {
					return errors.New("Can't access inner slice")
				}
				*(innerSlice[0].(*int)) = 1
				return nil
			},
			}
		}
		return &MockRow{mockScan: func(dest ...any) error {
			return errors.New("row not found")
		},
		}
	},
	}
	repo := &RecordingRepoImplement{conn: &mockDB}
	id, err := repo.Count(ctx)
	assert.Equal(t, id, expectedReturn)
	if err != nil {
		t.Errorf("Error testing Count method %v", err)
	}

}
