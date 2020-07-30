package dbs

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/rtntubmt97/profiler/internal/defines"
	"github.com/rtntubmt97/profiler/internal/utils"
	app "github.com/rtntubmt97/profiler/pkg/applications"
	k "github.com/rtntubmt97/profiler/pkg/kernel"
)

const pkgName = "dbs"

type profileMongoDb struct {
	client            *mongo.Client
	testDb            *mongo.Database
	profileCollection *mongo.Collection
}

var profileMongoDbInstance *profileMongoDb

func init() {
	var err error
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
		return
	}
	ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	testDb := client.Database("test")
	profileCollection := testDb.Collection("profile")

	profileMongoDbInstance = new(profileMongoDb)
	profileMongoDbInstance.client = client
	profileMongoDbInstance.testDb = testDb
	profileMongoDbInstance.profileCollection = profileCollection
}

func MongoDbInstance() *profileMongoDb {
	return profileMongoDbInstance
}

var profiler k.Profiler = app.HttpPageProfiler()

func (db *profileMongoDb) RetrieveProfile(id int64) (error, defines.Profile) {
	mark := k.CreateMark()
	defer profiler.Record("db.RetrieveProfile", mark)
	var profile defines.Profile
	err := db.profileCollection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&profile)
	return utils.WrapError(pkgName, "RetrieveProfile failed", err), profile
}

func (db *profileMongoDb) CreateProfile(profile defines.Profile) error {
	_, err := db.profileCollection.InsertOne(context.Background(),
		bson.M{"id": profile.Id, "name": profile.Name, "job": profile.Job})
	return utils.WrapError(pkgName, "CreateProfile failed", err)
}
