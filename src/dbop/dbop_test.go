package dbop

import (
	"testing"
)

func Test_GetSetCount(t *testing.T) {
	count, err := dbop.GetSetCount()
	if err != nil {
		t.Errorf("%s", err.Error())
	}
}

func Test_LockRedis(t *testing.T) {
    err := dbop.LockRedis()
    if err != nil {
        t.Errorf("%s", err.Error())
    }
}

func Test_UnlockRedis(t *testing.T){
     err := dbop.UnlockRedis()
     if err != nil {
     t.Errorf("%s", err.Error())
     }
}

func Test_tt(){
     return 
}