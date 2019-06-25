// +build ignore

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	pPath := strings.Replace(filepath.Join("C:\\", "Users", "Public", "env_windows_amd64", "5.13.0", "mingw73_64"), string(filepath.Separator), "\\", -1)
	if len(os.Args) >= 2 {
		pPath = os.Args[1]
	}
	if !strings.Contains(pPath, "5.13.0") {
		pPath = strings.Replace(filepath.Join(pPath, "5.13.0", "mingw73_64"), string(filepath.Separator), "\\", -1)
	}

	var oPath string
	for _, fn := range []string{
		"bin/Qt5Core.dll",
		"bin/qmake.exe",
	} {
		fn = filepath.Join("./5.13.0/mingw73_64/", fn)

		data, err := ioutil.ReadFile(fn)
		if err != nil {
			println("couldn't find", fn)
			continue
		}

		for _, path := range []string{"qt_prfxpath", "qt_epfxpath", "qt_hpfxpath"} {
			path += "="

			start := bytes.Index(data, []byte(path))
			if start == -1 {
				continue
			}

			end := bytes.IndexByte(data[start:], byte(0))
			if end == -1 {
				continue
			}

			if len(oPath) == 0 {
				oPath = string(data[start : start+end])
			}

			rep := append([]byte(path), []byte(pPath)...)
			if lendiff := end - len(rep); lendiff < 0 {
				end -= lendiff
			} else {
				rep = append(rep, bytes.Repeat([]byte{0}, lendiff)...)
			}
			data = bytes.Replace(data, data[start:start+end], rep, -1)
		}

		if err := ioutil.WriteFile(fn, data, 0644); err != nil {
			println("couldn't patch", fn)
		} else {
			println("patched", fn)
		}
	}

	fn := filepath.Join("./5.13.0/mingw73_64/bin/qtenv2.bat")
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		println("couldn't find", fn)
		return
	}
	oPath = strings.Split(oPath, "=")[1]
	data = bytes.Replace(data, []byte("cd /D "+oPath+"\r\n"), []byte(""), -1)
	data = bytes.Replace(data, []byte(oPath), []byte(pPath), -1)
	data = bytes.Replace(data, []byte(strings.Replace(strings.Replace(strings.Replace(oPath, "\\", "/", -1), "/5.13.0/", "/Tools/", -1), "73", "730", -1)), []byte(strings.Replace(strings.Replace(strings.Replace(strings.Replace(pPath, "\\", "/", -1), "/5.13.0/", "/Tools/", -1), "73", "730", -1), "_amd64/Tools/", "_amd64_Tools/", -1)), -1)
	if !bytes.Contains(data, []byte("-----")) {
		data = append(data, []byte("echo To export the current PATH to your default CMD or PS env run\r\necho ------------------------\r\necho setx PATH \"%%PATH%%\"\r\necho ------------------------\r\necho and re-open the command line window\r\n")...)
	}
	if err := ioutil.WriteFile(fn, data, 0644); err != nil {
		println("couldn't patch", fn)
	} else {
		println("patched", fn)
	}
}
