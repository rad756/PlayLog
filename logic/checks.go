package logic

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"fyne.io/fyne/v2/storage"
)

func PathExists(s string, MyApp MyApp) bool {
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

func IsInSyncModeAndServerInaccessible(MyApp MyApp) bool {
	if MyApp.App.Preferences().String("StorageMode") == "Sync" && !IsServerAccessible(fmt.Sprintf("http://%s:%s", MyApp.App.Preferences().String("IP"), MyApp.App.Preferences().String("Port"))) {
		return true
	} else {
		return false
	}
}

func ContainsComma(s string) bool {
	return strings.Contains(s, ",")
}

func IsServerAccessibleBoot(MyApp MyApp, ctx context.Context, cancel context.CancelFunc, callback func(MyApp, error)) {
	defer cancel()

	ip := MyApp.App.Preferences().String("IP")
	port := MyApp.App.Preferences().String("Port")

	d := &net.Dialer{}

	_, err := d.DialContext(ctx, "tcp", ip+":"+port)

	callback(MyApp, err)
}
