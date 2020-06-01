package validations

import (
	"fmt"
	"regexp"

	valid "github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
)

var (
	// ErrInvalidName it's used when a name contains invalid characters
	ErrInvalidName = errors.New("Invalid name Error")
	// ErrSessionNameTooLong it's used when a session name is longer tha MaxSessionNameLength
	ErrSessionNameTooLong = fmt.Errorf("Session Name too long. Maximum length %d", MaxSessionNameLength)
	// ErrSessionNameTooShort it's used when a session name is shorter tha MinSessionNameLength
	ErrSessionNameTooShort = fmt.Errorf("Session Name too short, minimum session name is %d", MinSessionNameLength)
	// ErrInvalidCharacterName it's trhown when a character name contains invalid characters
	ErrInvalidCharacterName = errors.New("Invalid Character Name")
	// ErrCharacterNameTooLong it's used when a character name is longer tha MaxCharacter NameLength
	ErrCharacterNameTooLong = fmt.Errorf("Character Name too long. Maximum length %d", MaxCharacterNameLength)
	// ErrCharacterNameTooShort it's used when a character name is shorter tha MinCharacterNameLength
	ErrCharacterNameTooShort = fmt.Errorf("Character Name too short, minimum character name is %d", MinCharacterNameLength)
	// ErrInvalidPlayerName it's trhown when a character name contains invalid characters
	ErrInvalidPlayerName = errors.New("Invalid Player Name")
	// ErrPlayerNameTooLong it's used when a session name is longer tha MaxPlayerNameLength
	ErrPlayerNameTooLong = fmt.Errorf("Player Name too long. Maximum length %d", MaxPlayerNameLength)
	// ErrPlayerNameTooShort it's used when a session name is shorter tha MinPlayerNameLength
	ErrPlayerNameTooShort = fmt.Errorf("Player Name too short, minimum player name is %d", MinPlayerNameLength)
	// ErrInvalidColor it's used when the color it's not valid hex color
	ErrInvalidColor = errors.New("Invalid Hex Color")
	// ErrInvalidDiceExpression it's used when an expression it's not a valid dice expression.
	ErrInvalidDiceExpression = errors.New("Invalid Dice Expression")
	// DiceExpressionRegexp regexp to match a valid dice expression
	DiceExpressionRegexp = regexp.MustCompile(DiceExpressionPattern)
	// ErrDiceExpressionTooLong it's used when a dice expresion length is longer tha MaxDiceExpressionLength
	ErrDiceExpressionTooLong = fmt.Errorf("Dice expression too long. Maximum length %d", MaxDiceExpressionLength)
	// InvalidTextRegexp regexp to match invalid texts
	InvalidTextRegexp = regexp.MustCompile(TextInvalidPattern)
	// ValidNameRegexp regexp to match valid names
	ValidNameRegexp = regexp.MustCompile(NameValidPattern)
	// ErrTextMessageTooLong it's used when a message length is longer tha MaxTextMessageLength
	ErrTextMessageTooLong = fmt.Errorf("Description too long. Maximum length %d", MaxTextMessageLength)
	// ErrDescriptionTooLong it's used when a description length is longer tha MaxDescriptionLength
	ErrDescriptionTooLong = fmt.Errorf("Description too long. Maximum length %d", MaxDescriptionLength)
)

const (
	// MinSessionNameLength minimum length for the session name
	MinSessionNameLength = 3
	// MaxSessionNameLength maximum length for the session name
	MaxSessionNameLength = 32
	// MinCharacterNameLength minimum length for the character name
	MinCharacterNameLength = 2
	// MaxCharacterNameLength maximum length for the character name
	MaxCharacterNameLength = 32
	// MinPlayerNameLength minimum length for the player name
	MinPlayerNameLength = 2
	// MaxPlayerNameLength maximum length for the player name
	MaxPlayerNameLength = 32
	// MaxDiceExpressionLength maximum length for a dice expression
	MaxDiceExpressionLength = 64
	// DiceExpressionPattern is
	DiceExpressionPattern = "^\\d*d\\d+(kl?\\d+|es\\d+|r\\d+|k\\d+|e|s\\d+)?$"
	// TextInvalidPattern pattern to valiate unaccetable charcters on text
	TextInvalidPattern = "[^\\d\\p{L} \\-\\.,:_\\*]"
	// NameValidPattern pattern to valiate accetable charcters on a name (character, player...)
	NameValidPattern = "[\\d\\p{L} \\-\\._]"
	// MaxTextMessageLength maximum length for a text in message
	MaxTextMessageLength = 256
	// MaxDescriptionLength maximum length for a description
	MaxDescriptionLength = 128
)

func init() {
	valid.TagMap["diceexpression"] = IsDiceExpression
}

// IsDiceExpression it's a govalidator style validtator that returns false when the received string
// it's not a valid dice expression, a true result doesn't ensure that the expression is valid, but this
// function should help to early find most of the invalid dice expressions
func IsDiceExpression(expr string) bool {
	return DiceExpressionRegexp.MatchString(expr)
}

// ValidateSessionName ensures a sesion name is valid, and removes non valid runes
func ValidateSessionName(name string) (string, error) {
	cleanSessionName := InvalidTextRegexp.ReplaceAllString(name, "")
	trimedCleanSessionName := valid.Trim(cleanSessionName, "") // Remove starting an tailing whitespace
	if len(trimedCleanSessionName) > MaxSessionNameLength {
		return trimedCleanSessionName, ErrSessionNameTooLong
	}
	if len(trimedCleanSessionName) < MinSessionNameLength {

		return trimedCleanSessionName, ErrSessionNameTooShort
	}
	return trimedCleanSessionName, nil
}

// ValidatePlayerName ensures a player name is valid, and removes the whitespace arround
func ValidatePlayerName(name string) (string, error) {
	trimedName := valid.Trim(name, "") // Remove starting an tailing whitespace
	if !ValidNameRegexp.MatchString(trimedName) {
		return trimedName, ErrInvalidPlayerName
	}
	if len(trimedName) > MaxPlayerNameLength {
		return trimedName, ErrPlayerNameTooLong
	}
	if len(trimedName) < MinPlayerNameLength {
		return trimedName, ErrPlayerNameTooShort
	}
	return trimedName, nil
}

// ValidateCharacterName ensures a character name is valid, and removes the whitespace arround
func ValidateCharacterName(name string) (string, error) {
	trimedName := valid.Trim(name, "") // Remove starting an tailing whitespace
	if !ValidNameRegexp.MatchString(trimedName) {
		return trimedName, ErrInvalidCharacterName
	}
	if len(trimedName) > MaxCharacterNameLength {
		return trimedName, ErrCharacterNameTooLong
	}
	if len(trimedName) < MinCharacterNameLength {
		return trimedName, ErrCharacterNameTooShort
	}
	return trimedName, nil
}

// ValidateColor validates if the given color it's a valid Hex color, and removes the whitespace arround
func ValidateColor(color string) (string, error) {
	trimedColor := valid.Trim(color, "")
	if !valid.IsHexcolor(trimedColor) {
		return trimedColor, ErrInvalidColor
	}
	return trimedColor, nil
}

// ValidateDescription validates the description it's a valid text
func ValidateDescription(desc string) (string, error) {
	validDescription := InvalidTextRegexp.ReplaceAllString(desc, "")
	trimedValidDescription := valid.Trim(validDescription, "")
	if len(trimedValidDescription) > MaxDescriptionLength {
		return trimedValidDescription, ErrDescriptionTooLong
	}
	return trimedValidDescription, nil
}

// CheckInvalidDiceExpression checks does only some basic checking on the dice expression to detect early cases of
// invalid dice expressions, this check doesn't ensure the expression is valid, but allows early detection of
// many invalid cases
func CheckInvalidDiceExpression(expr string) (string, error) {
	trimedExpr := valid.Trim(expr, "") // Remove starting an tailing whitespace
	if !IsDiceExpression(trimedExpr) {
		return trimedExpr, ErrInvalidDiceExpression
	}
	if len(trimedExpr) > MaxDiceExpressionLength {
		return trimedExpr, ErrDiceExpressionTooLong
	}
	return trimedExpr, nil
}

// CleanText removes invalid text charactes
func CleanText(text string) (string, error) {
	validText := InvalidTextRegexp.ReplaceAllString(text, "")
	if len(validText) > MaxTextMessageLength {
		return validText, ErrTextMessageTooLong
	}
	return validText, nil
}
