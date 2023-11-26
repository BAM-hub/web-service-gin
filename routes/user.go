package useerRouter

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type User struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
}

func GetUserByName(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	var user User
	name := c.Query("name")
	if(name == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name parameter required"})
		return
	}
	row := db.QueryRow("SELECT * FROM user WHERE name = ?", name)
	
	if err := row.Scan(&user.Id, &user.Name); err != nil {
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

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	if(user.Name == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name fieled is required"})
		return
	}
	
	stmt, err := db.Prepare("INSERT INTO user(name) VALUES(?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	row, insertErr := stmt.Exec(user.Name)

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
		  User{Name: user.Name, Id: userId}})
}

func GetUsers(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	var users []User
	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}