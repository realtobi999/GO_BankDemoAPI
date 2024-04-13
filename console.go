package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func clearConsole() {
	// Clearing the terminal screen based on the OS
	switch runtime.GOOS {
	case "linux", "darwin":
		// For Unix-like systems
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		// For Windows
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		// Unsupported OS
		fmt.Println("Unsupported OS. Cannot clear the terminal.")
	}
}

/***
 *    ______                _      _____              _                    
 *    | ___ \              | |    /  ___|            | |                   
 *    | |_/ /  __ _  _ __  | | __ \ `--.  _   _  ___ | |_   ___  _ __ ___  
 *    | ___ \ / _` || '_ \ | |/ /  `--. \| | | |/ __|| __| / _ \| '_ ` _ \ 
 *    | |_/ /| (_| || | | ||   <  /\__/ /| |_| |\__ \| |_ |  __/| | | | | |
 *    \____/  \__,_||_| |_||_|\_\ \____/  \__, ||___/ \__| \___||_| |_| |_|
 *                                         __/ |                           
 *                                        |___/                            
 */
func printASCII() {
	fmt.Println(``)
	fmt.Println(` ______                _      _____              _                    `)
	fmt.Println(` | ___ \              | |    /  ___|            | |                   `)
	fmt.Println(` | |_/ /  __ _  _ __  | | __ \ '--.  _   _  ___ | |_   ___  _ __ ___  `)
	fmt.Println(` | ___ \ / _' || '_ \ | |/ /  '--. \| | | |/ __|| __| / _ \| '_ ' _ \ `)
	fmt.Println(` | |_/ /| (_| || | | ||   <  /\__/ /| |_| |\__ \| |_ |  __/| | | | | |`)
	fmt.Println(` \____/  \__,_||_| |_||_|\_\ \____/  \__, ||___/ \__| \___||_| |_| |_|`)
	fmt.Println(`                                      __/ |                           `)
	fmt.Println(`                                     |___/                            `)
	fmt.Println(``)
}