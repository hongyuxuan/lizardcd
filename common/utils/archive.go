package utils

import (
	"context"
	"io"
	"os"
	"os/user"
	"strings"

	"github.com/mholt/archiver/v4"
)

func Unarchive(filename, dest string, uid, gid int) error {
	if dest != "" && !strings.HasSuffix(dest, "/") {
		dest += "/"
	}
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	format, _, err := archiver.Identify(filename, file)
	if err != nil {
		return err
	}
	if ex, ok := format.(archiver.Extractor); ok {
		ex.Extract(context.Background(), file, nil, func(ctx context.Context, f archiver.File) error {
			if f.FileInfo.IsDir() {
				os.MkdirAll(dest+f.NameInArchive, f.Mode())
				if curuser, _ := user.Current(); curuser.Name == "root" {
					os.Chown(dest+f.NameInArchive, uid, gid)
				}
				return nil
			}
			newFile, err := os.OpenFile(dest+f.NameInArchive, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer newFile.Close()
			rd, err := f.Open()
			if err != nil {
				return err
			}
			defer rd.Close()
			io.Copy(newFile, rd)
			if curuser, _ := user.Current(); curuser.Name == "root" {
				os.Chown(dest+f.NameInArchive, uid, gid)
			}
			return nil
		})
	}
	return nil
}
