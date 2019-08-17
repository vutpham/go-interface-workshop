package main

type event struct {
	id          string
	countryCode string
	publisherID string
	gamerID     string
}

type userInfo struct {
	gamerID     string
	startCounts int
	clickCounts int
	countryInfo countryCodeGetter
}

type dsp struct {
	dspID    string
	bidValue int
}

type userDspData struct {
	user       userInfo
	winningDSP dsp
}

type countryCodeGetter interface {
	getCountryCode() string
}

func (e event) getCountryCode() string {
	return e.countryCode
}

func (u userInfo) getCountryCode() string {
	return u.countryInfo.getCountryCode()
}

func (d userDspData) getCountryCode() string {
	return d.user.getCountryCode()
}

type userInfoGetter interface {
	retrieveUserInfo() userInfo
}

func (e event) retrieveUserInfo() userInfo {
	return getUserInfoWithCountryCode(e)
}

func (d userDspData) retrieveUserInfo() userInfo {
	return d.user
}

type dataPublisher interface {
	publish() error
}

func (d userDspData) publish() error {
	return nil
}

func (d userDspData) dspID() string {
	return d.winningDSP.dspID
}

func winningDSP(u userInfo, d []dsp) dsp {
	return d[0] // dspID will be either "China1" or "Global1"
}

func handleEvent(e userInfoGetter) userDspData {
	u := e.retrieveUserInfo()
	return u.callDSPs()
}

func getUserInfoWithCountryCode(ccg countryCodeGetter) userInfo {
	return userInfo{countryInfo: ccg}
}

func (u userInfo) callDSPs() userDspData {
	switch u.getCountryCode() {
	case "CN":
		dsp := handleCallToChinaDSPs(u)
		return aggregateUserDspData(u, dsp)
	default:
		dsp := handleCallToDSPs(u)
		return aggregateUserDspData(u, dsp)
	}
}

func handleCallToChinaDSPs(u userInfo) dsp {
	chineseDSPs := []dsp{
		{dspID: "China1"},
		{dspID: "China2"},
	}
	return winningDSP(u, chineseDSPs)
}

func handleCallToDSPs(u userInfo) dsp {
	globalDSPs := []dsp{
		{dspID: "Global1"},
		{dspID: "Global2"},
	}
	return winningDSP(u, globalDSPs)
}

func aggregateUserDspData(u userInfo, d dsp) userDspData {
	return userDspData{user: u, winningDSP: d}
}

func publishUserDSPData(udd dataPublisher) error {
	return udd.publish()
}

func createEventData() event {
	return event{}
}

func main() {
	e := createEventData()
	udd := handleEvent(e)
	publishUserDSPData(udd)
}
