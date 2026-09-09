package ui

import (
	"strconv"
	"strings"
	"time"
)

func (r *Registry) Stylesheet() string {
	if r == nil {
		return ""
	}
	var css strings.Builder
	writeThemeRule(&css, ":root", r.Default())
	for _, id := range r.order {
		writeThemeRule(&css, `[data-dreego-theme="`+id+`"]`, r.themes[id])
	}
	css.WriteString("@media (prefers-reduced-motion:reduce){:root{--dreego-motion-fast:0ms;--dreego-motion-normal:0ms;--dreego-motion-slow:0ms;}}")
	return css.String()
}

func writeThemeRule(css *strings.Builder, selector string, theme Theme) {
	css.WriteString(selector)
	css.WriteByte('{')
	writeToken(css, "color-canvas", theme.Colors.Canvas)
	writeToken(css, "color-surface", theme.Colors.Surface)
	writeToken(css, "color-surface-raised", theme.Colors.SurfaceRaised)
	writeToken(css, "color-text", theme.Colors.Text)
	writeToken(css, "color-text-muted", theme.Colors.TextMuted)
	writeToken(css, "color-border", theme.Colors.Border)
	writeToken(css, "color-accent", theme.Colors.Accent)
	writeToken(css, "color-accent-hover", theme.Colors.AccentHover)
	writeToken(css, "color-accent-text", theme.Colors.AccentText)
	writeToken(css, "color-focus", theme.Colors.Focus)
	writeToken(css, "color-success", theme.Colors.Success)
	writeToken(css, "color-warning", theme.Colors.Warning)
	writeToken(css, "color-danger", theme.Colors.Danger)
	writeToken(css, "type-font-sans", theme.Typography.FontSans)
	writeToken(css, "type-font-mono", theme.Typography.FontMono)
	writeToken(css, "type-base-size", theme.Typography.BaseSize)
	writeToken(css, "type-line-height", strconv.FormatFloat(theme.Typography.LineHeight, 'f', -1, 64))
	writeToken(css, "type-weight-regular", strconv.Itoa(theme.Typography.WeightRegular))
	writeToken(css, "type-weight-medium", strconv.Itoa(theme.Typography.WeightMedium))
	writeToken(css, "type-weight-strong", strconv.Itoa(theme.Typography.WeightStrong))
	writeToken(css, "shape-radius-small", theme.Shape.RadiusSmall)
	writeToken(css, "shape-radius-medium", theme.Shape.RadiusMedium)
	writeToken(css, "shape-radius-large", theme.Shape.RadiusLarge)
	writeToken(css, "shape-radius-pill", theme.Shape.RadiusPill)
	writeToken(css, "elevation-small", shadowCSS(theme.Elevation.Small))
	writeToken(css, "elevation-medium", shadowCSS(theme.Elevation.Medium))
	writeToken(css, "elevation-large", shadowCSS(theme.Elevation.Large))
	writeToken(css, "motion-fast", durationCSS(theme.Motion.Fast))
	writeToken(css, "motion-normal", durationCSS(theme.Motion.Normal))
	writeToken(css, "motion-slow", durationCSS(theme.Motion.Slow))
	writeToken(css, "motion-easing", theme.Motion.Easing)
	css.WriteByte('}')
}

func writeToken(css *strings.Builder, name, value string) {
	css.WriteString("--dreego-")
	css.WriteString(name)
	css.WriteByte(':')
	css.WriteString(value)
	css.WriteByte(';')
}

func shadowCSS(shadow Shadow) string {
	return strings.Join([]string{shadow.X, shadow.Y, shadow.Blur, shadow.Spread, shadow.Color}, " ")
}

func durationCSS(duration time.Duration) string {
	return strconv.FormatInt(duration.Milliseconds(), 10) + "ms"
}
