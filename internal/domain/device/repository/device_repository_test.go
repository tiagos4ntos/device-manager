package repository

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/tiagos4ntos/device-manager/internal/domain/device/entity"
)

func makeExpectedDeviceRecord() entity.Device {
	return entity.Device{
		ID:        uuid.MustParse("b44ecc02-872e-4c18-8d2a-ac09dfc4b49a"),
		Name:      "Galaxy S23 FE",
		Brand:     "Samsumg",
		State:     "available",
		CreatedAt: lo.Must(time.Parse(time.DateTime, "2025-08-31 15:01:02")),
	}
}

func Test_Create_Device(t *testing.T) {
	assert := assert.New(t)

	deviceCreateQuery := regexp.QuoteMeta(`
	INSERT INTO devices (id, name, brand, state)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at, deleted_at;`)

	expectedDevice := makeExpectedDeviceRecord()

	type args struct {
		context           context.Context
		deviceToBeCreated entity.Device
	}
	testArgs := args{
		context: context.TODO(),
		deviceToBeCreated: entity.Device{
			ID:    expectedDevice.ID,
			Name:  expectedDevice.Name,
			Brand: expectedDevice.Brand,
			State: expectedDevice.State,
		},
	}

	testCases := []struct {
		name         string
		sqlMock      func(mock sqlmock.Sqlmock)
		args         args
		wantedErr    error
		wantedResult entity.Device
	}{
		{
			name: "Create Device Success Case",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(deviceCreateQuery).
					WillBeClosed().
					ExpectQuery().
					WithArgs(expectedDevice.ID, expectedDevice.Name, expectedDevice.Brand, expectedDevice.State.String()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at"}).
						AddRow(uuid.MustParse("b44ecc02-872e-4c18-8d2a-ac09dfc4b49a"), lo.Must(time.Parse(time.DateTime, "2025-08-31 15:01:02")), nil, nil))
			},
			args:         testArgs,
			wantedErr:    nil,
			wantedResult: expectedDevice,
		},
		{
			name: "Create Device Fails on Prepare Statement",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(deviceCreateQuery).
					WillReturnError(fmt.Errorf("some database error"))
			},
			args:         testArgs,
			wantedErr:    fmt.Errorf("some database error"),
			wantedResult: testArgs.deviceToBeCreated,
		},
		{
			name: "Create Device Fails when insert",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(deviceCreateQuery).
					WillBeClosed().
					ExpectQuery().
					WillReturnError(fmt.Errorf("some database error"))
			},
			args:         testArgs,
			wantedErr:    fmt.Errorf("some database error"),
			wantedResult: testArgs.deviceToBeCreated,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoErrorf(err, "an error '%s' was nto expected when opening a stub database connection", err)

			deviceRepository := NewDeviceRepository(db)

			tt.sqlMock(mock)

			device := tt.args.deviceToBeCreated

			err = deviceRepository.CreateDevice(tt.args.context, &device)

			assert.Equal(tt.wantedErr, err)
			assert.Equal(tt.wantedResult, device)

			mock.ExpectClose()

			err = db.Close()
			assert.NoErrorf(err, "db was not closed")

			err = mock.ExpectationsWereMet()
			assert.NoErrorf(err, "there were unfulfilled expectations")
		})
	}
}

func Test_Get_Device_ByID(t *testing.T) {
	assert := assert.New(t)

	deviceGetByIdQuery := regexp.QuoteMeta(`
	SELECT id, name, brand, state, created_at, updated_at, deleted_at
	FROM devices
	WHERE id = $1 AND deleted_at IS NULL;`)

	expectedDevice := makeExpectedDeviceRecord()

	type args struct {
		context  context.Context
		deviceID uuid.UUID
	}
	testArgs := args{
		context:  context.TODO(),
		deviceID: uuid.MustParse("b44ecc02-872e-4c18-8d2a-ac09dfc4b49a"),
	}

	testCases := []struct {
		name         string
		sqlMock      func(mock sqlmock.Sqlmock)
		args         args
		wantedErr    error
		wantedResult entity.Device
	}{
		{
			name: "Get Device By ID Success Case",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(deviceGetByIdQuery).
					WillBeClosed().
					ExpectQuery().
					WithArgs(testArgs.deviceID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "brand", "state", "created_at", "updated_at", "deleted_at"}).
						AddRow(uuid.MustParse("b44ecc02-872e-4c18-8d2a-ac09dfc4b49a"),
							"Galaxy S23 FE",
							"Samsumg",
							"available",
							lo.Must(time.Parse(time.DateTime, "2025-08-31 15:01:02")), nil, nil))
			},
			args:         testArgs,
			wantedErr:    nil,
			wantedResult: expectedDevice,
		},
		{
			name: "Get Device By ID Fails on Prepare Statement",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(deviceGetByIdQuery).
					WillReturnError(fmt.Errorf("some database error"))
			},
			args:         testArgs,
			wantedErr:    fmt.Errorf("some database error"),
			wantedResult: entity.Device{},
		},
		{
			name: "Get Device By ID Fails when insert",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(deviceGetByIdQuery).
					WillBeClosed().
					ExpectQuery().
					WillReturnError(fmt.Errorf("some database error"))
			},
			args:         testArgs,
			wantedErr:    fmt.Errorf("some database error"),
			wantedResult: entity.Device{},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoErrorf(err, "an error '%s' was nto expected when opening a stub database connection", err)

			deviceRepository := NewDeviceRepository(db)

			tt.sqlMock(mock)

			var device entity.Device
			device, err = deviceRepository.GetDeviceByID(tt.args.context, tt.args.deviceID)

			assert.Equal(tt.wantedErr, err)
			assert.Equal(tt.wantedResult, device)

			mock.ExpectClose()

			err = db.Close()
			assert.NoErrorf(err, "db was not closed")

			err = mock.ExpectationsWereMet()
			assert.NoErrorf(err, "there were unfulfilled expectations")
		})
	}
}

func Test_Fully_Update_Device(t *testing.T) {
	assert := assert.New(t)
	deviceUpdatedAt := lo.Must(time.Parse(time.DateTime, "2025-09-01 19:11:22"))

	updateDeviceQuery := regexp.QuoteMeta(`
	UPDATE devices SET 
		name = $2,
		brand = $3,
		state = $4,
		updated_at = now()
	WHERE id = $1 AND deleted_at IS NULL
	RETURNING id, created_at, updated_at, deleted_at;`)

	updatedDevice := makeExpectedDeviceRecord()
	updatedDevice.UpdatedAt = lo.ToPtr(deviceUpdatedAt)

	deviceToBeUpdated := makeExpectedDeviceRecord()

	type args struct {
		context        context.Context
		deviceToUpdate entity.Device
	}
	testArgs := args{
		context:        context.TODO(),
		deviceToUpdate: deviceToBeUpdated,
	}

	testCases := []struct {
		name         string
		sqlMock      func(mock sqlmock.Sqlmock)
		args         args
		wantedErr    error
		wantedResult entity.Device
	}{
		{
			name: "Update Device Success Case",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(updateDeviceQuery).
					WillBeClosed().
					ExpectQuery().
					WithArgs(deviceToBeUpdated.ID, deviceToBeUpdated.Name, deviceToBeUpdated.Brand, deviceToBeUpdated.State.String()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at"}).
						AddRow(updatedDevice.ID, updatedDevice.CreatedAt, deviceUpdatedAt, nil))
			},
			args:         testArgs,
			wantedErr:    nil,
			wantedResult: updatedDevice,
		},
		{
			name: "Update Device Fails on Prepare Statement",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(updateDeviceQuery).
					WillReturnError(fmt.Errorf("some database error"))
			},
			args:         testArgs,
			wantedErr:    fmt.Errorf("some database error"),
			wantedResult: testArgs.deviceToUpdate,
		},
		{
			name: "Update Device Fails when insert",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(updateDeviceQuery).
					WillBeClosed().
					ExpectQuery().
					WillReturnError(fmt.Errorf("some database error"))
			},
			args:         testArgs,
			wantedErr:    fmt.Errorf("some database error"),
			wantedResult: testArgs.deviceToUpdate,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoErrorf(err, "an error '%s' was nto expected when opening a stub database connection", err)

			deviceRepository := NewDeviceRepository(db)

			tt.sqlMock(mock)

			err = deviceRepository.FullyUpdateDevice(tt.args.context, &tt.args.deviceToUpdate)

			assert.Equal(tt.wantedErr, err)
			assert.Equal(tt.wantedResult, tt.args.deviceToUpdate)

			mock.ExpectClose()

			err = db.Close()
			assert.NoErrorf(err, "db was not closed")

			err = mock.ExpectationsWereMet()
			assert.NoErrorf(err, "there were unfulfilled expectations")
		})
	}
}

func Test_Update_Device_State(t *testing.T) {
	assert := assert.New(t)
	deviceUpdatedAt := lo.Must(time.Parse(time.DateTime, "2025-09-01 19:11:22"))
	newStatus := entity.InUse

	updateDeviceQuery := regexp.QuoteMeta(`
	UPDATE devices SET 
		state = $2,
		updated_at = now()
	WHERE id = $1 AND deleted_at IS NULL
	RETURNING id, name, brand, state, created_at, updated_at, deleted_at;`)

	updatedDevice := makeExpectedDeviceRecord()
	updatedDevice.UpdatedAt = lo.ToPtr(deviceUpdatedAt)
	updatedDevice.State = newStatus

	deviceToBeUpdated := makeExpectedDeviceRecord()

	type args struct {
		context        context.Context
		deviceToUpdate entity.Device
	}
	testArgs := args{
		context:        context.TODO(),
		deviceToUpdate: deviceToBeUpdated,
	}

	testCases := []struct {
		name         string
		sqlMock      func(mock sqlmock.Sqlmock)
		args         args
		wantedErr    error
		wantedResult entity.Device
	}{
		{
			name: "Update Device Status Success Case",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(updateDeviceQuery).
					WillBeClosed().
					ExpectQuery().
					WithArgs(deviceToBeUpdated.ID, newStatus).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "brand", "state", "created_at", "updated_at", "deleted_at"}).
						AddRow(uuid.MustParse("b44ecc02-872e-4c18-8d2a-ac09dfc4b49a"),
							"Galaxy S23 FE",
							"Samsumg",
							newStatus.String(),
							lo.Must(time.Parse(time.DateTime, "2025-08-31 15:01:02")),
							deviceUpdatedAt,
							nil))
			},
			args:         testArgs,
			wantedErr:    nil,
			wantedResult: updatedDevice,
		},
		{
			name: "Update Device Status Fails on Prepare Statement",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(updateDeviceQuery).
					WillReturnError(fmt.Errorf("some database error"))
			},
			args:         testArgs,
			wantedErr:    fmt.Errorf("some database error"),
			wantedResult: entity.Device{},
		},
		{
			name: "Update Device Status Fails when insert",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(updateDeviceQuery).
					WillBeClosed().
					ExpectQuery().
					WillReturnError(fmt.Errorf("some database error"))
			},
			args:         testArgs,
			wantedErr:    fmt.Errorf("some database error"),
			wantedResult: entity.Device{},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoErrorf(err, "an error '%s' was nto expected when opening a stub database connection", err)

			deviceRepository := NewDeviceRepository(db)

			tt.sqlMock(mock)

			device, err := deviceRepository.UpdateDeviceState(tt.args.context, tt.args.deviceToUpdate.ID, newStatus)

			assert.Equal(tt.wantedErr, err)
			assert.Equal(tt.wantedResult, device)

			mock.ExpectClose()

			err = db.Close()
			assert.NoErrorf(err, "db was not closed")

			err = mock.ExpectationsWereMet()
			assert.NoErrorf(err, "there were unfulfilled expectations")
		})
	}
}

func Test_Delete_Device(t *testing.T) {
	assert := assert.New(t)
	deletedDeviceID := uuid.MustParse("b44ecc02-872e-4c18-8d2a-ac09dfc4b49a")

	deleteDeviceQuery := regexp.QuoteMeta(`UPDATE devices SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL AND state <> 'in-use'`)

	type args struct {
		context  context.Context
		deviceID uuid.UUID
	}
	testArgs := args{
		context:  context.TODO(),
		deviceID: deletedDeviceID,
	}

	testCases := []struct {
		name      string
		sqlMock   func(mock sqlmock.Sqlmock)
		args      args
		wantedErr error
	}{
		{
			name: "Delete Device Success Case",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(deleteDeviceQuery).
					WillBeClosed().
					ExpectExec().
					WithArgs(testArgs.deviceID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args:      testArgs,
			wantedErr: nil,
		},
		{
			name: "Delete Device Fails when no rows are affected",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(deleteDeviceQuery).
					WillBeClosed().
					ExpectExec().
					WithArgs(testArgs.deviceID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			args:      testArgs,
			wantedErr: sql.ErrNoRows,
		},
		{
			name: "Delete Device Fails on rows affected",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(deleteDeviceQuery).
					WillBeClosed().
					ExpectExec().
					WithArgs(testArgs.deviceID).
					WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some database error")))
			},
			args:      testArgs,
			wantedErr: fmt.Errorf("some database error"),
		},
		{
			name: "Delete Device Fails on execute",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(deleteDeviceQuery).
					WillBeClosed().
					ExpectExec().
					WillReturnError(fmt.Errorf("some database error"))
			},
			args:      testArgs,
			wantedErr: fmt.Errorf("some database error"),
		},
		{
			name: "Delete Device Fails on preapre statement",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(deleteDeviceQuery).
					WillReturnError(fmt.Errorf("some error"))
			},
			args:      testArgs,
			wantedErr: fmt.Errorf("some error"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoErrorf(err, "an error '%s' was nto expected when opening a stub database connection", err)

			deviceRepository := NewDeviceRepository(db)

			tt.sqlMock(mock)

			err = deviceRepository.DeleteDevice(tt.args.context, tt.args.deviceID)

			assert.Equal(tt.wantedErr, err)

			mock.ExpectClose()

			err = db.Close()
			assert.NoErrorf(err, "db was not closed")

			err = mock.ExpectationsWereMet()
			assert.NoErrorf(err, "there were unfulfilled expectations")
		})
	}
}

func Test_List_Devices(t *testing.T) {
	assert := assert.New(t)
	createdAt := lo.Must(time.Parse(time.DateTime, "2025-08-31 15:01:02"))

	deviceListQuery := regexp.QuoteMeta(`SELECT id, name, brand, state, created_at, updated_at, deleted_at
	FROM devices
	WHERE deleted_at IS NULL
	AND ($1 IS NULL OR brand = $1)
  	AND ($2 IS NULL OR state = $2)
	ORDER BY name;`)

	type args struct {
		context context.Context
		params  map[string]any
	}
	testArgs := args{
		context: context.TODO(),
		params: map[string]any{
			"brand": nil,
			"state": nil,
		},
	}

	testCases := []struct {
		name         string
		sqlMock      func(mock sqlmock.Sqlmock)
		args         args
		wantedErr    error
		wantedResult []entity.Device
	}{
		{
			name: "List Devices Success Case",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(deviceListQuery).
					WillBeClosed().
					ExpectQuery().
					WithArgs(nil, nil).
					WillReturnRows(
						sqlmock.
							NewRows([]string{"id", "name", "brand", "state", "created_at", "updated_at", "deleted_at"}).
							AddRow(uuid.MustParse("b44ecc02-872e-4c18-8d2a-ac09dfc4b49a"), "Galaxy S23 FE", "Samsumg", entity.Available, createdAt, nil, nil).
							AddRow(uuid.MustParse("c60dceb7-60c8-4d74-8d7c-cd34a0b4ce19"), "IPhone 15", "Apple", entity.InUse, createdAt, createdAt, nil))
			},
			args:      testArgs,
			wantedErr: nil,
			wantedResult: []entity.Device{
				{
					ID:        uuid.MustParse("b44ecc02-872e-4c18-8d2a-ac09dfc4b49a"),
					Name:      "Galaxy S23 FE",
					Brand:     "Samsumg",
					State:     entity.Available,
					CreatedAt: createdAt,
				},
				{
					ID:        uuid.MustParse("c60dceb7-60c8-4d74-8d7c-cd34a0b4ce19"),
					Name:      "IPhone 15",
					Brand:     "Apple",
					State:     entity.InUse,
					CreatedAt: createdAt,
					UpdatedAt: &createdAt,
				},
			},
		},
		{
			name: "List Devices filtering Brand Success Case",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(deviceListQuery).
					WillBeClosed().
					ExpectQuery().
					WithArgs("Apple", nil).
					WillReturnRows(
						sqlmock.
							NewRows([]string{"id", "name", "brand", "state", "created_at", "updated_at", "deleted_at"}).
							AddRow(uuid.MustParse("c60dceb7-60c8-4d74-8d7c-cd34a0b4ce19"), "IPhone 15", "Apple", entity.InUse, createdAt, createdAt, nil))
			},
			args: args{
				context: context.TODO(),
				params: map[string]any{
					"brand": "Apple",
					"state": nil,
				},
			},
			wantedErr: nil,
			wantedResult: []entity.Device{
				{
					ID:        uuid.MustParse("c60dceb7-60c8-4d74-8d7c-cd34a0b4ce19"),
					Name:      "IPhone 15",
					Brand:     "Apple",
					State:     entity.InUse,
					CreatedAt: createdAt,
					UpdatedAt: &createdAt,
				},
			},
		},
		{
			name: "List Devices filtering Brand and State Success Case",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(deviceListQuery).
					WillBeClosed().
					ExpectQuery().
					WithArgs("Apple", "in-use").
					WillReturnRows(
						sqlmock.
							NewRows([]string{"id", "name", "brand", "state", "created_at", "updated_at", "deleted_at"}).
							AddRow(uuid.MustParse("a60dceb7-60c8-4d74-8d7c-cd34a0b4ce11"), "IPhone 16", "Apple", entity.InUse, createdAt, createdAt, nil))
			},
			args: args{
				context: context.TODO(),
				params: map[string]any{
					"brand": "Apple",
					"state": "in-use",
				},
			},
			wantedErr: nil,
			wantedResult: []entity.Device{
				{
					ID:        uuid.MustParse("a60dceb7-60c8-4d74-8d7c-cd34a0b4ce11"),
					Name:      "IPhone 16",
					Brand:     "Apple",
					State:     entity.InUse,
					CreatedAt: createdAt,
					UpdatedAt: &createdAt,
				},
			},
		},
		{
			name: "List Devices  Fails on Prepare Statement",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(deviceListQuery).
					WillReturnError(fmt.Errorf("some database error"))
			},
			args:         testArgs,
			wantedErr:    fmt.Errorf("some database error"),
			wantedResult: nil,
		},
		{
			name: "List Devices Fails when query",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(deviceListQuery).
					WillBeClosed().
					ExpectQuery().
					WillReturnError(fmt.Errorf("some database error"))
			},
			args:         testArgs,
			wantedErr:    fmt.Errorf("some database error"),
			wantedResult: nil,
		},
		{
			name: "List Devices Fails on Row Scan",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(deviceListQuery).
					WillBeClosed().
					ExpectQuery().
					WithArgs(nil, nil).
					WillReturnRows(
						sqlmock.
							NewRows([]string{"id"}).
							AddRow(uuid.MustParse("b44ecc02-872e-4c18-8d2a-ac09dfc4b49a")))
			},
			args:         testArgs,
			wantedErr:    fmt.Errorf("sql: expected %d destination arguments in Scan, not %d", 1, 7),
			wantedResult: nil,
		},
		{
			name: "List Devices Fails on rows.Error",
			sqlMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(deviceListQuery).
					WillBeClosed().
					ExpectQuery().
					WithArgs(nil, nil).
					WillReturnRows(
						sqlmock.
							NewRows([]string{"id", "name", "brand", "state", "created_at", "updated_at", "deleted_at"}).
							AddRow(uuid.MustParse("b44ecc02-872e-4c18-8d2a-ac09dfc4b49a"), "Galaxy S23 FE", "Samsumg", entity.Available, createdAt, nil, nil).
							AddRow(uuid.MustParse("c60dceb7-60c8-4d74-8d7c-cd34a0b4ce19"), "IPhone 15", "Apple", entity.InUse, createdAt, createdAt, nil).
							RowError(1, fmt.Errorf("some error")))
			},
			args:         testArgs,
			wantedErr:    fmt.Errorf("some error"),
			wantedResult: nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoErrorf(err, "an error '%s' was nto expected when opening a stub database connection", err)

			deviceRepository := NewDeviceRepository(db)

			tt.sqlMock(mock)

			devices, err := deviceRepository.ListDevices(tt.args.context, tt.args.params)

			assert.Equal(tt.wantedErr, err)
			assert.Equal(tt.wantedResult, devices)

			mock.ExpectClose()

			err = db.Close()
			assert.NoErrorf(err, "db was not closed")

			err = mock.ExpectationsWereMet()
			assert.NoErrorf(err, "there were unfulfilled expectations")
		})
	}
}
