package stripe

import "github.com/google/uuid"

type TaskDescriptor struct {
	Name       string    `json:"name" yaml:"name" bson:"name"`
	Main       bool      `json:"main,omitempty" yaml:"main" bson:"main"`
	Account    string    `json:"account,omitempty" yaml:"account" bson:"account"`
	TransferID uuid.UUID `json:"transferID,omitempty" yaml:"transferID" bson:"transferID"`
}
