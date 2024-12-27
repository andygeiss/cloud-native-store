package gcp

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	"github.com/andygeiss/cloud-native-store/internal/app/config"
)

type CloudSpanner struct {
	cfg  *config.Config
	path string
}

func NewCloudSpanner(cfg *config.Config) *CloudSpanner {
	// Set up a connection to the Spanner database.
	path := fmt.Sprintf("projects/%s/instances/%s/databases/%s",
		cfg.PortCloudSpanner.ProjectID,
		cfg.PortCloudSpanner.InstanceID,
		cfg.PortCloudSpanner.DatabaseID,
	)

	return &CloudSpanner{
		cfg:  cfg,
		path: path,
	}
}

func (a *CloudSpanner) Delete(ctx context.Context, key string) (err error) {
	// Create a Spanner client
	client, err := spanner.NewClient(ctx, a.path)
	if err != nil {
		return err
	}
	defer client.Close()

	// Create a read-write transaction for deleting the row.
	table := a.cfg.PortCloudSpanner.Table
	m := spanner.Delete(table, spanner.Key{key})
	_, err = client.Apply(ctx, []*spanner.Mutation{m})
	if err != nil {
		return err
	}

	return nil
}

func (a *CloudSpanner) Get(ctx context.Context, key string) (value string, err error) {
	// Create a Spanner client.
	client, err := spanner.NewClient(ctx, a.path)
	if err != nil {
		return "", err
	}
	defer client.Close()

	// Create a read-only transaction.
	ro := client.Single()
	defer ro.Close()

	// Read a row from the table.
	table := a.cfg.PortCloudSpanner.Table
	row, err := ro.ReadRow(ctx, table, spanner.Key{key}, []string{"Value"})
	if err != nil {
		return "", err
	}

	// Get the value from the row.
	var val string
	if err := row.Column(0, &val); err != nil {
		return "", err
	}
	return val, nil
}

func (a *CloudSpanner) Put(ctx context.Context, key, value string) (err error) {
	// Create a Spanner client
	client, err := spanner.NewClient(ctx, a.path)
	if err != nil {
		return err
	}
	defer client.Close()

	// Create a read-write transaction.
	table := a.cfg.PortCloudSpanner.Table
	m := spanner.InsertOrUpdate(table, []string{"Key", "Value"}, []interface{}{key, value})
	if _, err = client.Apply(ctx, []*spanner.Mutation{m}); err != nil {
		return err
	}

	// Apply the insert mutation.
	_, err = client.Apply(ctx, []*spanner.Mutation{m})
	if err != nil {
		return err
	}
	return nil
}
