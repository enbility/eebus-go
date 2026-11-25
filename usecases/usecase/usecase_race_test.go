package usecase

import (
	"sync"

	"github.com/enbility/spine-go/model"
)

// Two "device removed" events for the same remote device handled
// concurrently: indices are looked up under the lock, but removed later
// under a new lock, so the second removal may use stale indices.
func (s *UseCaseSuite) Test_ConcurrentRemoveDeviceScenarios() {
	localDevice := s.localEntity.Device()
	remoteDevice2, remoteEntities2 := addRemoteDevice("remoteDevice2", "remoteDevice2Ski", localDevice, s.T())
	localDevice.AddRemoteDeviceForSki("remoteDevice2Ski", remoteDevice2)

	for round := 0; round < 2000; round++ {
		s.uc.updateRemoteEntityScenarios(remoteEntities2[0], []model.UseCaseScenarioSupportType{1})
		s.uc.updateRemoteEntityScenarios(remoteEntities2[1], []model.UseCaseScenarioSupportType{1})

		var wg sync.WaitGroup
		for g := 0; g < 2; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				s.uc.removeDeviceFromAvailableEntityScenarios(remoteDevice2)
			}()
		}
		wg.Wait()
	}
}
