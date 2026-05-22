package filesys

import "mime/multipart"

type MIMETYpe string

func DetectMIME(file multipart.File) string
