package edit

import (
	"os"
	"os/exec"
	"runtime"
)

// OpenFile opens a file in the default editor - cross platform
func OpenFile(file string) error {
	// Get user's preferred editor from the environment
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = defaultEditor()
	}

	// Execute the editor with the file path
	cmd := exec.Command(editor, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func defaultEditor() string {
	switch runtime.GOOS {
	case "windows":
		return "notepad"
	case "darwin":
		return "vim"
	default: // Linux and other OS
		return "nano"
	}
}
