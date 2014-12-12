package dbop

import (
	"config"
	"container/list"
	"github.com/hoisie/redis"
	"srvlog"
)

var client redis.Client

// func init() {

// }
func InitDb(redisAddr string) {
	client.Addr = redisAddr
}

func GetSetCount() (count int, err error) {
	sortedSetName := config.GetSortedSetName()
	card, err1 := client.Zcard(sortedSetName)
	if err1 != nil {
		srvlog.EPrintf("%s", err.Error())
		err = err1
		return
	}
	return card, nil
}

func LockRedis() error {
	keyName := config.GetRedisKeyName()
	err := client.Set(keyName, []byte("lock"))
	if err != nil {
		srvlog.Printf("%s", err.Error())
		return err
	}
	return nil
}

func UnlockRedis() error {
	keyName := config.GetRedisKeyName()
	err := client.Set(keyName, []byte("unlock"))
	if err != nil {
		srvlog.Printf("%s", err.Error())
		return err
	}
	return nil
}

func GetLeastUsedKeys() (keys [][]byte, err error) {
	sortedSetName := config.GetSortedSetName()
	card, err1 := client.Zcard(sortedSetName)
	if err1 != nil {
		srvlog.EPrintf("%s", err.Error())
		err = err1
		return
	}
	delPercent := config.GetDelPercent()
	delCount := int(float32(card) * (delPercent / 100.00))
	data, err1 := client.Zrange(sortedSetName, 0, delCount)
	if err1 != nil {
		srvlog.EPrintf("%s", err.Error())
		err = err1
		return
	}

	keys = make([][]byte, 0)
	keys = data
	err = nil
	return
}

func DeleteLeastUsedKeys() error {
	//delete sorted element in sorted set
	//delete element in the hash buckets.
	sortedSetName := config.GetSortedSetName()
	card, err1 := client.Zcard(sortedSetName)
	if err1 != nil {
		srvlog.EPrintf("%s", err1.Error())
		return err1
	}
	delPercent := config.GetDelPercent()
	delCount := int(float32(card) * (delPercent / 100.0))
	keys, err1 := client.Zrange(sortedSetName, 0, delCount)
	if err1 != nil {
		srvlog.EPrintf("%s", err1.Error())
		return err1
	}
	_, err1 = client.Zremrangebyrank(sortedSetName, 0, delCount)
	if err1 != nil {
		srvlog.EPrintf("%s", err1.Error())
		return err1
	}

	hashName := config.GetHashName()
	for _, val := range keys {
		_, err1 := client.Hdel(hashName, string(val))
		if err1 != nil {
			srvlog.EPrintf("%s", err1.Error())
			return err1
		}
	}
	return nil
}

func GetLeastUsedFiles(keys [][]byte) (filepaths *list.List, err error) {
	filepaths = list.New()
	hashName := config.GetHashName()
	strKeys := make([]string, 0)
	for _, v := range keys {
		strKeys = append(strKeys, string(v))
	}

	data, err1 := client.Hmget(hashName, strKeys...)
	if err1 != nil {
		srvlog.EPrintf("%s", err.Error())
		err = err1
		return
	}
	//filepaths = make([]string, len(data))
	srvlog.Printf("filepaths length :%d, data len:%d, strkeys len:%d, keys len:%d",
		filepaths.Len(), len(data), len(strKeys), len(keys))

	for _, v := range data {
		srvlog.Printf("%s\n", string(v))
		//filepaths = append(filepaths, string(v))
		filepaths.PushBack(v)
	}
	srvlog.Printf("filepaths length :%d", filepaths.Len())

	err = nil
	return
}
