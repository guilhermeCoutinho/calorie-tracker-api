package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/go-pg/pg"
	"github.com/guilhermeCoutinho/api-studies/messages"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/topfreegames/extensions/pg/interfaces"
)

func getAuthenticatedUser(t *testing.T) string {
	userName := "userName"
	password := "MyPassword"

	createUserRequest := &messages.CreateUserRequest{
		Username:     userName,
		Password:     password,
		CalorieLimit: 100,
	}

	doRequest(t, http.MethodPost, "/users", nil, createUserRequest, &messages.BaseResponse{})

	loginResponse := &messages.LoginResponse{}
	doRequest(t, http.MethodPost, "/auth", nil, &messages.LoginRequest{
		Username: userName,
		Password: password,
	}, loginResponse)

	assert.NotEmpty(t, loginResponse.AccessToken)

	return loginResponse.AccessToken
}

func doRequest(t *testing.T, method string, path string, token *string, payload interface{}, responsePtr interface{}) int {
	payloadBytes, err := json.Marshal(payload)
	assert.Nil(t, err)

	req, err := http.NewRequest(method, URL+path, bytes.NewBuffer(payloadBytes))
	assert.Nil(t, err)

	if token != nil {
		req.Header.Set("authorization", *token)
	}

	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)

	defer resp.Body.Close()

	bufferSize := int64(1024 * 1024)
	response, err := ioutil.ReadAll(io.LimitReader(resp.Body, bufferSize))
	assert.Nil(t, err)

	err = json.Unmarshal(response, responsePtr)
	assert.Nil(t, err)

	return resp.StatusCode
}

// GetConfigFromPath returns the config from the path
func GetConfigFromPath(t testing.TB, path string) *viper.Viper {
	t.Helper()

	config, err := getConfig(path, "calorie-tracker", "yaml")
	if err != nil {
		t.Fatal(err)
	}

	return config
}

func GetPG(t testing.TB) func() {
	config := GetConfigFromPath(t, "../config/config.yaml")
	db, _ := connectToPG(config, "db")

	return func() {
		TruncateTables(t, db)
		db.Close()
	}
}

func connectToPG(config *viper.Viper, prefix string) (*pg.DB, error) {
	user := config.GetString(fmt.Sprintf("%s.user", prefix))
	pass := config.GetString(fmt.Sprintf("%s.pass", prefix))
	host := config.GetString(fmt.Sprintf("%s.host", prefix))
	database := config.GetString(fmt.Sprintf("%s.database", prefix))
	port := config.GetInt(fmt.Sprintf("%s.port", prefix))
	poolSize := config.GetInt(fmt.Sprintf("%s.poolSize", prefix))
	maxRetries := config.GetInt(fmt.Sprintf("%s.maxRetries", prefix))
	timeout := config.GetInt(fmt.Sprintf("%s.connectionTimeout", prefix))

	options := &pg.Options{
		Addr:       fmt.Sprintf("%s:%d", host, port),
		User:       user,
		Password:   pass,
		Database:   database,
		PoolSize:   poolSize,
		MaxRetries: maxRetries,
	}
	db := pg.Connect(options)
	err := waitForConnection(db, timeout)
	return db, err
}

func waitForConnection(db *pg.DB, timeout int) error {
	t := time.Duration(timeout) * time.Second
	timeoutTimer := time.NewTimer(t)
	defer timeoutTimer.Stop()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutTimer.C:
			return fmt.Errorf("timed out waiting for PostgreSQL to connect")
		case <-ticker.C:
			if isConnected(db) {
				return nil
			}
		}
	}
}

func isConnected(db *pg.DB) bool {
	res, err := db.Exec("select 1")
	if err != nil {
		return false
	}
	return res.RowsReturned() == 1
}

func TruncateTables(t testing.TB, db interfaces.DB) {
	t.Helper()

	_, err := db.Exec(`
	CREATE OR REPLACE FUNCTION truncate_tables(username IN VARCHAR) RETURNS void AS $$
	DECLARE
		statements CURSOR FOR
			SELECT tablename FROM pg_tables
			WHERE tableowner = username AND schemaname = 'public' AND tablename != 'gopg_migrations';
	BEGIN
		FOR stmt IN statements LOOP
			EXECUTE 'TRUNCATE TABLE ' || quote_ident(stmt.tablename) || ' CASCADE;';
		END LOOP;
	END;
	$$ LANGUAGE plpgsql;
	SELECT truncate_tables('calorie-tracker-user');
	`)

	if err != nil {
		t.Fatal(err)
	}
}

func getConfig(path, envPrefix, configType string) (*viper.Viper, error) {
	config := viper.New()
	config.SetConfigFile(path)
	config.SetConfigType(configType)
	config.SetEnvPrefix(envPrefix)
	config.AddConfigPath(".")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()

	// If a config file is found, read it in.
	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}

	return config, nil
}
