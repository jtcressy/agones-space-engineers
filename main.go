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

	savesDir := "/appdata/space-engineers/SpaceEngineersDedicated/Saves"
	if ok,err := isEmpty(savesDir); ok {
		message := fmt.Sprintf(">>> Saves directory is empty! (%v)\n", savesDir)
		message += fmt.Sprintf(">>> Shell into the pod, copy/download a save file, then restart.\n")
		
		deadline := time.Now().Add(time.Hour)
	
		for range time.Tick(1 * time.Second) {
			timeRemaining := getTimeRemaining(deadline)
	
			if timeRemaining.t <= 0 {
				fmt.Println("Countdown reached!")
				break
			}
	
			fmt.Printf(message + ">>> Holding the door open for you to shell in... Minutes: %d Seconds: %d\n", timeRemaining.m, timeRemaining.s)
		}
		fmt.Println("")
	} else {
		if err != nil {
			panic(err)
		}
		fmt.Printf(">>> Saves directory is not empty, continuing with startup! (%v)\n", savesDir)
	}

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

func isEmpty(name string) (bool, error) {
    f, err := os.Open(name)
    if err != nil {
        return false, err
    }
    defer f.Close()

    _, err = f.Readdirnames(1) // Or f.Readdir(1)
    if err == io.EOF {
        return true, nil
    }
    return false, err // Either not empty or error, suits both cases
}

type countdown struct {
	t int
	d int
	h int
	m int
	s int
}

func getTimeRemaining(t time.Time) countdown {
	currentTime := time.Now()
	difference := t.Sub(currentTime)

	total := int(difference.Seconds())
	days := int(total / (60 * 60 * 24))
	hours := int(total / (60 * 60) % 24)
	minutes := int(total/60) % 60
	seconds := int(total % 60)

	return countdown{
		t: total,
		d: days,
		h: hours,
		m: minutes,
		s: seconds,
	}
}