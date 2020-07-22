package defines

type Profile struct {
	Id   int64  `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	Job  string `json:"job" bson:"job"`
}

type ProfileDb interface {
	CreateProfile(Profile) error
	RetrieveProfile(id int64) (error, Profile)
	// UpdateProfile(Profile) error
	// DeleteProfile(id int64) (error, Profile)
}

type ResponseData struct {
	Description string      `json:"description"`
	Data        interface{} `json:"data"`
}
