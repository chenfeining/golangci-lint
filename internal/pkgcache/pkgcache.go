package pkgcache

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"runtime"
	"sort"
	"sync"

	"golang.org/x/tools/go/packages"

	"github.com/chenfeining/golangci-lint/internal/cache"
	"github.com/chenfeining/golangci-lint/pkg/logutils"
	"github.com/chenfeining/golangci-lint/pkg/timeutils"
)

type HashMode int

const (
	HashModeNeedOnlySelf HashMode = iota
	HashModeNeedDirectDeps
	HashModeNeedAllDeps
)

// Cache is a per-package data cache. A cached data is invalidated when
// package, or it's dependencies change.
type Cache struct {
	lowLevelCache *cache.Cache
	pkgHashes     sync.Map
	sw            *timeutils.Stopwatch
	log           logutils.Log  // not used now, but may be needed for future debugging purposes
	ioSem         chan struct{} // semaphore limiting parallel IO
}

func NewCache(sw *timeutils.Stopwatch, log logutils.Log) (*Cache, error) {
	c, err := cache.Default()
	if err != nil {
		return nil, err
	}
	return &Cache{
		lowLevelCache: c,
		sw:            sw,
		log:           log,
		ioSem:         make(chan struct{}, runtime.GOMAXPROCS(-1)),
	}, nil
}

func (c *Cache) Trim() {
	c.sw.TrackStage("trim", func() {
		c.lowLevelCache.Trim()
	})
}

func (c *Cache) Put(pkg *packages.Package, mode HashMode, key string, data any) error {
	var err error
	buf := &bytes.Buffer{}
	c.sw.TrackStage("gob", func() {
		err = gob.NewEncoder(buf).Encode(data)
	})
	if err != nil {
		return fmt.Errorf("failed to gob encode: %w", err)
	}

	var aID cache.ActionID

	c.sw.TrackStage("key build", func() {
		aID, err = c.pkgActionID(pkg, mode)
		if err == nil {
			subkey, subkeyErr := cache.Subkey(aID, key)
			if subkeyErr != nil {
				err = fmt.Errorf("failed to build subkey: %w", subkeyErr)
			}
			aID = subkey
		}
	})
	if err != nil {
		return fmt.Errorf("failed to calculate package %s action id: %w", pkg.Name, err)
	}
	c.ioSem <- struct{}{}
	c.sw.TrackStage("cache io", func() {
		err = c.lowLevelCache.PutBytes(aID, buf.Bytes())
	})
	<-c.ioSem
	if err != nil {
		return fmt.Errorf("failed to save data to low-level cache by key %s for package %s: %w", key, pkg.Name, err)
	}

	return nil
}

var ErrMissing = errors.New("missing data")

func (c *Cache) Get(pkg *packages.Package, mode HashMode, key string, data any) error {
	var aID cache.ActionID
	var err error
	c.sw.TrackStage("key build", func() {
		aID, err = c.pkgActionID(pkg, mode)
		if err == nil {
			subkey, subkeyErr := cache.Subkey(aID, key)
			if subkeyErr != nil {
				err = fmt.Errorf("failed to build subkey: %w", subkeyErr)
			}
			aID = subkey
		}
	})
	if err != nil {
		return fmt.Errorf("failed to calculate package %s action id: %w", pkg.Name, err)
	}

	var b []byte
	c.ioSem <- struct{}{}
	c.sw.TrackStage("cache io", func() {
		b, _, err = c.lowLevelCache.GetBytes(aID)
	})
	<-c.ioSem
	if err != nil {
		if cache.IsErrMissing(err) {
			return ErrMissing
		}
		return fmt.Errorf("failed to get data from low-level cache by key %s for package %s: %w", key, pkg.Name, err)
	}

	c.sw.TrackStage("gob", func() {
		err = gob.NewDecoder(bytes.NewReader(b)).Decode(data)
	})
	if err != nil {
		return fmt.Errorf("failed to gob decode: %w", err)
	}

	return nil
}

func (c *Cache) pkgActionID(pkg *packages.Package, mode HashMode) (cache.ActionID, error) {
	hash, err := c.packageHash(pkg, mode)
	if err != nil {
		return cache.ActionID{}, fmt.Errorf("failed to get package hash: %w", err)
	}

	key, err := cache.NewHash("action ID")
	if err != nil {
		return cache.ActionID{}, fmt.Errorf("failed to make a hash: %w", err)
	}
	fmt.Fprintf(key, "pkgpath %s\n", pkg.PkgPath)
	fmt.Fprintf(key, "pkghash %s\n", hash)

	return key.Sum(), nil
}

// packageHash computes a package's hash. The hash is based on all Go
// files that make up the package, as well as the hashes of imported
// packages.
func (c *Cache) packageHash(pkg *packages.Package, mode HashMode) (string, error) {
	type hashResults map[HashMode]string
	hashResI, ok := c.pkgHashes.Load(pkg)
	if ok {
		hashRes := hashResI.(hashResults)
		if _, ok := hashRes[mode]; !ok {
			return "", fmt.Errorf("no mode %d in hash result", mode)
		}
		return hashRes[mode], nil
	}

	hashRes := hashResults{}

	key, err := cache.NewHash("package hash")
	if err != nil {
		return "", fmt.Errorf("failed to make a hash: %w", err)
	}

	fmt.Fprintf(key, "pkgpath %s\n", pkg.PkgPath)
	for _, f := range pkg.CompiledGoFiles {
		c.ioSem <- struct{}{}
		h, fErr := cache.FileHash(f)
		<-c.ioSem
		if fErr != nil {
			return "", fmt.Errorf("failed to calculate file %s hash: %w", f, fErr)
		}
		fmt.Fprintf(key, "file %s %x\n", f, h)
	}
	curSum := key.Sum()
	hashRes[HashModeNeedOnlySelf] = hex.EncodeToString(curSum[:])

	imps := make([]*packages.Package, 0, len(pkg.Imports))
	for _, imp := range pkg.Imports {
		imps = append(imps, imp)
	}
	sort.Slice(imps, func(i, j int) bool {
		return imps[i].PkgPath < imps[j].PkgPath
	})

	calcDepsHash := func(depMode HashMode) error {
		for _, dep := range imps {
			if dep.PkgPath == "unsafe" {
				continue
			}

			depHash, depErr := c.packageHash(dep, depMode)
			if depErr != nil {
				return fmt.Errorf("failed to calculate hash for dependency %s with mode %d: %w", dep.Name, depMode, depErr)
			}

			fmt.Fprintf(key, "import %s %s\n", dep.PkgPath, depHash)
		}
		return nil
	}

	if err := calcDepsHash(HashModeNeedOnlySelf); err != nil {
		return "", err
	}

	curSum = key.Sum()
	hashRes[HashModeNeedDirectDeps] = hex.EncodeToString(curSum[:])

	if err := calcDepsHash(HashModeNeedAllDeps); err != nil {
		return "", err
	}
	curSum = key.Sum()
	hashRes[HashModeNeedAllDeps] = hex.EncodeToString(curSum[:])

	if _, ok := hashRes[mode]; !ok {
		return "", fmt.Errorf("invalid mode %d", mode)
	}

	c.pkgHashes.Store(pkg, hashRes)
	return hashRes[mode], nil
}
