package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Queries contains all prepared SQL queries.
type Queries struct {
	GetDashboardStats *sqlx.Stmt `query:"get-dashboard-stats"`

	InsertSubscriber                *sqlx.Stmt `query:"insert-subscriber"`
	UpsertSubscriber                *sqlx.Stmt `query:"upsert-subscriber"`
	UpsertBlacklistSubscriber       *sqlx.Stmt `query:"upsert-blacklist-subscriber"`
	GetSubscriber                   *sqlx.Stmt `query:"get-subscriber"`
	GetSubscribersByEmails          *sqlx.Stmt `query:"get-subscribers-by-emails"`
	GetSubscriberLists              *sqlx.Stmt `query:"get-subscriber-lists"`
	UpdateSubscriber                *sqlx.Stmt `query:"update-subscriber"`
	BlacklistSubscribers            *sqlx.Stmt `query:"blacklist-subscribers"`
	AddSubscribersToLists           *sqlx.Stmt `query:"add-subscribers-to-lists"`
	DeleteSubscriptions             *sqlx.Stmt `query:"delete-subscriptions"`
	UnsubscribeSubscribersFromLists *sqlx.Stmt `query:"unsubscribe-subscribers-from-lists"`
	DeleteSubscribers               *sqlx.Stmt `query:"delete-subscribers"`
	Unsubscribe                     *sqlx.Stmt `query:"unsubscribe"`

	// Non-prepared arbitrary subscriber queries.
	QuerySubscribers                       string `query:"query-subscribers"`
	QuerySubscribersTpl                    string `query:"query-subscribers-template"`
	DeleteSubscribersByQuery               string `query:"delete-subscribers-by-query"`
	AddSubscribersToListsByQuery           string `query:"add-subscribers-to-lists-by-query"`
	BlacklistSubscribersByQuery            string `query:"blacklist-subscribers-by-query"`
	DeleteSubscriptionsByQuery             string `query:"delete-subscriptions-by-query"`
	UnsubscribeSubscribersFromListsByQuery string `query:"unsubscribe-subscribers-from-lists-by-query"`

	CreateList      *sqlx.Stmt `query:"create-list"`
	GetLists        *sqlx.Stmt `query:"get-lists"`
	UpdateList      *sqlx.Stmt `query:"update-list"`
	UpdateListsDate *sqlx.Stmt `query:"update-lists-date"`
	DeleteLists     *sqlx.Stmt `query:"delete-lists"`

	CreateCampaign           *sqlx.Stmt `query:"create-campaign"`
	QueryCampaigns           *sqlx.Stmt `query:"query-campaigns"`
	GetCampaign              *sqlx.Stmt `query:"get-campaign"`
	GetCampaignForPreview    *sqlx.Stmt `query:"get-campaign-for-preview"`
	GetCampaignStats         *sqlx.Stmt `query:"get-campaign-stats"`
	GetCampaignStatus        *sqlx.Stmt `query:"get-campaign-status"`
	NextCampaigns            *sqlx.Stmt `query:"next-campaigns"`
	NextCampaignSubscribers  *sqlx.Stmt `query:"next-campaign-subscribers"`
	GetOneCampaignSubscriber *sqlx.Stmt `query:"get-one-campaign-subscriber"`
	UpdateCampaign           *sqlx.Stmt `query:"update-campaign"`
	UpdateCampaignStatus     *sqlx.Stmt `query:"update-campaign-status"`
	UpdateCampaignCounts     *sqlx.Stmt `query:"update-campaign-counts"`
	RegisterCampaignView     *sqlx.Stmt `query:"register-campaign-view"`
	DeleteCampaign           *sqlx.Stmt `query:"delete-campaign"`

	InsertMedia *sqlx.Stmt `query:"insert-media"`
	GetMedia    *sqlx.Stmt `query:"get-media"`
	DeleteMedia *sqlx.Stmt `query:"delete-media"`

	CreateTemplate     *sqlx.Stmt `query:"create-template"`
	GetTemplates       *sqlx.Stmt `query:"get-templates"`
	UpdateTemplate     *sqlx.Stmt `query:"update-template"`
	SetDefaultTemplate *sqlx.Stmt `query:"set-default-template"`
	DeleteTemplate     *sqlx.Stmt `query:"delete-template"`

	CreateLink        *sqlx.Stmt `query:"create-link"`
	RegisterLinkClick *sqlx.Stmt `query:"register-link-click"`

	// GetStats *sqlx.Stmt `query:"get-stats"`
}

// connectDB initializes a database connection.
func connectDB(host string, port int, user, pwd, dbName string, sslMode string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, pwd, dbName, sslMode))
	if err != nil {
		return nil, err
	}

	return db, nil
}

// compileSubscriberQueryTpl takes a arbitrary WHERE expressions
// to filter subscribers from the subscribers table and prepares a query
// out of it using the raw `query-subscribers-template` query template.
// While doing this, a readonly transaction is created and the query is
// dry run on it to ensure that it is indeed readonly.
func (q *Queries) compileSubscriberQueryTpl(exp string, db *sqlx.DB) (string, error) {
	tx, err := db.BeginTxx(context.Background(), &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return "", err
	}

	// Perform the dry run.
	if exp != "" {
		exp = " AND " + exp
	}
	stmt := fmt.Sprintf(q.QuerySubscribersTpl, exp)
	if _, err := tx.Exec(stmt, true, pq.Int64Array{}); err != nil {
		tx.Rollback()
		return "", err
	}

	return stmt, nil
}

// compileSubscriberQueryTpl takes a arbitrary WHERE expressions and a subscriber
// query template that depends on the filter (eg: delete by query, blacklist by query etc.)
// combines and executes them.
func (q *Queries) execSubscriberQueryTpl(exp, tpl string, listIDs []int64, db *sqlx.DB, args ...interface{}) error {
	// Perform a dry run.
	filterExp, err := q.compileSubscriberQueryTpl(exp, db)
	if err != nil {
		return err
	}

	if len(listIDs) == 0 {
		listIDs = pq.Int64Array{}
	}
	// First argument is the boolean indicating if the query is a dry run.
	a := append([]interface{}{false, pq.Int64Array(listIDs)}, args...)
	if _, err := db.Exec(fmt.Sprintf(tpl, filterExp), a...); err != nil {
		return err
	}

	return nil
}
