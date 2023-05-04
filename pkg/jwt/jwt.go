package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dozheiny/it-captal-task/config"
	user2 "github.com/dozheiny/it-captal-task/models/user"
	"github.com/dozheiny/it-captal-task/pkg/mongo"
	"github.com/dozheiny/it-captal-task/serializers"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewToken just initialize TokenDetails.
func NewToken() *TokenDetails {
	t := new(TokenDetails)
	return t
}

// Initialize gets user id and generate refresh token and access token.
func (td *TokenDetails) Initialize(userID primitive.ObjectID, refreshBuild bool) error {
	var err error

	// get SecretKey.
	secretKey, err := config.Get("SECRET_KEY")
	if err != nil {
		return err
	}

	// Build access token structure
	td.AtExpires = time.Now().Add(accessTokenLifeTime).Unix()
	td.AccessUuid = uuid.New().String()

	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["authorized"] = true
	accessTokenClaims["access_uuid"] = td.AccessUuid
	accessTokenClaims["userId"] = userID
	accessTokenClaims["exp"] = td.AtExpires

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	td.AccessToken, err = accessToken.SignedString([]byte(secretKey))
	if err != nil {
		return err
	}

	// Build refresh token structure
	if refreshBuild {
		td.RtExpires = time.Now().Add(accessTokenLifeTime).Unix()
		td.RefreshUuid = uuid.New().String()
		refreshTokenClaims := jwt.MapClaims{}
		refreshTokenClaims["refresh_uuid"] = td.RefreshUuid
		refreshTokenClaims["userId"] = userID
		refreshTokenClaims["exp"] = td.RtExpires
		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
		td.RefreshToken, err = refreshToken.SignedString([]byte(secretKey))
		if err != nil {
			return err
		}
	}

	return nil
}

// ExtractAccessMetadata Extract access token claims from request header
func ExtractAccessMetadata(req *fiber.Ctx) (jwt.MapClaims, error) {
	// Extract token string from header.
	bearToken := req.GetReqHeaders()["Authorization"]
	strArr := strings.Split(bearToken, " ")
	var tokenString string
	if len(strArr) == 2 {
		tokenString = strArr[1]
	}

	// Verify token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		secretKey, err := config.Get("SECRET_KEY")
		if err != nil {
			return nil, err
		}

		// Make sure that the token method conforms to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "Token is expired") {
			return nil, errors.New(accessTokenExpired)
		} else {
			return nil, errors.New(accessTokenIsWrong)
		}
	}

	// Converting to MapClaims
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New(accessTokenIsWrong)
}

// ExtractRefreshMetadata Extract refresh token claims
func ExtractRefreshMetadata(refreshString string) (jwt.MapClaims, error) {

	// get secret key.
	secretKey, err := config.Get("SECRET_KEY")
	if err != nil {
		return nil, err
	}

	// Verify token
	token, err := jwt.Parse(refreshString, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conforms to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "Token is expired") {
			return nil, errors.New(accessTokenExpired)
		} else {
			return nil, errors.New(accessTokenIsWrong)
		}
	}

	// Converting to MapClaims
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New(accessTokenIsWrong)
}

// AuthMiddleware will set user data in the context and redirect it to the next controller.
// that means we can use this function for every router and give us better user data (*ω*)
func AuthMiddleware(ctx *fiber.Ctx) error {

	// Extract and verify token metadata
	tokenAuth, err := ExtractAccessMetadata(ctx)
	if err != nil {

		return ctx.Status(http.StatusBadRequest).JSON(serializers.NewState().
			SetCode(http.StatusUnauthorized).SetMessage(err.Error()).SetStatus(false))
	}

	// Fetch userId and convert to primaryKey.
	idString := tokenAuth["userId"].(string)

	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(serializers.NewState().
			SetCode(http.StatusInternalServerError).
			SetMessage(internalServerError).
			SetDetails(err.Error()))

	}

	user := new(user2.Model)

	if err := mongo.FindOne(ctx.Context(), user, bson.D{{"_id", id}}); err != nil {
		if errors.Is(err, mongo.NoDocuments) {

			return ctx.Status(fiber.StatusUnauthorized).JSON(serializers.NewState().
				SetCode(fiber.StatusUnauthorized).SetMessage(accessTokenIsWrong).
				SetStatus(false).SetDetails(err.Error()))
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(serializers.NewState().
			SetCode(fiber.StatusInternalServerError).SetMessage(err.Error()).SetStatus(false))
	}

	// Set user in context.
	ctx.Locals("user", user)
	return ctx.Next()
}
