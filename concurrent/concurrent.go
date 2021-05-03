package concurrent

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/MAAARKIN/notification-async/model"
)

func worker(ctx context.Context, out chan<- string, src <-chan string, event string) {
	for {
		select {
		case url, ok := <-src: // you must check for readable state of the channel.
			if !ok {
				return
			}

			payload := map[string]interface{}{
				"event":       event,
				"destination": map[string]string{"email": url},
				"source":      event,
			}

			data, _ := json.Marshal(payload)

			//my http client post here
			time.Sleep(100 * time.Millisecond) //to simulate the api.

			out <- fmt.Sprintf("send event to %v, payload to send %v", url, string(data)) // do somethingg useful.
		case <-ctx.Done(): // if the context is cancelled, quit.
			fmt.Println("worker finish")
			return
		}
	}
}

func Start(ctx context.Context, op model.Options, numberOfWorkers int) {

	start := time.Now()

	csvfile, err := os.Open(op.Filename)
	if err != nil {
		log.Fatal(err)
	}
	defer csvfile.Close()

	scanner := bufio.NewScanner(csvfile)
	scanner.Split(bufio.ScanLines)

	// create the pair of input/output channels for the controller=>workers com.
	src := make(chan string)
	out := make(chan string)

	// use a waitgroup to manage synchronization
	var wg sync.WaitGroup

	// declare the workers
	for i := 0; i < numberOfWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(ctx, out, src, op.Event)
		}()
	}

	// read the csv and write it to src
	go func() {
		// timerAux := 0
		for scanner.Scan() {
			// timerAux = timerAux + 1
			// if timerAux == 20 {
			// 	time.Sleep(700 * time.Millisecond)
			// 	timerAux = 0
			// }
			src <- scanner.Text() // you might select on ctx.Done().
		}
		close(src) // close src to signal workers that no more job are incoming.
	}()

	// drain the output
	go func() {
		for res := range out {
			fmt.Println(res)
		}
	}()

	// wait for worker group to finish and close out
	wg.Wait()  // wait for writers to quit.
	close(out) // when you close(out) it breaks the below loop.

	fmt.Printf("\n%2fs", time.Since(start).Seconds())
	fmt.Println()
}
