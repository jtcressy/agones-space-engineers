package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	sdk "agones.dev/agones/sdks/go"
)

type interceptor struct {
	forward io.Writer
	intercept func(p []byte)
}

func (i *interceptor) Write(p []byte) (n int, err error) {
	if i.intercept != nil {
		i.intercept(p)
	}

	return i.forward.Write(p)
}

func main() {
	input := flag.String("i", "", "path to server_linux.sh or server_windows.bat")
	args := flag.String("args", "", "additional arguments to pass to the script")
	flag.Parse()

	argsList := strings.Split(strings.Trim(strings.TrimSpace(*args), "'"), " ")
	fmt.Println(">>> Connecting to Agones with the SDK")
	s, err := sdk.NewSDK()
	if err != nil {
		log.Fatalf(">>> Could not connect to sdk: %v", err)
	}

	fmt.Println(">>> Starting health checking")
	go doHealth(s)

	fmt.Println(">>> Starting wrapper for space-engineers!")
	fmt.Printf(">>> Path to space-engineers server script: %s %v\n", *input, argsList)

	// track references to listening count
	gameReady := false

	cmd := exec.Command(*input, argsList...) // #nosec
	cmd.Stderr = &interceptor{forward: os.Stderr}
	cmd.Stdout = &interceptor{
		forward: os.Stdout,
		intercept: func(p []byte) {
			if gameReady {
				return
			}

			str := strings.TrimSpace(string(p))
			// space-engineers will say "Game ready..." when ready.
			if count := strings.Count(str, "Game ready..."); count > 0 {
				gameReady = true
				fmt.Printf(">>> Found 'Game ready...' statement \n")
				fmt.Printf(">>> Moving to READY: %s \n", str)
				err = s.Ready()
				if err != nil {
					log.Fatalf("Could not send ready message")
				}
			}
		}}

	err = cmd.Start()
	if err != nil {
		log.Fatalf(">>> Error Starting Cmd %v", err)
	}
	err = cmd.Wait()
	log.Fatal(">>> space-engineers shutdown unexpectantly", err)
}

// doHealth sends the regular Health Pings
func doHealth(sdk *sdk.SDK) {
	tick := time.Tick(2 * time.Second)
	for {
		err := sdk.Health()
		if err != nil {
			log.Fatalf("[wrapper] Could not send health ping, %v", err)
		}
		<-tick
	}
}