package model

type Game struct {
	Id    GameId
	Label string
}

func (g Game) IsNull() bool {
	return (g.Id == "") && (g.Label == "")
}

type GameId = string

type InputOption struct {
	// if non-nil then this can be clicked
	ClickOptions *ClickOption

	// if non-nil then this can be toggled
	ToggleOptions *ToggleOption
}

type ClickOption struct {
	Key InputKey
}

type ToggleOption struct {
	Key InputKey
}

// input command. Only one of each type is allowed
type Input struct {
	Click  *ClickInput
	Toggle *ToggleInput
}

type ClickInput struct {
	Key InputKey

	// number of times to repeat the input, where Repeats == 1
	// indicates to click once
	Repeat int
}

type ToggleInput struct {
	Key InputKey

	// optional value indicating which toggle state the input should have, rather
	// than just being toggled
	Value *bool
}

type InputKey = string
