package main

import (
	"context"
	"fmt"
	"log"

	//"os"

	//"strings"

	//"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"rishik.com/constants"
	"rishik.com/db"
	"rishik.com/enums"
)

// Role reprresents the role dicument structure
type Role struct {
	Name        string   `bson:"name"`
	Permissions []string `bson:"permissions"`
}

func seedRoles() error {
	log.Printf("seeding roles started...")

	rolesCollection := db.Database.Collection("roles")
	session, err := db.Client.StartSession()
	if err != nil {
		log.Fatalf("Failed to start session: %v", err)
	}
	defer session.EndSession(context.Background())
	err = mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		if err := sc.StartTransaction(); err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}
		log.Printf("clearing existing roles...")
		if _, err := rolesCollection.DeleteMany(sc, bson.M{}); err != nil {
			sc.AbortTransaction(sc)
			return fmt.Errorf("faile to clear roles: %w", err)

		}
		for roleName, permissions := range constants.RolePermissions {
			var existingRole Role
			err := rolesCollection.FindOne(sc, bson.M{"name": roleName}).Decode(&existingRole)
			if err == mongo.ErrNoDocuments {
				newRole := Role{
					Name:        string(roleName),
					Permissions: convertPermissiion(permissions),
				}
				_, err := rolesCollection.InsertOne(sc, newRole)
				if err != nil {
					sc.AbortTransaction(sc)
					return fmt.Errorf("failed to insert role: %w", err)
				}
				log.Printf("Inserted role: %s with permissions: %v", roleName, permissions)
			} else if err != nil {
				sc.AbortTransaction(sc)
				return fmt.Errorf("failed to check existing role: %w", err)
			} else {
				log.Printf("Role %s already exists, skipping insertion.", roleName)
			}

		}
		if err := sc.CommitTransaction(sc); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
		log.Println("Transactions commited")
		return nil
	})
	if err != nil {
		return fmt.Errorf("error during  %w", err)
	}
	log.Printf("seeding commleted sucessfully")
	return nil

}
func convertPermissiion(permissions []enums.PermissionType) []string {
	var result []string
	for _, p := range permissions {
		result = append(result, string(p))
	}
	return result
}
func main() {

	uri := "mongodb+srv://Realizz:fq8rCmCCfe6XtEVo@cluster0.7ua0zvw.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

	err := db.InitDB(uri, "teamsync_db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	if err := seedRoles(); err != nil {
		log.Fatalf("Error running seed script: %v", err)
	}
}
