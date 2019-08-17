package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_retrieveUserInfo(t *testing.T) {
	t.Run("retrieveUserInfo:  Returns correct countryCode (MX)", func(t *testing.T) {
		e := &mockEvent{}
		e.On("getCountryCode").Return("MX")
		e.On("retrieveUserInfo").Return(userInfo{countryInfo: e})
		got := e.retrieveUserInfo()
		assert.NotNil(t, got)
		assert.Equal(t, "MX", got.getCountryCode())
		assert.IsType(t, userInfo{}, got)

		e.AssertExpectations(t)
	})
}

func Test_handleEvent(t *testing.T) {
	t.Run("handleEvent with countryCode (US)", func(t *testing.T) {
		e := &mockEvent{}
		e.On("getCountryCode").Return("US")
		e.On("retrieveUserInfo").Return(userInfo{countryInfo: e})

		got := handleEvent(e)
		assert.NotNil(t, got)
		assert.Equal(t, "US", got.getCountryCode())
		assert.Equal(t, "Global1", got.dspID())
		assert.IsType(t, userDspData{}, got)

		e.AssertExpectations(t)
	})

	t.Run("handleEvent with countryCode (CN)", func(t *testing.T) {
		e := &mockEvent{}
		e.On("getCountryCode").Return("CN")
		e.On("retrieveUserInfo").Return(userInfo{countryInfo: e})

		got := handleEvent(e)
		assert.NotNil(t, got)
		assert.Equal(t, "CN", got.getCountryCode())
		assert.Equal(t, "China1", got.dspID())
		assert.IsType(t, userDspData{}, got)

		e.AssertExpectations(t)
	})
}

func Test_userDspData_retrieveUserInfo(t *testing.T) {
	t.Run("retrieveUserInfo", func(t *testing.T) {
		e := &mockEvent{}
		e.On("getCountryCode").Return("VN")
		u := userInfo{countryInfo: e}
		udd := userDspData{user: u}
		assert.Equal(t, "VN", udd.getCountryCode())
		assert.IsType(t, userInfo{}, udd.retrieveUserInfo())
		e.AssertExpectations(t)
	})
}

func Test_publishUserDSPData(t *testing.T) {
	t.Run("publishUserDSPData:  udd.publish() mock with error", func(t *testing.T) {
		udd := &mockUserDSPData{}
		udd.On("publish").Return(errors.New("Error publisher user-dsp data"))
		err := publishUserDSPData(udd)
		assert.Error(t, err)
		assert.Equal(t, "Error publisher user-dsp data", err.Error())
		udd.AssertExpectations(t)
	})

	t.Run("publishUserDSPData:  udd invokes publish() and returns nil", func(t *testing.T) {
		udd := userDspData{}
		err := publishUserDSPData(udd)
		assert.NoError(t, err)
	})

}

func Test_userDspData_dspID(t *testing.T) {
	t.Run("dspID", func(t *testing.T) {
		d := dsp{dspID: "FakeID"}
		udd := userDspData{winningDSP: d}
		assert.Equal(t, d.dspID, udd.dspID())
	})
}

type mockEvent struct {
	mock.Mock
}

func (e *mockEvent) getCountryCode() string {
	args := e.Called()
	return args.String(0)
}

func (e *mockEvent) retrieveUserInfo() userInfo {
	args := e.Called()
	return args.Get(0).(userInfo)
}

type mockUserInfo struct {
	mock.Mock
}

func (u *mockUserInfo) callDSPs() userDspData {
	args := u.Called()
	return args.Get(0).(userDspData)
}

func (u *mockUserInfo) getCountryCode() string {
	args := u.Called()
	return args.String(0)
}

type mockUserDSPData struct {
	mock.Mock
}

func (u *mockUserDSPData) publish() error {
	args := u.Called()
	return args.Error(0)
}

// Dumb tests that I just added for code coverage :P
func Test_createEventData(t *testing.T) {
	t.Run("CreatesEventData", func(t *testing.T) {
		got := createEventData()
		assert.NotNil(t, got)
	})
}

func Test_event_getCountryCode(t *testing.T) {
	t.Run("getCountryCode() returns correct country code info", func(t *testing.T) {
		e := event{countryCode: "KR"}
		assert.Equal(t, "KR", e.getCountryCode())
	})
}

func Test_event_retrieveUserInfo(t *testing.T) {
	t.Run("retrieveUserInfo", func(t *testing.T) {
		e := &event{countryCode: "EU"}
		u := e.retrieveUserInfo()
		assert.Equal(t, "EU", u.getCountryCode())
	})

	t.Run("getUserInfoWithCountryCode", func(t *testing.T) {
		e := &event{countryCode: "BR"}
		u := getUserInfoWithCountryCode(e)
		assert.Equal(t, "BR", u.getCountryCode())
	})
}
