package logic

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func PathExists(s string, MyApp *MyApp) bool {
	path, _ := storage.Child(MyApp.App.Storage().RootURI(), s)
	exists, _ := storage.Exists(path)

	return exists
}

func IsNum(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	} else {
		return false
	}
}

func IsInSyncModeAndServerInaccessible(MyApp *MyApp) bool {
	if MyApp.App.Preferences().String("StorageMode") == "Sync" && !IsServerAccessible(fmt.Sprintf("http://%s:%s", MyApp.App.Preferences().String("IP"), MyApp.App.Preferences().String("Port"))) {
		return true
	} else {
		return false
	}
}

func ContainsComma(s string) bool {
	return strings.Contains(s, ",")
}

func IsServerAccessibleBoot(MyApp *MyApp, ctx context.Context, cancel context.CancelFunc, callback func(*MyApp, error)) {
	defer cancel()

	ip := MyApp.App.Preferences().String("IP")
	port := MyApp.App.Preferences().String("Port")

	d := &net.Dialer{}

	_, err := d.DialContext(ctx, "tcp", ip+":"+port)

	callback(MyApp, err)
}

func IsServerAccessibleSwitch(MyApp *MyApp, ctx context.Context, cancel context.CancelFunc, popup *widget.PopUp, callback func(*MyApp)) {
	defer cancel()

	ip := MyApp.App.Preferences().String("IP")
	port := MyApp.App.Preferences().String("Port")

	d := &net.Dialer{}

	_, err := d.DialContext(ctx, "tcp", ip+":"+port)

	//If errored will hide popup, if not hidden already
	if !popup.Hidden {
		popup.Hide()
	}

	if err != nil {
		dialog.ShowError(fmt.Errorf("Cannot Connect to Server"), MyApp.Win)
		return
	}

	var errStr []string

	if MyApp.App.Preferences().String("StorageMode") == "Local" {
		if !FileConflictCheck(MyApp) {
			MyApp.App.Preferences().SetString("StorageMode", "Sync")
			return
		} else {
			callback(MyApp)
			return
		}
	} else {
		errStr = append(errStr, "Cannot connect to server, check details or if server is running")
	}

	if MyApp.App.Preferences().String("StorageMode") == "Desync" {
		if !FileConflictCheck(MyApp) {
			MyApp.App.Preferences().SetString("StorageMode", "Sync")
			return
		} else {
			callback(MyApp)
			return
		}
	} else {
		errStr = append(errStr, "Cannot switch to Sync Mode, check server details or if server is running")
	}

	if len(errStr) != 0 {
		dialog.NewError(BuildError(errStr), MyApp.Win)
	}
}

func IsServerAccessibleChange(MyApp *MyApp, popup *widget.PopUp, testIP string, testPort string, ctx context.Context, cancel context.CancelFunc, callback func(*MyApp, error, *widget.PopUp, string, string, func(*MyApp)), callback2 func(*MyApp)) {
	defer cancel()

	d := &net.Dialer{}

	_, err := d.DialContext(ctx, "tcp", testIP+":"+testPort)

	callback(MyApp, err, popup, testIP, testPort, callback2)
}

func BuildError(errStr []string) error {
	if len(errStr) != 0 {
		return errors.New(strings.Join(errStr[:], "\n\n"))
	} else {
		return nil
	}
}
