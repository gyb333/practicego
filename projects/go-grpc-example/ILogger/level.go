package ILogger

import (
	"errors"
	"strings"
)



//ModuleLevels maintains log levels based on module
type ModuleLevels struct {
	levels map[string]Level
}

// GetLevel returns the log level for the given module.
func (l *ModuleLevels) GetLevel(module string) Level {
	level, exists := l.levels[module]
	if !exists {
		level, exists = l.levels[""]
		// no configuration exists, default to info
		if !exists {
			level = INFO
		}
	}
	return level
}

// SetLevel sets the log level for the given module.
func (l *ModuleLevels) SetLevel(module string, level Level) {
	if l.levels == nil {
		l.levels = make(map[string]Level)
	}
	l.levels[module] = level
}

// IsEnabledFor will return true if logging is enabled for the given module.
func (l *ModuleLevels) IsEnabledFor(module string, level Level) bool {
	return level <= l.GetLevel(module)
}


//Log level names in string
var levelNames = []string{
	"CRITICAL",
	"ERROR",
	"WARNING",
	"INFO",
	"DEBUG",
}

// ParseLevel returns the log level from a string representation.
func ParseLevel(level string) (Level, error) {
	for i, name := range levelNames {
		if strings.EqualFold(name, level) {
			return Level(i), nil
		}
	}
	return ERROR, errors.New("logger: invalid log level")
}

//ParseString returns String repressentation of given log level
func ParseString(level Level) string {
	return levelNames[level]
}
