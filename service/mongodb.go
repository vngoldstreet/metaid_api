package service

import (
	"context"
	"errors"
	"log"
	"time"
	"vietvd/mql-api/entity"
	"vietvd/mql-api/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteMongoByID(id string) error {
	client := repository.GetMongoClient()

	collection := client.Database("FXGHDatabases").Collection("id_generates")
	idObj, err := primitive.ObjectIDFromHex(id) // Replace with your actual document ID
	if err != nil {
		return err
	}

	// Create a filter to match the document by _id
	filter := bson.M{"_id": idObj}

	// Delete the document
	deleteResult, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	log.Printf("Deleted %d document(s)\n", deleteResult.DeletedCount)
	return nil
}

func UpdateMongoByID(id string, vps_name string) error {
	client := repository.GetMongoClient()
	collection := client.Database("FXGHDatabases").Collection("id_generates")

	idObj, err := primitive.ObjectIDFromHex(id) // Replace with your actual document ID
	if err != nil {
		return err
	}

	// Define the update document
	update := bson.M{
		"$set": bson.M{
			"vps_name":   vps_name,
			"updated_at": time.Now(),
		},
	}

	// Update the document
	updateResult, err := collection.UpdateByID(context.Background(), idObj, update)
	if err != nil {
		return err
	}

	log.Printf("Matched %d document(s) and modified %d document(s)\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	return nil
}

func SetMongoDB(mql_id string, vps_name string) error {
	client := repository.GetMongoClient()
	collection := client.Database("FXGHDatabases").Collection("id_generates")
	document := entity.IDGenerate{
		UpdateAt: time.Now(),
		MQLID:    mql_id,
		VPSName:  vps_name,
	}

	_, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		return err
	}
	return nil
}

func GetMongoDB(vps_name string) (entity.IDGenerate, error) {
	client := repository.GetMongoClient()
	collection := client.Database("FXGHDatabases").Collection("id_generates")
	filter := bson.M{"vps_name": vps_name} // Replace with your actual filter condition

	var result entity.IDGenerate
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity.IDGenerate{}, errors.New("no document found")
		} else {
			return entity.IDGenerate{}, err
		}
	}
	return result, nil
}

func GetAllMongoData(vps_name string) (entity.IDGenerate, error) {
	client := repository.GetMongoClient()
	collection := client.Database("FXGHDatabases").Collection("id_generates")
	var results []entity.IDGenerate
	filter := bson.M{"vps_name": vps_name} // Replace with your actual filter condition

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return entity.IDGenerate{}, err
	}

	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var result entity.IDGenerate
		err := cursor.Decode(&result)
		if err != nil {
			return entity.IDGenerate{}, err
		}
		results = append(results, result)
	}
	if err := cursor.Err(); err != nil {
		return entity.IDGenerate{}, err
	}
	len := len(results)
	if len == 0 {
		return entity.IDGenerate{}, nil
	}
	return results[len-1], nil
}

func DeleteAllDataByName(vps_name string) error {
	client := repository.GetMongoClient()
	collection := client.Database("FXGHDatabases").Collection("id_generates")
	filter := bson.M{"vps_name": vps_name} // Replace with your actual filter condition

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var result entity.IDGenerate
		err := cursor.Decode(&result)
		if err != nil {
			return err
		}

		if err := UpdateMongoByID(result.ID.Hex(), ""); err != nil {
			return err
		}
	}
	if err := cursor.Err(); err != nil {
		return err
	}
	return nil
}
