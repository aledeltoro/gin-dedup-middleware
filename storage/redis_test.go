package storage

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
)

func TestAddSet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client, mock := redismock.NewClientMock()

		key := "myset"
		member := "hello"
		expire := 30 * time.Second

		mock.MatchExpectationsInOrder(true)
		mock.ExpectSAdd(key, member).SetVal(1)
		mock.ExpectExpireNX(key, expire).SetVal(true)

		cacheStorage := &RedisCache{
			client: client,
		}

		err := cacheStorage.AddSet(context.Background(), key, expire, member)
		if err != nil {
			t.Error(err)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("no new set members", func(t *testing.T) {
		client, mock := redismock.NewClientMock()

		key := "myset"
		member := "hello"
		expire := 30 * time.Second

		mock.MatchExpectationsInOrder(true)
		mock.ExpectSAdd(key, member).SetVal(0)

		cacheStorage := &RedisCache{
			client: client,
		}

		err := cacheStorage.AddSet(context.Background(), key, expire, member)
		if err != nil {
			t.Error(err)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("connection error", func(t *testing.T) {
		client, mock := redismock.NewClientMock()

		key := "myset"
		member := "hello"
		expire := 30 * time.Second

		mock.MatchExpectationsInOrder(true)
		mock.ExpectSAdd(key, member).SetErr(errors.New("connection error"))

		cacheStorage := &RedisCache{
			client: client,
		}

		err := cacheStorage.AddSet(context.Background(), key, expire, member)
		if err != nil && err.Error() != "connection error" {
			t.Errorf("unexpected error received: %s", err.Error())
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Error(err)
		}
	})
}

func TestIsSetMember(t *testing.T) {
	client, mock := redismock.NewClientMock()

	key := "myset"
	member := "hello"

	mock.MatchExpectationsInOrder(true)
	mock.ExpectSIsMember(key, member).SetVal(true)

	cacheStorage := &RedisCache{
		client: client,
	}

	isMember, err := cacheStorage.IsSetMember(context.Background(), key, member)
	if err != nil && !isMember {
		t.Error(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
