package cache

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

var (
	db   *redis.Client
	mock redismock.ClientMock
	rds  *redisService
	key  = "key"
	val  = "val"
	exp  = time.Duration(0)
)

func TestMain(m *testing.M) {
	ctx = context.TODO()
	db, mock = redismock.NewClientMock()
	rds = &redisService{instance: db}

	os.Exit(m.Run())
}

func TestSetFails(t *testing.T) {
	mock.ExpectSet(key, val, exp).SetErr(errors.New("FAIL"))

	err := rds.Set(key, val, exp)

	if err == nil {
		t.Error("must return an error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestSetSuccess(t *testing.T) {
	mock.ExpectSet(key, val, exp).SetVal(val)

	err := rds.Set(key, val, exp)

	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestGetFails(t *testing.T) {
	mock.ExpectGet(key).SetErr(errors.New("no value found by provided key"))

	_, err := rds.Get(key)

	if err == nil {
		t.Error("must return an error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestGetSuccess(t *testing.T) {
	mock.ExpectSet(key, val, exp).SetVal(val)

	err := rds.Set(key, val, exp)

	if err != nil {
		t.Error(err)
	}

	mock.ExpectGet(key).SetVal(val)

	actual, err := rds.Get(key)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, val, actual)
}

func TestDelFails(t *testing.T) {
	mock.ExpectDel(key).SetErr(errors.New("key not found"))

	err := rds.Delete(key)

	if err == nil {
		t.Error("must return an error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestDelSuccess(t *testing.T) {
	mock.ExpectSet(key, val, exp).SetVal(val)

	err := rds.Set(key, val, exp)

	if err != nil {
		t.Error(err)
	}

	mock.ExpectDel(key).SetVal(1)

	err = rds.Delete(key)

	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
