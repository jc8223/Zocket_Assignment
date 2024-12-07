// imageanalysis/imageanalysis.go

package imageanalysis

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // Import PNG image format for image.Decode
	"io"
	"context"
	"net/http"
	"os"
	"path/filepath"
	// "github.com/gin-gonic/gin"
	 
	// "zocket-assignment/db"
	// "zocket-assignment/models"

	"github.com/nfnt/resize"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CompressAndStoreImages downloads and compresses product images, updating the database.
func CompressAndStoreImages(productID int, imageUrls []string) error {
	// Create a directory to store compressed images
	outputDir := fmt.Sprintf("compressed_images/%d", productID)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	  
	// Loop through each image URL
	for index, imageUrl := range imageUrls {
		// Download the image
		imageData, err := downloadImage(imageUrl)
		if err != nil {
			fmt.Printf("Error downloading image %d: %v\n", index+1, err)
			continue
		}

		// Compress the image
		compressedImage, err := compressImage(imageData, 100, 0) // Change compression parameters as needed
		if err != nil {
			fmt.Printf("Error compressing image %d: %v\n", index+1, err)
			continue
		}

		// Save the compressed image to a file
		outputPath := filepath.Join(outputDir, fmt.Sprintf("compressed_image%d.jpg", index+1))
		if err := saveImage(compressedImage, outputPath); err != nil {
			fmt.Printf("Error saving compressed image %d: %v\n", index+1, err)
			continue
		}

		updateErr  := updateDatabaseWithImagePath(productID, outputPath)
        if updateErr  != nil {
            fmt.Printf("Error updating database with image path: %v\n", updateErr )
            continue
        }
		// Update the database with the local file path of the compressed image
		// TODO: Replace this with your actual database update logic
		// For now, we'll print the local file path
		fmt.Printf("Compressed image path for product ID %d: %s\n", productID, outputPath)
	}

	fmt.Printf("Image analysis completed for product ID %d\n", productID)
	return nil
}

// downloadImage downloads an image from the given URL
func downloadImage(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %v", err)
	}
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}

// compressImage compresses an image with the given quality and resize parameters
func compressImage(imageData []byte, quality int, newSize uint) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	if newSize > 0 {
		img = resize.Resize(newSize, 0, img, resize.Lanczos3)
	}

	return img, nil
}

// saveImage saves an image to a file in JPEG format
func saveImage(img image.Image, outputPath string) error {
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer out.Close()

	err = jpeg.Encode(out, img, nil)
	if err != nil {
		return fmt.Errorf("failed to encode image: %v", err)
	}

	return nil
}


// Example: product creation handler in your API code
// func CreateProductHandler(c *gin.Context) {
//     var product models.Product
//     if err := c.ShouldBindJSON(&product); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     // Store the product details in the database
//     // TODO: Replace this with your actual database storage logic
//     // For now, we'll print the product details
//     fmt.Printf("Stored product details in the database: %+v\n", product)

//     // Perform image analysis
//     err := imageanalysis.CompressAndStoreImages(product.ProductID, product.ProductImages)
//     if err != nil {
//         fmt.Printf("Error during image analysis: %v\n", err)
//         // Handle the error as needed
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform image analysis"})
//         return
//     }

//     c.JSON(http.StatusCreated, product)
// }

// updateDatabaseWithImagePath updates the database with the local file path of the compressed image
func updateDatabaseWithImagePath(productID int, imagePath string) error {
	// Create a MongoDB client (you should reuse your existing client)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("your_mongo_uri_here"))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Get the product collection
	productCollection := client.Database("your_database_name_here").Collection("products")

	// Define the filter to find the product by ID
	filter := bson.M{"product_id": productID}

	// Define the update to set the compressed_product_images field
	update := bson.M{
		"$push": bson.M{
			"compressed_product_images": imagePath,
		},
	}

	// Update the document in the database
	_, err = productCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update database with image path: %v", err)
	}

	return nil
}