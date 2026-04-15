package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ViitoJooj/door/internal/domain"
	"github.com/cockroachdb/pebble"
)

func (s *Store) nextID(prefix string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := []byte("meta:" + prefix + ":seq")

	val, closer, err := s.db.Get(key)
	if closer != nil {
		defer closer.Close()
	}

	id := 1
	if errors.Is(err, pebble.ErrNotFound) {
		// first record, id stays at 1
	} else if err != nil {
		return 0, err
	} else {
		id, err = strconv.Atoi(string(val))
		if err != nil {
			return 0, err
		}
		id++
	}

	if err := s.db.Set(key, []byte(strconv.Itoa(id)), pebble.NoSync); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Store) CreateUser(user *domain.User) error {
	id, err := s.nextID("user")
	if err != nil {
		return err
	}
	user.ID = id

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	batch := s.db.NewBatch()
	defer batch.Close()

	batch.Set([]byte(fmt.Sprintf("user:id:%010d", id)), data, nil)
	batch.Set([]byte("user:email:"+user.Email), []byte(strconv.Itoa(id)), nil)
	batch.Set([]byte("user:username:"+user.Username), []byte(strconv.Itoa(id)), nil)

	return batch.Commit(pebble.NoSync)
}

func (s *Store) CreateApplication(application *domain.Application) error {
	id, err := s.nextID("app")
	if err != nil {
		return err
	}
	application.ID = id

	data, err := json.Marshal(application)
	if err != nil {
		return err
	}

	batch := s.db.NewBatch()
	defer batch.Close()

	batch.Set([]byte(fmt.Sprintf("app:id:%010d", id)), data, nil)
	batch.Set([]byte("app:url:"+application.Url), []byte(strconv.Itoa(id)), nil)
	batch.Set([]byte("app:country:"+application.Country), []byte(strconv.Itoa(id)), nil)

	return batch.Commit(pebble.NoSync)
}

func (s *Store) InsertRequestLog(log *domain.RequestLog) error {
	id, err := s.nextID("log")
	if err != nil {
		return err
	}
	log.ID = id
	log.CreatedAt = time.Now()

	data, err := json.Marshal(log)
	if err != nil {
		return err
	}

	return s.db.Set([]byte(fmt.Sprintf("log:id:%010d", id)), data, pebble.NoSync)
}
