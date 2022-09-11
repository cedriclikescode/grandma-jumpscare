package main

import (
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	spriteSheet      *ebiten.Image
	mplusNormalFont  font.Face
	grandmaJumpscare = false
	grandmaJumpscareStart time.Time
	audioContext *audio.Context
	audioPlayer  *audio.Player
)

type Game struct {}

func (g *Game) Update() error {
	return nil
}

func AngryColor(percent float64) color.RGBA {
	percent = 100 - percent
	if percent <= 50 {
		return color.RGBA{255, uint8(percent * 5.1), 0, 255}
	} else {
		return color.RGBA{uint8(255 - 5.1 * percent + 255), 255, 0, 255}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !grandmaJumpscare && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		grandmaJumpscareStart = time.Now()
		grandmaJumpscare = true
		ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
		audioPlayer.Rewind()
		audioPlayer.Play()
	}

	op := &ebiten.DrawImageOptions{}
	if grandmaJumpscare {
		frame := int(time.Since(grandmaJumpscareStart).Seconds() * 30)
		if frame > 33 {
			frame = 33
			grandmaJumpscare = false
			ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMinimum)
			ebiten.ScheduleFrame()
		}
		x := frame % 7 * 290
		y := int(float64(frame / 7) * 476)
		screen.DrawImage(spriteSheet.SubImage(image.Rect(x, y, 288 + x, 474 + y)).(*ebiten.Image), op)
		return
	}

	screen.DrawImage(spriteSheet.SubImage(image.Rect(0, 0, 288, 540)).(*ebiten.Image), op)
	ebitenutil.DrawRect(screen, 0, 477, 288, 580, color.Gray{127})
	grandmaText := "      Click to make \n grandma jumpscare!"
	textRectangle := text.BoundString(mplusNormalFont, grandmaText)
	text.Draw(screen, grandmaText, mplusNormalFont, int(278 - textRectangle.Max.X) / 2, 510, color.RGBA{0, 0, 0, 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 288, 570
}

func main() {
	var _ image.Image
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	spriteSheet, _, err = ebitenutil.NewImageFromFile("output.png")
	if err != nil {
		log.Fatal(err)
	}
	var dpi float64 = 72 * 1

	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    30,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Initialize audio context.
	audioContext = audio.NewContext(44100)

	// In this example, embedded resource "Jab_wav" is used.
	//
	// If you want to use a wav file, open this and pass the file stream to wav.Decode.
	// Note that file's Close() should not be closed here
	// since audio.Player manages stream state.
	//
	//     f, err := os.Open("jab.wav")
	//     if err != nil {
	//         return err
	//     }
	//
	//     d, err := wav.DecodeWithoutResampling(f)
	//     ...

	// Decode wav-formatted data and retrieve decoded PCM stream.
	//d, err := wav.DecodeWithoutResampling(bytes.NewReader(raudio.Jab_wav))
	f, err := os.Open("grandma jumpscare.wav")
	if err != nil {
		log.Fatal(err)
	}
	d, err := wav.DecodeWithoutResampling(f)
	if err != nil {
		log.Fatal(err)
	}

	audioPlayer, err = audioContext.NewPlayer(d)
	if err != nil {
		log.Fatal(err)
	}




	//ebiten.SetScreenClearedEveryFrame(false)
	//ebiten.SetScreenFilterEnabled(false)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMinimum)
	ebiten.SetWindowSize(288, 570)
	ebiten.SetWindowTitle("Don't click on grandma!")
	rand.Seed(time.Now().UnixNano())
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
	ebiten.ScheduleFrame()
}
