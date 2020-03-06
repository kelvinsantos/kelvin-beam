package store

type Scooter struct {
	Id   string `json:"id,omitempty" bson:"id,omitempty`
	Lat  string `json:"lat,omitempty" bson:"lat,omitempty`
	Lng  string `json:"lng,omitempty" bson:"lng,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}
