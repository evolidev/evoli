package evoli

import (
	"github.com/cosmtrek/air/runner"
	"github.com/evolidev/evoli/framework/console"
	"github.com/evolidev/evoli/framework/use"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Start() {
	//console.Commands()
	//watch()
	console.Watch()
}

func watch() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var err error
	cfg, err := runner.InitConfig("airs.toml")
	if err != nil {
		log.Fatal(err)
		return
	}

	use.D(cfg)
	return

	r, err := runner.NewEngineWithConfig(cfg, true)
	if err != nil {
		log.Fatal(err)
		return
	}
	go func() {
		<-sigs
		r.Stop()
	}()

	defer func() {
		if e := recover(); e != nil {
			log.Fatalf("PANIC: %+v", e)
		}
	}()

	// kill after 5 seconds
	go func() {
		<-time.After(1 * time.Second)
		r.Stop()
	}()

	r.Run()

}
