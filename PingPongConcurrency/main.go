// Ping Pong Concurrency
//
// Two goroutines (ping and pong) keep sending messages into one shared
// channel. A reader goroutine prints those messages for 5 seconds, then
// signals main through a "done" channel so the program can exit.
// context.Context is used to tell the senders when to stop.
//
// ---------------------------------------------------------------------------
// sync.WaitGroup vs context.Context
// ---------------------------------------------------------------------------
//
// Both are used with goroutines, but they solve different problems:
//
//	WaitGroup = "wait until the goroutines are FINISHED"
//	Context   = "tell the goroutines they SHOULD STOP"
//
// | Aspect            | sync.WaitGroup                        | context.Context                                 |
// |-------------------|---------------------------------------|-------------------------------------------------|
// | Purpose           | Wait for a group of goroutines to end | Send cancellation / deadline / values down      |
// | Direction         | Workers -> waiter ("I'm done")        | Parent -> workers ("please stop")               |
// | Main API          | Add(n), Done(), Wait()                | WithCancel, WithTimeout, WithDeadline, Done()   |
// | Stops goroutines? | No, only waits for them               | No, only asks; goroutine must check ctx.Done()  |
// | Waits for them?   | Yes, Wait() blocks until counter == 0 | No, cancel() returns immediately                |
// | Timeouts          | Not built in                          | Built in (WithTimeout / WithDeadline)           |
// | Propagation       | None, one flat counter                | Tree: cancelling a parent cancels all children  |
// | Carries data      | No                                    | Yes, request-scoped values via WithValue        |
// | Typical use       | Fan-out work, wait for all results    | HTTP requests, DB calls, graceful shutdown      |
//
// They are usually used TOGETHER: context tells workers to stop, and a
// WaitGroup makes sure they actually stopped before main exits:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	var wg sync.WaitGroup
//	wg.Add(2)
//	go func() { defer wg.Done(); ping(ctx, ch) }()
//	go func() { defer wg.Done(); pong(ctx, ch) }()
//
//	// ... read from ch until ctx.Done() ...
//
//	wg.Wait() // ping and pong have really returned here
//
// In this program a plain "done" channel plays the WaitGroup role for the
// single reader goroutine. Nothing waits for ping/pong themselves, which
// is fine here only because main exits right after.
package main

import (
	"context"
	"fmt"
	"time"
)

// ping sends a "ping" message with the current time into ch about once
// per second, until ctx is cancelled.
func ping(ctx context.Context, ch chan string) {
	for {
		// select waits on whichever case is ready first.
		select {
		case <-ctx.Done():
			// Context cancelled -> stop the goroutine.
			return
		case ch <- fmt.Sprintf("ping: %v", time.Now()):
			// Message was received by the reader; pause before the next one.
			time.Sleep(1 * time.Second)
		}
	}
}

// pong is the same as ping but sends "pong" messages. Both share the same
// channel, so the order of ping/pong output is not guaranteed.
func pong(ctx context.Context, ch chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case ch <- fmt.Sprintf("pong: %v", time.Now()):
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {

	// Cancellable context passed to ping and pong. cancel() closes
	// ctx.Done(), which makes both goroutines return.
	// defer runs it when main returns.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Unbuffered channel: each send blocks until the reader receives it.
	pingerCh := make(chan string)

	// done is used by the reader goroutine to tell main it has finished.
	// struct{} takes no memory; only the signal matters.
	done := make(chan struct{})

	// Start the two senders.
	go ping(ctx, pingerCh)
	go pong(ctx, pingerCh)

	// Reader goroutine: prints messages until the 5 second timeout fires.
	go func() {

		// time.After returns a channel that receives a value once after 5s.
		timeout := time.After(5 * time.Second)

		for {
			select {
			case <-timeout:
				fmt.Println("Timeout reached")
				// NOTE: closing the channel while ping/pong may still try to
				// send on it can panic ("send on closed channel"). Safer:
				// call cancel() here instead, and let only the senders'
				// owner close the channel after they have stopped.
				close(pingerCh)
				// Tell main we are finished.
				done <- struct{}{}
				return
			case msg := <-pingerCh:
				fmt.Println(msg)
			}
		}

	}()

	// Block main until the reader goroutine signals it is done.
	<-done

	fmt.Println("done")
}
