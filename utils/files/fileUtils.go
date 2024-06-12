package filesUtils

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// 打包特定文件
func CreateTarArchiveFiles(sourceFiles []string, tarFilePath string) error {
	tarFile, err := os.Create(tarFilePath)
	if err != nil {
		return fmt.Errorf("can't create tar file %w", err)
	}
	defer tarFile.Close()

	tw := tar.NewWriter(tarFile)
	defer tw.Close()

	for _, sourceFile := range sourceFiles {
		file, err := os.Open(sourceFile)
		if err != nil {
			return fmt.Errorf("can't open  %s: %w", sourceFile, err)
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			return fmt.Errorf("can't get information %s: %w", sourceFile, err)
		}
		header := &tar.Header{
			Name: filepath.Base(sourceFile),
			Mode: int64(stat.Mode().Perm()),
			Size: stat.Size(),
		}

		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("can't write to tar header %s: %w", sourceFile, err)
		}

		_, err = io.Copy(tw, file)
		if err != nil {
			return fmt.Errorf("can't write file content to tar file %s: %w", sourceFile, err)
		}
	}

	return nil
}
