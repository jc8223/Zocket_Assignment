package api

import (
	"context"
	"net/http"
	"time"
	"zocket-assignment/db"
	"zocket-assignment/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
    // "go.mongodb.org/mongo-driver/mongo/options"

)

 
 
func CreateProductHandler(c *gin.Context)  {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	product.Created_At = time.Now()
	product.Updated_At = time.Now()

	productID, err := insertProduct(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product created successfully", "product_id": productID})
}

func insertProduct(product models.Product) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
 
    collection :=  Mongodb.Client.Database(db.DbName).Collection(db.ProductCollection)
 
    result, err := collection.InsertOne(ctx, product)
    if err != nil {
        return "", err
    }
 
    insertedID, ok := result.InsertedID.(primitive.ObjectID)
    if !ok {
        return "", err
    }

    return insertedID.Hex(), nil
}


// FetchProductHandler handles fetching a single product by ID
func FetchProductHandler(c *gin.Context) {
	productID := c.Param("id")

	// Convert the product ID string to an ObjectID
	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Connect to the "products" collection in the "product_management" database
	collection := Mongodb.Client.Database(db.DbName).Collection(db.ProductCollection)

	// Search for the product by ID
	var product models.Product
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// FetchAllProductsHandler handles fetching all products
func FetchAllProductsHandler(c *gin.Context) {
	// Connect to the "products" collection in the "product_management" database
	collection := Mongodb.Client.Database(db.DbName).Collection(db.ProductCollection)

	// Search for all products
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	defer cursor.Close(context.TODO())

	var products []models.Product
	err = cursor.All(context.TODO(), &products)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// UpdateProductHandler handles updating a product by ID
func UpdateProductHandler(c *gin.Context) {
	productID := c.Param("id")

	// Convert the product ID string to an ObjectID
	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Connect to the "products" collection in the "product_management" database
	collection := Mongodb.Client.Database(db.DbName).Collection(db.ProductCollection)

	// Bind JSON request to Product struct
	var updatedProduct models.Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set updated timestamp
	updatedProduct.Updated_At = time.Now()

	// Update the product in the database
	result, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		bson.D{
			{"$set", updatedProduct},
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}


func DeleteProductHandler(c *gin.Context) {
    productID := c.Param("id")

    // Convert the product ID string to an ObjectID
    objID, err := primitive.ObjectIDFromHex(productID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    // Connect to the "products" collection in the "product_management" database
    collection := Mongodb.Client.Database(db.DbName).Collection(db.ProductCollection)

    // Delete the product from the database
    result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
        return
    }

    if result.DeletedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

























 // User structure representing the User model
type User struct {
    ID         int       `json:"id"`
    Name       string    `json:"name"`
    Mobile     string    `json:"mobile"`
    Latitude   float64   `json:"latitude"`
    Longitude  float64   `json:"longitude"`
    Created_At  time.Time `json:"created_at"`
    Updated_At  time.Time `json:"updated_at"`
}

// Assuming you have a userCollection variable representing your MongoDB collection
// var userCollection *mongo.Collection

// CreateUserHandler handles the creation of a new user
func CreateUserHandler(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Set timestamps
    user.Created_At = time.Now()
    user.Updated_At = time.Now()

	
    // Insert the user into the database
    userDocument := bson.M{
        "id":         user.ID,
        "name":       user.Name,
        "mobile":     user.Mobile,
        "latitude":   user.Latitude,
        "longitude":  user.Longitude,
        "created_at": user.Created_At,
        "updated_at": user.Updated_At,
    }

    // Replace "userCollection" with your actual MongoDB collection
	collection :=  Mongodb.Client.Database(db.DbName).Collection(db.UserCollection)
 
    _, err := collection.InsertOne(context.Background(), userDocument)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    c.JSON(http.StatusCreated, user)
}

// FetchUserHandler retrieves a user by ID
func FetchUserHandler(c *gin.Context) {
    userID := c.Param("id")

    var user User

    // Retrieve the user from the database by ID
    // Replace "userCollection" with your actual MongoDB collection
    filter := bson.M{"id": userID}
	collection :=  Mongodb.Client.Database(db.DbName).Collection(db.UserCollection)
 
    err := collection.FindOne(context.Background(), filter).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
        return
    }

    c.JSON(http.StatusOK, user)
}
 
// UpdateUserHandler handles updating a user by ID
func UpdateUserHandler(c *gin.Context) {
    userID := c.Param("id")

    var updatedUser User
    if err := c.ShouldBindJSON(&updatedUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Set the updated timestamp
    updatedUser.Created_At = time.Now()
    updatedUser.Updated_At = time.Now()

    // Update the user in the database by ID
    // Replace "userCollection" with your actual MongoDB collection
    filter := bson.M{"id": userID}
    update := bson.M{
        "$set": bson.M{
            "name":        updatedUser.Name,
            "mobile":      updatedUser.Mobile,
            "latitude":    updatedUser.Latitude,
            "longitude":   updatedUser.Longitude,
            "created_at":  updatedUser.Created_At,
            "updated_at":  updatedUser.Updated_At,
        },
    }
	collection :=  Mongodb.Client.Database(db.DbName).Collection(db.UserCollection)
 
    _, err := collection.UpdateOne(context.Background(), filter, update)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }

    c.JSON(http.StatusOK, updatedUser)
}

// DeleteUserHandler handles deleting a user by ID
func DeleteUserHandler(c *gin.Context) {
    userID := c.Param("id")

    // Delete the user from the database by ID
    // Replace "userCollection" with your actual MongoDB collection
    filter := bson.M{"id": userID}
	collection :=  Mongodb.Client.Database(db.DbName).Collection(db.UserCollection)
 
    _, err := collection.DeleteOne(context.Background(), filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// FetchAllUsersHandler retrieves all users
func FetchAllUsersHandler(c *gin.Context) {
    var users []User

    // Retrieve all users from the database
    // Replace "userCollection" with your actual MongoDB collection
	collection :=  Mongodb.Client.Database(db.DbName).Collection(db.UserCollection)
 
    cursor, err := collection.Find(context.Background(), bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
        return
    }
    defer cursor.Close(context.Background())

    for cursor.Next(context.Background()) {
        var user User
        err := cursor.Decode(&user)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user"})
            return
        }
        users = append(users, user)
    }

    c.JSON(http.StatusOK, users)
}