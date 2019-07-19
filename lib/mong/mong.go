package mong

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

type Mgo struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	cfg        Config
	mux        sync.Mutex
}

func (m *Mgo) newCtx() (context.Context, context.CancelFunc) {
	if m.cfg.Timeout == 0 {
		return context.WithCancel(context.Background())
	}
	return context.WithTimeout(context.Background(), m.cfg.Timeout)
}

func (m *Mgo) new() error {
	option := options.Client()
	option.ApplyURI("mongodb://" + m.cfg.Host)
	option.SetMaxPoolSize(10)

	client, err := mongo.Connect(context.Background(), option)
	if err != nil {
		return err
	}

	m.client = client
	m.database = m.client.Database(m.cfg.Database)
	m.collection = m.database.Collection(m.cfg.Collection)
	return nil
}

// 连接
func (m *Mgo) Connect() error {

	m.mux.Lock()
	defer m.mux.Unlock()

	if m.client == nil {
		return m.new()
	}
	if err := m.client.Ping(context.Background(), nil); err != nil {
		return m.new()
	}
	return nil
}

// 单条写入
func (m *Mgo) InsertOne(document interface{}) (*mongo.InsertOneResult, error) {
	ctx, cancel := m.newCtx()
	defer cancel()
	return m.collection.InsertOne(ctx, document)
}

// 批量写入
func (m *Mgo) InsertMany(documents []interface{}) (*mongo.InsertManyResult, error) {
	ctx, cancel := m.newCtx()
	defer cancel()
	return m.collection.InsertMany(ctx, documents)
}

// 批量查询
func (m *Mgo) Find(filter interface{}) (*mongo.Cursor, error) {
	ctx, cancel := m.newCtx()
	defer cancel()
	return m.collection.Find(ctx, filter)
}

// 断开连接
func (m *Mgo) Disconnect() error {

	m.mux.Lock()
	defer m.mux.Unlock()

	if m.client == nil {
		return nil
	}
	return m.client.Disconnect(context.Background())
}

