package eBay

import (
        "testing"
        "time"
        "fmt"
)

func TestCacheCategory(t *testing.T) {
        c, err := NewDBCache("localhost", "test", "abc", "dusell")
        if err != nil { t.Errorf("NewDBCache: %v", err); return }

        cat := &Category{
        CategoryId: "123456",
        CategoryName: fmt.Sprintf("test category: %v", time.Nanoseconds()),
        }
        err = c.CacheCategory(cat)
        if err != nil { t.Errorf("c.CacheCategory: %v", err); return }

        cat2, err := c.GetCategory(cat.CategoryId)
        if err != nil { t.Errorf("c.GetCategory(%s): %v", cat.CategoryId, err); return }
        if cat2 == nil { t.Errorf("c.GetCategory(%s): no category returned", cat.CategoryId); return }
        if cat.CategoryId != cat2.CategoryId { t.Errorf("wrong category: %v != %v", cat, cat2); return }
        if cat.CategoryName != cat2.CategoryName { t.Errorf("wrong category: %v != %v", cat, cat2); return }
}

func TestCacheItem(t *testing.T) {
        
}
