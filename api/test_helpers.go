package api

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/nivb52/hotel-rent/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DEFAULT_testdb_uri = "mongodb://localhost:27017"
const dbname = "hotel-rent-testing"

type testdb struct {
	db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.Store.User.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func SetupTest(t *testing.T) *testdb {
	err := godotenv.Load("../.env", "../.env.test.local")
	if err != nil {
		log.Fatal("Error loading .env files")
	}
	envDBuri := os.Getenv("TESTDB_CONNECTION_STRING")
	dburi := db.GetDBUri(envDBuri, DEFAULT_testdb_uri)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client, dbname)
	store := db.Store{
		User:    db.NewMongoUserStore(client, dbname),
		Hotel:   hotelStore,
		Room:    db.NewMongoRoomStore(client, hotelStore, dbname),
		Booking: db.NewMongoBookingStore(client, dbname),
	}

	return &testdb{
		Store: store,
	}
}

func (tdb *testdb) afterAll(t *testing.T) {
	defer tdb.teardown(t)
}

// Data Compare -- START
type FieldValueType int

const (
	String FieldValueType = iota
	Int
	Float32
	Bool
	Date
	// Add more field types as needed
)

type TestByField struct {
	ErrStr    string
	ValueType FieldValueType
}

type DataCompare struct {
	Fields   map[string]TestByField
	source   any
	testData any
}

func NewDataCompare(feildDefenitions map[string]TestByField) *DataCompare {
	return &DataCompare{
		Fields: feildDefenitions,
	}
}

func (o *DataCompare) SetSource(prop any) error {
	o.source = prop
	return nil
}

func (o *DataCompare) SetTestData(prop any) error {
	o.testData = prop
	return nil
}

func (o *DataCompare) Compare(t *testing.T) {
	// var errors []error = make([]error, len( o.Fields))

	rActual := reflect.ValueOf(o.testData)
	rSource := reflect.ValueOf(o.source)
	for key, field := range o.Fields {
		actual := reflect.Indirect(rActual).FieldByName(key)
		expected := reflect.Indirect(rSource).FieldByName(key)
		switch field.ValueType {
		case String:
			if expected.String() != actual.String() {
				if len(field.ErrStr) > 0 {
					t.Errorf(field.ErrStr, expected, actual)
				} else {
					t.Errorf("expected %s to be %s but found %s", key, expected, actual)
				}
			}

		case Int:
			if !expected.Equal(actual) && expected.Int() != actual.Int() {
				if len(field.ErrStr) > 0 {
					t.Errorf(field.ErrStr, expected, actual)
				} else {
					t.Errorf("expected %s to be %d but found %d", key, expected.Int(), actual.Int())
				}
			}

		case Float32:
			if !expected.Equal(actual) {
				if len(field.ErrStr) > 0 {
					t.Errorf(field.ErrStr, expected, actual)
				} else {
					t.Errorf("expected %s to be %f but found %f", key, expected.Float(), actual.Float())
				}
			}

		case Bool:
			if !expected.Equal(actual) {
				if len(field.ErrStr) > 0 {
					t.Errorf(field.ErrStr, expected, actual)
				} else {
					t.Errorf("expected %s to be %s but found %s", key, expected, actual)
				}
			}

		case Date:
			if !expected.Equal(actual) && expected.Interface().(time.Time).Format("01/02/2006") != actual.Interface().(time.Time).Format("01/02/2006") {
				formattedActual := actual.Interface().(time.Time).Format("01/02/2006")
				formattedExpected := expected.Interface().(time.Time).Format("01/02/2006")
				if len(field.ErrStr) > 0 {
					t.Errorf(field.ErrStr, formattedExpected, formattedActual)
				} else {
					t.Errorf("expected %s to be %s but found %s", key, formattedExpected, formattedActual)
				}
			}

		default:
			t.Errorf("Type not implemented")
		}
	}
}

// Data Compare -- END
