package service

import (
	"path"
	"strings"

	"github.com/zwh8800/md-blog-gen/util"
)

const (
	IndexSetKey   = "indices:"
	TagSetKey     = "tags:"
	NoteSetKey    = "notes:"
	ArchiveSetKey = "archives:"
	SearchSetKey  = "searches:"
	CacheKey      = "cache:"
)

func AddCache(p, data string) error {
	key := CacheKey + p
	switch {
	case strings.Index(p, util.GetPageBase()) == 0:
		if err := redisClient.SAdd(IndexSetKey, key).Err(); err != nil {
			return err
		}
	case strings.Index(p, util.GetTagBase()) == 0:
		if err := redisClient.SAdd(TagSetKey, key).Err(); err != nil {
			return err
		}
	case strings.Index(p, util.GetNoteBase()) == 0:
		if err := redisClient.SAdd(NoteSetKey, key).Err(); err != nil {
			return err
		}
	case strings.Index(p, util.GetArchiveBase()) == 0:
		if err := redisClient.SAdd(ArchiveSetKey, key).Err(); err != nil {
			return err
		}
	case strings.Index(p, util.GetSearchBase()) == 0:
	case strings.Index(p, path.Join("/api", util.GetSearchBase())) == 0:
		if err := redisClient.SAdd(SearchSetKey, key).Err(); err != nil {
			return err
		}
	}
	if err := redisClient.Set(key, data, 0).Err(); err != nil {
		return err
	}
	return nil
}

func FindCache(path string) (string, error) {
	return redisClient.Get(CacheKey + path).Result()
}

func RemoveIndexCache() error {
	const homePage = CacheKey + "/"
	if err := redisClient.Del(homePage).Err(); err != nil {
		return err
	}
	return deleteAll(IndexSetKey)
}

func RemoveTagCache(tag string) error {
	key := CacheKey + util.GetTagBase()
	if err := redisClient.Del(key).Err(); err != nil {
		return err
	}
	key = CacheKey + util.GetTagNameUrl(tag)
	if err := redisClient.Del(key).Err(); err != nil {
		return err
	}
	return redisClient.SRem(TagSetKey, key).Err()
}

func RemoveAllTagCache() error {
	key := CacheKey + util.GetTagBase()
	if err := redisClient.Del(key).Err(); err != nil {
		return err
	}
	return deleteAll(TagSetKey)
}

func RemoveNoteCache(notename string) error {
	key := CacheKey + util.GetNoteUrlByNotename(notename)
	if err := redisClient.Del(key).Err(); err != nil {
		return err
	}
	return redisClient.SRem(TagSetKey, key).Err()
}

func RemoveArchiveCache(year, month int64) error {
	key := CacheKey + util.GetArchiveBase()
	if err := redisClient.Del(key).Err(); err != nil {
		return err
	}
	key = CacheKey + util.GetArchiveMonthUrl(year, month)
	if err := redisClient.Del(key).Err(); err != nil {
		return err
	}
	return redisClient.SRem(TagSetKey, key).Err()
}

func RemoveSearchCache() error {
	return deleteAll(SearchSetKey)
}

func deleteAll(key string) error {
	for {
		count, err := redisClient.SCard(key).Result()
		if err != nil {
			return err
		}
		if count == 0 {
			return nil
		}
		key, err := redisClient.SPop(key).Result()
		if err != nil {
			return err
		}
		if err := redisClient.Del(key).Err(); err != nil {
			return err
		}
	}
}

func RemoveAllCache() {
	redisClient.FlushAll()
}
