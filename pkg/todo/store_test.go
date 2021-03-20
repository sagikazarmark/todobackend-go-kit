package todo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryStore_StoreAnItem(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		store := NewInMemoryStore()

		item := Item{
			ID:    "id",
			Title: "Store me!",
		}

		err := store.Store(context.Background(), item)
		require.NoError(t, err)

		assert.Equal(t, item, store.items[item.ID])
	})

	t.Run("OverwriteAnExistingItem", func(t *testing.T) {
		t.Parallel()

		store := NewInMemoryStore()

		store.items["id"] = Item{
			ID:        "id",
			Title:     "Store me first!",
			Completed: true,
		}

		item := Item{
			ID:    "id",
			Title: "Store me!",
		}

		err := store.Store(context.Background(), item)
		require.NoError(t, err)

		assert.Equal(t, item, store.items[item.ID])
	})
}

func TestInMemoryStore_ListAllItems(t *testing.T) {
	store := NewInMemoryStore()

	store.items["id"] = Item{
		ID:        "id",
		Title:     "Store me first!",
		Completed: true,
	}

	store.items["id2"] = Item{
		ID:        "id2",
		Title:     "Store me second!",
		Completed: true,
	}

	items, err := store.GetAll(context.Background())
	require.NoError(t, err)

	assert.Equal(t, []Item{store.items["id"], store.items["id2"]}, items)
}

func TestInMemoryStore_DeleteAllItems(t *testing.T) {
	store := NewInMemoryStore()

	store.items["id"] = Item{
		ID:        "id",
		Title:     "Store me first!",
		Completed: true,
	}

	store.items["id2"] = Item{
		ID:        "id2",
		Title:     "Store me second!",
		Completed: true,
	}

	err := store.DeleteAll(context.Background())
	require.NoError(t, err)

	assert.Len(t, store.items, 0)
}

func TestInMemoryStore_GetAnItem(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		store := NewInMemoryStore()

		const id = "id"

		store.items[id] = Item{
			ID:    id,
			Title: "Store me!",
		}

		item, err := store.GetOne(context.Background(), id)
		require.NoError(t, err)

		assert.Equal(t, store.items[id], item)
	})

	t.Run("NotFound", func(t *testing.T) {
		t.Parallel()

		store := NewInMemoryStore()

		_, err := store.GetOne(context.Background(), "id")
		require.Error(t, err)

		require.IsType(t, NotFoundError{}, err)

		e := err.(NotFoundError) // nolint:errorlint
		assert.Equal(t, "id", e.ID)
	})
}

func TestInMemoryStore_DeleteAnItem(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		store := NewInMemoryStore()

		store.items["id"] = Item{
			ID:        "id",
			Title:     "Store me first!",
			Completed: true,
		}

		store.items["id2"] = Item{
			ID:        "id2",
			Title:     "Store me second!",
			Completed: true,
		}

		err := store.DeleteOne(context.Background(), "id2")
		require.NoError(t, err)

		assert.Len(t, store.items, 1)
		assert.Contains(t, store.items, "id")
	})

	t.Run("Idempotent", func(t *testing.T) {
		t.Parallel()

		store := NewInMemoryStore()

		store.items["id"] = Item{
			ID:        "id",
			Title:     "Store me first!",
			Completed: true,
		}

		err := store.DeleteOne(context.Background(), "id2")
		require.NoError(t, err)

		assert.Len(t, store.items, 1)
		assert.Contains(t, store.items, "id")
	})
}
