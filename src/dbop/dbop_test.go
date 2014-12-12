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
