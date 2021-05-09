package token

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Tokens struct {
	tokens map[string]mqtt.Token
}

func NewTokens() *Tokens {
	return &Tokens{
		tokens: make(map[string]mqtt.Token),
	}
}

func (ts *Tokens) Append(name string, t mqtt.Token) {
	ts.tokens[name] = t
}

// Wait will wait indefinitely for the Token to complete, ie the Publish
// to be sent and confirmed receipt from the broker.
func (ts *Tokens) Wait() bool {
	result := true

	for _, t := range ts.tokens {
		result = result && t.Wait()
	}

	return result
}

// WaitTimeout takes a time.Duration to wait for the flow associated with the
// Token to complete, returns true if it returned before the timeout or
// returns false if the timeout occurred. In the case of a timeout the Token
// does not have an error set in case the caller wishes to wait again.
func (ts *Tokens) WaitTimeout(d time.Duration) bool {
	result := atomic.Value{}
	result.Store(true)

	wg := &sync.WaitGroup{}
	wg.Add(len(ts.tokens))

	for _, t := range ts.tokens {
		t := t
		go func() {
			defer wg.Done()
			result.Store(result.Load().(bool) && t.WaitTimeout(d))
		}()
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-done:
		return result.Load().(bool)
	case <-time.After(d):
		return false
	}
}

// Done returns a channel that is closed when the flow associated
// with the Token completes. Clients should call Error after the
// channel is closed to check if the flow completed successfully.
//
// Done is provided for use in select statements. Simple use cases may
// use Wait or WaitTimeout.
func (ts *Tokens) Done() <-chan struct{} {
	ch := make(chan struct{})

	go func() {
		for _, t := range ts.tokens {
			<-t.Done()
		}

		close(ch)
	}()

	return ch
}

type Errors []error

func (es Errors) Error() string {
	result := ""

	for _, e := range []error(es) {
		result += e.Error()
	}

	return result
}

func (ts *Tokens) Error() error {
	errors := make([]error, 0, len(ts.tokens))

	for name, token := range ts.tokens {
		errors = append(errors, fmt.Errorf("%s: %w", name, token.Error()))
	}

	return Errors(errors)
}
