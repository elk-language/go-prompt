//go:build unix
// +build unix

package prompt

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/elk-language/go-prompt/debug"
)

func (p *Prompt) handleSignals(exitCh chan os.Signal, winSizeCh chan *WinSize, stop chan struct{}) {
	in := p.reader
	sigCh := make(chan os.Signal, 1)
	defer signal.Stop(sigCh)
	signal.Notify(
		sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGWINCH,
	)

	for {
		select {
		case <-stop:
			debug.Log("stop handleSignals")
			return
		case s := <-sigCh:
			switch s {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT: // kill -SIGQUIT XXXX
				debug.Log(fmt.Sprintf("Catch %s", s))
				exitCh <- s
			case syscall.SIGWINCH:
				debug.Log("Catch SIGWINCH")

				select {
				case winSizeCh <- in.GetWinSize():
				default:
				}
			}
		}
	}
}
