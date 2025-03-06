package file

import (
	"github.com/fsnotify/fsnotify"
	"log"
)

func Callback(event fsnotify.Event) {
	file := event.Name
	op := event.Op

	switch op {
	case fsnotify.Create:
		log.Println("创建文件:", file)

	case fsnotify.Write:
		log.Println("写入文件:", file)

	case fsnotify.Remove:
		log.Println("删除文件:", file)

	case fsnotify.Rename:
		log.Println("重命名文件:", file)

	default:
		log.Println("未知操作:", file, op.String())
	}
}
