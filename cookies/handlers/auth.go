package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ralphexp/recipes-cookie/models"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type AuthHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

func NewAuthHandler(ctx context.Context, collection *mongo.Collection) *AuthHandler {
	return &AuthHandler{
		collection: collection,
		ctx:        ctx,
	}
}

// swagger:operation POST /signin auth signIn
// Login with username and password
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
//	'401':
//	    description: Invalid credentials
func (handler *AuthHandler) SignInHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h := sha256.New()
	io.Copy(h, strings.NewReader(user.Password))
	sha256sum := hex.EncodeToString(h.Sum(nil))

	fmt.Printf("user: %s, %s, %s\n", user.Username, user.Password, sha256sum)

	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"username": user.Username,
		"password": sha256sum,
	})

	if cur.Err() != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	sessionToken := xid.New().String()
	session := sessions.Default(c)
	session.Set("username", user.Username)
	session.Set("token", sessionToken)

	// write to redis and HTTP response
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "User signed in"})
}

// swagger:operation POST /refresh auth refresh
// Refresh token
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
//	'401':
//	    description: Invalid credentials
func (handler *AuthHandler) RefreshHandler(c *gin.Context) {
	session := sessions.Default(c)
	sessionToken := session.Get("token")
	sessionUser := session.Get("username")
	if sessionToken == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session cookie"})
		return
	}

	sessionToken = xid.New().String()
	session.Set("username", sessionUser.(string))
	session.Set("token", sessionToken)

	// write to redis and HTTP response
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "New session issued"})
}

// swagger:operation POST /signout auth signOut
// Signing out
// ---
// responses:
//
//	'200':
//	    description: Successful operation
func (handler *AuthHandler) SignOutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "Signed out..."})
}

func (handler *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// return a session.Session interface (which is sesssions.session object)
		session := sessions.Default(c)
		// check if gin.Context has sessions.Session object, if not
		// 1) call session.store.Get(s.request, "recipes_api") to get the Session Object
		// then assign the session object to a Session interface

		// for detail, see redistore.go:
		// func (s *RediStore) Get(r *http.Request, "recipes_api") (*sessions.Session, error)

		// 2) call session.Get("token"), get data from map
		sessionToken := session.Get("token")
		if sessionToken == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Not logged",
			})
			c.Abort()
		}
		c.Next()
	}
}
