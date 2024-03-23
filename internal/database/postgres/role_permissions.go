package postgres

import (
	"fmt"
	"github.com/lib/pq"
)

func (db *DB) GetPatientIdPermissions(roleId int, path string) ([]string, error) {
	fmt.Println(path)
	stmt := `SELECT actions FROM role_permissions WHERE role_id=$1 AND paths=$2 LIMIT 1;`

	row := db.DB.QueryRow(stmt, roleId, path)

	var textArray pq.StringArray
	err := row.Scan(&textArray)

	if err != nil {
		return nil, err
	}

	permissions := []string(textArray)

	return permissions, nil
}
