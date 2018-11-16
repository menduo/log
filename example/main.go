package main

import "github.com/menduo/log"

func main() {
	log.SetLevelByString("warn")
	log.Warning("warning. this will be logged")
	log.Debug("debug. this will not be logged")

	logger := log.New()
	logger.SetLevelByString("debug")
	logger.Debug("debug, will be logged")

	// output
	// 2018/11/17 01:20:28 main.go:7: [W] warning. this will be logged
	// 2018/11/17 01:20:28 main.go:12: [D] debug, will be logged
}
