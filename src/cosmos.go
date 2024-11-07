package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

func startCosmos(writeOutput func(msg string)) error {

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, proceeding without it")
	}

	// <create_client>
	endpoint, found := os.LookupEnv("CONFIGURATION__AZURECOSMOSDB__ENDPOINT")
	if !found {
		panic("Azure Cosmos DB for Table account endpoint not set.")
	}

	log.Println("ENDPOINT: ", endpoint)

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	service, err := aztables.NewServiceClient(endpoint, credential, nil)
	if err != nil {
		panic(err)
	}
	// </create_client>

	writeOutput("Current Status:\tStarting...")

	tableName, found := os.LookupEnv("CONFIGURATION__AZURECOSMOSDB__TABLENAME")
	if !found {
		tableName = "cosmicworks-products"
	}

	log.Println("TABLE: ", endpoint)

	table := service.NewClient(tableName)

	writeOutput(fmt.Sprintf("Get table:\t%s", tableName))

	{
		entity := aztables.EDMEntity{
			Entity: aztables.Entity{
				RowKey:       "aaaaaaaa-0000-1111-2222-bbbbbbbbbbbb",
				PartitionKey: "gear-surf-surfboards",
			},
			Properties: map[string]any{
				"Name":      "Yamba Surfboard",
				"Quantity":  12,
				"Price":     850.00,
				"Clearance": false,
			},
		}

		context := context.TODO()

		bytes, err := json.Marshal(entity)
		if err != nil {
			panic(err)
		}

		_, err = table.UpsertEntity(context, bytes, nil)
		if err != nil {
			panic(err)
		}

		writeOutput(fmt.Sprintf("Upserted entity:\t%v", entity))
	}

	{
		entity := aztables.EDMEntity{
			Entity: aztables.Entity{
				RowKey:       "bbbbbbbb-1111-2222-3333-cccccccccccc",
				PartitionKey: "gear-surf-surfboards",
			},
			Properties: map[string]any{
				"Name":      "Kiama Classic Surfboard",
				"Quantity":  25,
				"Price":     790.00,
				"Clearance": true,
			},
		}

		context := context.TODO()

		bytes, err := json.Marshal(entity)
		if err != nil {
			panic(err)
		}

		_, err = table.UpsertEntity(context, bytes, nil)
		if err != nil {
			panic(err)
		}

		writeOutput(fmt.Sprintf("Upserted entity:\t%v", entity))
	}

	{
		context := context.TODO()

		rowKey := "aaaaaaaa-0000-1111-2222-bbbbbbbbbbbb"
		partitionKey := "gear-surf-surfboards"

		response, err := table.GetEntity(context, partitionKey, rowKey, nil)
		if err != nil {
			panic(err)
		}

		var entity aztables.EDMEntity
		err = json.Unmarshal(response.Value, &entity)
		if err != nil {
			panic(err)
		}

		writeOutput(fmt.Sprintf("Read item row key:\t%s", entity.RowKey))
	}

	{
		filter := "PartitionKey eq 'gear-surf-surfboards'"

		options := &aztables.ListEntitiesOptions{
			Filter: &filter,
		}

		pager := table.NewListEntitiesPager(options)

		context := context.TODO()

		for pager.More() {
			response, err := pager.NextPage(context)
			if err != nil {
				panic(err)
			}
			for _, entityBytes := range response.Entities {
				var entity aztables.EDMEntity
				err := json.Unmarshal(entityBytes, &entity)
				if err != nil {
					panic(err)
				}

				writeOutput(fmt.Sprintf("Found item:\t%s\t%s", entity.Properties["Name"], entity.RowKey))
			}
		}
	}

	return nil
}
