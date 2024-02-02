package controllers

import (
	"context"
	"fmt"
	"golang-Restaurant-Management/database"
	"golang-Restaurant-Management/helper"
	"golang-Restaurant-Management/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc{
	return func(c *gin.Context){
       var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	   recordPerPage, err:= strconv.Atoi(c.Query("recordPerPage"))
	   if err!=nil || recordPerPage<1{
		recordPerPage=10
	   }

       page, err:= strconv.Atoi(c.Query("page"))
	   if err!=nil || page<1{
		page=1
	   }

       startIndex := (page-1)*recordPerPage
	   startIndex, err= strconv.Atoi(c.Query("startIndex"))

	   matchStage:=bson.D{{"$match",bson.D{{}}}}
	   groupStage:=bson.D{{"$group", bson.D{{"_id", bson.D{{"_id","null"}}},{"total_count",bson.D{{"$sum", 1}}}, {"data",bson.D{{"$push", "$$ROOT"}}}}}}
	   projectStage:= bson.D{
		{
	   "$project",bson.D{
		{"_id",0},
		{"total_count",1},
		{"user_items",bson.D{{"$slice",[]interface{}{"$data", startIndex,recordPerPage}}}}}}}
        
		result,err:=userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
      
        defer cancel()

		if err!=nil{
        c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while listing users"})
		return
		}	

		var allUsers []bson.M

		if err = result.All(ctx, &allUsers); err!=nil{
			
			log.Fatal(err)
			return
		}

		c.JSON(http.StatusOK,allUsers[0])
	}
}

func GetUser() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel= context.WithTimeout(context.Background(), 100*time.Second)

		userId:=c.Param("user_id")

		var user models.User

		err:=userCollection.FindOne(ctx, bson.M{"user_id":userId}).Decode(&user)
       
		defer cancel()
		 
		if err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while listing user details"})
		}

		c.JSON(http.StatusOK,user)
	}
}

func SignUp() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User

		if err:= c.BindJSON(&user); err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}

		validationErr:=validate.Struct(user)
		if validationErr!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validationErr.Error()})
			return
		}

		countEmail, err:=userCollection.CountDocuments(ctx, bson.M{"email":user.Email})
		defer cancel()

		if err!=nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while checking mail"})
			return
		}

		password:=HashPassword(*user.Password)

		user.Password = &password

		countPhone, err:= userCollection.CountDocuments(ctx, bson.M{"phone":user.Phone})

		defer cancel()

		if err!=nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while checking phone number"})
			return
		}
        
		if countEmail>0 || countPhone>0{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"this email or phone number already exists"})
			return
		}

		user.Created_at, _ =time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ =time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, refreshToken,_:=helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *&user.User_id)
		user.Token = &token
		user.Refresh_Token = &refreshToken

		resultInsertion, insertErr:=userCollection.InsertOne(ctx, user)
		if insertErr!=nil{
			msg:=fmt.Sprintf("user was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}
        defer cancel()
		c.JSON(http.StatusOK,resultInsertion)

	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err:= c.BindJSON(&user); err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}

		err:=userCollection.FindOne(ctx, bson.M{"email":user.Email}).Decode(&foundUser)
		defer cancel()
		if err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"user not found"})
			return
		}

		passwordIsValid, msg:= VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid!=true{
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}

		token, refreshToken, _:=helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *&foundUser.User_id)
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		c.JSON(http.StatusOK, foundUser)

	}
}

func HashPassword(password string) string{
    bytes, err:=bcrypt.GenerateFromPassword([]byte(password),14)
	if err!=nil{
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providePassword string) (bool,string){
	err:= bcrypt.CompareHashAndPassword([]byte(providePassword),[]byte(userPassword))
	check:=true
	msg:=""

	if err!=nil{
		msg=fmt.Sprintf("login or password is incorrect")
		check=false
	}
	return check, msg
}