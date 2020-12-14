package srclient

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
We are going to test the client meant to be used for testing
*/

var srClient MockSchemaRegistryClient
var schema = "{  " +
	"\"type\": \"record\",  " +
	"\"namespace\": \"com.mycorp.mynamespace\",  " +
	"\"name\": \"value_cdc_fake_2\", " +
	"\"doc\": \"Sample schema to help you get started.\", " +
	"\"fields\": [    " +
	"{   " +
	"\"name\": \"aField\"," +
	"\"type\": \"int\", " +
	"\"doc\": \"The int type is a 32-bit signed integer.\"   " +
	"}" +
	"]" +
	"}"

var schema2 = "{  " +
	"\"type\": \"record\",  " +
	"\"namespace\": \"com.mycorp.mynamespace\",  " +
	"\"name\": \"value_cdc_fake_2\", " +
	"\"doc\": \"Sample schema to help you get started.\", " +
	"\"fields\": [    " +
	"{   " +
	"\"name\": \"bField\"," +
	"\"type\": \"int\", " +
	"\"doc\": \"The int type is a 32-bit signed integer.\"   " +
	"}" +
	"]" +
	"}"

/*
We will use init to register some schemas to run our test with.
Due to this, the function that tests MockSchemaRegistryClient.CreateSchema will actually just assert that
the values creates by init are correct and expected.
*/
func init() {
	srClient = CreateMockSchemaRegistryClient("mock://testingUrl")

	// Test Schema creation
	_, _ = srClient.CreateSchema("test1", schema, Avro)
	// Test version upgrades and more registration
	_, _ = srClient.CreateSchema("test1", schema2, Avro)

}

func TestMockSchemaRegistryClient_CreateSchema(t *testing.T) {

	/*
	 Assert Schemas are registered with proper IDs and Versions
	 By virtue of this test, we also test MockSchemaRegistryClient.GetSchema
	*/
	schemaReg1, _ := srClient.GetSchema(1)
	assert.Equal(t, schema, schemaReg1.schema)
	assert.Equal(t, 1, schemaReg1.version)
	schemaReg2, _ := srClient.GetSchema(2)
	assert.Equal(t, schema2, schemaReg2.schema)
	assert.Equal(t, 2, schemaReg2.version)

	// Test registering already registered schema
	_, err := srClient.CreateSchema("test1", schema, Avro)
	assert.EqualError(t, err, "POST \"mock://testingUrl/subjects/test1/versions\": Schema already registered with id 1")

}

func TestMockSchemaRegistryClient_GetLatestSchema(t *testing.T) {

	latest, err := srClient.GetLatestSchema("test1")
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	} else {
		assert.Equal(t, schema2, latest.schema)
	}
}

func TestMockSchemaRegistryClient_GetSchemaVersions(t *testing.T) {
	versions, _ := srClient.GetSchemaVersions("test1")
	assert.Equal(t, 2, len(versions))

}

func TestMockSchemaRegistryClient_GetSchemaByVersion(t *testing.T) {

	oldVersion, _ := srClient.GetSchemaByVersion("test1", 1)
	assert.Equal(t, schema, oldVersion.schema)

}

func TestMockSchemaRegistryClient_GetSubjects(t *testing.T) {
	allSubjects, _ := srClient.GetSubjects()
	sort.Strings(allSubjects)
	assert.Equal(t, allSubjects, []string{"test1"})
}
