package test

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"sumwhere/models"
	"testing"
	"time"
)

func client() *redis.Client {
	opt, err := redis.ParseURL("redis://:@1.215.236.26:53379")
	if err != nil {
		panic(err)
	}

	// Create client as usually.
	client := redis.NewClient(opt)
	return client
}
func TestRedisTopTripPlace(t *testing.T) {
	// Create client as usually.
	//client := client()

	//assert.Equal(t,"",client.HGet("trip1","user2").Val())
	//result, err := client.Get("helloworld").Result()
	//assert.NoError(t,err)
	//
	//t.Log(result)

	//assert.Equal(t,"",client.Get("trip:57").Val())

	//assert.Equal(t,int64(0),client.Exists("trip:57").Val())

	//
	//t.Log(client.Get("trip:57").Val())
	//
	//client.Set("trip:57",nil,0)

}

func TestHset(t *testing.T) {
	client := client()

	trip := models.Trip{
		Id:          0,
		UserId:      0,
		MatchTypeId: 0,
		Concept:     "",
		TripTypeId:  0,
		GenderType:  "",
		StartDate:   time.Time{},
		EndDate:     time.Time{},
		CreateAt:    time.Time{},
		UpdateAt:    time.Time{},
		DeleteAt:    time.Time{},
	}
	_, err := json.Marshal(&trip)
	assert.NoError(t, err)

	//assert.Equal(t,true,client.HSet("user:8","res",byte).Val())

	for k, v := range client.HGetAll("user:8").Val() {
		t.Log(k)
		t.Log(v)
	}
	//t.Log(client.HGetAll("user:8").Val())
	//t.Log(client.HMGet("8","57"))

}
