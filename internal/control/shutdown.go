package control

import (
	"context"
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

// Shutdown turns off the app gracefully
func Shutdown(systemWideCancel context.CancelFunc) {

	systemWideCancel()

	fmt.Println(aurora.Yellow("\nShutting down.... ").BgBlue())
	os.Exit(0)
}
