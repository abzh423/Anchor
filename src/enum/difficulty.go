package enum

import "fmt"

// Difficulty is a Minecraft difficulty level, defines either as an integer or
// a string. Use String() on a value to get the equivalent string value.
type Difficulty int

const (
	// DifficultyPeaceful is the peaceful difficulty (0)
	DifficultyPeaceful Difficulty = iota
	// DifficultyEasy is the easy difficulty (1)
	DifficultyEasy
	// DifficultyNormal is the normal difficulty (2)
	DifficultyNormal
	// DifficultyHard is the hard difficulty (3)
	DifficultyHard
)

// String converts the difficulty into an enumerable string in lowercase format.
func (d Difficulty) String() string {
	switch d {
	case DifficultyPeaceful:
		return "peaceful"
	case DifficultyEasy:
		return "easy"
	case DifficultyNormal:
		return "normal"
	case DifficultyHard:
		return "hard"
	default:
		panic(fmt.Errorf("unknown difficulty type: %d", d))
	}
}

// ParseDifficulty attempts to parse the difficulty value, whether it is a
// string, integer or Difficulty type itself. It will return an error if the
// parsed difficulty value is out of range or of unknown type.
// TODO test files for this method
func ParseDifficulty(value interface{}) (Difficulty, error) {
	switch v := value.(type) {
	case Difficulty:
		return v, nil
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		{
			if v.(Difficulty) < DifficultyPeaceful || v.(Difficulty) > DifficultyHard {
				return 0, fmt.Errorf("unknown difficulty type: %v", v)
			}

			return v.(Difficulty), nil
		}
	case string:
		switch v {
		case "peaceful":
			return DifficultyPeaceful, nil
		case "easy":
			return DifficultyEasy, nil
		case "normal":
			return DifficultyNormal, nil
		case "hard":
			return DifficultyHard, nil
		default:
			return 0, fmt.Errorf("failed to parse difficulty: %s", v)
		}
	default:
		return 0, fmt.Errorf("unknown type for difficulty: %T", value)
	}
}
