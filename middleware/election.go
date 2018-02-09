package middleware

import (
	"os"
	"io"
	"github.com/gin-gonic/gin"
	"github.com/ringtail/leadership"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store/etcd"
	"time"
	"go.uber.org/zap"
)

var DefaultErrorWriter io.Writer = os.Stderr
var sugar *zap.SugaredLogger
var isLeader, canProxy bool
var leaderEndpoint string

func init() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar = logger.Sugar()
	isLeader = false
	canProxy = false
	etcd.Register()
	go participate()
}


func Election() gin.HandlerFunc {
	return ElectionAndProxy(DefaultErrorWriter)
}

func ElectionAndProxy(out io.Writer) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if isLeader {

			}else if canProxy == true {

			}
		}()
		c.Next()
	}
}

func participate() {
	// Create a store using pkg/store.
	client, err := libkv.NewStore(store.ETCD, []string{"127.0.0.1:2379"}, &store.Config{})
	if err != nil {
		panic(err)
	}

	waitTime := 10 * time.Second
	underwood := leadership.NewCandidate(client, "service/swarm/leader", "underwood", 15*time.Second)

	go func() {
		for {
			run(underwood)
			time.Sleep(waitTime)
			// retry
		}
	}()
}

func run(candidate *leadership.Candidate) {
	electedCh, errCh := candidate.RunForElection()
	for {
		select {
		case isElected := <-electedCh:
			if isElected {
				// Do something
				isLeader = true
				canProxy = false
				sugar.Infof("I'm leader now")
			} else {
				isLeader = false
				canProxy = true
				// Do something else
				leaderEndpoint =
				sugar.Infof("I'm follower now")
			}

		case err := <-errCh:
			sugar.Errorf("Failed to select a leader, Because of %s", err.Error())
			isLeader = false
			canProxy = false
			return
		}
	}
}
