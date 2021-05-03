package basic

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MAAARKIN/notification-async/model"
)

func worker(url string, event string) {
	payload := map[string]interface{}{
		"event":       event,
		"destination": map[string]string{"email": url},
		"source":      event,
	}

	data, _ := json.Marshal(payload)

	//my http client post here
	time.Sleep(100 * time.Millisecond) //to simulate the api.

	out := fmt.Sprintf("send event to %v, payload to send %v", url, string(data)) // do somethingg useful.
	fmt.Println(out)
}

func Start(ctx context.Context, op model.Options) {
	start := time.Now()

	csvfile, err := os.Open(op.Filename)
	if err != nil {
		log.Fatal(err)
	}
	defer csvfile.Close()

	scanner := bufio.NewScanner(csvfile)
	scanner.Split(bufio.ScanLines)

	// timerAux := 0
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			fmt.Println("worker finish")
			os.Exit(0)
		default:
			// timerAux = timerAux + 1
			// if timerAux == 20 {
			// 	time.Sleep(700 * time.Millisecond)
			// 	timerAux = 0
			// }
			worker(scanner.Text(), op.Event)
		}
	}

	fmt.Printf("\n%2fs", time.Since(start).Seconds())
	fmt.Println()
}
