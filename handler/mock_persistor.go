package handler

import "github.com/mrkaspa/geoserver/models"

type mockPersistor struct {
	persistAndFind chan *models.StrokeNear
	persist        chan *models.Stroke
	usersFound     chan []models.Stroke
	saved          chan bool
}

func (p mockPersistor) PersistAndFind() chan *models.StrokeNear {
	return p.persistAndFind
}

func (p mockPersistor) Persist() chan *models.Stroke {
	return p.persist
}

func (p mockPersistor) UsersFound() chan []models.Stroke {
	return p.usersFound
}

func (p mockPersistor) Saved() chan bool {
	return p.saved
}

func newMockPersistor() models.Persistance {
	mock := &mockPersistor{
		persistAndFind: make(chan *models.StrokeNear, 1),
		persist:        make(chan *models.Stroke, 1),
		usersFound:     make(chan []models.Stroke, 1),
		saved:          make(chan bool, 1),
	}
	go mock.run()
	return mock
}

func (m *mockPersistor) run() {
	for {
		select {
		case <-m.persist:
			m.saved <- true
		case <-m.persistAndFind:
			m.saved <- true
			m.usersFound <- []models.Stroke{
				{
					UserID:   "a2",
					Info:     "a2",
					Location: []float64{-79.38066843, 43.65483486},
				},
			}
		}
	}
}
