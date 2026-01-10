// Package themes provides a theming system for the ArloCode TUI application.
// It wraps catppuccin/go and provides convenient access to themed lipgloss styles.
package themes

import (
	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/lipgloss"
)

// ColorPalette defines all the colors needed for a theme
type ColorPalette struct {
	// Base colors (backgrounds)
	Base     string
	Mantle   string
	Crust    string
	Surface0 string
	Surface1 string
	Surface2 string

	// Overlay colors
	Overlay0 string
	Overlay1 string
	Overlay2 string

	// Text colors
	Text     string
	Subtext0 string
	Subtext1 string

	// Accent colors
	Rosewater string
	Flamingo  string
	Pink      string
	Mauve     string
	Red       string
	Maroon    string
	Peach     string
	Yellow    string
	Green     string
	Teal      string
	Sky       string
	Sapphire  string
	Blue      string
	Lavender  string
}

// Theme wraps a color palette and provides lipgloss-compatible colors
type Theme struct {
	palette ColorPalette
	name    string
}

// fromCatppuccin creates a ColorPalette from a Catppuccin flavor
func fromCatppuccin(f catppuccin.Flavor) ColorPalette {
	return ColorPalette{
		Base:      f.Base().Hex,
		Mantle:    f.Mantle().Hex,
		Crust:     f.Crust().Hex,
		Surface0:  f.Surface0().Hex,
		Surface1:  f.Surface1().Hex,
		Surface2:  f.Surface2().Hex,
		Overlay0:  f.Overlay0().Hex,
		Overlay1:  f.Overlay1().Hex,
		Overlay2:  f.Overlay2().Hex,
		Text:      f.Text().Hex,
		Subtext0:  f.Subtext0().Hex,
		Subtext1:  f.Subtext1().Hex,
		Rosewater: f.Rosewater().Hex,
		Flamingo:  f.Flamingo().Hex,
		Pink:      f.Pink().Hex,
		Mauve:     f.Mauve().Hex,
		Red:       f.Red().Hex,
		Maroon:    f.Maroon().Hex,
		Peach:     f.Peach().Hex,
		Yellow:    f.Yellow().Hex,
		Green:     f.Green().Hex,
		Teal:      f.Teal().Hex,
		Sky:       f.Sky().Hex,
		Sapphire:  f.Sapphire().Hex,
		Blue:      f.Blue().Hex,
		Lavender:  f.Lavender().Hex,
	}
}

// Available themes - Catppuccin
var (
	Mocha     = Theme{palette: fromCatppuccin(catppuccin.Mocha), name: "Mocha"}
	Frappe    = Theme{palette: fromCatppuccin(catppuccin.Frappe), name: "Frappé"}
	Macchiato = Theme{palette: fromCatppuccin(catppuccin.Macchiato), name: "Macchiato"}
	Latte     = Theme{palette: fromCatppuccin(catppuccin.Latte), name: "Latte"}

	// TokyoNight variants
	TokyoNight = Theme{
		name: "Tokyo Night",
		palette: ColorPalette{
			// Backgrounds
			Base:     "#1a1b26",
			Mantle:   "#16161e",
			Crust:    "#13131a",
			Surface0: "#292e42",
			Surface1: "#3b4261",
			Surface2: "#414868",
			// Overlays
			Overlay0: "#545c7e",
			Overlay1: "#565f89",
			Overlay2: "#737aa2",
			// Text
			Text:     "#c0caf5",
			Subtext0: "#a9b1d6",
			Subtext1: "#9aa5ce",
			// Accents
			Rosewater: "#f7768e",
			Flamingo:  "#ff007c",
			Pink:      "#bb9af7",
			Mauve:     "#9d7cd8",
			Red:       "#f7768e",
			Maroon:    "#db4b4b",
			Peach:     "#ff9e64",
			Yellow:    "#e0af68",
			Green:     "#9ece6a",
			Teal:      "#1abc9c",
			Sky:       "#7dcfff",
			Sapphire:  "#2ac3de",
			Blue:      "#7aa2f7",
			Lavender:  "#b4f9f8",
		},
	}

	TokyoNightStorm = Theme{
		name: "Tokyo Night Storm",
		palette: ColorPalette{
			// Backgrounds
			Base:     "#24283b",
			Mantle:   "#1f2335",
			Crust:    "#1a1e30",
			Surface0: "#292e42",
			Surface1: "#3b4261",
			Surface2: "#414868",
			// Overlays
			Overlay0: "#545c7e",
			Overlay1: "#565f89",
			Overlay2: "#737aa2",
			// Text
			Text:     "#c0caf5",
			Subtext0: "#a9b1d6",
			Subtext1: "#9aa5ce",
			// Accents
			Rosewater: "#f7768e",
			Flamingo:  "#ff007c",
			Pink:      "#bb9af7",
			Mauve:     "#9d7cd8",
			Red:       "#f7768e",
			Maroon:    "#db4b4b",
			Peach:     "#ff9e64",
			Yellow:    "#e0af68",
			Green:     "#9ece6a",
			Teal:      "#1abc9c",
			Sky:       "#7dcfff",
			Sapphire:  "#2ac3de",
			Blue:      "#7aa2f7",
			Lavender:  "#b4f9f8",
		},
	}

	TokyoNightDay = Theme{
		name: "Tokyo Night Day",
		palette: ColorPalette{
			// Backgrounds (light theme)
			Base:     "#e1e2e7",
			Mantle:   "#d5d6db",
			Crust:    "#c9cad0",
			Surface0: "#c4c5cb",
			Surface1: "#b4b5b9",
			Surface2: "#9699a3",
			// Overlays
			Overlay0: "#848795",
			Overlay1: "#6c6f7f",
			Overlay2: "#5a5d6b",
			// Text
			Text:     "#3760bf",
			Subtext0: "#4c5067",
			Subtext1: "#5a5d6b",
			// Accents
			Rosewater: "#f52a65",
			Flamingo:  "#d23d7c",
			Pink:      "#9854f1",
			Mauve:     "#7847bd",
			Red:       "#f52a65",
			Maroon:    "#c64343",
			Peach:     "#b15c00",
			Yellow:    "#8c6c3e",
			Green:     "#587539",
			Teal:      "#118c74",
			Sky:       "#0f91b3",
			Sapphire:  "#007197",
			Blue:      "#2e7de9",
			Lavender:  "#68a3dc",
		},
	}

	// Default theme
	Current = TokyoNight
)

// AllThemes returns a list of all available themes
func AllThemes() []Theme {
	return []Theme{
		Mocha, Frappe, Macchiato, Latte,
		TokyoNight, TokyoNightStorm, TokyoNightDay,
	}
}

// SetTheme changes the current theme
func SetTheme(t Theme) {
	Current = t
}

// Name returns the theme name
func (t Theme) Name() string {
	return t.name
}

// =============================================================================
// GLAMOUR STYLE - Returns a glamour StyleConfig using theme colors
// =============================================================================

// GlamourStyle returns a glamour ansi.StyleConfig using the theme's colors.
// This can be used with glamour.WithStyles() to render markdown with themed colors.
func (t Theme) GlamourStyle() ansi.StyleConfig {
	// Helper to create string pointers
	sp := func(s string) *string { return &s }
	// Helper to create bool pointers
	bp := func(b bool) *bool { return &b }
	// Helper to create uint pointers
	up := func(u uint) *uint { return &u }

	return ansi.StyleConfig{
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix: "\n",
				BlockSuffix: "\n",
				Color:       sp(t.palette.Text),
			},
			Margin: up(0),
		},
		BlockQuote: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color:  sp(t.palette.Overlay1),
				Italic: bp(true),
			},
			Indent:      up(1),
			IndentToken: sp("│ "),
		},
		List: ansi.StyleList{
			LevelIndent: 2,
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Color: sp(t.palette.Text),
				},
			},
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix: "\n",
				Color:       sp(t.palette.Mauve),
				Bold:        bp(true),
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				Color:           sp(t.palette.Base),
				BackgroundColor: sp(t.palette.Mauve),
				Bold:            bp(true),
			},
		},
		H2: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "## ",
				Color:  sp(t.palette.Pink),
				Bold:   bp(true),
			},
		},
		H3: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "### ",
				Color:  sp(t.palette.Lavender),
				Bold:   bp(true),
			},
		},
		H4: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "#### ",
				Color:  sp(t.palette.Blue),
			},
		},
		H5: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "##### ",
				Color:  sp(t.palette.Sapphire),
			},
		},
		H6: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "###### ",
				Color:  sp(t.palette.Teal),
				Bold:   bp(false),
			},
		},
		Text: ansi.StylePrimitive{
			Color: sp(t.palette.Text),
		},
		Strikethrough: ansi.StylePrimitive{
			CrossedOut: bp(true),
		},
		Emph: ansi.StylePrimitive{
			Italic: bp(true),
			Color:  sp(t.palette.Text),
		},
		Strong: ansi.StylePrimitive{
			Bold:  bp(true),
			Color: sp(t.palette.Text),
		},
		HorizontalRule: ansi.StylePrimitive{
			Color:  sp(t.palette.Overlay0),
			Format: "\n--------\n",
		},
		Item: ansi.StylePrimitive{
			BlockPrefix: "• ",
			Color:       sp(t.palette.Text),
		},
		Enumeration: ansi.StylePrimitive{
			BlockPrefix: ". ",
			Color:       sp(t.palette.Text),
		},
		Task: ansi.StyleTask{
			StylePrimitive: ansi.StylePrimitive{
				Color: sp(t.palette.Text),
			},
			Ticked:   "[✓] ",
			Unticked: "[ ] ",
		},
		Link: ansi.StylePrimitive{
			Color:     sp(t.palette.Blue),
			Underline: bp(true),
		},
		LinkText: ansi.StylePrimitive{
			Color: sp(t.palette.Sapphire),
			Bold:  bp(true),
		},
		Image: ansi.StylePrimitive{
			Color:     sp(t.palette.Lavender),
			Underline: bp(true),
		},
		ImageText: ansi.StylePrimitive{
			Color:  sp(t.palette.Overlay1),
			Format: "Image: {{.text}} →",
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				Color:           sp(t.palette.Peach),
				BackgroundColor: sp(t.palette.Surface0),
			},
		},
		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Color: sp(t.palette.Text),
				},
				Margin: up(2),
			},
			Chroma: &ansi.Chroma{
				Text: ansi.StylePrimitive{
					Color: sp(t.palette.Text),
				},
				Error: ansi.StylePrimitive{
					Color:           sp(t.palette.Text),
					BackgroundColor: sp(t.palette.Red),
				},
				Comment: ansi.StylePrimitive{
					Color: sp(t.palette.Overlay1),
				},
				CommentPreproc: ansi.StylePrimitive{
					Color: sp(t.palette.Peach),
				},
				Keyword: ansi.StylePrimitive{
					Color: sp(t.palette.Mauve),
				},
				KeywordReserved: ansi.StylePrimitive{
					Color: sp(t.palette.Mauve),
				},
				KeywordNamespace: ansi.StylePrimitive{
					Color: sp(t.palette.Pink),
				},
				KeywordType: ansi.StylePrimitive{
					Color: sp(t.palette.Yellow),
				},
				Operator: ansi.StylePrimitive{
					Color: sp(t.palette.Sky),
				},
				Punctuation: ansi.StylePrimitive{
					Color: sp(t.palette.Overlay2),
				},
				Name: ansi.StylePrimitive{
					Color: sp(t.palette.Text),
				},
				NameBuiltin: ansi.StylePrimitive{
					Color: sp(t.palette.Red),
				},
				NameTag: ansi.StylePrimitive{
					Color: sp(t.palette.Pink),
				},
				NameAttribute: ansi.StylePrimitive{
					Color: sp(t.palette.Yellow),
				},
				NameClass: ansi.StylePrimitive{
					Color:     sp(t.palette.Yellow),
					Underline: bp(true),
					Bold:      bp(true),
				},
				NameConstant: ansi.StylePrimitive{
					Color: sp(t.palette.Peach),
				},
				NameDecorator: ansi.StylePrimitive{
					Color: sp(t.palette.Pink),
				},
				NameException: ansi.StylePrimitive{
					Color: sp(t.palette.Maroon),
				},
				NameFunction: ansi.StylePrimitive{
					Color: sp(t.palette.Blue),
				},
				NameOther: ansi.StylePrimitive{
					Color: sp(t.palette.Text),
				},
				Literal: ansi.StylePrimitive{
					Color: sp(t.palette.Text),
				},
				LiteralNumber: ansi.StylePrimitive{
					Color: sp(t.palette.Peach),
				},
				LiteralDate: ansi.StylePrimitive{
					Color: sp(t.palette.Peach),
				},
				LiteralString: ansi.StylePrimitive{
					Color: sp(t.palette.Green),
				},
				LiteralStringEscape: ansi.StylePrimitive{
					Color: sp(t.palette.Pink),
				},
				GenericDeleted: ansi.StylePrimitive{
					Color: sp(t.palette.Red),
				},
				GenericEmph: ansi.StylePrimitive{
					Italic: bp(true),
				},
				GenericInserted: ansi.StylePrimitive{
					Color: sp(t.palette.Green),
				},
				GenericStrong: ansi.StylePrimitive{
					Bold: bp(true),
				},
				GenericSubheading: ansi.StylePrimitive{
					Color: sp(t.palette.Overlay1),
				},
				Background: ansi.StylePrimitive{
					BackgroundColor: sp(t.palette.Surface0),
				},
			},
		},
		Table: ansi.StyleTable{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Color: sp(t.palette.Text),
				},
			},
			CenterSeparator: sp("┼"),
			ColumnSeparator: sp("│"),
			RowSeparator:    sp("─"),
		},
		DefinitionList: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color: sp(t.palette.Text),
			},
		},
		DefinitionTerm: ansi.StylePrimitive{
			Color: sp(t.palette.Lavender),
			Bold:  bp(true),
		},
		DefinitionDescription: ansi.StylePrimitive{
			BlockPrefix: "\n→ ",
			Color:       sp(t.palette.Text),
		},
		HTMLBlock: ansi.StyleBlock{},
		HTMLSpan:  ansi.StyleBlock{},
	}
}

// GlamourStyle returns a glamour StyleConfig for the current theme
func GlamourStyle() ansi.StyleConfig {
	return Current.GlamourStyle()
}

// =============================================================================
// COLOR ACCESSORS - Returns lipgloss.Color for easy use with styles
// =============================================================================

// Base colors (backgrounds)
func (t Theme) Base() lipgloss.Color     { return lipgloss.Color(t.palette.Base) }
func (t Theme) Mantle() lipgloss.Color   { return lipgloss.Color(t.palette.Mantle) }
func (t Theme) Crust() lipgloss.Color    { return lipgloss.Color(t.palette.Crust) }
func (t Theme) Surface0() lipgloss.Color { return lipgloss.Color(t.palette.Surface0) }
func (t Theme) Surface1() lipgloss.Color { return lipgloss.Color(t.palette.Surface1) }
func (t Theme) Surface2() lipgloss.Color { return lipgloss.Color(t.palette.Surface2) }

// Overlay colors (for floating elements)
func (t Theme) Overlay0() lipgloss.Color { return lipgloss.Color(t.palette.Overlay0) }
func (t Theme) Overlay1() lipgloss.Color { return lipgloss.Color(t.palette.Overlay1) }
func (t Theme) Overlay2() lipgloss.Color { return lipgloss.Color(t.palette.Overlay2) }

// Text colors
func (t Theme) Text() lipgloss.Color     { return lipgloss.Color(t.palette.Text) }
func (t Theme) Subtext0() lipgloss.Color { return lipgloss.Color(t.palette.Subtext0) }
func (t Theme) Subtext1() lipgloss.Color { return lipgloss.Color(t.palette.Subtext1) }

// Accent colors
func (t Theme) Rosewater() lipgloss.Color { return lipgloss.Color(t.palette.Rosewater) }
func (t Theme) Flamingo() lipgloss.Color  { return lipgloss.Color(t.palette.Flamingo) }
func (t Theme) Pink() lipgloss.Color      { return lipgloss.Color(t.palette.Pink) }
func (t Theme) Mauve() lipgloss.Color     { return lipgloss.Color(t.palette.Mauve) }
func (t Theme) Red() lipgloss.Color       { return lipgloss.Color(t.palette.Red) }
func (t Theme) Maroon() lipgloss.Color    { return lipgloss.Color(t.palette.Maroon) }
func (t Theme) Peach() lipgloss.Color     { return lipgloss.Color(t.palette.Peach) }
func (t Theme) Yellow() lipgloss.Color    { return lipgloss.Color(t.palette.Yellow) }
func (t Theme) Green() lipgloss.Color     { return lipgloss.Color(t.palette.Green) }
func (t Theme) Teal() lipgloss.Color      { return lipgloss.Color(t.palette.Teal) }
func (t Theme) Sky() lipgloss.Color       { return lipgloss.Color(t.palette.Sky) }
func (t Theme) Sapphire() lipgloss.Color  { return lipgloss.Color(t.palette.Sapphire) }
func (t Theme) Blue() lipgloss.Color      { return lipgloss.Color(t.palette.Blue) }
func (t Theme) Lavender() lipgloss.Color  { return lipgloss.Color(t.palette.Lavender) }

// =============================================================================
// PRE-BUILT STYLES - Common UI element styles
// =============================================================================

// Styles holds pre-built lipgloss styles for common UI elements
type Styles struct {
	// App chrome
	App lipgloss.Style

	// Typography
	Title    lipgloss.Style
	Subtitle lipgloss.Style
	Text     lipgloss.Style
	Muted    lipgloss.Style
	Bold     lipgloss.Style

	// Interactive elements
	Input       lipgloss.Style
	InputBorder lipgloss.Style
	Cursor      lipgloss.Style
	Placeholder lipgloss.Style

	// Containers
	Modal      lipgloss.Style
	ModalTitle lipgloss.Style
	Card       lipgloss.Style

	// Status
	Success lipgloss.Style
	Warning lipgloss.Style
	Error   lipgloss.Style
	Info    lipgloss.Style

	// Hints/help
	Hint lipgloss.Style
}

// GetStyles returns pre-built styles for the given theme
func (t Theme) GetStyles() Styles {
	return Styles{
		// App chrome - full screen background
		App: lipgloss.NewStyle().
			Background(t.Base()),

		// Typography
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Mauve()),

		Subtitle: lipgloss.NewStyle().
			Foreground(t.Subtext1()),

		Text: lipgloss.NewStyle().
			Foreground(t.Text()),

		Muted: lipgloss.NewStyle().
			Foreground(t.Overlay1()),

		Bold: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Text()),

		// Interactive elements
		Input: lipgloss.NewStyle().
			Foreground(t.Text()),

		InputBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Surface2()).
			Padding(0, 1),

		Cursor: lipgloss.NewStyle().
			Foreground(t.Rosewater()),

		Placeholder: lipgloss.NewStyle().
			Foreground(t.Overlay0()),

		// Containers
		Modal: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Mauve()).
			Padding(1, 2).
			Background(t.Surface0()),

		ModalTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Pink()).
			MarginBottom(1),

		Card: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Surface1()).
			Padding(1, 2).
			Background(t.Mantle()),

		// Status colors
		Success: lipgloss.NewStyle().
			Foreground(t.Green()),

		Warning: lipgloss.NewStyle().
			Foreground(t.Yellow()),

		Error: lipgloss.NewStyle().
			Foreground(t.Red()),

		Info: lipgloss.NewStyle().
			Foreground(t.Blue()),

		// Hints
		Hint: lipgloss.NewStyle().
			Foreground(t.Overlay1()).
			Italic(true),
	}
}

// =============================================================================
// CONVENIENCE FUNCTIONS - Use current theme
// =============================================================================

// Base returns Base color from current theme
func Base() lipgloss.Color { return Current.Base() }

// Mantle returns Mantle color from current theme
func Mantle() lipgloss.Color { return Current.Mantle() }

// Crust returns Crust color from current theme
func Crust() lipgloss.Color { return Current.Crust() }

// Surface0 returns Surface0 color from current theme
func Surface0() lipgloss.Color { return Current.Surface0() }

// Surface1 returns Surface1 color from current theme
func Surface1() lipgloss.Color { return Current.Surface1() }

// Surface2 returns Surface2 color from current theme
func Surface2() lipgloss.Color { return Current.Surface2() }

// Overlay0 returns Overlay0 color from current theme
func Overlay0() lipgloss.Color { return Current.Overlay0() }

// Overlay1 returns Overlay1 color from current theme
func Overlay1() lipgloss.Color { return Current.Overlay1() }

// Overlay2 returns Overlay2 color from current theme
func Overlay2() lipgloss.Color { return Current.Overlay2() }

// Text returns Text color from current theme
func Text() lipgloss.Color { return Current.Text() }

// Subtext0 returns Subtext0 color from current theme
func Subtext0() lipgloss.Color { return Current.Subtext0() }

// Subtext1 returns Subtext1 color from current theme
func Subtext1() lipgloss.Color { return Current.Subtext1() }

// Rosewater returns Rosewater color from current theme
func Rosewater() lipgloss.Color { return Current.Rosewater() }

// Flamingo returns Flamingo color from current theme
func Flamingo() lipgloss.Color { return Current.Flamingo() }

// Pink returns Pink color from current theme
func Pink() lipgloss.Color { return Current.Pink() }

// Mauve returns Mauve color from current theme
func Mauve() lipgloss.Color { return Current.Mauve() }

// Red returns Red color from current theme
func Red() lipgloss.Color { return Current.Red() }

// Maroon returns Maroon color from current theme
func Maroon() lipgloss.Color { return Current.Maroon() }

// Peach returns Peach color from current theme
func Peach() lipgloss.Color { return Current.Peach() }

// Yellow returns Yellow color from current theme
func Yellow() lipgloss.Color { return Current.Yellow() }

// Green returns Green color from current theme
func Green() lipgloss.Color { return Current.Green() }

// Teal returns Teal color from current theme
func Teal() lipgloss.Color { return Current.Teal() }

// Sky returns Sky color from current theme
func Sky() lipgloss.Color { return Current.Sky() }

// Sapphire returns Sapphire color from current theme
func Sapphire() lipgloss.Color { return Current.Sapphire() }

// Blue returns Blue color from current theme
func Blue() lipgloss.Color { return Current.Blue() }

// Lavender returns Lavender color from current theme
func Lavender() lipgloss.Color { return Current.Lavender() }

// GetStyles returns styles for the current theme
func GetStyles() Styles {
	return Current.GetStyles()
}
