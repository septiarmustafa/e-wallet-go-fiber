package service

import (
	"context"
	"e-wallet/domain"
	"e-wallet/dto"
	"e-wallet/internal/util"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository  domain.UserRepository
	cacheRepository domain.CacheRepository
	emailService    domain.EmailService
}

func NewUser(userRepository domain.UserRepository, cacheRepository domain.CacheRepository, emailService domain.EmailService) domain.UserService {
	return &UserService{
		userRepository:  userRepository,
		cacheRepository: cacheRepository,
		emailService:    emailService,
	}
}

// Authenticate implements domain.UserService.
func (u *UserService) Authenticate(ctx context.Context, req dto.AuthReq) (dto.AuthRes, error) {
	user, err := u.userRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		return dto.AuthRes{}, err
	}
	if user == (domain.User{}) {
		return dto.AuthRes{}, domain.ErrAuthFailed
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return dto.AuthRes{}, domain.ErrAuthFailed
	}

	if user.EmailVerifiedAtDB.Valid {
		return dto.AuthRes{}, domain.ErrAuthFailed
	}

	token := util.GeneratorRandomString(16)

	userJson, _ := json.Marshal(user)
	_ = u.cacheRepository.Set("user:"+token, userJson)

	return dto.AuthRes{
		Token: token,
	}, nil
}

// ValidateToken implements domain.UserService.
func (u *UserService) ValidateToken(ctx context.Context, token string) (dto.UserData, error) {
	data, err := u.cacheRepository.Get("user:" + token)
	if err != nil {
		return dto.UserData{}, domain.ErrAuthFailed
	}

	var user domain.User
	_ = json.Unmarshal(data, &user)

	return dto.UserData{
		ID:       user.ID,
		FullName: user.FullName,
		Phone:    user.Phone,
		Username: user.Username,
	}, nil
}

// Register implements domain.UserService.
func (u *UserService) Register(ctx context.Context, req dto.UserRegisterReq) (dto.UserRegisterRes, error) {
	exist, err := u.userRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		return dto.UserRegisterRes{}, err
	}
	if exist != (domain.User{}) {
		return dto.UserRegisterRes{}, domain.ErrUsernameTaken
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)

	user := domain.User{
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
		Username: req.Username,
		Password: string(hashedPassword),
	}
	err = u.userRepository.Insert(ctx, &user)
	if err != nil {
		return dto.UserRegisterRes{}, err
	}

	otpCode := util.GeneratorRandomNumber(4)
	referenceID := util.GeneratorRandomString(16)

	fmt.Println("OTP:: ", otpCode)

	_ = u.emailService.Send(req.Email, "OTP Code", "OTP Anda "+otpCode)

	_ = u.cacheRepository.Set("otp:"+referenceID, []byte(otpCode))
	_ = u.cacheRepository.Set("user-ref:"+referenceID, []byte(user.Username))
	return dto.UserRegisterRes{
		ReferenceID: referenceID,
	}, nil
}

// ValidateOTP implements domain.UserService.
func (u *UserService) ValidateOTP(ctx context.Context, req dto.ValidateOtpReq) error {
	val, err := u.cacheRepository.Get("otp:" + req.ReferenceID)
	if err != nil {
		return domain.ErrOtpInvalid
	}

	otp := string(val)
	if otp != req.OTP {
		return domain.ErrOtpInvalid
	}

	val, err = u.cacheRepository.Get("user-ref:" + req.ReferenceID)
	if err != nil {
		return domain.ErrOtpInvalid
	}

	user, err := u.userRepository.FindByUsername(ctx, string(val))
	if err != nil {
		return err
	}
	user.EmailVerifiedAt = time.Now()
	_ = u.userRepository.Update(ctx, &user)

	return nil
}
