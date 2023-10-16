package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"mucy-core/core"
	"mucy-core/jsons"
	"reflect"
	"time"
)

var (
	client *redis.Client
	ctx    = context.Background()
)

type Conf struct {
	Host     string `mapstructure:"hots" json:"hots" yaml:"hots"`
	Port     string `mapstructure:"port" json:"port" yaml:"port"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`
}

func InitializeRedis(config *Conf) error {
	client = redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password, // no password set
		DB:       config.DB,       // use default DB
		PoolSize: 100,             // 连接池大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	core.App.Redis = client
	return nil
}

// 缓存数据
func Set(key string, value string, timeOut time.Duration) error {
	err := client.Set(ctx, key, value, timeOut).Err()
	if err != nil {
		return err
	}
	return nil
}

// 缓存数据
func SetObj(key string, value interface{}, timeOut time.Duration) error {
	err := client.Set(ctx, key, jsons.ToJsonStr(value), timeOut).Err()
	if err != nil {
		return err
	}
	return nil
}

// 根据Key取值
func Get(key string) (string, error) {
	val, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.New("key不存在")
	} else if err != nil {
		return "", err
	} else {
		return val, nil
	}
}

// 缓存数据
func Delete(key string) (int64, error) {
	num, err := client.Del(ctx, key).Result()
	return num, err
}

// 获取当前数据库key的数量
func GetSize() (int64, error) {
	num, err := client.DBSize(ctx).Result()
	return num, err
}

// 清空当前数据库
func FlushDB() (string, error) {
	res, err := client.FlushDB(ctx).Result()
	return res, err
}

// 清空所有数据库
func FlushAll() (string, error) {
	res, err := client.FlushAll(ctx).Result()
	return res, err
}

// 创建合集
func ZAdd(key string, value ...*redis.Z) error {
	err := client.ZAdd(ctx, key, value...).Err()
	client.Expire(ctx, key, 40*24*time.Hour)
	if err != nil {
		return err
	}
	return nil
}

// 合集新增数据
func SetZIncrBy(key string, member string) error {
	err := client.ZIncrBy(ctx, key, 1, member).Err()
	client.Expire(ctx, key, 35*24*time.Hour)
	if err != nil {
		return err
	}
	return nil
}

// GetRangeWithScores 获取分数最高
func GetRangeWithScores(key string, start int64, stop int64) []redis.Z {
	ret := client.ZRevRangeWithScores(ctx, key, start, stop).Val()
	return ret
}

func ZUnionStore(dest string, keys ...string) error {
	//var weights []float64
	store := &redis.ZStore{
		Keys: keys,
		//Weights:weights,
		Aggregate: "SUM",
	}
	err := client.ZUnionStore(ctx, dest, store).Err()
	client.Expire(ctx, dest, 12*time.Hour)
	return err
}

func reflectValSum(val reflect.Value, args ...reflect.Value) reflect.Value {
	kind := val.Kind()
	vi := val.Interface()
	for _, v := range args {
		switch kind {
		case reflect.Int:
			val.Set(reflect.ValueOf(vi.(int) + v.Interface().(int)))
		case reflect.Int8:
			val.Set(reflect.ValueOf(vi.(int8) + v.Interface().(int8)))
		case reflect.Int16:
			val.Set(reflect.ValueOf(vi.(int16) + v.Interface().(int16)))
		case reflect.Int32:
			val.Set(reflect.ValueOf(vi.(int32) + v.Interface().(int32)))
		case reflect.Int64:
			val.Set(reflect.ValueOf(vi.(int64) + v.Interface().(int64)))
		case reflect.Uint:
			val.Set(reflect.ValueOf(vi.(uint) + v.Interface().(uint)))
		case reflect.Uint8:
			val.Set(reflect.ValueOf(vi.(uint8) + v.Interface().(uint8)))
		case reflect.Uint16:
			val.Set(reflect.ValueOf(vi.(uint16) + v.Interface().(uint16)))
		case reflect.Uint32:
			val.Set(reflect.ValueOf(vi.(uint32) + v.Interface().(uint32)))
		case reflect.Uint64:
			val.Set(reflect.ValueOf(vi.(uint64) + v.Interface().(uint64)))
		case reflect.Float32:
			val.Set(reflect.ValueOf(vi.(float32) + v.Interface().(float32)))
		case reflect.Float64:
			val.Set(reflect.ValueOf(vi.(float64) + v.Interface().(float64)))
		}
	}
	return val
}

func StructFieldSum(val interface{}, args ...interface{}) {
	v := reflect.ValueOf(val).Elem()
	t := v.Type()
	num := v.NumField()
	for _, arg := range args {
		vv := reflect.ValueOf(arg).Elem()
		if t != vv.Type() {
			continue
		}
		for i := 0; i < num; i++ {
			//如果是下划线或小写字符开头则忽略
			if t.Field(i).Name[0] >= 95 {
				continue
			}
			reflectValSum(v.Field(i), vv.Field(i))
		}
	}
}
