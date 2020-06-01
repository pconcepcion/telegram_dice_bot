package validations

import (
	"strings"
	"testing"
)

func TestIsDiceExpression(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want bool
	}{
		// Valid ones
		{name: "IsDiceExpression 001", expr: "d2", want: true},
		{name: "IsDiceExpression 002", expr: "d100", want: true},
		{name: "IsDiceExpression 003", expr: "6d30000", want: true},
		{name: "IsDiceExpression 004", expr: "6d3", want: true},
		{name: "IsDiceExpression 005", expr: "6d6k3", want: true},
		{name: "IsDiceExpression 006", expr: "6d6kl3", want: true},
		{name: "IsDiceExpression 007", expr: "d6kl3", want: true},
		{name: "IsDiceExpression 008", expr: "d6e", want: true},
		{name: "IsDiceExpression 009", expr: "31d10es8", want: true},
		{name: "IsDiceExpression 010", expr: "31d10r2", want: true},
		// Invalidalid ones
		{name: "IsDiceExpression Error 001", expr: " d2", want: false},
		{name: "IsDiceExpression Error 002", expr: "d2 ", want: false},
		{name: "IsDiceExpression Error 003", expr: " d2 ", want: false},
		{name: "IsDiceExpression Error 004", expr: "d", want: false},
		{name: "IsDiceExpression Error 005", expr: "dk", want: false},
		{name: "IsDiceExpression Error 006", expr: "kd", want: false},
		{name: "IsDiceExpression Error 007", expr: "", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDiceExpression(tt.expr); got != tt.want {
				t.Errorf("IsDiceExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateSessionName(t *testing.T) {
	tests := []struct {
		name        string
		sessionName string
		want        string
		wantErr     error
	}{
		// Valid session names without errors
		{name: "ValidateSessionName 001", sessionName: "test session name", want: "test session name", wantErr: nil},
		{name: "ValidateSessionName 002", sessionName: "sesión de prueba", want: "sesión de prueba", wantErr: nil},
		{name: "ValidateSessionName 003", sessionName: "     sesión de prueba      ", want: "sesión de prueba", wantErr: nil},
		{name: "ValidateSessionName 004", sessionName: strings.Repeat("a", MaxSessionNameLength), want: strings.Repeat("a", MaxSessionNameLength), wantErr: nil},
		// Session names with errors
		{name: "ValidateSessionName Error 001", sessionName: "test session name too long, should give an error ", want: "test session name too long, should give an error", wantErr: ErrSessionNameTooLong},
		{name: "ValidateSessionName Error 002", sessionName: "a", want: "a", wantErr: ErrSessionNameTooShort},
		{name: "ValidateSessionName Error 003", sessionName: "      a        ", want: "a", wantErr: ErrSessionNameTooShort},
		{name: "ValidateSessionName Error 004", sessionName: "#####a######", want: "a", wantErr: ErrSessionNameTooShort},
		{name: "ValidateSessionName Error 005", sessionName: "#    a     #", want: "a", wantErr: ErrSessionNameTooShort},
		{name: "ValidateSessionName Error 006", sessionName: "", want: "", wantErr: ErrSessionNameTooShort},
		{name: "ValidateSessionName Error 007", sessionName: strings.Repeat("a", MaxSessionNameLength+1), want: strings.Repeat("a", MaxSessionNameLength+1), wantErr: ErrSessionNameTooLong},
		// len("á") = 2 so it goes beyond MaxSessionNameLength
		{name: "ValidateSessionName Error 007", sessionName: strings.Repeat("á", MaxSessionNameLength), want: strings.Repeat("á", MaxSessionNameLength), wantErr: ErrSessionNameTooLong},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateSessionName(tt.sessionName)
			if err != tt.wantErr {
				t.Errorf("ValidateSessionName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateSessionName() = #%v#, want #%v#", got, tt.want)
			}
		})
	}
}

func TestValidatePlayerName(t *testing.T) {
	tests := []struct {
		name       string
		playerName string
		want       string
		wantErr    error
	}{
		// Valid player names that doesn't return error
		{name: "ValidatePlayerName 001", playerName: "Bruenor", want: "Bruenor", wantErr: nil},
		{name: "ValidatePlayerName 002", playerName: "Bob", want: "Bob", wantErr: nil},
		{name: "ValidatePlayerName 003", playerName: "Mr. Salax", want: "Mr. Salax", wantErr: nil},
		{name: "ValidatePlayerName 004", playerName: "Bruenor Firesoul", want: "Bruenor Firesoul", wantErr: nil},
		{name: "ValidatePlayerName 005", playerName: "  Bob  ", want: "Bob", wantErr: nil},
		{name: "ValidatePlayerName 006", playerName: "Bob el guapo", want: "Bob el guapo", wantErr: nil},
		{name: "ValidatePlayerName 007", playerName: strings.Repeat("a", MaxPlayerNameLength), want: strings.Repeat("a", MaxPlayerNameLength), wantErr: nil},
		// Invalid player names that return error
		{name: "ValidatePlayerName Error 001", playerName: "Bruenor123456789012345678901234567890", want: "Bruenor123456789012345678901234567890", wantErr: ErrPlayerNameTooLong},
		{name: "ValidatePlayerName Error 002", playerName: "B", want: "B", wantErr: ErrPlayerNameTooShort},
		{name: "ValidatePlayerName Error 003", playerName: "", want: "", wantErr: ErrInvalidPlayerName},
		{name: "ValidatePlayerName Error 004", playerName: strings.Repeat("a", MaxPlayerNameLength+1), want: strings.Repeat("a", MaxPlayerNameLength+1), wantErr: ErrPlayerNameTooLong},
		// len ("á") = 2 so this goes beyond the max
		{name: "ValidatePlayerName 007", playerName: strings.Repeat("á", MaxPlayerNameLength), want: strings.Repeat("á", MaxPlayerNameLength), wantErr: ErrPlayerNameTooLong},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidatePlayerName(tt.playerName)
			if err != tt.wantErr {
				t.Errorf("ValidatePlayerName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidatePlayerName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateCharacterName(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name          string
		characterName string
		want          string
		wantErr       error
	}{
		// Valid character names that doesn't return error
		{name: "ValidateCharacterName 001", characterName: "Bruenor", want: "Bruenor", wantErr: nil},
		{name: "ValidateCharacterName 002", characterName: "Bob", want: "Bob", wantErr: nil},
		{name: "ValidateCharacterName 003", characterName: "Mr. Salax", want: "Mr. Salax", wantErr: nil},
		{name: "ValidateCharacterName 004", characterName: "Bruenor Firesoul", want: "Bruenor Firesoul", wantErr: nil},
		{name: "ValidateCharacterName 005", characterName: "  Bob  ", want: "Bob", wantErr: nil},
		{name: "ValidateCharacterName 006", characterName: "Bob el guapo", want: "Bob el guapo", wantErr: nil},
		{name: "ValidateCharacterName 007", characterName: strings.Repeat("a", MaxCharacterNameLength), want: strings.Repeat("a", MaxCharacterNameLength), wantErr: nil},
		// Invalid character names that return error
		{name: "ValidateCharacterName Error 001", characterName: "Bruenor123456789012345678901234567890", want: "Bruenor123456789012345678901234567890", wantErr: ErrCharacterNameTooLong},
		{name: "ValidateCharacterName Error 002", characterName: "Bruenor 1234567890 1234567890 1234567890", want: "Bruenor 1234567890 1234567890 1234567890", wantErr: ErrCharacterNameTooLong},
		{name: "ValidateCharacterName Error 003", characterName: strings.Repeat("a", MaxCharacterNameLength+1), want: strings.Repeat("a", MaxCharacterNameLength+1), wantErr: ErrCharacterNameTooLong},
		{name: "ValidateCharacterName Error 004", characterName: "B", want: "B", wantErr: ErrCharacterNameTooShort},
		{name: "ValidateCharacterName Error 005", characterName: "", want: "", wantErr: ErrInvalidCharacterName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateCharacterName(tt.characterName)
			if err != tt.wantErr {
				t.Errorf("ValidateCharacterName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateCharacterName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateColor(t *testing.T) {
	tests := []struct {
		name    string
		color   string
		want    string
		wantErr error
	}{
		// Valid colors, no errors
		{name: "ValidateColor 001", color: "#ffccaa", want: "#ffccaa", wantErr: nil},
		{name: "ValidateColor 002", color: "#fca", want: "#fca", wantErr: nil},
		{name: "ValidateColor 003", color: "  #fca  ", want: "#fca", wantErr: nil},
		{name: "ValidateColor 004", color: "#002233", want: "#002233", wantErr: nil},
		{name: "ValidateColor 005", color: "#ff0000", want: "#ff0000", wantErr: nil},
		{name: "ValidateColor 006", color: "#00ff00", want: "#00ff00", wantErr: nil},
		{name: "ValidateColor 007", color: "  #000  ", want: "#000", wantErr: nil},
		{name: "ValidateColor 008", color: "  #fff  ", want: "#fff", wantErr: nil},
		{name: "ValidateColor 009", color: "#000000", want: "#000000", wantErr: nil},
		{name: "ValidateColor 010", color: "#ffffff", want: "#ffffff", wantErr: nil},
		{name: "ValidateColor 011", color: "   #ffffff   ", want: "#ffffff", wantErr: nil},
		// Invalid Colors, give error
		{name: "ValidateColor Error 001", color: "#fcj", want: "#fcj", wantErr: ErrInvalidColor},
		{name: "ValidateColor Error 002", color: "#ffff", want: "#ffff", wantErr: ErrInvalidColor},
		{name: "ValidateColor Error 003", color: "cj", want: "cj", wantErr: ErrInvalidColor},
		{name: "ValidateColor Error 004", color: "", want: "", wantErr: ErrInvalidColor},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateColor(tt.color)
			if err != tt.wantErr {
				t.Errorf("ValidateColor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateDescription(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name    string
		desc    string
		want    string
		wantErr error
	}{
		// Valid descriptions without errors
		{name: "ValidateDescription 001", desc: "test description", want: "test description", wantErr: nil},
		{name: "ValidateDescription 002", desc: "descripción de prueba", want: "descripción de prueba", wantErr: nil},
		{name: "ValidateDescription 003", desc: "     descripción de prueba      ", want: "descripción de prueba", wantErr: nil},
		{name: "ValidateDescription 004", desc: "a", want: "a", wantErr: nil},
		{name: "ValidateDescription 005", desc: "      a        ", want: "a", wantErr: nil},
		{name: "ValidateDescription 006", desc: "#####a######", want: "a", wantErr: nil},
		{name: "ValidateDescription 007", desc: "#    a     #", want: "a", wantErr: nil},
		{name: "ValidateDescription 008", desc: "", want: "", wantErr: nil},
		{name: "ValidateDescription 009", desc: "test long description, should not give an error even if it is more than 64 characters",
			want: "test long description, should not give an error even if it is more than 64 characters", wantErr: nil},
		// it's --> its
		{name: "ValidateDescription 010", desc: "test long description, should not give an error even if it's more than 64 characters",
			want: "test long description, should not give an error even if its more than 64 characters", wantErr: nil},
		{name: "ValidateDescription 011", desc: strings.Repeat("a", MaxDescriptionLength), want: strings.Repeat("a", MaxDescriptionLength), wantErr: nil},
		// Session names with errors
		{name: "ValidateDescription Error 001", desc: strings.Repeat("test 1", 100), want: strings.Repeat("test 1", 100), wantErr: ErrDescriptionTooLong},
		{name: "ValidateDescription Error 002", desc: strings.Repeat("a", MaxDescriptionLength+1), want: strings.Repeat("a", MaxDescriptionLength+1), wantErr: ErrDescriptionTooLong},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateDescription(tt.desc)
			if err != tt.wantErr {
				t.Errorf("ValidateDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckInvalidDiceExpression(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		want    string
		wantErr error
	}{
		// Valid ones
		{name: "IsDiceExpression 001", expr: "d2", want: "d2", wantErr: nil},
		{name: "IsDiceExpression 002", expr: "d100", want: "d100", wantErr: nil},
		{name: "IsDiceExpression 003", expr: "6d30000", want: "6d30000", wantErr: nil},
		{name: "IsDiceExpression 004", expr: "6d3", want: "6d3", wantErr: nil},
		{name: "IsDiceExpression 005", expr: "6d6k3", want: "6d6k3", wantErr: nil},
		{name: "IsDiceExpression 006", expr: "6d6kl3", want: "6d6kl3", wantErr: nil},
		{name: "IsDiceExpression 007", expr: "d6kl3", want: "d6kl3", wantErr: nil},
		{name: "IsDiceExpression 008", expr: "d6e", want: "d6e", wantErr: nil},
		{name: "IsDiceExpression 009", expr: "31d10es8", want: "31d10es8", wantErr: nil},
		{name: "IsDiceExpression 010", expr: "31d10r2", want: "31d10r2", wantErr: nil},
		{name: "IsDiceExpression 011", expr: " d2", want: "d2", wantErr: nil},
		{name: "IsDiceExpression 012", expr: "d2 ", want: "d2", wantErr: nil},
		{name: "IsDiceExpression 013", expr: " d2 ", want: "d2", wantErr: nil},
		// Invalidalid ones
		{name: "IsDiceExpression Error 004", expr: "d", want: "d", wantErr: ErrInvalidDiceExpression},
		{name: "IsDiceExpression Error 005", expr: "dk", want: "dk", wantErr: ErrInvalidDiceExpression},
		{name: "IsDiceExpression Error 006", expr: "kd", want: "kd", wantErr: ErrInvalidDiceExpression},
		{name: "IsDiceExpression Error 007", expr: "", want: "", wantErr: ErrInvalidDiceExpression},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckInvalidDiceExpression(tt.expr)
			if err != tt.wantErr {
				t.Errorf("CheckInvalidDiceExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckInvalidDiceExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCleanText(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CleanText(tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("CleanText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CleanText() = %v, want %v", got, tt.want)
			}
		})
	}
}
