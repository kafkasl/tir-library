package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
)

func Message(status bool, message string) (map[string]interface{}) {
	return map[string]interface{} {"status" : status, "message" : message}
}

func Respond(w http.ResponseWriter, data map[string] interface{}, status int)  {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func DeleteCreatedEntities(db *gorm.DB) func() {
	type entity struct {
		table   string
		keyname string
		key     interface{}
	}
	var entries []entity
	hookName := "cleanupHook"

	// Setup the onCreate Hook
	db.Callback().Create().After("gorm:create").Register(hookName, func(scope *gorm.Scope) {
		fmt.Printf("Inserted entities of %s with %s=%v\n", scope.TableName(), scope.PrimaryKey(), scope.PrimaryKeyValue())
		entries = append(entries, entity{table: scope.TableName(), keyname: scope.PrimaryKey(), key: scope.PrimaryKeyValue()})
	})
	return func() {
		// Remove the hook once we're done
		defer db.Callback().Create().Remove(hookName)
		// Find out if the current db object is already a transaction
		_, inTransaction := db.CommonDB().(*sql.Tx)
		tx := db
		if !inTransaction {
			tx = db.Begin()
		}
		// Loop from the end. It is important that we delete the entries in the
		// reverse order of their insertion
		for i := len(entries) - 1; i >= 0; i-- {
			entry := entries[i]
			fmt.Printf("Deleting entities from '%s' table with key %v\n", entry.table, entry.key)
			tx.Table(entry.table).Where(entry.keyname+" = ?", entry.key).Delete("")
		}

		if !inTransaction {
			tx.Commit()
		}
	}
}