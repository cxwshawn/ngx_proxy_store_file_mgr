package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type FileMgrConfig struct {
	MaxFileLimit   int
	CheckInterval  int
	ExpireDays     int
	ErrorLog       bool
	AccessLog      bool
	SortedSetName  string
	HashName       string
	DelPercentOnce float32
	RedisAddr      string
	RoutineCount   int
	RedisLockName  string
}

var Defaultfmc *FileMgrConfig

func init() {
	Defaultfmc = &FileMgrConfig{100000,
		10 * 60, 7, true, false, "defset", "defhash",
		33.3333, "127.0.0.1:6379", 32, "cache"}
}

func ParseConfig(configFileName string) error {
	exePath, err1 := exec.LookPath(os.Args[0])
	if err1 != nil {
		fmt.Println(err1.Error())
		return err1
	}
	exeFullPath, err1 := filepath.Abs(exePath)
	if err1 != nil {
		fmt.Println(err1.Error())
		return err1
	}
	configFullPath := fmt.Sprintf("%s/%s", filepath.Dir(exeFullPath), configFileName)

	file, err1 := os.Open(configFullPath)
	defer file.Close()

	if err1 != nil {
		fmt.Println(err1.Error())
		return err1
	}
	data, err1 := ioutil.ReadAll(file)
	if err1 != nil {
		fmt.Println(err1.Error())
		return err1
	}
	if Defaultfmc == nil {
		fmt.Println("Defaultfmc is nil!")
		return errors.New("Default config object is nil.")
	}
	err1 = json.Unmarshal(data, Defaultfmc)
	if err1 != nil {
		fmt.Println(err1.Error())
		return err1
	}
	return nil
}

func GetCheckInterval() int {
	return Defaultfmc.CheckInterval
}

func GetMaxFileLimit() int {
	return Defaultfmc.MaxFileLimit
}

func IsLogError() bool {
	return Defaultfmc.ErrorLog
}

func IsLogAccess() bool {
	return Defaultfmc.AccessLog
}

func GetExpireDays() int {
	return Defaultfmc.ExpireDays
}

func GetSortedSetName() string {
	return Defaultfmc.SortedSetName
}

func GetHashName() string {
	return Defaultfmc.HashName
}

func GetDelPercent() float32 {
	return Defaultfmc.DelPercentOnce
}

func GetRedisAddr() string {
	return Defaultfmc.RedisAddr
}

func GetRoutineCount() int {
	return Defaultfmc.RoutineCount
}

func GetRedisKeyName() string {
	return Defaultfmc.RedisLockName
}
