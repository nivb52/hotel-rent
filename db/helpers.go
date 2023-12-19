package db

import "go.mongodb.org/mongo-driver/bson/primitive"

func CompareIDs(ID1 string, OID2 primitive.ObjectID) (bool, error) {
	OID1, err := primitive.ObjectIDFromHex(ID1)
	if err != nil {
		return false, err
	}

	return !OID1.IsZero() || OID1 == OID2, nil

}
