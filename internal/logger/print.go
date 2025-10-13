package logger

import (
	"fmt"
	"os"
)

func Hellof(format string, args ...any) {
	if opt.Silent {
		return
	}

	fmt.Printf("👋 "+format+"\n", args...)
}

func Debugf(format string, args ...any) {
	if opt.Silent {
		return
	}

	if !opt.Debug {
		return
	}

	fmt.Printf("⚙️ "+format+"\n", args...)
}

func Warningf(format string, args ...any) {
	if opt.Silent {
		return
	}

	message := "⚠️ WARNING: " + fmt.Sprintf(format, args...) + "\n"

	if _, err := fmt.Fprint(os.Stderr, message); err != nil {
		fmt.Print(message)
	}
}

func Errorf(format string, args ...any) {
	message := "🚫 ERROR: " + fmt.Sprintf(format, args...) + "\n"

	if _, err := fmt.Fprint(os.Stderr, message); err != nil {
		fmt.Print(message)
	}
}

func Successf(format string, args ...any) {
	if opt.Silent {
		return
	}

	fmt.Printf("👍 "+format+"\n", args...)
}
