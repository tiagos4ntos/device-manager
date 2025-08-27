package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagos4ntos/device-manager/internal/domain/device/entity"
)

func makeExpectedDeviceRecord() *entity.Device {
	return &entity.Device{
		ID:           uuid.MustParse("b44ecc02-872e-4c18-8d2a-ac09dfc4b49a"),
		Name:         "Galaxy S23 FE",
		Brand:        "Samsumg",
		Status:       "active",
		CreationTime: lo.Must(time.Parse(time.DateTime, "2025-08-31 15:01:02")),
	}
}

func Test_Create_Device(t *testing.T) {
	assert := assert.New(t)

	deviceCreateQuery := regexp.QuoteMeta(`
	INSERT INTO devices (id, name, brand, status)
	VALUES ($1, $2, $3, $4)
	RETURNING id, creation_time;`)

	expectedDevice := makeExpectedDeviceRecord()

	testCases := []struct {
		name              string
		sqlMock           func(mock sqlmock.Sqlmock)
		deviceToBeCreated entity.Device
		wantedErr         error
		wantedResult      *entity.Device
	}{
		{
			name: "Create Device Success Case",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(deviceCreateQuery).
					WillBeClosed().
					ExpectQuery().
					WithArgs(expectedDevice.ID, expectedDevice.Name, expectedDevice.Brand, expectedDevice.Status.String()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "creation_time"}).
						AddRow(uuid.MustParse("b44ecc02-872e-4c18-8d2a-ac09dfc4b49a"), lo.Must(time.Parse(time.DateTime, "2025-08-31 15:01:02"))))
			},
			deviceToBeCreated: entity.Device{
				ID:     expectedDevice.ID,
				Name:   expectedDevice.Name,
				Brand:  expectedDevice.Brand,
				Status: expectedDevice.Status,
			},
			wantedErr:    nil,
			wantedResult: expectedDevice,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			db, mockSql, err := sqlmock.New()
			require.NoError(t, err)

			tt.sqlMock(mockSql)
			deviceRepository := NewDeviceRepository(db)

			err = deviceRepository.CreateDevice(context.Background(), &tt.deviceToBeCreated)
			require.NoError(t, err)
			assert.Equal(tt.wantedResult, &tt.deviceToBeCreated)
		})
	}
}
