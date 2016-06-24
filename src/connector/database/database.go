package database

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/tebben/sensorthings-connector/src/connector/models"
)

var open bool
var connectorBucketName = "connectors"

type Database struct {
	bolt *bolt.DB
}

func (db *Database) Open(dbPath string) error {
	var err error
	dbFile := dbPath
	config := &bolt.Options{Timeout: 1 * time.Second}
	db.bolt, err = bolt.Open(dbFile, 0600, config)
	if err != nil {
		log.Fatal(err)
	}

	db.bolt.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(connectorBucketName))
		return nil
	})

	open = true
	return nil
}

func (db *Database) Close() {
	open = false
	db.bolt.Close()
}

// InsertConnector inserts or updates a connector in the database
func (db *Database) InsertConnector(connector *models.ConnectorBase) error {
	if !open {
		return fmt.Errorf("db must be opened before saving!")
	}
	err := db.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(connectorBucketName))
		enc, err := json.Marshal(connector)
		if err != nil {
			return fmt.Errorf("could not encode module %s: %s", connector.GetName(), err)
		}

		err = b.Put([]byte(connector.GetID()), enc)
		return err
	})
	return err
}

// GetConnectors loads all connectors from the database
func (db *Database) GetConnectors() ([]*models.ConnectorBase, error) {
	if !open {
		return nil, fmt.Errorf("db must be opened before reading!")
	}

	connectors := make([]*models.ConnectorBase, 0)
	err := db.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(connectorBucketName))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			con := &models.ConnectorBase{}
			err := json.Unmarshal(v, &con)
			if err != nil {
				log.Printf("Error loading connector fro db: %v", string(k[:]))
				continue
			}

			connectors = append(connectors, con)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Could not get connetors")
		return nil, err
	}

	return connectors, nil
}

// DeleteConnector removes an connector from the database
func (db *Database) DeleteConnector(id string) error {
	if !open {
		return fmt.Errorf("db must be opened before saving!")
	}

	err := db.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(connectorBucketName))
		err := b.Delete([]byte(id))
		return err
	})

	return err
}

// SaveConnectorState saves the running state of a module
func (db *Database) SaveConnectorState(id string, running bool) error {
	if !open {
		return fmt.Errorf("db must be opened before saving!")
	}

	err := db.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(connectorBucketName))
		c := b.Get([]byte(id))
		if c != nil {
			con := &models.ConnectorBase{}
			if err := json.Unmarshal(c, &con); err != nil {
				return err
			} else {
				con.Running = running
				enc, _ := json.Marshal(con)
				if err = b.Put([]byte(con.GetID()), enc); err != nil {
					return err
				}

			}
		}

		return nil
	})

	return err
}
