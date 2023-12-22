package shaker

import (
	"math/rand"
	"time"

	"github.com/peterhellberg/gfx"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type ScreenShaker interface {
	Shake()
	SetShakeIntensity(maxIntensity float64)
	Update() error
	GetOffsets() (dx, dy float64)
}

type Shaker struct {
	ShakeIntensity    float64
	MaxShakeIntensity float64
	ShakeDuration     float64
}

func New() ScreenShaker {
	return &Shaker{
		MaxShakeIntensity: 10.0,
		ShakeIntensity:    1.0,
		ShakeDuration:     1.0,
	}
}

func (s *Shaker) Shake() {
	s.ShakeIntensity = 0.0
}

// 10.0 = Very Intense, 1.0  = Non Existent
func (s *Shaker) SetShakeIntensity(maxIntensity float64) {
	s.MaxShakeIntensity = maxIntensity
	s.ShakeDuration = maxIntensity / 10.0
}

func (s *Shaker) Update() error {
	s.ShakeIntensity += 1 / 60.0 // 60 FPS fixed
	return nil
}

func (s *Shaker) GetOffsets() (dx, dy float64) {
	if s.ShakeIntensity < 1 {
		lerped := gfx.Lerp(s.ShakeDuration, 0, s.ShakeIntensity)
		amplitude := s.MaxShakeIntensity * lerped
		dx = amplitude * (2*rand.Float64() - 1)
		dy = amplitude * (2*rand.Float64() - 1)
	}
	return
}
