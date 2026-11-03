//go:build unix
// +build unix

package prompt

import (
	"io"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"
)

type testWriter struct {
	VT100Writer
}

func (w *testWriter) Flush() error {
	w.buffer = nil
	return nil
}

var _ Writer = &testWriter{}

// testReader queues scripted input and reports a fixed window size.
// Read returns io.EOF when the queue is empty so the prompt's poll loop can retry.
type testReader struct {
	mu      sync.Mutex
	inputs  [][]byte
	winSize *WinSize
}

func (r *testReader) Open() error {
	return nil
}

func (r *testReader) Close() error {
	return nil
}

func (r *testReader) GetWinSize() *WinSize {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.winSize != nil {
		return r.winSize
	}
	return &WinSize{Row: DefRowCount, Col: DefColCount}
}

func (r *testReader) Read(buf []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.inputs) == 0 {
		return 0, io.EOF
	}
	b := r.inputs[0]
	n := copy(buf, b)
	if n < len(b) {
		r.inputs[0] = b[n:]
		return n, nil
	}
	r.inputs = r.inputs[1:]
	return n, nil
}

func (r *testReader) appendInput(b []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.inputs = append(r.inputs, append([]byte(nil), b...))
}

func TestHandleSignals_WinSizeNonBlocking(t *testing.T) {
	reader := &testReader{
		winSize: &WinSize{Row: 30, Col: 100},
	}
	pt := &Prompt{
		reader: reader,
	}

	exitCh := make(chan os.Signal)
	// Unbuffered channel with no reader to simulate the main loop being busy.
	winSizeCh := make(chan *WinSize)
	stopCh := make(chan struct{})
	doneCh := make(chan struct{})

	go func() {
		pt.handleSignals(exitCh, winSizeCh, stopCh)
		close(doneCh)
	}()

	// Send SIGWINCH signal. In previous versions with blocking send, this would hang.
	if err := syscall.Kill(syscall.Getpid(), syscall.SIGWINCH); err != nil {
		t.Fatalf("failed to send SIGWINCH: %v", err)
	}

	// Allow a brief moment for the signal handler to process SIGWINCH.
	time.Sleep(50 * time.Millisecond)

	// Stopping handleSignals must succeed promptly without deadlock.
	stopCh <- struct{}{}

	select {
	case <-doneCh:
	case <-time.After(2 * time.Second):
		t.Fatal("handleSignals deadlocked on SIGWINCH with unbuffered winSizeCh")
	}
}

func TestHandleSignals_RapidConcurrentSignals(t *testing.T) {
	reader := &testReader{
		winSize: &WinSize{Row: 40, Col: 120},
	}
	pt := &Prompt{
		reader: reader,
	}

	exitCh := make(chan os.Signal)
	winSizeCh := make(chan *WinSize, 128)
	stopCh := make(chan struct{})
	doneCh := make(chan struct{})

	go func() {
		pt.handleSignals(exitCh, winSizeCh, stopCh)
		close(doneCh)
	}()

	// Read from winSizeCh in background.
	stopReader := make(chan struct{})
	readerDone := make(chan struct{})
	go func() {
		defer close(readerDone)
		for {
			select {
			case <-stopReader:
				return
			case <-winSizeCh:
			}
		}
	}()

	// Fire rapid SIGWINCH signals concurrently.
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				_ = syscall.Kill(syscall.Getpid(), syscall.SIGWINCH)
				time.Sleep(2 * time.Millisecond)
			}
		}()
	}
	wg.Wait()

	// Stop handleSignals
	stopCh <- struct{}{}

	select {
	case <-doneCh:
	case <-time.After(2 * time.Second):
		t.Fatal("handleSignals deadlocked during rapid concurrent SIGWINCH signals")
	}

	close(stopReader)
	<-readerDone
}

func TestPrompt_Run_ConcurrentSIGWINCH(t *testing.T) {
	reader := &testReader{
		winSize: &WinSize{Row: DefRowCount, Col: DefColCount},
	}
	writer := &testWriter{}
	renderer := NewRenderer()
	renderer.out = writer

	executed := make(chan string, 5)
	pt := &Prompt{
		reader:   reader,
		renderer: renderer,
		buffer:   NewBuffer(),
		executor: func(in string) {
			executed <- in
		},
		history:                NewHistory(),
		completion:             NewCompletionManager(6),
		keyBindMode:            EmacsKeyBind,
		executeOnEnterCallback: DefaultExecuteOnEnterCallback,
	}

	runDone := make(chan struct{})
	go func() {
		pt.Run()
		close(runDone)
	}()

	// Fire SIGWINCH signals concurrently while sending inputs.
	stopSignals := make(chan struct{})
	signalsDone := make(chan struct{})
	go func() {
		defer close(signalsDone)
		for {
			select {
			case <-stopSignals:
				return
			default:
				_ = syscall.Kill(syscall.Getpid(), syscall.SIGWINCH)
				time.Sleep(5 * time.Millisecond)
			}
		}
	}()

	// Queue input to execute a command: first the command text, then Enter.
	reader.appendInput([]byte("test_command"))
	reader.appendInput([]byte{0xa})

	select {
	case cmd := <-executed:
		if cmd != "test_command" {
			t.Fatalf("expected 'test_command', got %q", cmd)
		}
	case <-time.After(3 * time.Second):
		close(stopSignals)
		t.Fatal("prompt runner deadlocked while executing command with concurrent SIGWINCH")
	}

	// Send Control-D to exit the prompt runner.
	reader.appendInput([]byte{0x04})

	close(stopSignals)
	<-signalsDone

	select {
	case <-runDone:
	case <-time.After(3 * time.Second):
		t.Fatal("prompt runner deadlocked on shutdown with concurrent SIGWINCH")
	}
}
