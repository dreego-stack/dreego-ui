package ui

import "time"

type Theme struct {
	ID          string
	Name        string
	Description string
	Colors      Colors
	Typography  Typography
	Shape       Shape
	Elevation   Elevation
	Motion      Motion
}

type Colors struct {
	Canvas        string
	Surface       string
	SurfaceRaised string
	Text          string
	TextMuted     string
	Border        string
	Accent        string
	AccentHover   string
	AccentText    string
	Focus         string
	Success       string
	Warning       string
	Danger        string
}

type Typography struct {
	FontSans      string
	FontMono      string
	BaseSize      string
	LineHeight    float64
	WeightRegular int
	WeightMedium  int
	WeightStrong  int
}

type Shape struct {
	RadiusSmall  string
	RadiusMedium string
	RadiusLarge  string
	RadiusPill   string
}

type Shadow struct {
	X      string
	Y      string
	Blur   string
	Spread string
	Color  string
}

type Elevation struct {
	Small  Shadow
	Medium Shadow
	Large  Shadow
}

type Motion struct {
	Fast   time.Duration
	Normal time.Duration
	Slow   time.Duration
	Easing string
}

func WhiteTheme() Theme {
	return Theme{
		ID:          "white",
		Name:        "White",
		Description: "A calm light theme with cool neutral surfaces.",
		Colors: Colors{
			Canvas:        "#ffffff",
			Surface:       "#f4f7f7",
			SurfaceRaised: "#ffffff",
			Text:          "#17191c",
			TextMuted:     "#5f6872",
			Border:        "#d9dfe2",
			Accent:        "#006d77",
			AccentHover:   "#005a63",
			AccentText:    "#ffffff",
			Focus:         "#007c91",
			Success:       "#147a4f",
			Warning:       "#9a5b00",
			Danger:        "#b4232c",
		},
		Typography: defaultTypography(),
		Shape:      defaultShape(),
		Elevation: Elevation{
			Small:  Shadow{X: "0", Y: "0.0625rem", Blur: "0.1875rem", Spread: "0", Color: "#10182014"},
			Medium: Shadow{X: "0", Y: "0.375rem", Blur: "1.25rem", Spread: "-0.25rem", Color: "#1018201f"},
			Large:  Shadow{X: "0", Y: "1.25rem", Blur: "3rem", Spread: "-0.75rem", Color: "#1018202e"},
		},
		Motion: defaultMotion(),
	}
}

func BlackTheme() Theme {
	return Theme{
		ID:          "black",
		Name:        "Black",
		Description: "A deliberate true-black theme with raised near-black surfaces.",
		Colors: Colors{
			Canvas:        "#000000",
			Surface:       "#101214",
			SurfaceRaised: "#181b1e",
			Text:          "#f7f9fa",
			TextMuted:     "#a7b0b7",
			Border:        "#343a3f",
			Accent:        "#62d0c9",
			AccentHover:   "#8de0dc",
			AccentText:    "#001f22",
			Focus:         "#77dde5",
			Success:       "#55c991",
			Warning:       "#f1b95b",
			Danger:        "#ff7b83",
		},
		Typography: defaultTypography(),
		Shape:      defaultShape(),
		Elevation: Elevation{
			Small:  Shadow{X: "0", Y: "0.0625rem", Blur: "0.1875rem", Spread: "0", Color: "#00000052"},
			Medium: Shadow{X: "0", Y: "0.375rem", Blur: "1.25rem", Spread: "-0.25rem", Color: "#0000008f"},
			Large:  Shadow{X: "0", Y: "1.25rem", Blur: "3rem", Spread: "-0.75rem", Color: "#000000b8"},
		},
		Motion: defaultMotion(),
	}
}

func defaultTypography() Typography {
	return Typography{
		FontSans:      `-apple-system, BlinkMacSystemFont, "SF Pro Text", "Segoe UI", sans-serif`,
		FontMono:      `"SFMono-Regular", Consolas, "Liberation Mono", monospace`,
		BaseSize:      "1rem",
		LineHeight:    1.5,
		WeightRegular: 400,
		WeightMedium:  500,
		WeightStrong:  650,
	}
}

func defaultShape() Shape {
	return Shape{
		RadiusSmall:  "0.5rem",
		RadiusMedium: "0.875rem",
		RadiusLarge:  "1.25rem",
		RadiusPill:   "999rem",
	}
}

func defaultMotion() Motion {
	return Motion{
		Fast:   120 * time.Millisecond,
		Normal: 200 * time.Millisecond,
		Slow:   320 * time.Millisecond,
		Easing: "cubic-bezier(0.2,0,0,1)",
	}
}
