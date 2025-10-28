package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/render-test-server/internal/config"
	httpserver "github.com/render-test-server/internal/http-server"
	"github.com/render-test-server/internal/logger"
	"github.com/render-test-server/internal/smtptool"
)

func main() {
    fmt.Println("Hello, World!")

    logger := logger.SetupLogger()
    cfg := config.MustLoadConfig()

    logger.Infof("Start %s", time.Now().String())

    logger.Infof("%v", cfg)
    logger.WithError(smtptool.RunOnceToCheck(cfg.SMTPConfig)).Info("Calling the smtp execution")

    wg := new(sync.WaitGroup)	
	wg.Add(1)

    serverapp := httpserver.NewServerApp(cfg, logger, wg)

    go serverapp.Run()

    wg.Wait()
}