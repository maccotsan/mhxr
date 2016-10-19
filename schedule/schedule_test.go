package schedule

import (
	"testing"
)

func TestGetSchedule(t *testing.T) {

	_, err := GetEventSchedule()
	if err != nil {
		t.Errorf("%s", err)
	}

	_, err2 := CreateHTML()
	if err2 != nil {
		t.Errorf("%s", err2)
	}
}
