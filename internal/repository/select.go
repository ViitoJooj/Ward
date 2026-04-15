package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/ViitoJooj/door/internal/domain"
	"github.com/cockroachdb/pebble"
)

func upperBound(prefix []byte) []byte {
	end := make([]byte, len(prefix))
	copy(end, prefix)
	for i := len(end) - 1; i >= 0; i-- {
		end[i]++
		if end[i] != 0 {
			return end[:i+1]
		}
	}
	return nil
}

func (s *Store) getIDByIndex(indexKey string) (int, error) {
	val, closer, err := s.db.Get([]byte(indexKey))
	if closer != nil {
		defer closer.Close()
	}
	if errors.Is(err, pebble.ErrNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(val))
}

func (s *Store) FindUserByID(id int) (*domain.User, error) {
	val, closer, err := s.db.Get([]byte(fmt.Sprintf("user:id:%010d", id)))
	if closer != nil {
		defer closer.Close()
	}
	if errors.Is(err, pebble.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	user := &domain.User{}
	if err := json.Unmarshal(val, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) FindUserByEmail(email string) (*domain.User, error) {
	id, err := s.getIDByIndex("user:email:" + email)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, nil
	}
	return s.FindUserByID(id)
}

func (s *Store) FindUserByUsername(username string) (*domain.User, error) {
	id, err := s.getIDByIndex("user:username:" + username)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, nil
	}
	return s.FindUserByID(id)
}

func (s *Store) ListUsers() ([]*domain.User, error) {
	prefix := []byte("user:id:")
	iter, err := s.db.NewIter(&pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: upperBound(prefix),
	})
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var users []*domain.User
	for iter.First(); iter.Valid(); iter.Next() {
		user := &domain.User{}
		if err := json.Unmarshal(iter.Value(), user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, iter.Error()
}

func (s *Store) FindApplicationByID(id int) (*domain.Application, error) {
	val, closer, err := s.db.Get([]byte(fmt.Sprintf("app:id:%010d", id)))
	if closer != nil {
		defer closer.Close()
	}
	if errors.Is(err, pebble.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	app := &domain.Application{}
	if err := json.Unmarshal(val, app); err != nil {
		return nil, err
	}

	return app, nil
}

func (s *Store) FindApplicationByURL(url string) (*domain.Application, error) {
	id, err := s.getIDByIndex("app:url:" + url)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, nil
	}
	return s.FindApplicationByID(id)
}

func (s *Store) FindApplicationByCountry(country string) (*domain.Application, error) {
	id, err := s.getIDByIndex("app:country:" + country)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, nil
	}
	return s.FindApplicationByID(id)
}

func (s *Store) ListApplications() ([]*domain.Application, error) {
	prefix := []byte("app:id:")
	iter, err := s.db.NewIter(&pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: upperBound(prefix),
	})
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var applications []*domain.Application
	for iter.First(); iter.Valid(); iter.Next() {
		app := &domain.Application{}
		if err := json.Unmarshal(iter.Value(), app); err != nil {
			return nil, err
		}
		applications = append(applications, app)
	}

	return applications, iter.Error()
}

func (s *Store) ListRequestLogs() ([]*domain.RequestLog, error) {
	prefix := []byte("log:id:")
	iter, err := s.db.NewIter(&pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: upperBound(prefix),
	})
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var logs []*domain.RequestLog
	for iter.First(); iter.Valid(); iter.Next() {
		entry := &domain.RequestLog{}
		if err := json.Unmarshal(iter.Value(), entry); err != nil {
			return nil, err
		}
		logs = append(logs, entry)
	}

	if err := iter.Error(); err != nil {
		return nil, err
	}

	sort.Slice(logs, func(i, j int) bool {
		return logs[i].CreatedAt.After(logs[j].CreatedAt)
	})

	return logs, nil
}
