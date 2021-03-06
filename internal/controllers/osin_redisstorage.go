package controllers

// https://github.com/ShaleApps/osinredis

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/bartmika/osin-example/internal/models"
	"github.com/garyburd/redigo/redis"
	"github.com/google/uuid"
	"github.com/openshift/osin"
	"github.com/pkg/errors"
)

func init() {
	gob.Register(map[string]interface{}{})
	gob.Register(&osin.DefaultClient{})
	gob.Register(osin.AuthorizeData{})
	gob.Register(osin.AccessData{})
	gob.Register(models.UserLite{})
}

// OSINRedisStorage implements "github.com/RangelReale/osin".OSINRedisStorage
type OSINRedisStorage struct {
	pool      *redis.Pool
	keyPrefix string
}

func NewOSINRedisStorage() *OSINRedisStorage {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", ":6379")
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
	return &OSINRedisStorage{
		pool:      pool,
		keyPrefix: "prefix",
	}
}

// New initializes and returns a new OSINRedisStorage
func New(pool *redis.Pool, keyPrefix string) *OSINRedisStorage {
	return &OSINRedisStorage{
		pool:      pool,
		keyPrefix: keyPrefix,
	}
}

// Clone the storage if needed. For example, using mgo, you can clone the session with session.Clone
// to avoid concurrent access problems.
// This is to avoid cloning the connection at each method access.
// Can return itself if not a problem.
func (s *OSINRedisStorage) Clone() osin.Storage {
	return s
}

// Close the resources the OSINRedisStorage potentially holds (using Clone for example)
func (s *OSINRedisStorage) Close() {}

// CreateClient inserts a new client
func (s *OSINRedisStorage) CreateClient(client osin.Client) error {
	conn := s.pool.Get()
	if err := conn.Err(); err != nil {
		return err
	}

	defer conn.Close()

	payload, err := encode(client)
	if err != nil {
		return errors.Wrap(err, "failed to encode client")
	}

	_, err = conn.Do("SET", s.makeKey("client", client.GetId()), payload)
	return errors.Wrap(err, "failed to save client")
}

// GetClient gets a client by ID
func (s *OSINRedisStorage) GetClient(id string) (osin.Client, error) {
	conn := s.pool.Get()
	if err := conn.Err(); err != nil {
		return nil, err
	}

	defer conn.Close()

	var (
		rawClientGob interface{}
		err          error
	)

	if rawClientGob, err = conn.Do("GET", s.makeKey("client", id)); err != nil {
		return nil, errors.Wrap(err, "unable to GET client")
	}
	if rawClientGob == nil {
		return nil, nil
	}

	clientGob, _ := redis.Bytes(rawClientGob, err)

	var client osin.DefaultClient
	err = decode(clientGob, &client)
	return &client, errors.Wrap(err, "failed to decode client gob")
}

// UpdateClient updates a client
func (s *OSINRedisStorage) UpdateClient(client osin.Client) error {
	return errors.Wrap(s.CreateClient(client), "failed to update client")
}

// DeleteClient deletes given client
func (s *OSINRedisStorage) DeleteClient(client osin.Client) error {
	conn := s.pool.Get()
	if err := conn.Err(); err != nil {
		return err
	}

	defer conn.Close()

	_, err := conn.Do("DEL", s.makeKey("client", client.GetId()))
	return errors.Wrap(err, "failed to delete client")
}

// SaveAuthorize saves authorize data.
func (s *OSINRedisStorage) SaveAuthorize(data *osin.AuthorizeData) (err error) {
	conn := s.pool.Get()
	if err := conn.Err(); err != nil {
		return err
	}

	defer conn.Close()

	payload, err := encode(data)
	if err != nil {
		return errors.Wrap(err, "failed to encode data")
	}

	_, err = conn.Do("SETEX", s.makeKey("auth", data.Code), data.ExpiresIn, string(payload))
	return errors.Wrap(err, "failed to set auth")
}

// LoadAuthorize looks up AuthorizeData by a code.
// Client information MUST be loaded together.
// Optionally can return error if expired.
func (s *OSINRedisStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	conn := s.pool.Get()
	if err := conn.Err(); err != nil {
		return nil, err
	}

	defer conn.Close()

	var (
		rawAuthGob interface{}
		err        error
	)

	if rawAuthGob, err = conn.Do("GET", s.makeKey("auth", code)); err != nil {
		return nil, errors.Wrap(err, "unable to GET auth")
	}
	if rawAuthGob == nil {
		return nil, nil
	}

	authGob, _ := redis.Bytes(rawAuthGob, err)

	var auth osin.AuthorizeData
	err = decode(authGob, &auth)
	return &auth, errors.Wrap(err, "failed to decode auth")
}

// RemoveAuthorize revokes or deletes the authorization code.
func (s *OSINRedisStorage) RemoveAuthorize(code string) (err error) {
	conn := s.pool.Get()
	if err := conn.Err(); err != nil {
		return err
	}

	defer conn.Close()

	_, err = conn.Do("DEL", s.makeKey("auth", code))
	return errors.Wrap(err, "failed to delete auth")
}

// SaveAccess creates AccessData.
func (s *OSINRedisStorage) SaveAccess(data *osin.AccessData) (err error) {
	conn := s.pool.Get()
	if err := conn.Err(); err != nil {
		return err
	}

	defer conn.Close()

	payload, err := encode(data)
	if err != nil {
		return errors.Wrap(err, "failed to encode access")
	}

	accessID := uuid.NewString()

	if _, err := conn.Do("SETEX", s.makeKey("access", accessID), data.ExpiresIn, string(payload)); err != nil {
		return errors.Wrap(err, "failed to save access")
	}

	if _, err := conn.Do("SETEX", s.makeKey("access_token", data.AccessToken), data.ExpiresIn, accessID); err != nil {
		return errors.Wrap(err, "failed to register access token")
	}

	_, err = conn.Do("SETEX", s.makeKey("refresh_token", data.RefreshToken), data.ExpiresIn, accessID)
	return errors.Wrap(err, "failed to register refresh token")
}

// LoadAccess gets access data with given access token
func (s *OSINRedisStorage) LoadAccess(token string) (*osin.AccessData, error) {
	return s.loadAccessByKey(s.makeKey("access_token", token))
}

// RemoveAccess deletes AccessData with given access token
func (s *OSINRedisStorage) RemoveAccess(token string) error {
	return s.removeAccessByKey(s.makeKey("access_token", token))
}

// LoadRefresh gets access data with given refresh token
func (s *OSINRedisStorage) LoadRefresh(token string) (*osin.AccessData, error) {
	return s.loadAccessByKey(s.makeKey("refresh_token", token))
}

// RemoveRefresh deletes AccessData with given refresh token
func (s *OSINRedisStorage) RemoveRefresh(token string) error {
	return s.removeAccessByKey(s.makeKey("refresh_token", token))
}

func (s *OSINRedisStorage) removeAccessByKey(key string) error {
	conn := s.pool.Get()
	if err := conn.Err(); err != nil {
		return err
	}

	defer conn.Close()

	accessID, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return errors.Wrap(err, "failed to get access")
	}

	access, err := s.loadAccessByKey(key)
	if err != nil {
		return errors.Wrap(err, "unable to load access for removal")
	}

	if access == nil {
		return nil
	}

	accessKey := s.makeKey("access", accessID)
	if _, err := conn.Do("DEL", accessKey); err != nil {
		return errors.Wrap(err, "failed to delete access")
	}

	accessTokenKey := s.makeKey("access_token", access.AccessToken)
	if _, err := conn.Do("DEL", accessTokenKey); err != nil {
		return errors.Wrap(err, "failed to deregister access_token")
	}

	refreshTokenKey := s.makeKey("refresh_token", access.RefreshToken)
	_, err = conn.Do("DEL", refreshTokenKey)
	return errors.Wrap(err, "failed to deregister refresh_token")
}

func (s *OSINRedisStorage) loadAccessByKey(key string) (*osin.AccessData, error) {
	conn := s.pool.Get()
	if err := conn.Err(); err != nil {
		return nil, err
	}

	defer conn.Close()

	var (
		rawAuthGob interface{}
		err        error
	)

	if rawAuthGob, err = conn.Do("GET", key); err != nil {
		return nil, errors.Wrap(err, "unable to GET auth")
	}
	if rawAuthGob == nil {
		return nil, nil
	}

	accessID, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return nil, errors.Wrap(err, "unable to get access ID")
	}

	accessIDKey := s.makeKey("access", accessID)
	accessGob, err := redis.Bytes(conn.Do("GET", accessIDKey))
	if err != nil {
		return nil, errors.Wrap(err, "unable to get access gob")
	}

	var access osin.AccessData
	if err := decode(accessGob, &access); err != nil {
		return nil, errors.Wrap(err, "failed to decode access gob")
	}

	ttl, err := redis.Int(conn.Do("TTL", accessIDKey))
	if err != nil {
		return nil, errors.Wrap(err, "unable to get access TTL")
	}

	access.ExpiresIn = int32(ttl)

	access.Client, err = s.GetClient(access.Client.GetId())
	if err != nil {
		return nil, errors.Wrap(err, "unable to get client for access")
	}

	if access.AuthorizeData != nil && access.AuthorizeData.Client != nil {
		access.AuthorizeData.Client, err = s.GetClient(access.AuthorizeData.Client.GetId())
		if err != nil {
			return nil, errors.Wrap(err, "unable to get client for access authorize data")
		}
	}

	return &access, nil
}

func (s *OSINRedisStorage) makeKey(namespace, id string) string {
	return fmt.Sprintf("%s:%s:%s", s.keyPrefix, namespace, id)
}

func encode(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(v); err != nil {
		return nil, errors.Wrap(err, "unable to encode")
	}
	return buf.Bytes(), nil
}

func decode(data []byte, v interface{}) error {
	err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(v)
	return errors.Wrap(err, "unable to decode")
}
