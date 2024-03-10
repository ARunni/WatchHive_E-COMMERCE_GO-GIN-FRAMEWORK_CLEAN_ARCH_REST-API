package usecase

import (
	"WatchHive/pkg/config"
	mockhelper "WatchHive/pkg/helper/mock"
	mockRepository "WatchHive/pkg/repository/mock"
	"WatchHive/pkg/utils/models"

	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_GetAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockRepository.NewMockUserRepository(ctrl)
	cfg := config.Config{}
	adminRepo := mockRepository.NewMockAdminRepository(ctrl)
	wallet := mockRepository.NewMockWalletRepository(ctrl)
	helper := mockhelper.NewMockHelper(ctrl)

	userUseCase := NewUserUseCase(userRepo, adminRepo, cfg, helper, wallet)
	testingData := map[string]struct {
		input   int
		stub    func(*mockRepository.MockUserRepository, *mockhelper.MockHelper, int)
		want    []models.AddressInfoResponse
		wantErr error
	}{
		"success": {
			input: 1,
			stub: func(ur *mockRepository.MockUserRepository, mh *mockhelper.MockHelper, i int) {
				ur.EXPECT().GetAllAddress(i).Times(1).Return([]models.AddressInfoResponse{}, nil)
			},
			want:    []models.AddressInfoResponse{},
			wantErr: nil,
		},
		"failed": {
			input: 1,
			stub: func(ur *mockRepository.MockUserRepository, mh *mockhelper.MockHelper, i int) {
				ur.EXPECT().GetAllAddress(i).Times(1).Return([]models.AddressInfoResponse{}, errors.New("error in getting addresses"))
			},
			want:    []models.AddressInfoResponse{},
			wantErr: errors.New("error in getting addresses"),
		},
	}
	for _, test := range testingData {
		test.stub(userRepo, helper, test.input)
		result, err := userUseCase.GetAllAddress(test.input)
		assert.Equal(t, test.want, result)
		assert.Equal(t, test.wantErr, err)
	}
}
