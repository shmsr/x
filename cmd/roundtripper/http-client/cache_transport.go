package main

import (
	"bufio"
	"bytes"
	"errors"
	"net/http"
	"net/http/httputil"
	"sync"
)

const defaultCacheSize = 1000

var (
	ErrEmptyCache = errors.New("cache is empty")
	ErrInitCache  = errors.New("cache is not initialized")
)

type cacheTransport struct {
	// rt is the original RoundTripper
	rt http.RoundTripper

	cache struct {
		// mu protects data map from concurrent access/ modifications
		mu sync.RWMutex
		// data hold req.URL as key and response as value
		data map[string]string
	}
}

// getCacheTransport returns a pointer to cacheTransport
func getCacheTransport(size ...int64) *cacheTransport {
	c := new(cacheTransport)
	c.rt = http.DefaultTransport
	if len(size) == 0 {
		c.initCache(defaultCacheSize)
		return c
	}
	c.initCache(size[0])
	return c
}

// getRequestURL returns the URL to acess (client requests)
func getRequestURL(r *http.Request) string {
	if r != nil {
		return r.URL.String()
	}
	return ""
}

func getCachedResponse(b []byte, req *http.Request) (*http.Response, error) {
	// The req parameter optionally specifies the Request that corresponds to this
	// Response. If nil, a GET request is assumed.
	buf := bytes.NewBuffer(b)
	return http.ReadResponse(bufio.NewReader(buf), req)
}

// initCache initializes cache with cacheSize len
func (c *cacheTransport) initCache(size int64) {
	c.cache.data = make(map[string]string, size)
}

// Set makes a entry to the cache
func (c *cacheTransport) Set(req *http.Request, value string) error {
	c.cache.mu.Lock()
	defer c.cache.mu.Unlock()
	if c.cache.data == nil {
		return ErrInitCache
	}
	c.cache.data[getRequestURL(req)] = value
	return nil
}

// Get fetches a entry from cache, if available
func (c *cacheTransport) Get(req *http.Request) (string, error) {
	c.cache.mu.RLock()
	v, ok := c.cache.data[getRequestURL(req)]
	c.cache.mu.RUnlock()
	if ok {
		return v, nil
	}
	return "", ErrEmptyCache
}

// Detaches from older references and points to newly allocated map
// GC will cleanup the older cache which isn't referenced anymore
func (c *cacheTransport) Clear(size ...int64) {
	if len(size) == 0 {
		c.initCache(defaultCacheSize)
		return
	}
	c.initCache(size[0])
}

// RoundTripper interface should implement RoundTrip method
// type RoundTripper interface{
// 	RoundTrip(*Request) (*Response, error)
// }

// compile-time safety check
var _ http.RoundTripper = (*cacheTransport)(nil)

// RoundTrip first tries the cache, and if the response is not cached,
// the request is relayed to server, and if successful the response is then
// added to cache
func (c *cacheTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if v, err := c.Get(req); err == nil {
		return getCachedResponse([]byte(v), req)
	}
	resp, err := c.rt.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	buf, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}
	if err := c.Set(req, string(buf)); err != nil {
		// In case of error adding entry to cache; response is successfully
		// returned but with an error
		return resp, err
	}
	return resp, nil
}
