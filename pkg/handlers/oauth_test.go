/****************************************************************************
 * Copyright 2020, Optimizely, Inc. and contributors                        *
 *                                                                          *
 * Licensed under the Apache License, Version 2.0 (the "License");          *
 * you may not use this file except in compliance with the License.         *
 * You may obtain a copy of the License at                                  *
 *                                                                          *
 *    http://www.apache.org/licenses/LICENSE-2.0                            *
 *                                                                          *
 * Unless required by applicable law or agreed to in writing, software      *
 * distributed under the License is distributed on an "AS IS" BASIS,        *
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. *
 * See the License for the specific language governing permissions and      *
 * limitations under the License.                                           *
 ***************************************************************************/

// Package handlers //
package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/optimizely/agent/config"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

type OAuthTestSuite struct {
	suite.Suite
	handler *OAuthHandler
	mux     *chi.Mux
	secret  string
}

func (s *OAuthTestSuite) SetupTest() {
	s.secret = "RW+Uo/7z4ag9hAb10w8LIZFRFaSwS4nt1/l+uVgChIQ="
	authConfig := config.ServiceAuthConfig{
		Clients: []config.OAuthClientCredentials{
			{
				ID:         "optly_user",
				SecretHash: "JDJhJDEyJDNDOG12LmNCNzlHaHhGcEJtLzZZQk9VLnRneEpGTTlnTXozb2kyNS9ERzhJTDZOZkpGa0ND",
				SDKKeys:    []string{"123"},
			},
		},
		HMACSecrets: []string{"hmac_seekrit"},
		TTL:         30 * time.Minute,
	}
	s.handler = NewOAuthHandler(&authConfig)

	mux := chi.NewMux()
	mux.Post("/api/token", s.handler.CreateAPIAccessToken)
	mux.Post("/admin/token", s.handler.CreateAdminAccessToken)
	s.mux = mux
}

func (s *OAuthTestSuite) TestGetAPIAccessTokenMissingGrantType() {
	bodyValues := url.Values{}
	bodyValues.Set("client_id", "optly")
	bodyValues.Set("client_secret", s.secret)
	req := httptest.NewRequest("POST", "/api/token", strings.NewReader(bodyValues.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	s.mux.ServeHTTP(rec, req)
	s.Equal(http.StatusBadRequest, rec.Code)
}

func (s *OAuthTestSuite) TestGetAdminAccessTokenMissingGrantType() {
	bodyValues := url.Values{}
	bodyValues.Set("client_id", "optly")
	bodyValues.Set("client_secret", s.secret)
	req := httptest.NewRequest("POST", "/admin/token", strings.NewReader(bodyValues.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	s.mux.ServeHTTP(rec, req)
	s.Equal(http.StatusBadRequest, rec.Code)
}

func (s *OAuthTestSuite) TestGetAPIAccessTokenUnsupportedGrantType() {
	bodyValues := url.Values{}
	bodyValues.Set("grant_type", "authorization")
	bodyValues.Set("client_id", "optly_user")
	bodyValues.Set("client_secret", s.secret)
	req := httptest.NewRequest("POST", "/api/token", strings.NewReader(bodyValues.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	s.mux.ServeHTTP(rec, req)
	s.Equal(http.StatusBadRequest, rec.Code)
}

func (s *OAuthTestSuite) TestGetAdminAccessTokenUnsupportedGrantType() {
	bodyValues := url.Values{}
	bodyValues.Set("grant_type", "authorization")
	bodyValues.Set("client_id", "optly_user")
	bodyValues.Set("client_secret", s.secret)
	req := httptest.NewRequest("POST", "/admin/token", strings.NewReader(bodyValues.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	s.mux.ServeHTTP(rec, req)
	s.Equal(http.StatusBadRequest, rec.Code)
}

func (s *OAuthTestSuite) TestGetAPIAccessTokenMissingClientId() {
	bodyValues := url.Values{}
	bodyValues.Set("grant_type", "client_credentials")
	bodyValues.Set("client_secret", s.secret)
	req := httptest.NewRequest("POST", "/api/token", strings.NewReader(bodyValues.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	s.mux.ServeHTTP(rec, req)
	s.Equal(http.StatusUnauthorized, rec.Code)
}

func (s *OAuthTestSuite) TestGetAdminAccessTokenMissingClientId() {
	bodyValues := url.Values{}
	bodyValues.Set("grant_type", "client_credentials")
	bodyValues.Set("client_secret", s.secret)
	req := httptest.NewRequest("POST", "/admin/token", strings.NewReader(bodyValues.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	s.mux.ServeHTTP(rec, req)
	s.Equal(http.StatusUnauthorized, rec.Code)
}

func (s *OAuthTestSuite) TestGetAPIAccessTokenMissingClientSecret() {
	bodyValues := url.Values{}
	bodyValues.Set("grant_type", "client_credentials")
	bodyValues.Set("client_id", "optly_user")
	req := httptest.NewRequest("POST", "/api/token", strings.NewReader(bodyValues.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	s.mux.ServeHTTP(rec, req)
	s.Equal(http.StatusUnauthorized, rec.Code)
}

func (s *OAuthTestSuite) TestGetAdminAccessTokenMissingClientSecret() {
	bodyValues := url.Values{}
	bodyValues.Set("grant_type", "client_credentials")
	bodyValues.Set("client_id", "optly_user")
	req := httptest.NewRequest("POST", "/admin/token", strings.NewReader(bodyValues.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	s.mux.ServeHTTP(rec, req)
	s.Equal(http.StatusUnauthorized, rec.Code)
}

func (s *OAuthTestSuite) TestGetAPIAccessTokenInvalidClientId() {
	bodyValues := url.Values{}
	bodyValues.Set("grant_type", "client_credentials")
	bodyValues.Set("client_id", "not_an_optly_user")
	bodyValues.Set("client_secret", s.secret)
	req := httptest.NewRequest("POST", "/api/token", strings.NewReader(bodyValues.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	s.mux.ServeHTTP(rec, req)
	s.Equal(http.StatusUnauthorized, rec.Code)
}

func (s *OAuthTestSuite) TestGetAdminAccessTokenInvalidClientId() {
	bodyValues := url.Values{}
	bodyValues.Set("grant_type", "client_credentials")
	bodyValues.Set("client_id", "not_an_optly_user")
	bodyValues.Set("client_secret", s.secret)
	req := httptest.NewRequest("POST", "/admin/token", strings.NewReader(bodyValues.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	s.mux.ServeHTTP(rec, req)
	s.Equal(http.StatusUnauthorized, rec.Code)
}

func (s *OAuthTestSuite) TestGetAPIAccessTokenInvalidClientSecret() {
	bodyValues := url.Values{}
	bodyValues.Set("grant_type", "client_credentials")
	bodyValues.Set("client_id", "optly_user")
	bodyValues.Set("client_secret", "GpDgQx7w8Hb3ibD6K+T77S/0kHgr9qoRxsEpC2lDv08=")
	req := httptest.NewRequest("POST", "/api/token", strings.NewReader(bodyValues.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	s.mux.ServeHTTP(rec, req)
	s.Equal(http.StatusUnauthorized, rec.Code)
}

func (s *OAuthTestSuite) TestGetAdminAccessTokenInvalidClientSecret() {
	bodyValues := url.Values{}
	bodyValues.Set("grant_type", "client_credentials")
	bodyValues.Set("client_id", "optly_user")
	bodyValues.Set("client_secret", "GpDgQx7w8Hb3ibD6K+T77S/0kHgr9qoRxsEpC2lDv08=")
	req := httptest.NewRequest("POST", "/admin/token", strings.NewReader(bodyValues.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	s.mux.ServeHTTP(rec, req)
	s.Equal(http.StatusUnauthorized, rec.Code)
}

func (s *OAuthTestSuite) TestGetAPIAccessTokenInvalidBody() {
	req := httptest.NewRequest("POST", "/api/token", strings.NewReader("fjklM<>CXM><&&^&%"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	s.mux.ServeHTTP(rec, req)
	s.Equal(http.StatusBadRequest, rec.Code)
}

func (s *OAuthTestSuite) TestGetAPIAccessTokenSuccess() {
	bodyValues := url.Values{}
	bodyValues.Set("grant_type", "client_credentials")
	bodyValues.Set("client_id", "optly_user")
	bodyValues.Set("client_secret", s.secret)
	req := httptest.NewRequest("POST", "/api/token", strings.NewReader(bodyValues.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	s.mux.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var actual tokenResponse
	err := json.Unmarshal(rec.Body.Bytes(), &actual)
	s.NoError(err)
	s.Equal("bearer", actual.TokenType)
	s.NotEmpty(actual.AccessToken)
	s.NotEmpty(actual.ExpiresIn)
}

func (s *OAuthTestSuite) TestGetAdminAccessTokenSuccess() {
	bodyValues := url.Values{}
	bodyValues.Set("grant_type", "client_credentials")
	bodyValues.Set("client_id", "optly_user")
	bodyValues.Set("client_secret", s.secret)
	req := httptest.NewRequest("POST", "/admin/token", strings.NewReader(bodyValues.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	s.mux.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var actual tokenResponse
	err := json.Unmarshal(rec.Body.Bytes(), &actual)
	s.NoError(err)
	s.Equal("bearer", actual.TokenType)
	s.NotEmpty(actual.AccessToken)
	s.NotEmpty(actual.ExpiresIn)
}

func TestOAuthTestSuite(t *testing.T) {
	suite.Run(t, new(OAuthTestSuite))
}

type OAuthDisabledTestSuite struct {
	suite.Suite
	handler *OAuthHandler
	mux     *chi.Mux
}

func (s *OAuthDisabledTestSuite) SetupTest() {
	authConfig := config.ServiceAuthConfig{
		Clients:     make([]config.OAuthClientCredentials, 0),
		HMACSecrets: make([]string, 0),
		TTL:         0,
	}
	s.handler = NewOAuthHandler(&authConfig)

	mux := chi.NewMux()
	mux.Post("/api/token", s.handler.CreateAPIAccessToken)
	mux.Post("/admin/token", s.handler.CreateAdminAccessToken)
	s.mux = mux
}

func (s *OAuthDisabledTestSuite) TestGetAdminAccessTokenDisabled() {
	bodyValues := url.Values{}
	bodyValues.Set("grant_type", "client_credentials")
	bodyValues.Set("client_id", "optly_user")
	bodyValues.Set("client_secret", "client_seekrit")
	req := httptest.NewRequest("POST", "/admin/token", strings.NewReader(bodyValues.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	s.mux.ServeHTTP(rec, req)
	s.Equal(http.StatusUnauthorized, rec.Code)
}

func (s *OAuthDisabledTestSuite) TestGetAPIAccessTokenDisabled() {
	bodyValues := url.Values{}
	bodyValues.Set("grant_type", "client_credentials")
	bodyValues.Set("client_id", "optly_user")
	bodyValues.Set("client_secret", "client_seekrit")
	req := httptest.NewRequest("POST", "/api/token", strings.NewReader(bodyValues.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	s.mux.ServeHTTP(rec, req)
	s.Equal(http.StatusUnauthorized, rec.Code)
}

func TestOAuthDisabledTestSuite(t *testing.T) {
	suite.Run(t, new(OAuthDisabledTestSuite))
}

type OAuthMissingHMACSecretTestSuite struct {
	suite.Suite
	handler *OAuthHandler
	mux     *chi.Mux
	secret  string
}

func (s *OAuthMissingHMACSecretTestSuite) SetupTest() {
	s.secret = "RW+Uo/7z4ag9hAb10w8LIZFRFaSwS4nt1/l+uVgChIQ="
	config := config.ServiceAuthConfig{
		Clients: []config.OAuthClientCredentials{
			{
				ID:         "optly_user",
				SecretHash: "JDJhJDEyJDNDOG12LmNCNzlHaHhGcEJtLzZZQk9VLnRneEpGTTlnTXozb2kyNS9ERzhJTDZOZkpGa0ND",
				SDKKeys:    []string{"123"},
			},
		},
		// No HMACSecrets provided - this configuration is invalid
		HMACSecrets: []string{},
		TTL:         30 * time.Minute,
	}
	s.handler = NewOAuthHandler(&config)

	mux := chi.NewMux()
	mux.Post("/api/token", s.handler.CreateAPIAccessToken)
	mux.Post("/admin/token", s.handler.CreateAdminAccessToken)
	s.mux = mux
}

func (s *OAuthMissingHMACSecretTestSuite) TestInvalidConfig() {
	s.Nil(s.handler)
}

func TestOAuthMissingHMACSecretTestSuite(t *testing.T) {
	suite.Run(t, new(OAuthMissingHMACSecretTestSuite))
}
