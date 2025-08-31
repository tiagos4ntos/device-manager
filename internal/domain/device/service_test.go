package device

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tiagos4ntos/device-manager/internal/domain/device/entity"
	"github.com/tiagos4ntos/device-manager/internal/domain/device/errors"
	"github.com/tiagos4ntos/device-manager/internal/domain/mocks"
)

var errDatabaseGeneric = fmt.Errorf("some database error")

func Test_List_Device(t *testing.T) {
	type args struct {
		context context.Context
	}

	testArgs := args{
		context: context.TODO(),
	}

	tests := []struct {
		name                 string
		testArgs             args
		wantRepositoryResult []entity.Device
		wantRepositoryErr    error
		wantErr              error
	}{
		{
			name:     "List Success Case",
			testArgs: testArgs,
			wantRepositoryResult: []entity.Device{
				{
					ID:    uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Name:  "Galaxy S21",
					Brand: "Samsung",
					State: entity.InUse,
				},
				{
					ID:    uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Name:  "Iphone 13",
					Brand: "Apple",
					State: entity.Available,
				},
			},
			wantRepositoryErr: nil,
			wantErr:           nil,
		},
		{
			name:                 "List Repository Error Case",
			testArgs:             testArgs,
			wantRepositoryResult: nil,
			wantRepositoryErr:    errDatabaseGeneric,
			wantErr:              errors.NewDeviceError(errors.ErrInternal, "something went wrong while listing devices", errDatabaseGeneric),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockRepo := mocks.NewMockDeviceRepository(mockCtrl)
			service := NewDeviceService(mockRepo)

			mockRepo.
				EXPECT().
				ListDevices(tt.testArgs.context).
				Return(tt.wantRepositoryResult, tt.wantRepositoryErr).
				AnyTimes()

			devices, err := service.List(tt.testArgs.context)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantRepositoryResult, devices)
		})
	}
}

func Test_GetByID_Device(t *testing.T) {
	type args struct {
		context  context.Context
		deviceID uuid.UUID
	}

	testArgs := args{
		context:  context.TODO(),
		deviceID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
	}

	tests := []struct {
		name                 string
		testArgs             args
		wantRepositoryResult entity.Device
		wantRepositoryErr    error
		wantErr              error
	}{
		{
			name:     "GetByID Success Case",
			testArgs: testArgs,
			wantRepositoryResult: entity.Device{
				ID:    testArgs.deviceID,
				Name:  "Galaxy S21",
				Brand: "Samsung",
				State: entity.InUse,
			},
			wantRepositoryErr: nil,
			wantErr:           nil,
		},
		{
			name:                 "GetByID Device Repository Error Case",
			testArgs:             testArgs,
			wantRepositoryResult: entity.Device{},
			wantRepositoryErr:    errDatabaseGeneric,
			wantErr:              errors.NewDeviceError(errors.ErrInternal, "something went wrong while retrieving device", errDatabaseGeneric),
		},
		{
			name:                 "GetByID Device Not Found Case",
			testArgs:             testArgs,
			wantRepositoryResult: entity.Device{},
			wantRepositoryErr:    sql.ErrNoRows,
			wantErr:              errors.NewDeviceError(errors.ErrNotFound, "device not found", sql.ErrNoRows),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockRepo := mocks.NewMockDeviceRepository(mockCtrl)
			service := NewDeviceService(mockRepo)

			mockRepo.
				EXPECT().
				GetDeviceByID(tt.testArgs.context, gomock.Any()).
				Return(tt.wantRepositoryResult, tt.wantRepositoryErr).
				AnyTimes()

			device, err := service.GetByID(tt.testArgs.context, tt.testArgs.deviceID)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantRepositoryResult, device)
		})
	}
}

func Test_Create_Device(t *testing.T) {
	type args struct {
		context context.Context
	}

	testArgs := args{
		context: context.TODO(),
	}

	tests := []struct {
		name                string
		testArgs            args
		device              entity.Device
		wantedRepositoryErr error
		wantErr             error
	}{
		{
			name:     "Create Device Success Case",
			testArgs: testArgs,
			device: entity.Device{
				Name:  "Galaxy S21",
				Brand: "Samsung",
				State: entity.InUse,
			},
			wantedRepositoryErr: nil,
			wantErr:             nil,
		},
		{
			name:     "Create Device Repository Error Case",
			testArgs: testArgs,
			device: entity.Device{
				Name:  "iPhone 13",
				Brand: "Apple",
				State: entity.Available,
			},
			wantedRepositoryErr: errDatabaseGeneric,
			wantErr:             errors.NewDeviceError(errors.ErrInternal, "something went wrong while creating device", errDatabaseGeneric),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockRepo := mocks.NewMockDeviceRepository(mockCtrl)
			service := NewDeviceService(mockRepo)

			mockRepo.
				EXPECT().
				CreateDevice(tt.testArgs.context, gomock.Any()).
				Return(tt.wantedRepositoryErr).
				AnyTimes()

			device, err := service.Create(tt.testArgs.context, tt.device)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.device.Name, device.Name)
			assert.Equal(t, tt.device.Brand, device.Brand)
			assert.Equal(t, tt.device.State, device.State)
			assert.NotEqual(t, uuid.Nil, device.ID)
		})
	}
}

func Test_Update_Device(t *testing.T) {
	type args struct {
		context context.Context
	}

	testArgs := args{
		context: context.TODO(),
	}

	tests := []struct {
		name                    string
		testArgs                args
		device                  entity.Device
		updateOnlyStatus        bool
		wantedRepoGetByIdResult entity.Device
		wantedRepoGetByIdError  error
		wantedRepoUpdateErr     error
		wantErr                 error
	}{
		{
			name:     "Update Device Success Case",
			testArgs: testArgs,
			device: entity.Device{
				ID:    uuid.MustParse("215f759c-aa0f-494f-84ba-0d706dd6d59a"),
				Name:  "Galaxy S21 Updated",
				Brand: "Samsung",
				State: entity.InUse,
			},
			updateOnlyStatus: false,
			wantedRepoGetByIdResult: entity.Device{
				ID:    uuid.MustParse("215f759c-aa0f-494f-84ba-0d706dd6d59a"),
				Name:  "Galaxy S21",
				Brand: "Samsung",
				State: entity.Available,
			},
			wantedRepoGetByIdError: nil,
			wantedRepoUpdateErr:    nil,
			wantErr:                nil,
		},
		{
			name:     "Update Device Error Retrieving Device Case",
			testArgs: testArgs,
			device: entity.Device{
				ID:    uuid.MustParse("215f759c-aa0f-494f-84ba-0d706dd6d59a"),
				Name:  "Galaxy S21 Updated",
				Brand: "Samsung",
				State: entity.InUse,
			},
			updateOnlyStatus:        false,
			wantedRepoGetByIdResult: entity.Device{},
			wantedRepoGetByIdError:  errDatabaseGeneric,
			wantedRepoUpdateErr:     nil,
			wantErr:                 errors.NewDeviceError(errors.ErrNotFound, "something went wrong while retrieving device", errDatabaseGeneric),
		},
		{
			name:     "Update Only Device Not Occur When In Use Case",
			testArgs: testArgs,
			device: entity.Device{
				ID:    uuid.MustParse("215f759c-aa0f-494f-84ba-0d706dd6d59a"),
				Name:  "Galaxy S21 Updated",
				Brand: "Samsung",
				State: entity.InUse,
			},
			updateOnlyStatus: true,
			wantedRepoGetByIdResult: entity.Device{
				ID:    uuid.MustParse("215f759c-aa0f-494f-84ba-0d706dd6d59a"),
				Name:  "Galaxy S21 Updated",
				Brand: "Samsung",
				State: entity.InUse,
			},
			wantedRepoGetByIdError: nil,
			wantedRepoUpdateErr:    nil,
			wantErr:                errors.NewDeviceError(errors.ErrInvalid, "device is in use and cannot be updated", fmt.Errorf("device state is the same: %s", entity.InUse)),
		},
		{
			name:     "Update Only Device State with Success",
			testArgs: testArgs,
			device: entity.Device{
				ID:    uuid.MustParse("215f759c-aa0f-494f-84ba-0d706dd6d59a"),
				Name:  "Galaxy S21 Updated",
				Brand: "Samsung",
				State: entity.Available,
			},
			updateOnlyStatus: true,
			wantedRepoGetByIdResult: entity.Device{
				ID:    uuid.MustParse("215f759c-aa0f-494f-84ba-0d706dd6d59a"),
				Name:  "Galaxy S21 Updated",
				Brand: "Samsung",
				State: entity.InUse,
			},
			wantedRepoGetByIdError: nil,
			wantedRepoUpdateErr:    nil,
			wantErr:                nil,
		},
		{
			name:     "Update Only Device State with Expected Error",
			testArgs: testArgs,
			device: entity.Device{
				ID:    uuid.MustParse("215f759c-aa0f-494f-84ba-0d706dd6d59a"),
				Name:  "Galaxy S21 Updated",
				Brand: "Samsung",
				State: entity.Available,
			},
			updateOnlyStatus: true,
			wantedRepoGetByIdResult: entity.Device{
				ID:    uuid.MustParse("215f759c-aa0f-494f-84ba-0d706dd6d59a"),
				Name:  "Galaxy S21 Updated",
				Brand: "Samsung",
				State: entity.InUse,
			},
			wantedRepoGetByIdError: nil,
			wantedRepoUpdateErr:    errDatabaseGeneric,
			wantErr:                errors.NewDeviceError(errors.ErrInternal, "something went wrong while update device state", fmt.Errorf("error updating device state: %v", errDatabaseGeneric)),
		},
		{
			name:     "Update Only Device State with Expected Error",
			testArgs: testArgs,
			device: entity.Device{
				ID:    uuid.MustParse("215f759c-aa0f-494f-84ba-0d706dd6d59a"),
				Name:  "Galaxy S21 Updated",
				Brand: "Samsung",
				State: entity.InUse,
			},
			updateOnlyStatus: false,
			wantedRepoGetByIdResult: entity.Device{
				ID:    uuid.MustParse("215f759c-aa0f-494f-84ba-0d706dd6d59a"),
				Name:  "Galaxy S21 Updated",
				Brand: "Samsung",
				State: entity.Available,
			},
			wantedRepoGetByIdError: nil,
			wantedRepoUpdateErr:    errDatabaseGeneric,
			wantErr:                errors.NewDeviceError(errors.ErrInternal, "something went wrong while fully update device", errDatabaseGeneric),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockRepo := mocks.NewMockDeviceRepository(mockCtrl)
			service := NewDeviceService(mockRepo)

			mockRepo.
				EXPECT().
				GetDeviceByID(tt.testArgs.context, gomock.Any()).
				Return(tt.wantedRepoGetByIdResult, tt.wantedRepoGetByIdError).
				AnyTimes()

			if tt.updateOnlyStatus {
				mockRepo.EXPECT().
					UpdateDeviceState(tt.testArgs.context, tt.device.ID, tt.device.State).
					Return(tt.device, tt.wantedRepoUpdateErr).
					AnyTimes()
			} else {
				mockRepo.EXPECT().
					FullyUpdateDevice(tt.testArgs.context, gomock.Any()).
					Return(tt.wantedRepoUpdateErr).
					AnyTimes()
			}

			updatedDevice, err := service.Update(tt.testArgs.context, tt.device)
			assert.Equal(t, tt.wantErr, err)
			if err == nil {
				assert.Equal(t, tt.device.Name, updatedDevice.Name)
				assert.Equal(t, tt.device.Brand, updatedDevice.Brand)
				assert.Equal(t, tt.device.State, updatedDevice.State)
			}
		})
	}
}

func Test_Delete_Device(t *testing.T) {
	type args struct {
		context context.Context
	}

	testArgs := args{
		context: context.TODO(),
	}

	tests := []struct {
		name                string
		testArgs            args
		deviceId            uuid.UUID
		wantedRepositoryErr error
		wantErr             error
	}{
		{
			name:                "Delete Device Success Case",
			testArgs:            testArgs,
			deviceId:            uuid.MustParse("215f759c-aa0f-494f-84ba-0d706dd6d59a"),
			wantedRepositoryErr: nil,
			wantErr:             nil,
		},
		{
			name:                "Delete Device Repository Error Case",
			testArgs:            testArgs,
			deviceId:            uuid.MustParse("215f759c-aa0f-494f-84ba-0d706dd6d59a"),
			wantedRepositoryErr: errDatabaseGeneric,
			wantErr:             errors.NewDeviceError(errors.ErrInternal, "something went wrong while delete device", errDatabaseGeneric),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockRepo := mocks.NewMockDeviceRepository(mockCtrl)
			service := NewDeviceService(mockRepo)

			mockRepo.
				EXPECT().
				DeleteDevice(tt.testArgs.context, tt.deviceId).
				Return(tt.wantedRepositoryErr).
				AnyTimes()

			err := service.Delete(tt.testArgs.context, tt.deviceId)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
