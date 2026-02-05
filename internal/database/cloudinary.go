package database

import (
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

// NewCloudinary returns a configured Cloudinary instance
func NewCloudinary() (*cloudinary.Cloudinary, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("missing Cloudinary credentials in environment variables")
	}

	// Build Cloudinary URL
	cldURL := fmt.Sprintf("cloudinary://%s:%s@%s", apiKey, apiSecret, cloudName)

	cld, err := cloudinary.NewFromURL(cldURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary: %w", err)
	}

	return cld, nil
}
