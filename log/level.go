package log

import "log/slog"

func SetLevel(opts *slog.HandlerOptions, level_name string) error {

	// parse log level from string into proper value
	var level slog.Level
	err := level.UnmarshalText([]byte(level_name))
	if err != nil {
		return err
	}
	opts.Level = level

	return nil
}
