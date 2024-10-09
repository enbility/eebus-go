package cevc

func (s *CemCEVCSuite) Test_Heartbeat() {
	s.sut.StartHeartbeat()

	s.sut.StopHeartbeat()
}

func (s *CemCEVCSuite) Test_OperatingState() {
	s.sut.SetOperatingState(true)
}
