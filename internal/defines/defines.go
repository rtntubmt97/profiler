package defines

type Profile struct {
	Id   int64  `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	Job  string `json:"job" bson:"job"`
}

type ProfileDb interface {
	// createProfile(Profile) error
	retrieveProfile(id int64) (error, Profile)
	// updateProfile(Profile) error
	// deleteProfile(id int64) (error, Profile)
}

type Param map[string]interface{}

type Result struct {
	data interface{}
	error
}

type Handler func(Param) Result
