package tomate

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/0xAX/notificator"
)

// Pomodoro contains the key durations
type Pomodoro struct {
	Focus         time.Duration `json:"focus"`
	ShortPause    time.Duration `json:"short_pause"`
	LongPause     time.Duration `json:"long_pause"`
	NumberSprints int           `json:"number_sprints"`
	ntf           *notificator.Notificator
}

// New builds a new Pomodoro from the durations of the different steps in minutes
func New(focus, short, long, nsprints int) Pomodoro {
	return Pomodoro{
		time.Duration(focus) * time.Minute,
		time.Duration(short) * time.Minute,
		time.Duration(long) * time.Minute,
		nsprints,
		notificator.New(notificator.Options{}),
	}
}

// Launch starts a timer
func (p *Pomodoro) launch(d time.Duration) time.Duration {
	start := time.Now()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case <-time.After(d):
		return time.Since(start)
	case <-sigs:
		return time.Since(start)
	}
}

func (p *Pomodoro) sprint(w io.Writer) {
	if _, err := fmt.Fprintln(w, "Starting sprint"); err != nil {
		log.Fatalf("Error writing start of sprint: %s", err.Error())
	}
	d := p.launch(p.Focus)
	_, err := fmt.Fprintf(
		w, "End of sprint after %v\n", d,
	)
	if err != nil {
		log.Fatalf("Error writing end of sprint: %s\n", err.Error())
	}
	p.ntf.Push("End of sprint", "Take some rest", "", notificator.UR_NORMAL)
}

func (p *Pomodoro) pause(w io.Writer, long bool) {
	fmt.Fprintln(w, "Starting pause")
	var d time.Duration
	if long {
		d = p.launch(p.LongPause)
	} else {
		d = p.launch(p.ShortPause)
	}
	_, err := fmt.Fprintf(
		w, "End of pause after %v\n", d,
	)
	if err != nil {
		log.Fatalf("Error writing end of pause: %s\n", err.Error())
	}
	p.ntf.Push("End of pause", "Back to Work", "", notificator.UR_NORMAL)
}

// Loop starts looping between the three states
func (p *Pomodoro) Loop(w io.Writer) {
	fmt.Fprintf(
		w,
		"--------Starting tomate--------\nshort: %v - long: %v - num sprints: %v\n",
		p.ShortPause, p.LongPause, p.NumberSprints,
	)
	for {
		for sprint := 0; sprint < p.NumberSprints-1; sprint++ {
			p.sprint(w)
			p.pause(w, false)
		}
		p.sprint(w)
		p.pause(w, true)
	}
}
