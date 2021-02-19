// A dummy music player
// Controls available to the user: "play", "pause", "prev", "next", "exit"

package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/ini.v1"
)

// Add controls.
const (
	exit = iota
	stop
	play
	pause
	previous
	next
)

// Song structure.
type Song struct {
	Name     string
	Duration time.Duration
}

// MP3Player structure.
type MP3Player struct {
	Play         bool
	Playing      int
	DurationLeft int
	Songs        map[int]Song

	// internal: channels are for internal use
	channel struct {
		stop             chan bool
		pause            chan bool
		chanDurationLeft chan int
	}
}

// AddMP3Songs adds songs.
func AddMP3Songs(cfg *ini.File, songs []string, s map[int]Song) {
	var duration time.Duration
	for index, song := range songs {
		if d, err := cfg.Section("songs").Key(song).Int(); err == nil {
			duration = time.Duration(d) * time.Second
		}
		s[index+1] = Song{
			Name:     song,
			Duration: duration,
		}
	}
}

// Buttons function listens for inputs from the user.
func Buttons(control chan<- int) {
	pressed := ""
	for {
		fmt.Scanln(&pressed)
		switch pressed {
		case "play":
			control <- play
		case "pause":
			control <- pause
		case "prev":
			control <- previous
		case "next":
			control <- next
		case "exit":
			control <- exit
			return
		default:
			fmt.Printf("P: Available controls: [play, pause, prev, next, exit]\n")
		}
	}
}

// Play controls the play/pause for songs
func Play(name string, d time.Duration, controlCh chan<- int, pauseCh, stopCh <-chan bool, durationLeftCh chan<- int) {
	fmt.Printf("P: Playing song: %s, seek: %v\n", name, int(d.Seconds()))
	tick := time.NewTicker(1 * time.Second)
	duration := int(d.Seconds())
	for {
		select {
		case <-tick.C:
			duration--
			if duration <= 0 {
				controlCh <- stop
				tick.Stop()
				fmt.Printf("P: Finished song: %s\n", name)
				return
			}
		case p, ok := <-pauseCh:
			if p && ok {
				tick.Stop()
				durationLeftCh <- duration
				return
			}
		case s, ok := <-stopCh:
			if s && ok {
				tick.Stop()
				return
			}
		}
	}
}

// MusicPlayerController is the controller for music player.
func MusicPlayerController(mp3 MP3Player, control chan int) {
	var dur time.Duration
	for c := range control {
		switch c {
		// Play
		case play:
			if !mp3.Play {
				mp3.Play = true
				switch mp3.Playing {
				case 0:
					mp3.Playing = 1
					dur = mp3.Songs[mp3.Playing].Duration
				default:
					if mp3.DurationLeft == 0 {
						dur = mp3.Songs[mp3.Playing].Duration
					} else {
						dur = time.Duration(mp3.DurationLeft) * time.Second
					}
				}
				go Play(
					mp3.Songs[mp3.Playing].Name,
					dur,
					control,
					mp3.channel.pause,
					mp3.channel.stop,
					mp3.channel.chanDurationLeft,
				)
			} else {
				fmt.Printf("P: Song is already playing\n")
			}
		// Pause
		case pause:
			if mp3.Play {
				mp3.Play = false
				mp3.channel.pause <- true
				mp3.DurationLeft = <-mp3.channel.chanDurationLeft
				fmt.Printf("P: Paused at seek: %v\n", mp3.DurationLeft)
			} else {
				fmt.Printf("P: No song is playing, currently\n")
			}
		// Previous
		case previous:
			switch mp3.Playing {
			case 1:
				mp3.Playing = len(mp3.Songs)
			default:
				mp3.Playing--
			}
			mp3.DurationLeft = 0
			if mp3.Play {
				mp3.channel.stop <- true
				mp3.Play = false
				control <- play
			}
			fmt.Printf("P: Select previous song: %s\n", mp3.Songs[mp3.Playing].Name)
		// Next
		case next:
			switch mp3.Playing {
			case len(mp3.Songs):
				mp3.Playing = 1
			default:
				mp3.Playing++
			}
			mp3.DurationLeft = 0
			if mp3.Play {
				mp3.channel.stop <- true
				mp3.Play = false
				control <- play
			}
			fmt.Printf("P: Select next song: %s\n", mp3.Songs[mp3.Playing].Name)
		// Stop (Finished Playing)
		case stop:
			mp3.Play = false
			mp3.DurationLeft = 0
		// Exit Player
		case exit:
			close(control)
			close(mp3.channel.stop)
			close(mp3.channel.pause)
			close(mp3.channel.chanDurationLeft)
			fmt.Printf("P: Exit\n")
			return
		}
	}
}

// Implements the Stringer interface
func (s Song) String() string {
	return fmt.Sprintf("=> Song: %s, Duration: %ds", s.Name, s.Duration/time.Second)
}

func main() {
	var internalPlaylist bool
	cfg, err := ini.Load("songs.ini")
	if err != nil {
		fmt.Println("P: Playing from the internal playlist")
		internalPlaylist = true
	}

	var (
		mp3     = new(MP3Player)
		control = make(chan int, 1)
	)

	if !internalPlaylist {
		// Read from the configuration file
		songs := cfg.Section("songs").KeyStrings()
		mp3.Songs = make(map[int]Song, len(songs))
		AddMP3Songs(cfg, songs, mp3.Songs)
	} else {
		// Feed from the internal playlist
		mp3.Songs = make(map[int]Song, 3)
		mp3.Songs[1] = Song{Name: "highway to hell", Duration: 5 * time.Second}
		mp3.Songs[2] = Song{Name: "stairway to heaven", Duration: 10 * time.Second}
		mp3.Songs[3] = Song{Name: "cardigan", Duration: 30 * time.Second}
	}

	if len(mp3.Songs) == 0 {
		fmt.Println("P: No songs available")
		os.Exit(0)
	}

	fmt.Printf("P: Playlist:\n")
	for i := 1; i <= len(mp3.Songs); i++ {
		fmt.Println(mp3.Songs[i])
	}

	mp3.channel.chanDurationLeft = make(chan int)
	mp3.channel.stop, mp3.channel.pause = make(chan bool), make(chan bool)

	go Buttons(control)
	MusicPlayerController(*mp3, control)
}
