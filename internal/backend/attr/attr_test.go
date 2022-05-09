package attr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSimpleSave(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	store, err := New(db)
	if err != nil {
		t.Fatal(err)
	}

	if err := store.Save(&FileAttr{
		Path:  "/dummy",
		Key:   "foo",
		Value: "bar",
	}); err != nil {
		t.Fatal(err)
	}

	attr, err := store.Take(&FileAttr{
		Path: "/dummy",
		Key:  "foo",
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &FileAttr{
		Path:  "/dummy",
		Key:   "foo",
		Value: "bar",
	}, attr)
}
func TestSimpleNotFound(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	store, err := New(db)
	if err != nil {
		t.Fatal(err)
	}

	if err := store.Save(&FileAttr{
		Path:  "/path1",
		Key:   "foo",
		Value: "bar",
	}); err != nil {
		t.Fatal(err)
	}

	_, err = store.Take(&FileAttr{
		Path: "/path2",
		Key:  "foo",
	})
	assert.True(t, IsNotFound(err))

	_, err = store.Take(&FileAttr{
		Path: "/path1",
		Key:  "bar",
	})
	assert.True(t, IsNotFound(err))
}