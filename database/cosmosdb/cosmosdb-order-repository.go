package cosmosdb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/mitchellh/mapstructure"
	"github.com/perocha/goutils/pkg/telemetry"
)

// CosmosDB repository
type CosmosdbRepository struct {
	client    ClientInterface
	database  DatabaseClientInterface
	container ContainerClientInterface
}

// Initialize CosmosDB repository using the provided connection string
func NewCosmosdbRepository(ctx context.Context, endPoint, connectionString, databaseName, containerName string) (*CosmosdbRepository, error) {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	// Create a new default azure credential
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		telemetryClient.TrackException(ctx, "CosmosdbRepository::NewCosmosdbRepository::Error creating default azure credential", err, telemetry.Critical, nil, true)
		return nil, err
	}

	// Create a new CosmosDB client
	clientOptions := azcosmos.ClientOptions{
		EnableContentResponseOnWrite: true,
	}
	client, err := azcosmos.NewClient(endPoint, credential, &clientOptions)
	if err != nil {
		telemetryClient.TrackException(ctx, "CosmosdbRepository::NewCosmosdbRepository::Error creating client", err, telemetry.Critical, nil, true)
		return nil, err
	}
	cosmosClient := &CosmosClient{client: client}

	// Retrieve database
	database, err := client.NewDatabase(databaseName)
	if err != nil {
		return nil, err
	}
	databaseClient := &CosmosDatabase{database: database}

	// Create a new container
	container, err := database.NewContainer(containerName)
	if err != nil {
		return nil, err
	}
	containerClient := &CosmosContainer{container: container}

	return &CosmosdbRepository{
		client:    cosmosClient,
		database:  databaseClient,
		container: containerClient,
	}, nil
}

// Creates a new document in CosmosDB
func (r *CosmosdbRepository) CreateDocument(ctx context.Context, partitionKey string, document interface{}) error {
	startTime := time.Now()
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	// Convert document to map[string]interface{}
	documentMap := make(map[string]interface{})
	if err := mapstructure.Decode(document, &documentMap); err != nil {
		return err
	}

	// Convert document to json
	docJson, err := json.Marshal(document)
	if err != nil {
		return err
	}

	// Create partition key
	pk := azcosmos.NewPartitionKeyString(partitionKey)

	// Create an item
	_, err = r.container.CreateItem(ctx, pk, docJson, nil)
	if err != nil {
		return err
	}

	// Construct telemetry properties from document content
	telemetryProps := make(map[string]string)
	for key, value := range documentMap {
		telemetryProps[key] = fmt.Sprintf("%v", value)
	}

	telemetryClient.TrackDependency(ctx, "CosmosdbRepository", "CreateDocument", "CosmosDB", r.client.Endpoint(), true, startTime, time.Now(), telemetryProps, true)

	return nil
}

// Updates an existing document in CosmosDB
func (r *CosmosdbRepository) UpdateDocument(ctx context.Context, partitionKey string, id string, document interface{}) error {
	startTime := time.Now()
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	// Create partition key
	pk := azcosmos.NewPartitionKeyString(partitionKey)

	// Convert document to json
	docJson, err := json.Marshal(document)
	if err != nil {
		return err
	}

	// Update an item
	_, err = r.container.UpsertItem(ctx, pk, docJson, nil)
	if err != nil {
		return err
	}

	telemetryClient.TrackDependency(ctx, "CosmosdbRepository", "UpdateDocument", "CosmosDB", r.client.Endpoint(), true, startTime, time.Now(), nil, true)

	return nil
}

// Deletes an document from CosmosDB
func (r *CosmosdbRepository) DeleteDocument(ctx context.Context, partitionKey string, id string) error {
	startTime := time.Now()
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	// Create partition key
	pk := azcosmos.NewPartitionKeyString(partitionKey)

	// Delete an item
	_, err := r.container.DeleteItem(ctx, pk, id, nil)
	if err != nil {
		return err
	}

	properties := map[string]string{
		"Id":           id,
		"PartitionKey": partitionKey,
	}
	telemetryClient.TrackDependency(ctx, "CosmosdbRepository", "DeleteDocument", "CosmosDB", r.client.Endpoint(), true, startTime, time.Now(), properties, true)

	return nil
}

// Retrieves an document from CosmosDB
func (r *CosmosdbRepository) GetDocument(ctx context.Context, partitionKey string, id string) (interface{}, error) {
	startTime := time.Now()
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	// Create partition key
	pk := azcosmos.NewPartitionKeyString(partitionKey)

	// Retrieve an item
	item, err := r.container.ReadItem(ctx, pk, id, nil)
	if err != nil {
		return nil, err
	}

	// Convert item to document
	var readDoc map[string]interface{}
	err = json.Unmarshal(item.Value, &readDoc)
	if err != nil {
		return nil, err
	}

	telemetryClient.TrackDependency(ctx, "CosmosdbRepository", "GetDocument", "CosmosDB", r.client.Endpoint(), true, startTime, time.Now(), nil, true)

	return readDoc, nil
}
