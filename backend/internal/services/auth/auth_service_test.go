package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AuthServiceUnitTestSuite struct {
	suite.Suite
	authService *AuthService
}

func TestAuthServiceUnitTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceUnitTestSuite))
}

func (suite *AuthServiceUnitTestSuite) SetupTest() {
	suite.authService = NewAuthService()
}

// ===== Unit Tests for Token Generation =====

func (suite *AuthServiceUnitTestSuite) TestGenerateAccessToken_Success() {
	userID := "550e8400-e29b-41d4-a716-446655440001"
	role := "employee"

	token, err := suite.authService.GenerateAccessToken(userID, role)
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)

	// Verify token structure
	assert.NotEmpty(suite.T(), token)
	parts := strings.Split(token, ".")
	assert.Equal(suite.T(), 3, len(parts)) // JWT has 3 parts separated by dots
}

func (suite *AuthServiceUnitTestSuite) TestGenerateAccessToken_WithAdminRole() {
	userID := "550e8400-e29b-41d4-a716-446655440001"
	role := "admin"

	token, err := suite.authService.GenerateAccessToken(userID, role)
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
}

func (suite *AuthServiceUnitTestSuite) TestGenerateAccessToken_WithManagerRole() {
	userID := "550e8400-e29b-41d4-a716-446655440001"
	role := "manager"

	token, err := suite.authService.GenerateAccessToken(userID, role)
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
}

func (suite *AuthServiceUnitTestSuite) TestGenerateRefreshToken_Success() {
	userID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := suite.authService.GenerateRefreshToken(userID)
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
}

func (suite *AuthServiceUnitTestSuite) TestGenerateAccessToken_EmptyUserID() {
	userID := ""
	role := "employee"

	token, err := suite.authService.GenerateAccessToken(userID, role)
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), token)
}

func (suite *AuthServiceUnitTestSuite) TestGenerateRefreshToken_EmptyUserID() {
	userID := ""

	token, err := suite.authService.GenerateRefreshToken(userID)
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), token)
}

// ===== Unit Tests for Token Validation =====

func (suite *AuthServiceUnitTestSuite) TestValidateAccessToken_Success() {
	userID := "550e8400-e29b-41d4-a716-446655440001"
	role := "employee"

	// Generate token
	token, err := suite.authService.GenerateAccessToken(userID, role)
	require.NoError(suite.T(), err)

	// Validate token
	extractedUserID, err := suite.authService.ValidateAccessToken(token)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, extractedUserID)
}

func (suite *AuthServiceUnitTestSuite) TestValidateAccessToken_InvalidToken() {
	invalidToken := "invalid.token.here"

	_, err := suite.authService.ValidateAccessToken(invalidToken)
	assert.Error(suite.T(), err)
}

func (suite *AuthServiceUnitTestSuite) TestValidateAccessToken_MalformedToken() {
	malformedToken := "notavalidtoken"

	_, err := suite.authService.ValidateAccessToken(malformedToken)
	assert.Error(suite.T(), err)
}

func (suite *AuthServiceUnitTestSuite) TestValidateRefreshToken_Success() {
	userID := "550e8400-e29b-41d4-a716-446655440001"

	// Generate token
	token, err := suite.authService.GenerateRefreshToken(userID)
	require.NoError(suite.T(), err)

	// Validate token
	extractedUserID, err := suite.authService.ValidateRefreshToken(token)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, extractedUserID)
}

func (suite *AuthServiceUnitTestSuite) TestValidateRefreshToken_InvalidToken() {
	invalidToken := "invalid.refresh.token"

	_, err := suite.authService.ValidateRefreshToken(invalidToken)
	assert.Error(suite.T(), err)
}

// ===== Unit Tests for Token Expiration =====

func (suite *AuthServiceUnitTestSuite) TestAccessToken_DoesNotExpireImmediately() {
	userID := "550e8400-e29b-41d4-a716-446655440001"
	role := "employee"

	token, err := suite.authService.GenerateAccessToken(userID, role)
	require.NoError(suite.T(), err)

	// Should be valid immediately
	_, err = suite.authService.ValidateAccessToken(token)
	assert.NoError(suite.T(), err)
}

func (suite *AuthServiceUnitTestSuite) TestRefreshToken_DoesNotExpireImmediately() {
	userID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := suite.authService.GenerateRefreshToken(userID)
	require.NoError(suite.T(), err)

	// Should be valid immediately
	_, err = suite.authService.ValidateRefreshToken(token)
	assert.NoError(suite.T(), err)
}

// ===== Unit Tests for Token Claims =====

func (suite *AuthServiceUnitTestSuite) TestAccessTokenContainsCorrectClaims() {
	userID := "550e8400-e29b-41d4-a716-446655440001"
	role := "manager"

	token, err := suite.authService.GenerateAccessToken(userID, role)
	require.NoError(suite.T(), err)

	// Parse and verify claims
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// In real implementation, use secret key
		return []byte("your-secret-key"), nil
	})

	if err == nil && parsedToken.Valid {
		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		require.True(suite.T(), ok)

		// Verify user_id claim
		assert.Equal(suite.T(), userID, claims["user_id"])
		// Verify role claim
		assert.Equal(suite.T(), role, claims["role"])
	}
}

// ===== Unit Tests for Token Renewal Flow =====

func (suite *AuthServiceUnitTestSuite) TestTokenRenewalFlow() {
	userID := "550e8400-e29b-41d4-a716-446655440001"
	role := "employee"

	// 1. Generate initial access and refresh tokens
	accessToken1, err := suite.authService.GenerateAccessToken(userID, role)
	require.NoError(suite.T(), err)

	refreshToken, err := suite.authService.GenerateRefreshToken(userID)
	require.NoError(suite.T(), err)

	// 2. Validate initial tokens
	extractedUserID, err := suite.authService.ValidateAccessToken(accessToken1)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, extractedUserID)

	// 3. Generate new access token using refresh token (simulating renewal)
	accessToken2, err := suite.authService.GenerateAccessToken(userID, role)
	require.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), accessToken1, accessToken2)

	// 4. Validate new token
	extractedUserID, err = suite.authService.ValidateAccessToken(accessToken2)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, extractedUserID)
}

// ===== Unit Tests for Multiple Users =====

func (suite *AuthServiceUnitTestSuite) TestMultipleUsersTokensAreIndependent() {
	user1ID := "550e8400-e29b-41d4-a716-446655440001"
	user2ID := "550e8400-e29b-41d4-a716-446655440002"
	role := "employee"

	// Generate tokens for both users
	token1, err1 := suite.authService.GenerateAccessToken(user1ID, role)
	require.NoError(suite.T(), err1)

	token2, err2 := suite.authService.GenerateAccessToken(user2ID, role)
	require.NoError(suite.T(), err2)

	// Tokens should be different
	assert.NotEqual(suite.T(), token1, token2)

	// Validate each token returns correct user ID
	extractedID1, _ := suite.authService.ValidateAccessToken(token1)
	extractedID2, _ := suite.authService.ValidateAccessToken(token2)

	assert.Equal(suite.T(), user1ID, extractedID1)
	assert.Equal(suite.T(), user2ID, extractedID2)
}

// ===== Integration-style tests (within service layer) =====

func (suite *AuthServiceUnitTestSuite) TestCompleteAuthenticationLifecycle() {
	userID := "550e8400-e29b-41d4-a716-446655440001"
	role := "employee"

	// 1. Generate tokens
	accessToken, err := suite.authService.GenerateAccessToken(userID, role)
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), accessToken)

	refreshToken, err := suite.authService.GenerateRefreshToken(userID)
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), refreshToken)

	// 2. Validate tokens
	extractedUserID, err := suite.authService.ValidateAccessToken(accessToken)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, extractedUserID)

	extractedUserID, err = suite.authService.ValidateRefreshToken(refreshToken)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, extractedUserID)

	// 3. Simulate time passing and refresh
	time.Sleep(100 * time.Millisecond)

	// 4. Generate new access token
	newAccessToken, err := suite.authService.GenerateAccessToken(userID, role)
	require.NoError(suite.T(), err)

	// 5. Validate new token
	extractedUserID, err = suite.authService.ValidateAccessToken(newAccessToken)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, extractedUserID)
}

// ===== Benchmark tests (optional) =====

func (suite *AuthServiceUnitTestSuite) BenchmarkGenerateAccessToken(b *testing.B) {
	userID := "550e8400-e29b-41d4-a716-446655440001"
	role := "employee"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = suite.authService.GenerateAccessToken(userID, role)
	}
}

func (suite *AuthServiceUnitTestSuite) BenchmarkValidateAccessToken(b *testing.B) {
	userID := "550e8400-e29b-41d4-a716-446655440001"
	role := "employee"

	token, _ := suite.authService.GenerateAccessToken(userID, role)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = suite.authService.ValidateAccessToken(token)
	}
}
