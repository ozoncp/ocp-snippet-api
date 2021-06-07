package saver

import (
	"fmt"
	"time"

	"github.com/ozoncp/ocp-snippet-api/internal/flusher"
	"github.com/ozoncp/ocp-snippet-api/internal/models"
)

type Saver interface {
	Save(snippets []models.Snippet) error
	Init()
	Close() error
}

type saver struct {
	flusher  flusher.Flusher
	snippets []models.Snippet

	saveInterval uint
	doneCh       chan struct{}
}

// NewSaver возвращает Save с поддержкой периодического сохранения.
// capacity - capacity слайса сохраняемых сущностей.
// flusher - экземпляр интерфейса flusher.Flusher, осуществляющий сохранения.
// saveInterval - интервал сохранения в секундах.
func NewSaver(capacity uint, flusher flusher.Flusher, saveInterval uint) Saver {
	return &saver{
		flusher:      flusher,
		snippets:     make([]models.Snippet, 0, capacity),
		saveInterval: saveInterval,
	}
}

func (s saver) Save(snippets []models.Snippet) error {
	snippets, err := s.flusher.Flush(snippets)

	return err
}

func (s saver) Init() {
	go s.initSaverLoop()
}

func (s saver) Close() error {
	close(s.doneCh)
	return s.Save(s.snippets)
}

func (s saver) initSaverLoop() {
	ticker := time.NewTicker(time.Second * time.Duration(s.saveInterval))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := s.Save(s.snippets); err != nil {
				fmt.Printf("Error while saving snippets: %s", err.Error()) // TO DE FIXED: надо бы писать в лог...
			}
		case <-s.doneCh:
			return
		}
	}
}
