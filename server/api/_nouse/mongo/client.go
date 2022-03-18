package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func insertBsonD(col *mongo.Collection) error {
	bsonD := bson.D{
		{Key: "str1", Value: "abc"},
		{Key: "num1", Value: 1},
		{Key: "str2", Value: "xyz"},
		{Key: "num2", Value: bson.A{2, 3, 4}},
		{Key: "subdoc", Value: bson.D{{Key: "str", Value: "subdoc"}, {Key: "num", Value: 987}}},
		{Key: "date", Value: time.Now()},
	}
	_, err := col.InsertOne(context.Background(), bsonD)
	return err
}

func insertBsonM(col *mongo.Collection) error {
	bsonM := bson.M{
		"str1":   "efg",
		"num1":   11,
		"str2":   "opq",
		"num2":   bson.A{12, 13, 14},
		"subdoc": bson.M{"str": "subdoc", "num": 987},
		"date":   time.Now(),
	}
	for i := 0; i < 10; i++ {
		_, err := col.InsertOne(context.Background(), bsonM)
		if err != nil {
			return err
		}
	}
	return nil
}

type myType struct {
	Str1   string
	Num1   int
	Str2   string
	Num2   []int
	Subdoc struct {
		Str string
		Num int
	}
	Date time.Time
}

func insertStruct(col *mongo.Collection) error {
	doc := myType{
		"hij",
		21,
		"rst",
		[]int{22, 23, 24},
		struct {
			Str string
			Num int
		}{"subdoc", 987},
		time.Now(),
	}
	_, err := col.InsertOne(context.Background(), doc)
	return err
}

func insertNextBsonM(col *mongo.Collection) error {
	var num int32
	var err error
	if num, err = findMaxNum(col); err != nil {
		return err
	}
	// log.Printf(">>> findedNum: %d\n", num)
	num++
	// log.Printf(">>> addedNum: %d\n", num)
	bsonM := bson.M{
		"str1":   "efg",
		"num1":   num,
		"str2":   "opq",
		"num2":   bson.A{12, 13, 14},
		"subdoc": bson.M{"str": "subdoc", "num": 987},
		"date":   time.Now(),
	}
	log.Printf("> Inserting: %d\n", num)
	_, err = col.InsertOne(context.Background(), bsonM)
	if err != nil {
		return err
	}
	return nil
}

func findMaxNum(col *mongo.Collection) (int32, error) {

	var ret int32
	options := options.Find()
	options.SetSort(bson.D{bson.E{Key: "num1", Value: -1}})
	options.SetLimit(1)
	cursor, err := col.Find(context.Background(), bson.D{}, options)
	if err != nil {
		return ret, err
	}
	for cursor.Next(context.Background()) {
		elem := bson.M{}
		err := cursor.Decode(&elem)
		if err != nil {
			return ret, err
		}
		ret = elem["num1"].(int32)
		// log.Printf(">> Max: %d\n", ret)
		// results = append(results, *elem)
		// lastValue = elem.Lookup("_id")
	}
	// log.Printf(">> Returning: %d\n", ret)
	return ret, nil
}

func mainMain() error {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/?connect=direct"))
	if err != nil {
		return err
	}
	if err = client.Connect(context.Background()); err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	col := client.Database("test").Collection("col")
	// if err = insertBsonD(col); err != nil {
	//     return err
	// }
	// if err = insertBsonM(col); err != nil {
	//     return err
	// }
	// if err = insertStruct(col); err != nil {
	//     return err
	// }
	for i := 0; i < 10000; i++ {
		// if _, err = findMaxNum(col); err != nil {
		//     return err
		// }
		if err = insertNextBsonM(col); err != nil {
			return err
		}
	}
	return nil
}

// func main() {
//     if err := mainMain(); err != nil {
//         log.Fatal(err)
//     }
//     log.Println("normal end.")
// }
