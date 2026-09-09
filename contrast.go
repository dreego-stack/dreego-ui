package ui

import (
	"fmt"
	"math"
	"strconv"
)

func validateContrast(colors Colors) error {
	pairs := []struct {
		name       string
		foreground string
		background string
		minimum    float64
	}{
		{"Text on Canvas", colors.Text, colors.Canvas, 4.5},
		{"Text on Surface", colors.Text, colors.Surface, 4.5},
		{"Text on SurfaceRaised", colors.Text, colors.SurfaceRaised, 4.5},
		{"TextMuted on Canvas", colors.TextMuted, colors.Canvas, 4.5},
		{"AccentText on Accent", colors.AccentText, colors.Accent, 4.5},
		{"AccentText on AccentHover", colors.AccentText, colors.AccentHover, 4.5},
		{"Focus on Canvas", colors.Focus, colors.Canvas, 3},
		{"Focus on Surface", colors.Focus, colors.Surface, 3},
		{"Success on Canvas", colors.Success, colors.Canvas, 4.5},
		{"Warning on Canvas", colors.Warning, colors.Canvas, 4.5},
		{"Danger on Canvas", colors.Danger, colors.Canvas, 4.5},
	}
	for _, pair := range pairs {
		ratio := contrastRatio(pair.foreground, pair.background)
		if ratio < pair.minimum {
			return invalidTheme(
				"Colors."+pair.name,
				fmt.Sprintf("contrast %.2f:1 is below %.1f:1", ratio, pair.minimum),
			)
		}
	}
	return nil
}

func contrastRatio(foreground, background string) float64 {
	first := relativeLuminance(foreground)
	second := relativeLuminance(background)
	lighter := math.Max(first, second)
	darker := math.Min(first, second)
	return (lighter + 0.05) / (darker + 0.05)
}

func relativeLuminance(color string) float64 {
	red := colorChannel(color[1:3])
	green := colorChannel(color[3:5])
	blue := colorChannel(color[5:7])
	return 0.2126*red + 0.7152*green + 0.0722*blue
}

func colorChannel(value string) float64 {
	parsed, _ := strconv.ParseUint(value, 16, 8)
	channel := float64(parsed) / 255
	if channel <= 0.04045 {
		return channel / 12.92
	}
	return math.Pow((channel+0.055)/1.055, 2.4)
}
