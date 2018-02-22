package pip

import (
	"testing"
)

func newCache(t *testing.T, size int) *QueryCache {
	qc, err := NewQueryCache(size)
	if err != nil {
		t.Fatalf("Error creating query cache of size %d: %s", size, err)
	}

	return qc
}

func TestCacheGet(t *testing.T) {
	qc := newCache(t, 100)

	domain := "cnn.com"
	categories := "News"
	qc.Add(domain, categories)

	res, ok := qc.Get(domain)
	if !ok {
		t.Errorf("Expected ok is '%T', but got '%T'", true, false)
	} else if res != categories {
		t.Errorf("Expected res is '%s', but got '%s'", categories, res)
	}
}

func TestBadCacheSize(t *testing.T) {
	_, err := NewQueryCache(0)

	if err == nil {
		t.Errorf("Expected NewQueryCache(0) returns err, but got nil")
	}
}

func TestCachePurge(t *testing.T) {
	qc := newCache(t, 100)

	domain := "cnn.com"
	categories := "News"
	qc.Add(domain, categories)

	qcLen := qc.Len()
	if qcLen != 1 {
		t.Errorf("Expected len is '%d', but got '%d'", 1, qcLen)
	}

	qc.Purge()

	qcLen = qc.Len()
	if qcLen != 0 {
		t.Errorf("Expected len is '%d', but got '%d'", 0, qcLen)
	}

}

type testEntry struct {
	domain, categories string
}

var (
	testEntry1 = testEntry{domain: "cnn.com", categories: "News"}
	testEntry2 = testEntry{domain: "google.com", categories: "Search Engine"}
	testEntry3 = testEntry{domain: "espn.com", categories: "Sports"}
)

func TestCacheEvict1(t *testing.T) {
	testEntries := []testEntry{
		testEntry1, testEntry1, testEntry1,
		testEntry2, testEntry2, testEntry2, testEntry2,
		testEntry3,
	}

	const cacheSize = 2
	qc := newCache(t, cacheSize)

	for _, e := range testEntries {
		if _, ok := qc.Get(e.domain); !ok {
			qc.Add(e.domain, e.categories)
		}
	}

	qcLen := qc.Len()
	if qcLen != cacheSize {
		t.Errorf("Expected len is '%d', but got '%d'", cacheSize, qcLen)
	}

	if _, ok := qc.Get(testEntry1.domain); ok {
		t.Errorf("Expect testEntry1 to be evicted but got result")
	}

	if val, ok := qc.Get(testEntry2.domain); !ok {
		t.Errorf("Expect testEntry2 to be in cache but got no result")
	} else {
		if val.(string) != testEntry2.categories {
			t.Errorf("Expect value in cache for testEntry2 to be '%s' but got '%s'", testEntry2.categories, val.(string))
		}
	}

	if val, ok := qc.Get(testEntry3.domain); !ok {
		t.Errorf("Expect testEntry3 to be in cache but got no result")
	} else {
		if val.(string) != testEntry3.categories {
			t.Errorf("Expect value in cache for testEntry3 to be '%s' but got '%s'", testEntry3.categories, val.(string))
		}
	}
}

func TestCacheEvict2(t *testing.T) {
	testEntries := []testEntry{
		testEntry1, testEntry1, testEntry1,
		testEntry2, testEntry2, testEntry2, testEntry2,
		testEntry3, testEntry3, testEntry1,
	}

	const cacheSize = 2
	qc := newCache(t, cacheSize)

	for _, e := range testEntries {
		if _, ok := qc.Get(e.domain); !ok {
			qc.Add(e.domain, e.categories)
		}
	}

	qcLen := qc.Len()
	if qcLen != cacheSize {
		t.Errorf("Expected len is '%d', but got '%d'", cacheSize, qcLen)
	}

	if val, ok := qc.Get(testEntry1.domain); !ok {
		t.Errorf("Expect testEntry1 to be in cache but got no result")
	} else {
		if val.(string) != testEntry1.categories {
			t.Errorf("Expect value in cache for testEntry2 to be '%s' but got '%s'", testEntry1.categories, val.(string))
		}
	}

	if _, ok := qc.Get(testEntry2.domain); ok {
		t.Errorf("Expect testEntry2 to be evicted but got result")
	}

	if val, ok := qc.Get(testEntry3.domain); !ok {
		t.Errorf("Expect testEntry3 to be in cache but got no result")
	} else {
		if val.(string) != testEntry3.categories {
			t.Errorf("Expect value in cache for testEntry3 to be '%s' but got '%s'", testEntry3.categories, val.(string))
		}
	}
}
