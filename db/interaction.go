package db

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/tofa-project/client-daemon/lib"
)

// Retrieves all records from 'apps' table formatted
func GetApps() []j {
	rows := Query("SELECT * FROM apps;")
	defer rows.Close()
	response := make([]j, 0)

	for rows.Next() {
		var id, data string

		err := rows.Scan(&id, &data)
		if err != nil {
			fmt.Println("Some error occurred: ", err, "Possibly corrupted data skipping...")
			continue
		}

		dataMap := Decrypt(data)

		dataMap["id"] = id
		response = append(response, dataMap)
	}

	return response
}

// Retrieves row data by its ID formatted
func GetAppByID(id string) j {
	var i, data string

	res := Instance.QueryRow("SELECT * FROM apps WHERE id = ?", id)
	err := res.Scan(&i, &data)
	if err != nil {
		return nil
	}

	dataMap := Decrypt(data)
	dataMap["id"] = i

	return dataMap
}

// Creates a new application with afferent data and inserts it in DB.
// Returns newly created app's ID
func MakeApp(data j) string {
	dataEnc := Encrypt(data)
	idBytes := md5.Sum([]byte(lib.GenerateRandomString(32)))
	idHash := hex.EncodeToString(idBytes[:])

	Exec("INSERT INTO apps VALUES(?, ?);", idHash, dataEnc)

	return idHash
}

// Updates an app based on ID. Overwrites existing data
func UpdateApp(appID string, data j) {
	dataEnc := Encrypt(data)

	Exec("UPDATE apps SET data = ? WHERE id = ?", dataEnc, appID)
}

// Permanently deletes app from database
func DeleteApp(appID string) {
	Exec("DELETE FROM apps WHERE id = ?", appID)
}
