package upload

import (
	"os"
	"path"
	"log"
	"fmt"
	//"strings"
	"mime/multipart"
	"ggin/pkg/setting"
	"strings"
	"ggin/pkg/util"
	"ggin/pkg/file"
	"ggin/pkg/logging"
)

// 上传所需检查：大小，格式，文件夹，权限！！！


func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

func GetImageFullUrl(name string) string {
	return setting.AppSetting.ImagePrefixUrl + "/" + GetImagePath() + name
}

func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)
	return fileName + ext
}

func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

func CheckImageExt(fileName string) bool {		// 图片格式
	ext := file.GetExt(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}
	return false
}

func CheckImageSize(f multipart.File) bool {		// 图片大小
	size, err := file.GetSize(f)
	if err != nil {
		log.Println(err)
		logging.Warn(err)
		return false
	}
	return size <= setting.AppSetting.ImageMaxSize
}

func CheckImage(src string) error {		// 检查上传图片所需（权限、文件夹）
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err : %v", err)
	}

	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := file.CheckPermission(src)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}
