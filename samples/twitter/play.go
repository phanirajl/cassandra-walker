package main

import (
	"github.com/gocql/gocql"
	"github.com/jozn/cassandra-walker/samples/twitter/xc"
)

func main() {
	// create cassandra session
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "twitter"
	cluster.Consistency = gocql.One
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Create
	twitt1 := xc.Twitt{
		Body:       "Hello World",
		CreateTime: 1566000000,
		TwiitId:    1,
		UserId:     1,
	}

	err := twitt1.Save(session)

	// Delete one object
	twitt1.Delete(session)

	//////////////// For Selector
	twitts, err := xc.NewTwitt_Selector().UserId_Eq(1).Limit(5).GetRows(session) // returns and array of twitts ( []*twitt ,err )

	twitt, err := xc.NewTwitt_Selector().UserId_Eq(1).Limit(5).GetRows(session) // returns a single twitt ( *twitt ,err )

	//can use clustering columns too
	twitts, err = xc.NewTwitt_Selector().UserId_Eq(1).And_TwiitId_In(1, 25, 68).GetRows(session)

	//can select just some columns, it will returns *[]Twitt, with just selected columns sets
	twitts, err = xc.NewTwitt_Selector().Select_UserId().Select_Body().UserId_Eq(1).And_TwiitId_In(1, 25, 68).Limit(12).GetRows(session)

	//for when need to use filtering
	twitts, err = xc.NewTwitt_Selector().UserId_LT_Filtering(100).Limit(10).AllowFiltering().GetRows(session)

	//////////////// For Updater
	err = xc.NewTwitt_Updater().
		Body("new twitt text").UserId_Eq(1).And_TwiitId_In(1, 2, 3).Update(session)

	//////////////// For Deleter
	err = xc.NewTwitt_Deleter().UserId_Eq(1).And_TwiitId_In(1, 2, 3).Delete(session)
	err = xc.NewTwitt_Deleter().UserId_Eq(1).Delete(session)

	_ = err
	_ = twitts
	_ = twitt
}

/* log output - this is produced CQL queries to cassandra:

2018/09/18 22:35:54 CQL:  [insert into twitter.twitt (body,create_time,twiit_id,user_id) values (?,?,?,?)  [Hello World 1566000000 1 1]]
2018/09/18 22:35:54 CQL:  [DELETE FROM twitter.twitt WHERE  user_id = ? And twiit_id = ?  [1 1]]
2018/09/18 22:35:54 CQL:  [SELECT * FROM twitter.twitt WHERE  user_id = ?  LIMIT 5 [1]]
2018/09/18 22:35:54 CQL:  [SELECT * FROM twitter.twitt WHERE  user_id = ?  LIMIT 5 [1]]
2018/09/18 22:35:54 CQL:  [SELECT * FROM twitter.twitt WHERE  user_id = ? And twiit_id IN (?,?,?)  [1 1 25 68]]
2018/09/18 22:35:54 CQL:  [SELECT user_id, body FROM twitter.twitt WHERE  user_id = ? And twiit_id IN (?,?,?)  LIMIT 12 [1 1 25 68]]
2018/09/18 22:35:54 CQL:  [SELECT * FROM twitter.twitt WHERE  user_id < ?  LIMIT 10  ALLOW FILTERING [100]]
2018/09/18 22:35:54 CQL:  [UPDATE twitter.twitt SET body = ?  WHERE  user_id = ? And twiit_id IN (?,?,?)  [new twitt text 1 1 2 3]]
2018/09/18 22:35:54 CQL:  [DELETE FROM twitter.twitt WHERE  user_id = ? And twiit_id IN (?,?,?)  [1 1 2 3]]
2018/09/18 22:35:54 CQL:  [DELETE FROM twitter.twitt WHERE  user_id = ?  [1]]

*/
