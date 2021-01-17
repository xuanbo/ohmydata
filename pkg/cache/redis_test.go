package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/xuanbo/ohmydata/pkg/cache"
	"github.com/xuanbo/ohmydata/pkg/entity"
)

func init() {
	// 初始化缓存
	if err := cache.Init(); err != nil {
		panic(err)
	}
}

func TestCache(t *testing.T) {
	user := &entity.User{
		ID:       "1",
		Name:     "1",
		Username: "1",
		Password: "1",
	}
	if err := cache.Set(context.TODO(), "test:1", user, time.Minute); err != nil {
		t.Error(err)
		return
	}
	var u entity.User
	if err := cache.Get(context.TODO(), "test:1", &u); err != nil {
		t.Error(err)
		return
	}
	t.Logf("user: %v", u)

	if err := cache.Set(context.TODO(), "test", "1", time.Minute); err != nil {
		t.Error(err)
		return
	}
	var v interface{}
	if err := cache.Get(context.TODO(), "test", &v); err != nil {
		t.Error(err)
		return
	}
	t.Logf("s: %v", v)

	if err := cache.Get(context.TODO(), "test:2", &u); err != nil {
		t.Error(err)
		return
	}
	t.Logf("user: %v", u)
}

func TestDelEntries(t *testing.T) {
	if err := cache.DelMatch(context.TODO(), "id*"); err != nil {
		t.Error(err)
		return
	}
}
