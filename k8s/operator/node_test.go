package operator

import (
	"sync"
	"testing"
	"time"

	v1 "github.com/allinbits/starport-operator/api/v1"
	"github.com/stretchr/testify/require"
)

func TestDataRaceInNewNode(t *testing.T) {
	nc1 := NodeConfiguration{
		Name:               "test1",
		DockerImage:        "dockerImage",
		DockerImageVersion: "main",
		Namespace:          "namespace",
		JoinConfig:         &v1.JoinConfig{},
	}

	nc2 := NodeConfiguration{
		Name:               "test2",
		DockerImage:        "dockerImage",
		DockerImageVersion: "main",
		Namespace:          "namespace",
		JoinConfig:         &v1.JoinConfig{},
	}

	timeout := 10 * time.Second

	wg := sync.WaitGroup{}

	f := func(n NodeConfiguration, w *sync.WaitGroup) {
		after := time.After(timeout)

		for {
			select {
			case <-after:
				wg.Done()
				return
			default:
				_, err := NewNode(n)
				require.NoError(t, err)
			}
		}
	}

	wg.Add(2)

	go f(nc1, &wg)
	go f(nc2, &wg)

	wg.Wait()
}
