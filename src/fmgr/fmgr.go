package main

import (
	"config"
	"container/list"
	"dbop"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"srvlog"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

var configFileName *string = flag.String("config", "server.json", "server config file name")

var stopCh chan bool
var stopFlag uint32
var mutex sync.Mutex

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
}
func testDBOp() {
	count, err := dbop.GetSetCount()
	if err != nil {
		fmt.Println("%s", err.Error())
	}
	fmt.Println("count is:", count)
	keys, err := dbop.GetLeastUsedKeys()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range keys {
		fmt.Println(string(v))
	}
	// filepaths, err := dbop.GetLeastUsedFiles(keys)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// for _, v := range filepaths {
	// 	fmt.Println(v)
	// }
}

func stopHandler(w http.ResponseWriter, req *http.Request) {
	srvlog.Println("stop handler...")

	if req.Method != "GET" {
		return
	}
	w.WriteHeader(http.StatusOK)
	stopCh <- true
}

func restAPI() {
	http.HandleFunc("/stop", stopHandler)

	s := &http.Server{
		Addr:           fmt.Sprintf("127.0.0.1:10000"),
		Handler:        nil,
		ReadTimeout:    100 * time.Millisecond,
		WriteTimeout:   100 * time.Millisecond,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

var filepaths *list.List

func getFiles() error {
	defer func() {
		if err, ok := recover().(error); ok {
			srvlog.EPrintln("WARN: panic %v", err)
			srvlog.EPrintln(string(debug.Stack()))
		}
	}()
	srvlog.Printf("getFiles...")
	data, err := dbop.GetLeastUsedKeys()
	if err != nil {
		srvlog.EPrintf("%s", err.Error())
		return err
	}
	filepaths, err = dbop.GetLeastUsedFiles(data)
	if err != nil {
		srvlog.EPrintf("%s", err.Error())
		return err
	}
	srvlog.Printf("filepaths length :%d", filepaths.Len())
	elem := filepaths.Front()
	for {
		if elem != nil {
			srvlog.Println(string(elem.Value.([]byte)))
		} else {
			break
		}
		elem = elem.Next()
	}
	return nil
}

func deleteFiles(ch chan bool) {
	defer func(ch chan bool) {
		if err, ok := recover().(error); ok {
			srvlog.EPrintln("WARN: panic %v", err)
			srvlog.EPrintln(string(debug.Stack()))
		}
		srvlog.Printf("deleteFiles send over to channel.")
		ch <- true
	}(ch)

	srvlog.Printf("deleteFiles...")

	for {
		stop := atomic.LoadUint32(&stopFlag)
		if stop == 1 {
			break
		}

		mutex.Lock()
		if filepaths.Len() == 0 {
			mutex.Unlock()
			break
		}
		pElem := filepaths.Front()
		filepath := filepaths.Remove(pElem).([]byte)
		mutex.Unlock()
		//todo:delete file.
		srvlog.Printf("delete file %s", string(filepath))
		// srvlog.Printf("delete file %s", filepath)
		err := os.Remove(string(filepath))
		if err != nil {
			srvlog.EPrintf("%s", err.Error())
		}
	}
	srvlog.Printf("deleteFiles end...")
}

func NeedLRU() (lru bool, err error) {
	fs := &syscall.Statfs_t{}
	err = syscall.Statfs("/", fs)
	if err != nil {
		srvlog.EPrintf("%s", err.Error())
		lru = false
		return
	}

	all := fs.Blocks * uint64(fs.Bsize)
	free := fs.Bfree * uint64(fs.Bsize)
	leftPercent := int32(float64(free) / float64(all) * 100.0)
	if leftPercent < 20 {
		return true, nil
	}
	return false, nil
}

func main() {
	err := config.ParseConfig(*configFileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	srvlog.InitLog("Nginx_File_Manager")
	dbop.InitDb(config.GetRedisAddr())

	stopFlag = 0
	processing := false
	stopChs := make([]chan bool, config.GetRoutineCount())
	for i := 0; i < config.GetRoutineCount(); i++ {
		stopChs[i] = make(chan bool)
	}
	timeout := make(chan bool)
	stopCh = make(chan bool)
	loopBreak := false

	go restAPI()

	for {
		go func() {
			time.Sleep(time.Duration(config.GetCheckInterval() * int(time.Second)))
			timeout <- true
			//close(timeout)
		}()
		select {
		case <-stopCh:
			srvlog.Printf("received stop command!")
			atomic.SwapUint32(&stopFlag, 1)
			if processing == true {
				for _, ch := range stopChs {
					<-ch
				}
			}
			loopBreak = true
		case <-timeout:
			srvlog.Printf("timeout, start to process.")
			lru, _ := NeedLRU()
			if lru == true {
				err = dbop.LockRedis()
				if err != nil {
					continue
				}

				processing = true
				err = getFiles()
				if err != nil {
					dbop.UnlockRedis()
					processing = false
					continue
				}
				for i := 0; i < config.GetRoutineCount(); i++ {
					go deleteFiles(stopChs[i])
				}
				// for _, ch := range stopChs {
				// 	<-ch
				// }
				for i := 0; i < config.GetRoutineCount(); i++ {
					<-stopChs[i]
				}
				dbop.DeleteLeastUsedKeys()
				dbop.UnlockRedis()

				srvlog.Printf("timeout process end.")
				processing = false
			}
		}
		if loopBreak == true {
			break
		}
	}
	srvlog.Println("nginx file manager quitting... ")

	close(timeout)
	close(stopCh)
	for _, val := range stopChs {
		close(val)
	}
}
