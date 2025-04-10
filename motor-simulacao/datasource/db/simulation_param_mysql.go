package db

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	simulationParamRepository struct {
		mongoClient *mongo.Client
	}
	SimulationParamRepository interface {
		FindParams() ([]Document, error)
	}
	Document struct{
		Class string `bson:"classe"`
		Rate float64 `bson:"taxa"`
	}
)

func NewSimulationParamRepository(mongoClient *mongo.Client) SimulationParamRepository {
	return &simulationParamRepository{
		mongoClient: mongoClient,
	}
}

func (r *simulationParamRepository) FindParams() ([]Document, error) {
	// collection := r.mongoClient.Database("simlador_credito").Collection("parametros_classe_usuarios")

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// cursor, err := collection.Find(ctx, bson.D{}, nil)
	// if err != nil {
	// 	if err == mongo.ErrNoDocuments {
	// 		return nil, errors.New("document not found")
	// 	}
	// 	log.Printf("Error finding document: %v", err)
	// 	return nil, err
	// }

	// defer cursor.Close(ctx)

	
	// var result []Document
	// for ;cursor.Next(ctx); {
	// 	var doc Document
	// 	if err := cursor.Decode(&doc); err != nil {
	// 		log.Printf("Error decoding document: %v", err)
	// 		return nil, err
	// 	}
	// 	result = append(result, doc)
	// }

	// println(result)

	result := []Document{
		{
			Class: "25-",
			Rate: 0.05,
		},
		{
			Class: "26-40",
			Rate: 0.03,
		},
		{
			Class: "41-60",
			Rate: 0.02,
		},
		{
			Class: "61+",
			Rate: 0.04,
		},
	}

	return result, nil
}
