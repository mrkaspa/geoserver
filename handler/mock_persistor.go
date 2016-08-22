package handler

import "github.com/mrkaspa/geoserver/models"

type mockPersistor struct {
	persistAndFind chan models.PersistAndFindWithResponse
	persist        chan models.PersistWithResponse
	findStrokes    chan models.FindStrokesWithResponse
}

func (p mockPersistor) PersistAndFind() chan models.PersistAndFindWithResponse {
	return p.persistAndFind
}

func (p mockPersistor) Persist() chan models.PersistWithResponse {
	return p.persist
}

func (p mockPersistor) FindStrokes() chan models.FindStrokesWithResponse {
	return p.findStrokes
}

func newMockPersistor() models.Persistance {
	mock := &mockPersistor{
		persistAndFind: make(chan models.PersistAndFindWithResponse, 1),
		persist:        make(chan models.PersistWithResponse, 1),
		findStrokes:    make(chan models.FindStrokesWithResponse, 1),
	}
	go mock.run()
	return mock
}

func (m *mockPersistor) run() {
	for {
		select {
		case pwr := <-m.persist:
			pwr.Response <- true
		case pfwr := <-m.persistAndFind:
			pfwr.SavedResponse <- true
			pfwr.UsersResponse <- []models.Stroke{
				{
					UserID:   "a2",
					Info:     "a2",
					Location: []float64{-79.38066843, 43.65483486},
				},
			}
		case fs := <-m.findStrokes:
			fs.Response <- []models.Stroke{
				{
					UserID:   "a1",
					Info:     "a1",
					Location: []float64{-79.38066843, 43.65483486},
				},
			}
		}
	}
}
