package ui

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

var ErrInvalidTheme = errors.New("dreego-ui: invalid theme")

var (
	themeIDPattern = regexp.MustCompile(`^[a-z][a-z0-9-]{0,62}$`)
	colorPattern   = regexp.MustCompile(`^#[0-9a-fA-F]{6}([0-9a-fA-F]{2})?$`)
	opaquePattern  = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
	lengthPattern  = regexp.MustCompile(`^(0|-?[0-9]+(\.[0-9]+)?(px|rem|em))$`)
	fontPattern    = regexp.MustCompile(`^[a-zA-Z0-9 ,"'_\-]+$`)
	easingPattern  = regexp.MustCompile(`^(linear|ease|ease-in|ease-out|ease-in-out|cubic-bezier\((-?[0-9]+(\.[0-9]+)?),(-?[0-9]+(\.[0-9]+)?),(-?[0-9]+(\.[0-9]+)?),(-?[0-9]+(\.[0-9]+)?)\))$`)
)

func (t Theme) Validate() error {
	if !themeIDPattern.MatchString(t.ID) || isReservedThemeID(t.ID) {
		return invalidTheme("ID", "must start with a letter and contain only lowercase letters, numbers, or hyphens")
	}
	if strings.TrimSpace(t.Name) == "" || utf8.RuneCountInString(t.Name) > 80 {
		return invalidTheme("Name", "must contain between 1 and 80 characters")
	}
	if utf8.RuneCountInString(t.Description) > 240 {
		return invalidTheme("Description", "must not exceed 240 characters")
	}
	if err := validateColors(t.Colors); err != nil {
		return err
	}
	if err := validateTypography(t.Typography); err != nil {
		return err
	}
	if err := validateShape(t.Shape); err != nil {
		return err
	}
	if err := validateElevation(t.Elevation); err != nil {
		return err
	}
	return validateMotion(t.Motion)
}

func isReservedThemeID(id string) bool {
	return id == "auto" || id == "default" || id == "system"
}

func validateColors(colors Colors) error {
	values := []struct{ name, value string }{
		{"Canvas", colors.Canvas}, {"Surface", colors.Surface},
		{"SurfaceRaised", colors.SurfaceRaised}, {"Text", colors.Text},
		{"TextMuted", colors.TextMuted}, {"Border", colors.Border},
		{"Accent", colors.Accent}, {"AccentHover", colors.AccentHover},
		{"AccentText", colors.AccentText}, {"Focus", colors.Focus},
		{"Success", colors.Success}, {"Warning", colors.Warning},
		{"Danger", colors.Danger},
	}
	for _, item := range values {
		if !opaquePattern.MatchString(item.value) {
			return invalidTheme("Colors."+item.name, "must be an opaque six-digit hex color")
		}
	}
	return validateContrast(colors)
}

func validateTypography(typography Typography) error {
	if !fontPattern.MatchString(typography.FontSans) {
		return invalidTheme("Typography.FontSans", "contains unsupported characters")
	}
	if !fontPattern.MatchString(typography.FontMono) {
		return invalidTheme("Typography.FontMono", "contains unsupported characters")
	}
	if !lengthPattern.MatchString(typography.BaseSize) {
		return invalidTheme("Typography.BaseSize", "must be a supported CSS length")
	}
	if typography.LineHeight < 1 || typography.LineHeight > 2.5 {
		return invalidTheme("Typography.LineHeight", "must be between 1 and 2.5")
	}
	weights := []struct {
		name  string
		value int
	}{
		{"WeightRegular", typography.WeightRegular},
		{"WeightMedium", typography.WeightMedium},
		{"WeightStrong", typography.WeightStrong},
	}
	for _, item := range weights {
		if item.value < 100 || item.value > 900 {
			return invalidTheme("Typography."+item.name, "must be between 100 and 900")
		}
	}
	return nil
}

func validateShape(shape Shape) error {
	values := []struct{ name, value string }{
		{"RadiusSmall", shape.RadiusSmall}, {"RadiusMedium", shape.RadiusMedium},
		{"RadiusLarge", shape.RadiusLarge}, {"RadiusPill", shape.RadiusPill},
	}
	for _, item := range values {
		if !lengthPattern.MatchString(item.value) || strings.HasPrefix(item.value, "-") {
			return invalidTheme("Shape."+item.name, "must be a non-negative supported CSS length")
		}
	}
	return nil
}

func validateElevation(elevation Elevation) error {
	for _, item := range []struct {
		name   string
		shadow Shadow
	}{{"Small", elevation.Small}, {"Medium", elevation.Medium}, {"Large", elevation.Large}} {
		if err := validateShadow(item.shadow); err != nil {
			return fmt.Errorf("%w: Elevation.%s: %v", ErrInvalidTheme, item.name, err)
		}
	}
	return nil
}

func validateShadow(shadow Shadow) error {
	values := []struct{ name, value string }{
		{"X", shadow.X}, {"Y", shadow.Y}, {"Blur", shadow.Blur}, {"Spread", shadow.Spread},
	}
	for _, item := range values {
		if !lengthPattern.MatchString(item.value) {
			return fmt.Errorf("%s must be a supported CSS length", item.name)
		}
	}
	if strings.HasPrefix(shadow.Blur, "-") {
		return errors.New("Blur must be non-negative")
	}
	if !colorPattern.MatchString(shadow.Color) {
		return errors.New("Color must be a six- or eight-digit hex color")
	}
	return nil
}

func validateMotion(motion Motion) error {
	if motion.Fast < time.Millisecond || motion.Normal < time.Millisecond || motion.Slow < time.Millisecond {
		return invalidTheme("Motion", "durations must be at least one millisecond")
	}
	if motion.Fast > motion.Normal || motion.Normal > motion.Slow {
		return invalidTheme("Motion", "durations must be ordered from fast to slow")
	}
	if !easingPattern.MatchString(motion.Easing) {
		return invalidTheme("Motion.Easing", "must be a supported easing value")
	}
	return nil
}

func invalidTheme(field, reason string) error {
	return fmt.Errorf("%w: %s %s", ErrInvalidTheme, field, reason)
}
