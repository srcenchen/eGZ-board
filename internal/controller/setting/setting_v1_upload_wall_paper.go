package setting

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"eGZ-Board/api/setting/v1"
)

var wallpaperTypes = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".webp": "image/webp",
}

func (c *ControllerV1) UploadWallPaper(ctx context.Context, req *v1.UploadWallPaperReq) (res *v1.UploadWallPaperRes, err error) {
	extension := strings.ToLower(filepath.Ext(req.File.Filename))
	expectedType, allowed := wallpaperTypes[extension]
	if !allowed {
		return nil, fmt.Errorf("wallpaper must be jpg, jpeg, png, or webp")
	}
	upload, err := req.File.Open()
	if err != nil {
		return nil, fmt.Errorf("open wallpaper upload: %w", err)
	}
	buffer := make([]byte, 512)
	n, readErr := io.ReadFull(upload, buffer)
	if readErr != nil && readErr != io.ErrUnexpectedEOF {
		upload.Close()
		return nil, fmt.Errorf("read wallpaper upload: %w", readErr)
	}
	if err := upload.Close(); err != nil {
		return nil, fmt.Errorf("close wallpaper upload: %w", err)
	}
	detectedType := http.DetectContentType(buffer[:n])
	if detectedType != expectedType {
		return nil, fmt.Errorf("wallpaper content type %s does not match extension %s", detectedType, extension)
	}

	fileName, err := req.File.Save("./resource/upload", true)
	if err != nil {
		return nil, fmt.Errorf("save wallpaper upload: %w", err)
	}
	if filepath.Dir(fileName) != "." {
		_ = os.Remove(filepath.Join("./resource/upload", fileName))
		return nil, fmt.Errorf("invalid wallpaper filename")
	}
	return &v1.UploadWallPaperRes{FileName: fileName}, nil
}
