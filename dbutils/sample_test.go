package dbutils

import (
	"context"
	"coresamples/ent/enttest"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"coresamples/ent/sampleidgenerate"
	"coresamples/ent/orderinfo"
	
)

func TestGetdailyCollectionSamples(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	ctx := context.Background()

	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	tomorrow := now.Add(24 * time.Hour)

	// Create test patient - first create without assigning the result
	client.Patient.Create().
		ExecX(ctx)

	// Then get the created patient in a separate step
	patient := client.Patient.Query().
		FirstX(ctx)

	// Create test samples
	client.Sample.Create().
		SetPatient(patient).
		SetAccessionID("TEST001").
		SetSampleReceivedTime(now).
		SetDelayedHours(0).
		ExecX(ctx)

	client.Sample.Create().
		SetPatient(patient).
		SetAccessionID("TEST002").
		SetSampleReceivedTime(now.Add(-1 * time.Hour)).
		SetDelayedHours(0).
		ExecX(ctx)

	client.Sample.Create().
		SetPatient(patient).
		SetAccessionID("TEST003").
		SetSampleReceivedTime(yesterday).
		SetDelayedHours(0).
		ExecX(ctx)

	tests := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
		want      int // expected number of samples
	}{
		{
			name:      "Should find samples within time range",
			startTime: yesterday,
			endTime:   tomorrow,
			want:      3,
		},
		{
			name:      "Should find no samples outside time range",
			startTime: tomorrow,
			endTime:   tomorrow.Add(24 * time.Hour),
			want:      0,
		},
		{
			name:      "Should find samples in shorter time range",
			startTime: now.Add(-2 * time.Hour),
			endTime:   now.Add(1 * time.Hour),
			want:      2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			samples, err := GetdailyCollectionSamples(tt.startTime, tt.endTime, client, ctx)

			assert.NoError(t, err)

			// Assert correct number of samples returned
			assert.Equal(t, tt.want, len(samples))

			// Verify samples are within the time range
			for _, s := range samples {
				fullSample := client.Sample.GetX(ctx, s.ID)
				assert.True(t, fullSample.SampleReceivedTime.After(tt.startTime) || fullSample.SampleReceivedTime.Equal(tt.startTime))
				assert.True(t, fullSample.SampleReceivedTime.Before(tt.endTime))
			}
		})
	}
}

func TestGetdailyCollectionSamples_EdgeCases(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	ctx := context.Background()

	now := time.Now()
	client.Patient.Create().
		ExecX(ctx)

	// Then get the created patient in a separate step
	patient := client.Patient.Query().
		FirstX(ctx)

	// Create first sample
	client.Sample.Create().
		SetPatient(patient).
		SetAccessionID("EDGE-001").
		SetSampleReceivedTime(now).
		SetDelayedHours(0).
		ExecX(ctx)

	// Create second sample
	client.Sample.Create().
		SetPatient(patient).
		SetAccessionID("EDGE-002").
		SetSampleReceivedTime(now.Add(1 * time.Hour)).
		SetDelayedHours(0).
		ExecX(ctx)

	tests := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
		want      int
	}{
		{
			name:      "Exact time range match",
			startTime: now,
			endTime:   now.Add(1 * time.Hour),
			want:      1, // Should only include the start time sample
		},
		{
			name:      "Zero duration time range",
			startTime: now,
			endTime:   now,
			want:      0,
		},
		{
			name:      "Reversed time range",
			startTime: now.Add(1 * time.Hour),
			endTime:   now,
			want:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			samples, err := GetdailyCollectionSamples(tt.startTime, tt.endTime, client, ctx)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, len(samples))
		})
	}

	t.Run("Context cancelled", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		samples, err := GetdailyCollectionSamples(now, now.Add(1*time.Hour), client, cancelledCtx)
		assert.Error(t, err)
		assert.Nil(t, samples)
	})

	t.Run("Nil client", func(t *testing.T) {
		samples, err := GetdailyCollectionSamples(now, now.Add(1*time.Hour), nil, ctx)
		assert.Error(t, err)
		assert.Nil(t, samples)
	})
}

func TestGetdailyCheckNonReceivedSamples(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	ctx := context.Background()

	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	tomorrow := now.Add(24 * time.Hour)

	// Create test patient - first create without assigning the result
	client.Patient.Create().
		ExecX(ctx)

	// Then get the created patient in a separate step
	patient := client.Patient.Query().
		FirstX(ctx)

	// Create test order
	client.OrderInfo.Create().
		SetOrderCreateTime(now).
		SetOrderConfirmationNumber("TEST-ORDER-123").
		ExecX(ctx)

	// Get first order
	order1 := client.OrderInfo.Query().
		Where(orderinfo.OrderConfirmationNumberEQ("TEST-ORDER-123")).
		OnlyX(ctx)

	client.OrderInfo.Create().
		SetOrderCreateTime(now).
		SetOrderConfirmationNumber("TEST-ORDER-456").
		ExecX(ctx)

	// Get second order
	order2 := client.OrderInfo.Query().
		Where(orderinfo.OrderConfirmationNumberEQ("TEST-ORDER-456")).
		OnlyX(ctx)

	// Create test samples
	client.Sample.Create().
		SetPatient(patient).
		SetOrder(order1).
		SetAccessionID("TEST001").
		SetDelayedHours(0).
		ExecX(ctx)

	client.Sample.Create().
		SetPatient(patient).
		SetOrder(order2).
		SetAccessionID("TEST002").
		SetSampleReceivedTime(now).
		SetDelayedHours(0).
		ExecX(ctx)

	tests := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
		want      int // expected number of samples
	}{
		{
			name:      "Should find non-received samples within time range",
			startTime: yesterday,
			endTime:   tomorrow,
			want:      1,
		},
		{
			name:      "Should find no samples outside time range",
			startTime: tomorrow,
			endTime:   tomorrow.Add(24 * time.Hour),
			want:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			samples, err := GetdailyCheckNonReceivedSamples(tt.startTime, tt.endTime, client, ctx)

			assert.NoError(t, err)

			assert.Equal(t, tt.want, len(samples))

			// If samples were found, verify they are non-received
			for _, s := range samples {
				fullSample := client.Sample.GetX(ctx, s.ID)
				assert.True(t, fullSample.SampleReceivedTime.IsZero(), "expected sample received time to be zero")
				assert.NotEmpty(t, fullSample.AccessionID)
			}
		})
	}
}

func TestGetdailyCheckNonReceivedSamples_EdgeCases(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	ctx := context.Background()

	now := time.Now()
	// Create test patient - first create without assigning the result
	client.Patient.Create().
		ExecX(ctx)

	// Then get the created patient in a separate step
	patient := client.Patient.Query().
		FirstX(ctx)

	// Create orders with edge case timestamps
	client.OrderInfo.Create().
		SetOrderCreateTime(now).
		SetOrderConfirmationNumber("EXACT-START").
		ExecX(ctx)

	exactStartOrder := client.OrderInfo.Query().
		Where(orderinfo.OrderConfirmationNumberEQ("EXACT-START")).
		OnlyX(ctx)

	// Create samples for edge case testing
	client.Sample.Create().
		SetPatient(patient).
		SetOrder(exactStartOrder).
		SetAccessionID("EDGE-001").
		SetDelayedHours(0).
		ExecX(ctx)

	tests := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
		want      int
	}{
		{
			name:      "Exact time range match",
			startTime: now,
			endTime:   now.Add(1 * time.Hour),
			want:      1, // Should only include the start time sample
		},
		{
			name:      "Zero duration time range",
			startTime: now,
			endTime:   now,
			want:      0,
		},
		{
			name:      "Reversed time range",
			startTime: now.Add(1 * time.Hour),
			endTime:   now,
			want:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			samples, err := GetdailyCheckNonReceivedSamples(tt.startTime, tt.endTime, client, ctx)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, len(samples))
		})
	}
}

func TestGenerateSampleID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	ctx := context.Background()

	t.Run("Basic functionality", func(t *testing.T) {
		sampleID, err := GenerateSampleID(client, ctx)
		assert.NoError(t, err)
		assert.NotNil(t, sampleID)
		assert.Greater(t, sampleID.ID, 0)
	})

	t.Run("Multiple generations are unique", func(t *testing.T) {
		// Generate multiple IDs and ensure they're unique
		idMap := make(map[int]bool)
		for i := 0; i < 5; i++ {
			sampleID, err := GenerateSampleID(client, ctx)
			assert.NoError(t, err)
			assert.NotNil(t, sampleID)
			
			// Verify ID is unique
			assert.False(t, idMap[sampleID.ID], "Generated ID should be unique")
			idMap[sampleID.ID] = true
		}
	})

	t.Run("Cancelled context", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()
		
		_, err := GenerateSampleID(client, cancelledCtx)
		assert.Error(t, err)
	})

	t.Run("Nil client", func(t *testing.T) {
		_, err := GenerateSampleID(nil, ctx)
		assert.Error(t, err)
	})
}

func TestGetBarcodeForSampleID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	ctx := context.Background()

	// Setup test data - first create the entity without assigning the result
	client.SampleIDGenerate.Create().
		SetBarcode("TEST123456").
		ExecX(ctx)

	// Then get the ID in a separate step
	sampleIDGen := client.SampleIDGenerate.Query().
		Where(sampleidgenerate.BarcodeEQ("TEST123456")).
		OnlyX(ctx)

	tests := []struct {
		name      string
		sampleID  int
		wantCode  string
		wantError bool
	}{
		{
			name:      "Valid sample ID",
			sampleID:  sampleIDGen.ID,
			wantCode:  "TEST123456",
			wantError: false,
		},
		{
			name:      "Non-existent sample ID",
			sampleID:  99999,
			wantCode:  "",
			wantError: true,
		},
		{
			name:      "Invalid sample ID (zero)",
			sampleID:  0,
			wantCode:  "",
			wantError: true,
		},
		{
			name:      "Invalid sample ID (negative)",
			sampleID:  -1,
			wantCode:  "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			barcode, err := GetBarcodeForSampleID(tt.sampleID, client, ctx)

			if tt.wantError {
				assert.Error(t, err)
				assert.Empty(t, barcode)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantCode, barcode)
		})
	}

	t.Run("Context cancelled", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		barcode, err := GetBarcodeForSampleID(sampleIDGen.ID, client, cancelledCtx)
		assert.Error(t, err)
		assert.Empty(t, barcode)
	})

	t.Run("Nil client", func(t *testing.T) {
		barcode, err := GetBarcodeForSampleID(sampleIDGen.ID, nil, ctx)
		assert.Error(t, err)
		assert.Empty(t, barcode)
	})
}

