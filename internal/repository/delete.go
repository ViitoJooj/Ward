package repository

import (
	"fmt"

	"github.com/cockroachdb/pebble"
)

func (s *Store) DeleteUserByID(id int) error {
	old, err := s.FindUserByID(id)
	if err != nil {
		return err
	}

	batch := s.db.NewBatch()
	defer batch.Close()

	batch.Delete([]byte(fmt.Sprintf("user:id:%010d", id)), nil)
	if old != nil {
		batch.Delete([]byte("user:email:"+old.Email), nil)
		batch.Delete([]byte("user:username:"+old.Username), nil)
	}

	return batch.Commit(pebble.NoSync)
}

func (s *Store) DeleteApplicationByID(id int) error {
	old, err := s.FindApplicationByID(id)
	if err != nil {
		return err
	}

	batch := s.db.NewBatch()
	defer batch.Close()

	batch.Delete([]byte(fmt.Sprintf("app:id:%010d", id)), nil)
	if old != nil {
		batch.Delete([]byte("app:url:"+old.Url), nil)
		batch.Delete([]byte("app:country:"+old.Country), nil)
	}

	return batch.Commit(pebble.NoSync)
}
