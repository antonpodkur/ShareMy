package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/antonpodkur/ShareMy/config"
	"github.com/antonpodkur/ShareMy/internal/auth"
	"github.com/antonpodkur/ShareMy/internal/models"
	"github.com/antonpodkur/ShareMy/pkg/db"
	"github.com/antonpodkur/ShareMy/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type authUsecase struct {
    cfg *config.Config
    mongoClient *mongo.Client
    ctx context.Context
}

func NewAuthUsecase(cfg *config.Config, mongoClient *mongo.Client, ctx context.Context) auth.Usecase {
    return &authUsecase{
        cfg: cfg,
        mongoClient: mongoClient,
        ctx: ctx,
    }
}

func (au *authUsecase) SignUp(user *models.SignUpInput) (*models.UserDBResponse, error) {
    usersCollection := db.OpenCollection(au.mongoClient, au.cfg, "users")

    user.CreatedAt = time.Now()
    user.UpdatedAt = user.CreatedAt
    user.Email = strings.ToLower(user.Email)
    user.PasswordConfirm = ""
    user.Verified = true
    user.Role = "user"

    hashedPassword, _ := utils.HashPassword(user.Password)
    user.Password = hashedPassword

    res, err := usersCollection.InsertOne(au.ctx, &user)
    if err != nil {
        if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
            return nil, errors.New("user with that email already exists")
        }
        return nil, err
    }

    opt := options.Index()
    opt.SetUnique(true)
    index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: opt}

    if _, err := usersCollection.Indexes().CreateOne(au.ctx, index); err != nil {
        return nil, errors.New("could not create index for email")
    }

    var newUser *models.UserDBResponse
    query := bson.M{"_id": res.InsertedID}

    err = usersCollection.FindOne(au.ctx, query).Decode(&newUser)
    if err == nil {
        return nil, err
    }

    return newUser, nil
}

func (au *authUsecase) SignIn(_ *models.SignInInput) (*models.UserDBResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (au *authUsecase) GetUserById(id string) (*models.UserDBResponse, error) {
    usersCollection := db.OpenCollection(au.mongoClient, au.cfg, "users") 
    oid, _ := primitive.ObjectIDFromHex(id)

    var user *models.UserDBResponse

    query := bson.M{ "_id": oid }
    err := usersCollection.FindOne(au.ctx, query).Decode(&user)

    if err != nil {
        if err == mongo.ErrNoDocuments {
            return &models.UserDBResponse{}, err
        }
        return nil, err
    }

    return user, nil
}

func (au *authUsecase) GetUserByEmail(email string) (*models.UserDBResponse, error) {
    usersCollection := db.OpenCollection(au.mongoClient, au.cfg, "users") 

    var user *models.UserDBResponse

    query := bson.M{ "email": strings.ToLower(email) }
    err := usersCollection.FindOne(au.ctx, query).Decode(&user)

    if err != nil {
        if err == mongo.ErrNoDocuments {
            return &models.UserDBResponse{}, err
        }
        return nil, err
    }

    return user, nil

}

