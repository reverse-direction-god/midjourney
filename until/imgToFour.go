package until

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func splitImage_____TWO(imagePath string) ([]image.Image, error) {
	var images []image.Image
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 读取图像
	m, err := png.Decode(file)
	if err != nil {
		return nil, err
	}
	rgba := m.(*image.RGBA)
	//第一个
	subImage := rgba.SubImage(image.Rect(0, 0, m.Bounds().Max.X/2, m.Bounds().Max.Y/2)).(*image.RGBA)
	//第二个
	subImage1 := rgba.SubImage(image.Rect(m.Bounds().Max.X/2, 0, m.Bounds().Max.X, m.Bounds().Max.Y/2)).(*image.RGBA)
	//第三个
	subImage2 := rgba.SubImage(image.Rect(0, m.Bounds().Max.Y/2, m.Bounds().Max.X/2, m.Bounds().Max.Y)).(*image.RGBA)
	//第四个
	subImage3 := rgba.SubImage(image.Rect(m.Bounds().Max.X/2, m.Bounds().Max.Y/2, m.Bounds().Max.X, m.Bounds().Max.Y)).(*image.RGBA)
	images = append(images, subImage)
	images = append(images, subImage1)
	images = append(images, subImage2)
	images = append(images, subImage3)
	return images, nil
}
func splitImage(imagePath string) ([]image.Image, error) {
	// 打开图像文件
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 读取图像
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// 获取图像的宽度和高度
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	// 计算分割后每张图的宽度和高度
	newWidth := width / 2
	newHeight := height / 2

	// 分割图像
	var images []image.Image
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			// 定义每张图的矩形区域
			rect := image.Rect(j*newWidth, i*newHeight, (j+1)*newWidth, (i+1)*newHeight)

			// 裁剪图像
			subImage := img.(*image.RGBA).SubImage(rect)

			// 将裁剪后的图像添加到列表中
			images = append(images, subImage)
		}
	}

	return images, nil
}
func OssEd(by []byte, fileName string) error {
	// 创建OSS Client实例
	client, err := oss.New("oss-cn-shanghai.aliyuncs.com", "*", "*")
	if err != nil {
		log.Println(err)
		return err
	}

	// 获取Bucket实例
	bucket, err := client.Bucket("*")
	if err != nil {
		log.Println(err)
		return err
	}

	err = bucket.PutObject(fileName, bytes.NewReader(by))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func imageToBytes(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func saveImage(imagePath string, img image.Image) error {

	// 根据文件扩展名选择相应的编码器
	switch filepath.Ext(imagePath) {
	case ".jpg", ".jpeg":
		by, _ := imageToBytes(img)
		return OssEd(by, imagePath)
		// return jpeg.Encode(file, img, nil)
	case ".png":
		by, _ := imageToBytes(img)
		return OssEd(by, imagePath)
		// return png.Encode(file, img)
	default:
		return fmt.Errorf("unsupported image format")
	}
}

func ImgToFour(path string, userId int) (string, error) {
	// 图像文件路径
	userIdStr := strconv.Itoa(userId)

	fileName := strconv.Itoa(int(time.Now().Unix())) + userIdStr
	imagePath := path + ".png"
	var filesName string
	ok, err := func() (string, error) {
		// 分割图像
		croppedImages, err := splitImage_____TWO(imagePath)
		if err != nil {
			log.Println("Error:", err)
			return "", err
		}

		// 保存分割后的图像
		for i, img := range croppedImages {
			imagePath := fmt.Sprintf(fileName+"_%d.png", i+1)

			err := saveImage(imagePath, img)
			if err != nil {
				log.Println(err)
				return "", err
			}
			filesName = filesName + imagePath + "|"
		}
		return filesName, nil
	}()

	if ok == "" {
		ok, err = func() (string, error) {
			// 分割图像
			croppedImages, err := splitImage_____TWO(imagePath)
			if err != nil {
				log.Println("Error:", err)
				return "", err
			}

			// 保存分割后的图像
			for i, img := range croppedImages {
				imagePath := fmt.Sprintf(fileName+"_%d.png", i+1)

				err := saveImage(imagePath, img)
				if err != nil {
					log.Println(err)
					return "", err
				}
				filesName = filesName + imagePath + "|"
			}
			return filesName, nil
		}()
	}
	if ok == "" {
		ok, err = func() (string, error) {
			// 分割图像
			croppedImages, err := splitImage_____TWO(imagePath)
			if err != nil {
				log.Println("Error:", err)
				return "", err
			}

			// 保存分割后的图像
			for i, img := range croppedImages {
				imagePath := fmt.Sprintf(fileName+"_%d.png", i+1)

				err := saveImage(imagePath, img)
				if err != nil {
					log.Println(err)
					return "", err
				}
				filesName = filesName + imagePath + "|"
			}
			return filesName, nil
		}()
	}

	return ok, err
}
