package nginx_parser

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func SaveVisitorToDB(visitors []Visitor) error {
	db, err := NewDb()
	if err != nil {
		return err
	}
	defer db.Close()

	for _, visitor := range visitors {
		fmt.Printf("Visitor: %+v\n", visitor)

		id := uuid.NewString()
		query := `INSERT INTO visitor (id, host_name, request_ip, request_ip_location, request_time, request_method, 
			  request_uri, user_agent, response_status, created_at, created_by, is_deleted) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err = db.Exec(query,
			id,
			"BMCM",
			visitor.IP,
			visitor.IPLocation,
			visitor.RequestTime,
			visitor.RequestMethod,
			visitor.RequestURI,
			visitor.UserAgent,
			visitor.ResponseStatus,
			time.Now(),
			"system",
			false)
		if err != nil {
			return err
		}

		// fmt.Printf("Visitor saved to DB: %+v\n", visitor)
	}

	return nil
}
