package useerRouter

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)


func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(passowrd string, hash string) (bool) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passowrd))
	return err == nil
}

func GetUserByName(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	var user UserResponse
	name := c.Query("name")
	if(name == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name parameter required"})
		return
	}
	row := db.QueryRow("SELECT id, name, email, created_at FROM user WHERE name = ?", name)
	
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	var user User 
	hash, hashErr := HashPassword(user.Password) 
	if hashErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	user.Password = hash
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSOn"})
		return
	}

	if(user.Name == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name fieled is required"})
		return
	}

	if(user.Email == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email fieled is required"})
		return
	}

	if(user.Password == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password fieled is required"})
		return
	}
	
	stmt, err := db.Prepare("INSERT INTO user(name, email, password) VALUES(?,?,?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	row, insertErr := stmt.Exec(user.Name, user.Email, hash)

	if insertErr != nil {
		if insertErr.(*mysql.MySQLError).Number == 1062 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}
	
	defer stmt.Close()

	userId, err := row.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user created successfully", "user":
		  User{Name: user.Name, Id: userId, Email: user.Email}})
}

func GetUsers(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	var users []UserResponse
	rows, err := db.Query("SELECT id, name, email, created_at FROM user")
	if err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error, " + err.Error()})
		return
	}

	for rows.Next() {
		var user UserResponse
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error" + err.Error()})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}