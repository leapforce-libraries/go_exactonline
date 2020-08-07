package exactonline

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/iterator"
)

const tableRefreshToken string = "leapforce.refreshtokens"
const api string = "ExactOnline"

// BigQueryGetRefreshToken get refreshtoken from BigQuery
//
func (eo *ExactOnline) GetTokenFromBigQuery() error {
	fmt.Println("***GetTokenFromBigQuery***")
	// create client
	bqClient, err := eo.BigQuery.CreateClient()
	if err != nil {
		fmt.Println("\nerror in BigQueryCreateClient")
		return err
	}

	ctx := context.Background()

	//sql := "SELECT Value FROM `" + BIGQUERY_DATASET + "." + BIGQUERY_TABLENAME + "` WHERE key = '" + key + "'"
	//sql := "SELECT refreshtoken AS RefreshToken FROM `" + tableRefreshToken + "` WHERE client_id = '" + eo.ClientID + "'"
	sql := fmt.Sprintf("SELECT refreshtoken AS RefreshToken FROM `%s` WHERE api = '%s' AND client_id = '%s'", tableRefreshToken, api, eo.ClientID)

	//fmt.Println(sql)

	q := bqClient.Query(sql)
	it, err := q.Read(ctx)
	if err != nil {
		return err
	}

	token := new(Token)

	for {
		err := it.Next(token)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		break
	}
	//fmt.Println(token)

	/*
		if r.Value == "" {
			return nil, err
		}*/

	//token := new(Token)
	if eo.Token == nil {
		eo.Token = new(Token)
	}

	eo.Token.TokenType = "bearer"
	eo.Token.Expiry = time.Now().Add(-10 * time.Second)
	eo.Token.RefreshToken = token.RefreshToken
	eo.Token.AccessToken = ""

	//eo.Token = token

	return nil
}

// BigQuerySaveToken saves refreshtoken to BigQuery
//
func (eo *ExactOnline) SaveTokenToBigQuery() error {
	// create client
	bqClient, err := eo.BigQuery.CreateClient()
	if err != nil {
		fmt.Println("\nerror in BigQueryCreateClient")
		return err
	}

	//fmt.Println("[save]", eo.Token.RefreshToken[0:20])

	ctx := context.Background()

	sql := "MERGE `" + tableRefreshToken + "` AS TARGET " +
		"USING  (SELECT '" + api + "' AS api,'" + eo.ClientID + "' AS client_id,'" + eo.Token.RefreshToken + "' AS refreshtoken) AS SOURCE " +
		" ON TARGET.api = SOURCE.api " +
		" AND TARGET.client_id = SOURCE.client_id " +
		"WHEN MATCHED THEN " +
		"	UPDATE " +
		"	SET refreshtoken = SOURCE.refreshtoken " +
		"WHEN NOT MATCHED BY TARGET THEN " +
		"	INSERT (api, client_id, refreshtoken) " +
		"	VALUES (SOURCE.api, SOURCE.client_id, SOURCE.refreshtoken)"

	q := bqClient.Query(sql)

	job, err := q.Run(ctx)
	if err != nil {
		return err
	}

	for {
		status, err := job.Status(ctx)
		if err != nil {
			return err
		}
		if status.Done() {
			if status.Err() != nil {
				return status.Err()
				//log.Fatalf("Job failed with error %v", status.Err())
			}
			break
		}
		time.Sleep(1 * time.Second)
	}

	return nil
}
