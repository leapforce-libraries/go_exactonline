package exactonline

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/iterator"
)

// BigQueryGetRefreshToken get refreshtoken from BigQuery
//
func (eo *ExactOnline) GetTokenFromBigQuery() error {
	// create client
	bqClient, err := eo.BigQuery.CreateClient()
	if err != nil {
		fmt.Println("\nerror in BigQueryCreateClient")
		return err
	}

	ctx := context.Background()

	//sql := "SELECT Value FROM `" + BIGQUERY_DATASET + "." + BIGQUERY_TABLENAME + "` WHERE key = '" + key + "'"
	sql := "SELECT Value AS RefreshToken FROM `" + eo.BigQueryDataset + "." + eo.BigQueryTablename + "` WHERE key = '" + eo.RefreshTokenKey + "'"

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

	sql := "MERGE `" + eo.BigQueryDataset + "." + eo.BigQueryTablename + "` AS TARGET " +
		"USING  (select '" + eo.RefreshTokenKey + "' AS key,'" + eo.Token.RefreshToken + "' AS value) AS SOURCE " +
		" ON TARGET.key = SOURCE.key " +
		"WHEN MATCHED THEN " +
		"	UPDATE " +
		"	SET value = SOURCE.value " +
		"WHEN NOT MATCHED BY TARGET THEN " +
		"	INSERT (key, value) " +
		"	VALUES (SOURCE.key, SOURCE.value)"

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
