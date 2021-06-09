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
	Close()
}

type saver struct {
	flusher  flusher.Flusher
	snippets []models.Snippet

	saveInterval time.Duration
	doneCh       chan struct{}
}

// NewSaver возвращает Save с поддержкой периодического сохранения.
// capacity - capacity слайса сохраняемых сущностей.
// flusher - экземпляр интерфейса flusher.Flusher, осуществляющий сохранения.
// saveInterval - интервал сохранения
func NewSaver(capacity uint, flusher flusher.Flusher, saveInterval time.Duration) Saver {
	return &saver{
		flusher:      flusher,
		snippets:     make([]models.Snippet, 0, capacity),
		saveInterval: saveInterval,
	}
}

func (s saver) Save(snippets []models.Snippet) error {
	resSnippets, err := s.flusher.Flush(snippets)
	s.snippets = resSnippets

	return err
}

func (s saver) Init() {
	go s.initSaverLoop()
}

func (s saver) Close() {
	close(s.doneCh)
}

func (s saver) initSaverLoop() {
	ticker := time.NewTicker(s.saveInterval)
	defer ticker.Stop()

	save := func() {
		if err := s.Save(s.snippets); err != nil {
			fmt.Printf("Error while saving snippets: %s", err.Error()) // TO DE FIXED: надо бы писать в лог...
		}
	}
	for {
		select {
		case <-ticker.C:
			save()
		case <-s.doneCh:
			save()
			return
		}
	}
}
