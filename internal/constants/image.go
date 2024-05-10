package constants

var UploadDir = "uploads"

const (
	MaxFileSize    = 10 << 20 // 10 MB limit
	AllowedFileExt = ".jpg|.jpeg|.png|.gif"
)

// use magic bytes

// IsAllowedFileExtension checks if the file extension is allowed
func IsAllowedFileExtension(ext string) bool {
	allowedExts := map[string]struct{}{
		".jpg":  {},
		".jpeg": {},
		".png":  {},
		".gif":  {},
	}

	_, ok := allowedExts[ext]
	return ok
}
