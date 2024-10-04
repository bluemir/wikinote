package metadata_test

import (
	"context"
	"testing"

	"github.com/bluemir/wikinote/internal/backend/metadata"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSave(t *testing.T) {
	ctx := context.Background()
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		t.Fatal(err)
		return
	}

	store, err := metadata.New(ctx, db)
	if err != nil {
		t.Fatal(err)
		return
	}

	if err := store.Save(ctx, "test", "foo", "bar"); err != nil {
		t.Fatal(err)
		return
	}

	v, err := store.Take(ctx, "test", "foo")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "bar", v)
}

func TestFind(t *testing.T) {

	ctx := context.Background()
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		t.Fatal(err)
		return
	}

	store, err := metadata.New(ctx, db)
	if err != nil {
		t.Fatal(err)
		return
	}

	data := []metadata.StoreItem{
		{"test-1", "a", "b"},
		{"test-1", "c", "d"},
		{"test-1", "e", "f"},
		{"test-2", "a", "b"},
		{"test-2", "g", "h"},
	}

	for _, item := range data {
		if err := store.Save(ctx, item.Path, item.Key, item.Value); err != nil {
			t.Fatal(err)
			return
		}
	}

	{
		items, err := store.FindByLabels(ctx, map[string]string{"a": "b"})
		if err != nil {
			t.Fatal(err)
			return
		}
		assert.Len(t, items, 2)
	}
	{
		items, err := store.FindByLabels(ctx, map[string]string{"c": "d"})
		if err != nil {
			t.Fatal(err)
			return
		}
		assert.Len(t, items, 1)
		assert.Equal(t, "test-1", items[0].Path)
	}
	{
		items, err := store.FindByLabels(ctx, map[string]string{"g": "h"})
		if err != nil {
			t.Fatal(err)
			return
		}
		assert.Len(t, items, 1)
		assert.Equal(t, "test-2", items[0].Path)
	}
}
