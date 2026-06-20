package auth

import (
	"context"
	"errors"
	"time"

	"github.com/kaikaa/api/internal/token"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// errors ที่ handler ใช้ map เป็น HTTP status
var (
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type Service struct {
	users     *mongo.Collection
	shops     *mongo.Collection
	jwtSecret string
}

func NewService(db *mongo.Database, jwtSecret string) *Service {
	return &Service{
		users:     db.Collection("users"),
		shops:     db.Collection("shops"),
		jwtSecret: jwtSecret,
	}
}

// Register สร้าง shop + user (password hash ด้วย bcrypt) แล้วคืน token
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	count, err := s.users.CountDocuments(ctx, bson.M{"email": req.Email})
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	userID := primitive.NewObjectID()
	shopID := primitive.NewObjectID()

	if _, err := s.shops.InsertOne(ctx, Shop{
		ID: shopID, OwnerID: userID, Name: req.ShopName, CreatedAt: now,
	}); err != nil {
		return nil, err
	}

	if _, err := s.users.InsertOne(ctx, User{
		ID: userID, Email: req.Email, PasswordHash: string(hash), ShopID: shopID, CreatedAt: now,
	}); err != nil {
		return nil, err
	}

	tok, err := token.Generate(s.jwtSecret, userID.Hex(), shopID.Hex())
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: tok,
		Shop:  ShopInfo{ID: shopID.Hex(), Name: req.ShopName},
	}, nil
}

// Login ตรวจ email + password แล้วคืน token (ไม่บอกว่า email หรือ password ผิด เพื่อความปลอดภัย)
func (s *Service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	var user User
	err := s.users.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	var shop Shop
	if err := s.shops.FindOne(ctx, bson.M{"_id": user.ShopID}).Decode(&shop); err != nil {
		return nil, err
	}

	tok, err := token.Generate(s.jwtSecret, user.ID.Hex(), user.ShopID.Hex())
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: tok,
		Shop:  ShopInfo{ID: shop.ID.Hex(), Name: shop.Name},
	}, nil
}
