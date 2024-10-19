package handlers

import (
	"go_proj_example/internal/config"
	"go_proj_example/internal/database"
	"net/http"
	"path/filepath"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var cfg = config.LoadConfig()
var JwtSecret = []byte(cfg.Jwt_key)

type Claims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}

func LoadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	baseLayout, err := filepath.Glob(templatesDir + "/base/base.html")
	if err != nil {
		panic(err.Error())
	}

	authLayouts, err := filepath.Glob(templatesDir + "/auth/*.html")
	if err != nil {
		panic(err.Error())
	}
	restLayouts, err := filepath.Glob(templatesDir + "/*.html")
	if err != nil {
		panic(err.Error())
	}
	for _, templateFile := range authLayouts {
		r.AddFromFiles(filepath.Base(templateFile), templateFile)
	}
	for _, templateFile := range restLayouts {
		layoutCopy := make([]string, len(baseLayout))
		copy(layoutCopy, baseLayout)
		layoutCopy = append(layoutCopy, templateFile)
		r.AddFromFiles(filepath.Base(templateFile), layoutCopy...)
	}
	return r
}

func HomePage(c *gin.Context) {
	c.HTML(200, "home.html", nil)
}

func LoginPage(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{})
}

func SignupPage(c *gin.Context) {
	c.HTML(200, "signup.html", gin.H{})
}

func NotFound(c *gin.Context) {
	c.HTML(200, "404.html", gin.H{})
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func Register(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")
	secret := c.PostForm("secret")

	var user database.User
	result := database.DB.Where("Login = ?", login).First(&user)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "There is already a user with such name"})
		return
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	err = database.CreateUser(login, hashedPassword, secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func generateJWT(login string) (string, error) {
	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &Claims{
		Login: login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

func Login(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")

	var user database.User
	result := database.DB.Where("Login = ?", login).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := generateJWT(login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie("jwtToken", token, 3600, "/", "", false, false)
	c.Redirect(http.StatusFound, "/")
}

func SignOut(c *gin.Context) {
	c.SetCookie("jwtToken", "", -1, "/", "", false, false)
	c.Redirect(http.StatusFound, "/login")
}
