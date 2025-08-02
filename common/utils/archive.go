package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
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
		return fmt.Errorf("identify file format failed:%w", err)
	}
	if ex, ok := format.(archiver.Extractor); ok {
		return ex.Extract(context.Background(), file, nil, func(ctx context.Context, f archiver.File) error {
			filePath := filepath.Join(dest, f.NameInArchive)
			if f.IsDir() {
				if err := os.MkdirAll(filePath, f.Mode()); err != nil {
					return fmt.Errorf("mkdir %s failed: %w", filePath, err)
				}
				if curuser, _ := user.Current(); curuser.Name == "root" {
					if err := os.Chown(filePath, uid, gid); err != nil {
						return fmt.Errorf("chown %d:%d %s failed", uid, gid, filePath)
					}
				}
				return nil
			}
			// create extract file's directory
			if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				return fmt.Errorf("mkdir %s failed: %w", filepath.Dir(filePath), err)
			}
			// extract file
			newFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
			if err != nil {
				return fmt.Errorf("open newfile %s failed: %w", filePath, err)
			}
			defer newFile.Close()
			rd, err := f.Open()
			if err != nil {
				return fmt.Errorf("open archived file %s failed: %w", f.Name(), err)
			}
			defer rd.Close()
			if _, err = io.Copy(newFile, rd); err != nil {
				return fmt.Errorf("copy to newfile from archived file failed: %w", err)
			}
			if curuser, _ := user.Current(); curuser.Name == "root" {
				if err := os.Chown(filePath, uid, gid); err != nil {
					return fmt.Errorf("chown %d:%d %s failed", uid, gid, filePath)
				}
			}
			return nil
		})
	}
	return nil
}
