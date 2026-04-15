package repository

import (
	"encoding/json"
	"fmt"

	"github.com/ViitoJooj/door/internal/domain"
	"github.com/cockroachdb/pebble"
)

func (s *Store) UpdateUser(user *domain.User) error {
	old, err := s.FindUserByID(user.ID)
	if err != nil {
		return err
	}

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	batch := s.db.NewBatch()
	defer batch.Close()

	if old != nil && old.Email != user.Email {
		batch.Delete([]byte("user:email:"+old.Email), nil)
	}
	if old != nil && old.Username != user.Username {
		batch.Delete([]byte("user:username:"+old.Username), nil)
	}

	batch.Set([]byte(fmt.Sprintf("user:id:%010d", user.ID)), data, nil)
	batch.Set([]byte("user:email:"+user.Email), []byte(fmt.Sprintf("%d", user.ID)), nil)
	batch.Set([]byte("user:username:"+user.Username), []byte(fmt.Sprintf("%d", user.ID)), nil)

	return batch.Commit(pebble.NoSync)
}

func (s *Store) UpdateApplication(application *domain.Application) error {
	old, err := s.FindApplicationByID(application.ID)
	if err != nil {
		return err
	}

	data, err := json.Marshal(application)
	if err != nil {
		return err
	}

	batch := s.db.NewBatch()
	defer batch.Close()

	if old != nil && old.Url != application.Url {
		batch.Delete([]byte("app:url:"+old.Url), nil)
	}
	if old != nil && old.Country != application.Country {
		batch.Delete([]byte("app:country:"+old.Country), nil)
	}

	batch.Set([]byte(fmt.Sprintf("app:id:%010d", application.ID)), data, nil)
	batch.Set([]byte("app:url:"+application.Url), []byte(fmt.Sprintf("%d", application.ID)), nil)
	batch.Set([]byte("app:country:"+application.Country), []byte(fmt.Sprintf("%d", application.ID)), nil)

	return batch.Commit(pebble.NoSync)
}
