package Store

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

type Store struct {
	Db *sql.DB
}

func (s Store) Save(tweet NewTweet) (NewTweet, error) {

	if err := s.Db.QueryRow("INSERT INTO tweets(message,account_id) VALUES($1,$2) RETURNING id;", tweet.Message, tweet.AccountId).Scan(&tweet.Id); err != nil {
		return NewTweet{}, err
	}

	return tweet, nil
}
func (s Store) GetByAccountID(account_id int) (SomeSlice []NewTweet, err error) {
	query := "SELECT id, message, account_id FROM tweets WHERE account_id = $1;"
	row, err := s.Db.Query(query, account_id)
	if err != nil {
		fmt.Println(err)
		return
	}
	tweet := NewTweet{}
	for row.Next() {
		row.Scan(&tweet.Id, &tweet.Message, &tweet.AccountId)
		SomeSlice = append(SomeSlice, tweet)
	}

	return SomeSlice, nil
}
