package main

import (
	_ "go.uber.org/goleak"
)

// TODO: investigate why netpoll leaves pending goroutines
//  https://github.com/modelgateway/Glide/issues/33
//func TestMain(m *testing.M) {
//	goleak.VerifyTestMain(m)
//}
