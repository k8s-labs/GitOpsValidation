package logging

import "testing"

func TestInitLogger(t *testing.T) {
	InitLogger()
	if Logger == nil {
		t.Error("Logger should not be nil after InitLogger")
	}
	SyncLogger()
}
