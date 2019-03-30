package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Jguer/yay/v9/generic"
)

// waitLock will lock yay checking the status of db.lck until it does not exist
func waitLock() {
	if _, err := os.Stat(filepath.Join(pacmanConf.DBPath, "db.lck")); err != nil {
		return
	}

	fmt.Println(generic.Bold(generic.Yellow(generic.SmallArrow)), filepath.Join(pacmanConf.DBPath, "db.lck"), "is present.")

	fmt.Print(generic.Bold(generic.Yellow(generic.SmallArrow)), " There may be another Pacman instance running. Waiting...")

	for {
		time.Sleep(3 * time.Second)
		if _, err := os.Stat(filepath.Join(pacmanConf.DBPath, "db.lck")); err != nil {
			fmt.Println()
			return
		}
	}
}

func passToPacman(args *arguments) *exec.Cmd {
	argArr := make([]string, 0)

	if args.needRoot() {
		argArr = append(argArr, "sudo")
	}

	argArr = append(argArr, config.PacmanBin)
	argArr = append(argArr, cmdArgs.formatGlobals()...)
	argArr = append(argArr, args.formatArgs()...)
	if config.NoConfirm {
		argArr = append(argArr, "--noconfirm")
	}

	argArr = append(argArr, "--config", config.PacmanConf)
	argArr = append(argArr, "--")
	argArr = append(argArr, args.targets...)

	if args.needRoot() {
		waitLock()
	}
	return exec.Command(argArr[0], argArr[1:]...)
}

func passToMakepkg(dir string, args ...string) *exec.Cmd {
	if config.NoConfirm {
		args = append(args)
	}

	mflags := strings.Fields(config.MFlags)
	args = append(args, mflags...)

	if config.MakepkgConf != "" {
		args = append(args, "--config", config.MakepkgConf)
	}

	cmd := exec.Command(config.MakepkgBin, args...)
	cmd.Dir = dir
	return cmd
}

func passToGit(dir string, _args ...string) *exec.Cmd {
	gitflags := strings.Fields(config.GitFlags)
	args := []string{"-C", dir}
	args = append(args, gitflags...)
	args = append(args, _args...)

	cmd := exec.Command(config.GitBin, args...)
	return cmd
}

func isTty() bool {
	cmd := exec.Command("test", "-t", "1")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	return err == nil
}
