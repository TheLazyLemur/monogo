package monogo

import rl "github.com/gen2brain/raylib-go/raylib"

// SoundSource plays audio clips
type SoundSource struct {
	BaseComponent
	sound       rl.Sound
	music       rl.Music
	isMusic     bool
	Volume      float32
	Loop        bool
	PlayOnStart bool
}

// NewSoundSource creates a sound source from a file (short clips)
func NewSoundSource(path string) *SoundSource {
	return &SoundSource{
		sound:   rl.LoadSound(path),
		isMusic: false,
		Volume:  1.0,
	}
}

// NewMusicSource creates a music source from a file (streams, for longer audio)
func NewMusicSource(path string) *SoundSource {
	return &SoundSource{
		music:   rl.LoadMusicStream(path),
		isMusic: true,
		Volume:  1.0,
	}
}

func (s *SoundSource) Start() {
	if s.PlayOnStart {
		s.Play()
	}
}

func (s *SoundSource) Update(dt float32) {
	if s.isMusic {
		rl.UpdateMusicStream(s.music)
	}
}

// Play starts playback
func (s *SoundSource) Play() {
	if s.isMusic {
		rl.PlayMusicStream(s.music)
		if s.Loop {
			s.music.Looping = true
		}
	} else {
		rl.PlaySound(s.sound)
	}
}

// Stop stops playback
func (s *SoundSource) Stop() {
	if s.isMusic {
		rl.StopMusicStream(s.music)
	} else {
		rl.StopSound(s.sound)
	}
}

// Pause pauses playback
func (s *SoundSource) Pause() {
	if s.isMusic {
		rl.PauseMusicStream(s.music)
	} else {
		rl.PauseSound(s.sound)
	}
}

// Resume resumes playback
func (s *SoundSource) Resume() {
	if s.isMusic {
		rl.ResumeMusicStream(s.music)
	} else {
		rl.ResumeSound(s.sound)
	}
}

// IsPlaying returns true if currently playing
func (s *SoundSource) IsPlaying() bool {
	if s.isMusic {
		return rl.IsMusicStreamPlaying(s.music)
	}
	return rl.IsSoundPlaying(s.sound)
}

// SetVolume sets volume (0.0 to 1.0)
func (s *SoundSource) SetVolume(volume float32) {
	s.Volume = volume
	if s.isMusic {
		rl.SetMusicVolume(s.music, volume)
	} else {
		rl.SetSoundVolume(s.sound, volume)
	}
}

// Unload frees audio resources - call when done
func (s *SoundSource) Unload() {
	if s.isMusic {
		rl.UnloadMusicStream(s.music)
	} else {
		rl.UnloadSound(s.sound)
	}
}
