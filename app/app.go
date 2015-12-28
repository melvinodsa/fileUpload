package app

import (
	"fmt"

	"github.com/miketheprogrammer/go-thrust/lib/commands"
	"github.com/miketheprogrammer/go-thrust/lib/connection"
	"github.com/miketheprogrammer/go-thrust/thrust"
)

//StartFileUpload starts the GUI of the application for the front end.
func StartFileUpload(url string) {

	thrust.InitLogger()
	// thrust.Start() must always come before any bindings are created.
	thrust.Start()

	thrustWindow := thrust.NewWindow(thrust.WindowOptions{
		RootUrl: url,
	})
	thrustWindow.Show()
	thrustWindow.Maximize()
	thrustWindow.Focus()

	onclose, err := thrust.NewEventHandler("closed", func(cr commands.EventResult) {
		fmt.Println("Close Event Occured")
		connection.CleanExit()
	})
	fmt.Println(onclose)
	if err != nil {
		fmt.Println(err)
		connection.CleanExit()
	}
	// In lieu of something like an http server, we need to lock this thread
	// in order to keep it open, and keep the process running.
	// Dont worry we use runtime.Gosched :)
	thrust.LockThread()
}
