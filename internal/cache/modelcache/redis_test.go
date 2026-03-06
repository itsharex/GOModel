package modelcache

import (
	"context"
	"testing"
	"time"

	"gomodel/internal/cache"
)

func TestRedisModelCache_GetSet(t *testing.T) {
	store := cache.NewMapStore()
	defer store.Close()
	c := NewRedisModelCacheWithStore(store, "test:models", time.Hour)
	defer c.Close()

	ctx := context.Background()
	got, err := c.Get(ctx)
	if err != nil {
		t.Fatalf("Get empty: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil for empty cache, got %v", got)
	}

	mc := &ModelCache{
		UpdatedAt: time.Now(),
		Providers: map[string]CachedProvider{
			"openai": {
				ProviderType: "openai",
				OwnedBy:      "openai",
				Models: []CachedModel{
					{ID: "gpt-4", Created: 123},
				},
			},
		},
	}
	if err := c.Set(ctx, mc); err != nil {
		t.Fatalf("Set: %v", err)
	}

	got, err = c.Get(ctx)
	if err != nil {
		t.Fatalf("Get after Set: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil ModelCache")
	}
	if len(got.Providers) != 1 {
		t.Errorf("Providers: got %d entries, want 1", len(got.Providers))
	}
	p, ok := got.Providers["openai"]
	if !ok {
		t.Fatal("expected openai in Providers")
	}
	if p.ProviderType != "openai" {
		t.Errorf("ProviderType: got %s, want openai", p.ProviderType)
	}
	if len(p.Models) != 1 {
		t.Errorf("Models: got %d entries, want 1", len(p.Models))
	}
	if p.Models[0].ID != "gpt-4" {
		t.Errorf("Model ID: got %s, want gpt-4", p.Models[0].ID)
	}
}

func TestRedisModelCache_DefaultKeyAndTTL(t *testing.T) {
	store := cache.NewMapStore()
	defer store.Close()
	c := NewRedisModelCacheWithStore(store, "", 0)
	defer c.Close()

	ctx := context.Background()
	mc := &ModelCache{
		UpdatedAt: time.Now(),
		Providers: map[string]CachedProvider{},
	}
	if err := c.Set(ctx, mc); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := c.Get(ctx)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil ModelCache")
	}
}
