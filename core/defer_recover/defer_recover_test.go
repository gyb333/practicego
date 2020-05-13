package defer_recover_test

import (
	"../defer_recover"
	"testing"
)

func TestDeferFunc(t *testing.T) {
	defer_recover.DeferFunc()
}

func TestPanicRecover(t *testing.T) {
	defer_recover.PanicRecover()
}