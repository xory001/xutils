package xzip

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

func ZipFileToTarGz(srcFile string, destFile string) error {
	return ZipFileToTarGzWithLevel(srcFile, destFile, gzip.DefaultCompression)
}

func ZipFileToTarGzWithLevel(srcFile string, destFile string, zipLevel int) error {
	file, err := os.OpenFile(destFile, os.O_CREATE|os.O_WRONLY, 0644)
	if nil != err {
		return err
	}
	defer file.Close()
	gzw, err := gzip.NewWriterLevel(file, zipLevel)
	if nil != err {
		return err
	}
	defer gzw.Close()
	tw := tar.NewWriter(gzw)
	defer tw.Close()

	fileInfo, err := os.Stat(srcFile)
	if nil != err {
		return err
	}
	if fileInfo.IsDir() {
		sliceFile, err := os.ReadDir(srcFile)
		if nil != err {
			return err
		}
		for _, file := range sliceFile {
			fileInfo, err := os.Stat(filepath.Join(srcFile, file.Name()))
			if nil != err {
				return err
			}
			if fileInfo.IsDir() {
				err = zipDir(filepath.Join(srcFile, fileInfo.Name()), fileInfo.Name(), fileInfo, tw)
				if nil != err {
					return err
				}
			} else {
				err = zipFile(filepath.Join(srcFile, file.Name()), fileInfo.Name(), fileInfo, tw)
				if nil != err {
					return err
				}
			}
		}
	} else {
		err = zipFile(srcFile, fileInfo.Name(), fileInfo, tw)
		if nil != err {
			return err
		}
	}

	return nil
}

func ZipFilesToTarGz(sliceSrcFile []string, destFile string) error {
	return ZipFilesToTarGzWithLevel(sliceSrcFile, destFile, gzip.DefaultCompression)
}

// ZipFilesToTarGzWithLevel zipLevel: -1 - 9, 0: no zip; 9: best zip, but slowest; -1: default level
func ZipFilesToTarGzWithLevel(sliceSrcDir []string, destFile string, zipLevel int) error {
	file, err := os.OpenFile(destFile, os.O_CREATE|os.O_WRONLY, 0644)
	if nil != err {
		return err
	}
	defer file.Close()
	gzw, err := gzip.NewWriterLevel(file, zipLevel)
	if nil != err {
		return err
	}
	defer gzw.Close()
	tw := tar.NewWriter(gzw)
	defer tw.Close()
	for _, srcFile := range sliceSrcDir {
		fileInfo, err := os.Stat(srcFile)
		if nil != err {
			return err
		}
		if fileInfo.IsDir() {
			err = zipDir(srcFile, fileInfo.Name(), fileInfo, tw)
			if nil != err {
				return err
			}
		} else {
			err = zipFile(srcFile, fileInfo.Name(), fileInfo, tw)
			if nil != err {
				return err
			}
		}
	}

	return nil
}

func zipDir(srcDir string, destRelativeDir string, srcFileInfo os.FileInfo, tw *tar.Writer) error {
	//add dir info to tar.gz
	hdr, err := tar.FileInfoHeader(srcFileInfo, "syslink")
	hdr.Name = destRelativeDir
	hdr.Format = tar.FormatGNU
	if nil != err {
		return err
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}

	sliceFile, err := os.ReadDir(srcDir)
	if nil != err {
		return err
	}
	for _, file := range sliceFile {
		fileInfo, err := os.Stat(filepath.Join(srcDir, file.Name()))
		if nil != err {
			return err
		}
		if fileInfo.IsDir() {
			err = zipDir(filepath.Join(srcDir, fileInfo.Name()), filepath.Join(destRelativeDir, fileInfo.Name()), fileInfo, tw)
			if nil != err {
				return err
			}
		} else {
			err = zipFile(filepath.Join(srcDir, file.Name()), filepath.Join(destRelativeDir, fileInfo.Name()), fileInfo, tw)
			if nil != err {
				return err
			}
		}
	}
	return nil
}

func zipFile(srcFile string, destRelativeDir string, srcFileInfo os.FileInfo, tw *tar.Writer) error {
	hdr, err := tar.FileInfoHeader(srcFileInfo, "syslink")
	hdr.Name = destRelativeDir
	hdr.Format = tar.FormatGNU
	if nil != err {
		return err
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	fileSrc, err := os.Open(srcFile)
	if nil != err {
		return err
	}
	_, err = io.Copy(tw, fileSrc)
	return err
}

func UnZipTarGzFileToDir(srcFile string, destDir string) error {
	fw, err := os.Open(srcFile)
	if nil != err {
		return err
	}
	defer fw.Close()

	gw, err := gzip.NewReader(fw)
	if nil != err {
		return err
	}
	defer gw.Close()

	tw := tar.NewReader(gw)

	_, err = os.Stat(destDir)
	if os.IsNotExist(err) {
		os.MkdirAll(destDir, 0644)
	}

	for {
		hdr, err := tw.Next()
		if nil == err {
			if hdr.FileInfo().IsDir() {
				os.MkdirAll(filepath.Join(destDir, hdr.Name), 0466)
			} else {
				destFile, err := os.Create(filepath.Join(destDir, hdr.Name))
				if nil != err {
					return err
				}
				_, err = io.Copy(destFile, tw)
			}
		} else {
			break
		}
	}

	return nil
}
