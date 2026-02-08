package services

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"path/filepath"
	"strings"
	"time"

	"labkoding.my.id/kasir-api/external"
	"labkoding.my.id/kasir-api/models"
	"labkoding.my.id/kasir-api/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) GetAllProducts(name string) ([]models.Product, error) {
	return s.repo.GetAllProducts(name)
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	return s.repo.CreateProduct(product)
}

func (s *ProductService) GetProductByID(id string) (*models.Product, error) {
	return s.repo.GetProductByID(id)
}

func (s *ProductService) UpdateProduct(product *models.Product) error {
	return s.repo.UpdateProduct(product)
}

func (s *ProductService) DeleteProduct(id string) error {
	return s.repo.DeleteProduct(id)
}

// UploadProductImage uploads an image reader to configured R2 and returns the public URL.
// filename is used to preserve extension when generating the storage key.
// Images are automatically compressed before upload.
func (s *ProductService) UploadProductImage(ctx context.Context, r io.Reader, filename, contentType string) (string, error) {
	// Compress image before upload
	compressed, finalContentType, err := compressImage(r, contentType)
	if err != nil {
		return "", fmt.Errorf("failed to compress image: %w", err)
	}

	ext := filepath.Ext(filename)
	// Force .jpg extension if we converted to JPEG
	if finalContentType == "image/jpeg" && ext != ".jpg" && ext != ".jpeg" {
		ext = ".jpg"
	}
	key := fmt.Sprintf("products/%d%s", time.Now().UnixNano(), ext)
	url, err := external.UploadObject(ctx, key, compressed, finalContentType)
	if err != nil {
		return "", err
	}
	return url, nil
}

// compressImage compresses an image to reduce file size
// Supports JPEG and PNG. Returns compressed image as reader, content type, and error
func compressImage(r io.Reader, contentType string) (io.Reader, string, error) {
	// Read into buffer so we can decode config first, then decode again
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, r); err != nil {
		return nil, "", fmt.Errorf("failed to read image: %w", err)
	}

	// Check image dimensions before full decode to prevent decompression bombs
	config, format, err := image.DecodeConfig(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image config: %w", err)
	}

	// Reject images that are too large (e.g., max 10000x10000 = 100 megapixels)
	const maxWidth = 10000
	const maxHeight = 10000
	const maxPixels = 100_000_000 // 100 megapixels

	if config.Width > maxWidth || config.Height > maxHeight {
		return nil, "", fmt.Errorf("image dimensions too large: %dx%d (max: %dx%d)",
			config.Width, config.Height, maxWidth, maxHeight)
	}

	totalPixels := config.Width * config.Height
	if totalPixels > maxPixels {
		return nil, "", fmt.Errorf("image has too many pixels: %d (max: %d)",
			totalPixels, maxPixels)
	}

	// Now safe to decode the full image
	img, _, err := image.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Get image dimensions
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Resize if too large (max 1200px on longest side)
	maxDimension := 1200
	if width > maxDimension || height > maxDimension {
		if width > height {
			height = height * maxDimension / width
			width = maxDimension
		} else {
			width = width * maxDimension / height
			height = maxDimension
		}
		// Clamp to at least 1 pixel to avoid zero dimensions
		if width < 1 {
			width = 1
		}
		if height < 1 {
			height = 1
		}
		img = resize(img, width, height)
	}

	// Encode to buffer
	var encodedBuf bytes.Buffer
	var finalContentType string

	// Convert PNG to JPEG for better compression, keep JPEG as JPEG
	if strings.ToLower(format) == "png" {
		// Flatten PNG onto white background to handle transparency
		flatImg := image.NewRGBA(img.Bounds())
		// Fill with white background
		for y := flatImg.Bounds().Min.Y; y < flatImg.Bounds().Max.Y; y++ {
			for x := flatImg.Bounds().Min.X; x < flatImg.Bounds().Max.X; x++ {
				flatImg.Set(x, y, image.White)
			}
		}
		// Draw original image on top
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
				flatImg.Set(x, y, img.At(x, y))
			}
		}
		// Convert to JPEG with 85% quality
		err = jpeg.Encode(&encodedBuf, flatImg, &jpeg.Options{Quality: 85})
		finalContentType = "image/jpeg"
	} else {
		// JPEG: re-encode with 85% quality
		err = jpeg.Encode(&encodedBuf, img, &jpeg.Options{Quality: 85})
		finalContentType = "image/jpeg"
	}

	if err != nil {
		return nil, "", fmt.Errorf("failed to encode image: %w", err)
	}

	return &encodedBuf, finalContentType, nil
}

// resize implements simple nearest-neighbor image resizing
func resize(img image.Image, newWidth, newHeight int) image.Image {
	oldBounds := img.Bounds()
	oldWidth := oldBounds.Dx()
	oldHeight := oldBounds.Dy()

	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// Map new coordinates to old coordinates
			oldX := x * oldWidth / newWidth
			oldY := y * oldHeight / newHeight
			newImg.Set(x, y, img.At(oldBounds.Min.X+oldX, oldBounds.Min.Y+oldY))
		}
	}

	return newImg
}
