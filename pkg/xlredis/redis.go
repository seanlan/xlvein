package xlredis

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	incrByExpScript = `local current = redis.call("INCRBYFLOAT", KEYS[1], ARGV[1])
if (ARGV[1] - current == 0)
then
    if redis.call("expire", KEYS[1], ARGV[2]) == 0
    then
        return -1
	else
		return current
    end
else
	return current
end
`
)

type Redis struct {
	prefix string
	Client redis.UniversalClient
}

func NewClient(uri, username, password, prefix string, db int) (client *Redis, err error) {
	addrs := strings.Split(uri, ",")
	opt := &redis.UniversalOptions{
		Addrs:    addrs,
		DB:       db,
		Username: username,
		Password: password,
	}
	c := redis.NewUniversalClient(opt)
	return &Redis{prefix: prefix, Client: c}, c.Ping(context.TODO()).Err()
}

func (i *Redis) BuildKey(key string) string {
	if len(i.prefix) > 0 {
		return i.prefix + ":" + key
	}
	return key
}

// GetLock 获取锁(redis)
func (i *Redis) GetLock(ctx context.Context, key string, expiration time.Duration) bool {
	return i.Client.SetNX(ctx, i.BuildKey(key), 1, expiration).Val()
}

//ReleaseLock 释放锁(redis)
func (i *Redis) ReleaseLock(ctx context.Context, key string) bool {
	return i.Client.Del(ctx, i.BuildKey(key)).Val() > 0
}

func (i *Redis) Do(ctx context.Context, args ...interface{}) *redis.Cmd {
	return i.Client.Do(ctx, args...)
}

func (i *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return i.Client.Set(ctx, i.BuildKey(key), value, expiration)
}

func (i *Redis) Get(ctx context.Context, key string) *redis.StringCmd {
	return i.Client.Get(ctx, i.BuildKey(key))
}

func (i *Redis) Del(ctx context.Context, key string) *redis.IntCmd {
	return i.Client.Del(ctx, i.BuildKey(key))
}

func (i *Redis) Exists(ctx context.Context, key string) *redis.IntCmd {
	return i.Client.Exists(ctx, i.BuildKey(key))
}

func (i *Redis) Incr(ctx context.Context, key string) *redis.IntCmd {
	return i.Client.Incr(ctx, i.BuildKey(key))
}

func (i *Redis) Decr(ctx context.Context, key string) *redis.IntCmd {
	return i.Client.Decr(ctx, i.BuildKey(key))
}

func (i *Redis) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	return i.Client.IncrBy(ctx, i.BuildKey(key), value)
}

// IncrByExpire 首次添加则设置过期时间
func (i *Redis) IncrByExpire(ctx context.Context, key string, value float64, expire time.Duration) (int64, error) {
	return i.Client.Eval(ctx, incrByExpScript, []string{i.BuildKey(key)}, value, expire.Seconds()).Int64()
}

func (i *Redis) DecrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	return i.Client.DecrBy(ctx, i.BuildKey(key), value)
}

func (i *Redis) IncrByFloat(ctx context.Context, key string, value float64) *redis.FloatCmd {
	return i.Client.IncrByFloat(ctx, i.BuildKey(key), value)
}

func (i *Redis) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return i.Client.Expire(ctx, i.BuildKey(key), expiration)
}

func (i *Redis) ExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	return i.Client.ExpireAt(ctx, i.BuildKey(key), tm)
}

func (i *Redis) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	return i.Client.Keys(ctx, i.BuildKey(pattern))
}

func (i *Redis) Pipeline() redis.Pipeliner {
	return i.Client.Pipeline()
}

func (i *Redis) Pipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return i.Client.Pipelined(ctx, fn)
}

func (i *Redis) TxPipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return i.Client.TxPipelined(ctx, fn)
}

func (i *Redis) TxPipeline() redis.Pipeliner {
	return i.Client.TxPipeline()
}

func (i *Redis) Command(ctx context.Context) *redis.CommandsInfoCmd {
	return i.Client.Command(ctx)
}

func (i *Redis) ClientGetName(ctx context.Context) *redis.StringCmd {
	return i.Client.ClientGetName(ctx)
}

func (i *Redis) Echo(ctx context.Context, message interface{}) *redis.StringCmd {
	return i.Client.Echo(ctx, message)
}

func (i *Redis) Ping(ctx context.Context) *redis.StatusCmd {
	return i.Client.Ping(ctx)
}

func (i *Redis) Quit(ctx context.Context) *redis.StatusCmd {
	return i.Client.Quit(ctx)
}

func (i *Redis) SetBit(ctx context.Context, key string, offset int64, value int) *redis.IntCmd {
	return i.Client.SetBit(ctx, i.BuildKey(key), offset, value)
}

func (i *Redis) GetBit(ctx context.Context, key string, offset int64) *redis.IntCmd {
	return i.Client.GetBit(ctx, i.BuildKey(key), offset)
}

func (i *Redis) BitCount(ctx context.Context, key string, bitCount *redis.BitCount) *redis.IntCmd {
	return i.Client.BitCount(ctx, i.BuildKey(key), bitCount)
}

func (i *Redis) HSet(ctx context.Context, key, field string, value interface{}) *redis.IntCmd {
	return i.Client.HSet(ctx, i.BuildKey(key), field, value)
}

func (i *Redis) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	return i.Client.HGet(ctx, i.BuildKey(key), field)
}

func (i *Redis) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	return i.Client.HGetAll(ctx, i.BuildKey(key))
}

func (i *Redis) HExists(ctx context.Context, key, field string) *redis.BoolCmd {
	return i.Client.HExists(ctx, i.BuildKey(key), field)
}

func (i *Redis) HDel(ctx context.Context, key, field string) *redis.IntCmd {
	return i.Client.HDel(ctx, i.BuildKey(key), field)
}

func (i *Redis) HLen(ctx context.Context, key string) *redis.IntCmd {
	return i.Client.HLen(ctx, i.BuildKey(key))
}

func (i *Redis) HKeys(ctx context.Context, key string) *redis.StringSliceCmd {
	return i.Client.HKeys(ctx, i.BuildKey(key))
}

func (i *Redis) HVals(ctx context.Context, key string) *redis.StringSliceCmd {
	return i.Client.HVals(ctx, i.BuildKey(key))
}

func (i *Redis) BLPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	newKeys := make([]string, len(keys))
	for d, key := range keys {
		newKeys[d] = i.BuildKey(key)
	}
	return i.Client.BLPop(ctx, timeout, newKeys...)
}

func (i *Redis) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	newKeys := make([]string, len(keys))
	for d, key := range keys {
		newKeys[d] = i.BuildKey(key)
	}
	return i.Client.BRPop(ctx, timeout, newKeys...)
}

func (i *Redis) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) *redis.StringCmd {
	return i.Client.BRPopLPush(ctx, i.BuildKey(source), i.BuildKey(destination), timeout)
}

func (i *Redis) LIndex(ctx context.Context, key string, index int64) *redis.StringCmd {
	return i.Client.LIndex(ctx, i.BuildKey(key), index)
}

func (i *Redis) LInsert(ctx context.Context, key, op string, pivot, value interface{}) *redis.IntCmd {
	return i.Client.LInsert(ctx, i.BuildKey(key), op, pivot, value)
}

func (i *Redis) LInsertBefore(ctx context.Context, key string, pivot, value interface{}) *redis.IntCmd {
	return i.Client.LInsertBefore(ctx, i.BuildKey(key), pivot, value)
}

func (i *Redis) LInsertAfter(ctx context.Context, key string, pivot, value interface{}) *redis.IntCmd {
	return i.Client.LInsertAfter(ctx, i.BuildKey(key), pivot, value)
}

func (i *Redis) LLen(ctx context.Context, key string) *redis.IntCmd {
	return i.Client.LLen(ctx, i.BuildKey(key))
}

func (i *Redis) LPop(ctx context.Context, key string) *redis.StringCmd {
	return i.Client.LPop(ctx, i.BuildKey(key))
}

func (i *Redis) LPopCount(ctx context.Context, key string, count int) *redis.StringSliceCmd {
	return i.Client.LPopCount(ctx, i.BuildKey(key), count)
}

func (i *Redis) LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return i.Client.LPush(ctx, i.BuildKey(key), values...)
}

func (i *Redis) LPushX(ctx context.Context, key string, value interface{}) *redis.IntCmd {
	return i.Client.LPushX(ctx, i.BuildKey(key), value)
}

func (i *Redis) LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return i.Client.LRange(ctx, i.BuildKey(key), start, stop)
}

func (i *Redis) LRem(ctx context.Context, key string, count int64, value interface{}) *redis.IntCmd {
	return i.Client.LRem(ctx, i.BuildKey(key), count, value)
}

func (i *Redis) LSet(ctx context.Context, key string, index int64, value interface{}) *redis.StatusCmd {
	return i.Client.LSet(ctx, i.BuildKey(key), index, value)
}

func (i *Redis) LTrim(ctx context.Context, key string, start, stop int64) *redis.StatusCmd {
	return i.Client.LTrim(ctx, i.BuildKey(key), start, stop)
}

func (i *Redis) RPop(ctx context.Context, key string) *redis.StringCmd {
	return i.Client.RPop(ctx, i.BuildKey(key))
}

func (i *Redis) RPopCount(ctx context.Context, key string, count int) *redis.StringSliceCmd {
	return i.Client.RPopCount(ctx, i.BuildKey(key), count)
}

func (i *Redis) RPopLPush(ctx context.Context, source, destination string) *redis.StringCmd {
	return i.Client.RPopLPush(ctx, i.BuildKey(source), i.BuildKey(destination))
}

func (i *Redis) RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return i.Client.RPush(ctx, i.BuildKey(key), values...)
}

func (i *Redis) RPushX(ctx context.Context, key string, value interface{}) *redis.IntCmd {
	return i.Client.RPushX(ctx, i.BuildKey(key), value)
}

func (i *Redis) LMove(ctx context.Context, source, destination, srcpos, destpos string) *redis.StringCmd {
	return i.Client.LMove(ctx, i.BuildKey(source), i.BuildKey(destination), srcpos, destpos)
}

func (i *Redis) BLMoveBLMove(ctx context.Context, source, destination, srcpos, destpos string, timeout time.Duration) *redis.StringCmd {
	return i.Client.BLMove(ctx, i.BuildKey(source), i.BuildKey(destination), srcpos, destpos, timeout)
}

func (i *Redis) SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return i.Client.SAdd(ctx, i.BuildKey(key), members...)
}

func (i *Redis) SCard(ctx context.Context, key string) *redis.IntCmd {
	return i.Client.SCard(ctx, i.BuildKey(key))
}

func (i *Redis) SDiff(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	newKeys := make([]string, len(keys))
	for d, key := range keys {
		newKeys[d] = i.BuildKey(key)
	}
	return i.Client.SDiff(ctx, newKeys...)
}

func (i *Redis) SDiffStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	newKeys := make([]string, len(keys))
	for d, key := range keys {
		newKeys[d] = i.BuildKey(key)
	}
	return i.Client.SDiffStore(ctx, i.BuildKey(destination), newKeys...)
}

func (i *Redis) SInter(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	newKeys := make([]string, len(keys))
	for d, key := range keys {
		newKeys[d] = i.BuildKey(key)
	}
	return i.Client.SInter(ctx, newKeys...)
}

func (i *Redis) SInterStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	newKeys := make([]string, len(keys))
	for d, key := range keys {
		newKeys[d] = i.BuildKey(key)
	}
	return i.Client.SInterStore(ctx, i.BuildKey(destination), newKeys...)
}

func (i *Redis) SIsMember(ctx context.Context, key string, member interface{}) *redis.BoolCmd {
	return i.Client.SIsMember(ctx, i.BuildKey(key), member)
}

func (i *Redis) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	return i.Client.SMembers(ctx, i.BuildKey(key))
}

func (i *Redis) SMove(ctx context.Context, source, destination string, member interface{}) *redis.BoolCmd {
	return i.Client.SMove(ctx, i.BuildKey(source), i.BuildKey(destination), member)
}

func (i *Redis) SPop(ctx context.Context, key string) *redis.StringCmd {
	return i.Client.SPop(ctx, i.BuildKey(key))
}

func (i *Redis) SPopN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	return i.Client.SPopN(ctx, i.BuildKey(key), count)
}

func (i *Redis) SRandMember(ctx context.Context, key string) *redis.StringCmd {
	return i.Client.SRandMember(ctx, i.BuildKey(key))
}

func (i *Redis) SRandMemberN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	return i.Client.SRandMemberN(ctx, i.BuildKey(key), count)
}

func (i *Redis) SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return i.Client.SRem(ctx, i.BuildKey(key), members...)
}

func (i *Redis) SUnion(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	newKeys := make([]string, len(keys))
	for d, key := range keys {
		newKeys[d] = i.BuildKey(key)
	}
	return i.Client.SUnion(ctx, newKeys...)
}

func (i *Redis) SUnionStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	newKeys := make([]string, len(keys))
	for d, key := range keys {
		newKeys[d] = i.BuildKey(key)
	}
	return i.Client.SUnionStore(ctx, i.BuildKey(destination), newKeys...)
}

func (i *Redis) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return i.Client.SScan(ctx, i.BuildKey(key), cursor, match, count)
}

func (i *Redis) SScanMap(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return i.Client.SScan(ctx, i.BuildKey(key), cursor, match, count)
}

func (i *Redis) ZAdd(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return i.Client.ZAdd(ctx, i.BuildKey(key), members...)
}

func (i *Redis) ZAddNX(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return i.Client.ZAddNX(ctx, i.BuildKey(key), members...)
}

func (i *Redis) ZAddXX(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return i.Client.ZAddXX(ctx, i.BuildKey(key), members...)
}

func (i *Redis) ZAddCh(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return i.Client.ZAddCh(ctx, i.BuildKey(key), members...)
}

func (i *Redis) ZAddNXCh(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return i.Client.ZAddNXCh(ctx, i.BuildKey(key), members...)
}

func (i *Redis) ZAddXXCh(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return i.Client.ZAddXXCh(ctx, i.BuildKey(key), members...)
}

func (i *Redis) ZAddArgs(ctx context.Context, key string, args redis.ZAddArgs) *redis.IntCmd {
	return i.Client.ZAddArgs(ctx, i.BuildKey(key), args)
}

func (i *Redis) ZAddArgsIncr(ctx context.Context, key string, args redis.ZAddArgs) *redis.FloatCmd {
	return i.Client.ZAddArgsIncr(ctx, i.BuildKey(key), args)
}

func (i *Redis) ZIncr(ctx context.Context, key string, member *redis.Z) *redis.FloatCmd {
	return i.Client.ZIncr(ctx, i.BuildKey(key), member)
}

func (i *Redis) ZIncrNX(ctx context.Context, key string, member *redis.Z) *redis.FloatCmd {
	return i.Client.ZIncrNX(ctx, i.BuildKey(key), member)
}

func (i *Redis) ZIncrXX(ctx context.Context, key string, member *redis.Z) *redis.FloatCmd {
	return i.Client.ZIncrXX(ctx, i.BuildKey(key), member)
}

func (i *Redis) ZCard(ctx context.Context, key string) *redis.IntCmd {
	return i.Client.ZCard(ctx, i.BuildKey(key))
}

func (i *Redis) ZCount(ctx context.Context, key, min, max string) *redis.IntCmd {
	return i.Client.ZCount(ctx, i.BuildKey(key), min, max)
}

func (i *Redis) ZLexCount(ctx context.Context, key, min, max string) *redis.IntCmd {
	return i.Client.ZLexCount(ctx, i.BuildKey(key), min, max)
}

func (i *Redis) ZIncrBy(ctx context.Context, key string, increment float64, member string) *redis.FloatCmd {
	return i.Client.ZIncrBy(ctx, i.BuildKey(key), increment, member)
}

func (i *Redis) ZMScore(ctx context.Context, key string, members ...string) *redis.FloatSliceCmd {
	return i.Client.ZMScore(ctx, i.BuildKey(key), members...)
}

func (i *Redis) ZPopMax(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	return i.Client.ZPopMax(ctx, i.BuildKey(key), count...)
}

func (i *Redis) ZPopMin(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	return i.Client.ZPopMin(ctx, i.BuildKey(key), count...)
}

func (i *Redis) ZRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return i.Client.ZRange(ctx, i.BuildKey(key), start, stop)
}

func (i *Redis) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	return i.Client.ZRangeWithScores(ctx, i.BuildKey(key), start, stop)
}

func (i *Redis) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return i.Client.ZRangeByScore(ctx, i.BuildKey(key), opt)
}

func (i *Redis) ZRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return i.Client.ZRangeByLex(ctx, i.BuildKey(key), opt)
}

func (i *Redis) ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	return i.Client.ZRangeByScoreWithScores(ctx, i.BuildKey(key), opt)
}

func (i *Redis) ZRank(ctx context.Context, key, member string) *redis.IntCmd {
	return i.Client.ZRank(ctx, i.BuildKey(key), member)
}

func (i *Redis) ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return i.Client.ZRem(ctx, i.BuildKey(key), members...)
}

func (i *Redis) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *redis.IntCmd {
	return i.Client.ZRemRangeByRank(ctx, i.BuildKey(key), start, stop)
}

func (i *Redis) ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd {
	return i.Client.ZRemRangeByScore(ctx, i.BuildKey(key), min, max)
}

func (i *Redis) ZRemRangeByLex(ctx context.Context, key, min, max string) *redis.IntCmd {
	return i.Client.ZRemRangeByLex(ctx, i.BuildKey(key), min, max)
}

func (i *Redis) ZRevRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return i.Client.ZRevRange(ctx, i.BuildKey(key), start, stop)
}

func (i *Redis) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	return i.Client.ZRevRangeWithScores(ctx, i.BuildKey(key), start, stop)
}

func (i *Redis) ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return i.Client.ZRevRangeByScore(ctx, i.BuildKey(key), opt)
}

func (i *Redis) ZRevRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return i.Client.ZRevRangeByLex(ctx, i.BuildKey(key), opt)
}

func (i *Redis) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	return i.Client.ZRevRangeByScoreWithScores(ctx, i.BuildKey(key), opt)
}

func (i *Redis) ZRevRank(ctx context.Context, key, member string) *redis.IntCmd {
	return i.Client.ZRevRank(ctx, i.BuildKey(key), member)
}

func (i *Redis) ZScore(ctx context.Context, key, member string) *redis.FloatCmd {
	return i.Client.ZScore(ctx, i.BuildKey(key), member)
}

func (i *Redis) ZUnionStore(ctx context.Context, dest string, store *redis.ZStore) *redis.IntCmd {
	return i.Client.ZUnionStore(ctx, i.BuildKey(dest), store)
}

func (i *Redis) ZRandMember(ctx context.Context, key string, count int, withScores bool) *redis.StringSliceCmd {
	return i.Client.ZRandMember(ctx, i.BuildKey(key), count, withScores)
}

func (i *Redis) ZDiff(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	newKeys := make([]string, len(keys))
	for d, key := range keys {
		newKeys[d] = i.BuildKey(key)
	}
	return i.Client.ZDiff(ctx, newKeys...)
}

func (i *Redis) ZDiffWithScores(ctx context.Context, keys ...string) *redis.ZSliceCmd {
	newKeys := make([]string, len(keys))
	for d, key := range keys {
		newKeys[d] = i.BuildKey(key)
	}
	return i.Client.ZDiffWithScores(ctx, newKeys...)
}
