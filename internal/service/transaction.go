package service

import (
	"context"
	"e-wallet/domain"
	"e-wallet/dto"
	"e-wallet/internal/util"
	"encoding/json"
	"fmt"
	"time"
)

type transactionService struct {
	accountRespository    domain.AccountRepository
	transactionRepository domain.TransactionRepository
	cacheRepository       domain.CacheRepository
	notificationService   domain.NotificationService
	// notificationRepository domain.NotificationRepository
	// hub                    *dto.Hub
}

func NewTransaction(accountRespository domain.AccountRepository, transactionRepository domain.TransactionRepository, cacheRepository domain.CacheRepository, notificationService domain.NotificationService) domain.TransactionService {
	return &transactionService{
		accountRespository:    accountRespository,
		transactionRepository: transactionRepository,
		cacheRepository:       cacheRepository,
		notificationService:   notificationService,
		// notificationRepository domain.NotificationRepository
		// hub                    *dto.Hub
	}
}

// TransferInquiry implements domain.TransactionService.
func (t *transactionService) TransferInquiry(ctx context.Context, req dto.TransferInquiryReq) (dto.TransferInquiryRes, error) {

	// find my account
	user := ctx.Value("x-user").(dto.UserData)
	myAccount, err := t.accountRespository.FindByUserID(ctx, user.ID)
	if err != nil {
		return dto.TransferInquiryRes{}, err
	}

	if myAccount == (domain.Account{}) {
		return dto.TransferInquiryRes{}, domain.ErrAccountNotFound
	}

	// find dof account
	dofAcount, err := t.accountRespository.FindByAccountNumber(ctx, req.AccountNumber)
	if err != nil {
		return dto.TransferInquiryRes{}, err
	}

	if dofAcount == (domain.Account{}) {
		return dto.TransferInquiryRes{}, domain.ErrAccountNotFound
	}

	// validate my account balance < transfer amount
	if myAccount.Balance < req.Amount {
		return dto.TransferInquiryRes{}, domain.ErrInsufficientBalance
	}

	inquiryKey := util.GeneratorRandomNumber(32)

	// store inquiry
	jsonData, _ := json.Marshal(req)
	_ = t.cacheRepository.Set(inquiryKey, jsonData)

	return dto.TransferInquiryRes{
		InquiryKey: inquiryKey,
	}, nil
}

// TransferExecute implements domain.TransactionService.
func (t *transactionService) TransferExecute(ctx context.Context, req dto.TransferExecuteReq) error {
	// get inquiry
	val, err := t.cacheRepository.Get(req.InquiryKey)
	if err != nil {
		return domain.ErrInquiryNotFound
	}

	var reqInquiry dto.TransferInquiryReq
	_ = json.Unmarshal(val, &reqInquiry)

	if reqInquiry == (dto.TransferInquiryReq{}) {
		return domain.ErrInquiryNotFound
	}

	// check my account
	user := ctx.Value("x-user").(dto.UserData)
	myAccount, err := t.accountRespository.FindByUserID(ctx, user.ID)
	if err != nil {
		return err
	}

	// check account number dof / receiver
	dofAcount, err := t.accountRespository.FindByAccountNumber(ctx, reqInquiry.AccountNumber)
	if err != nil {
		return err
	}

	// debit for sof / sender
	debitTransaction := domain.Transaction{
		Account_id:          myAccount.ID,
		SofNumber:           myAccount.AccountNumber,
		DofNumber:           dofAcount.AccountNumber,
		TransactionType:     "D",
		Amount:              reqInquiry.Amount,
		TransactionDatetime: time.Now(),
	}

	err = t.transactionRepository.Insert(ctx, &debitTransaction)
	if err != nil {
		return err
	}

	// credit for dof / receiver
	creditTransaction := domain.Transaction{
		Account_id:          dofAcount.ID,
		SofNumber:           myAccount.AccountNumber,
		DofNumber:           dofAcount.AccountNumber,
		TransactionType:     "C",
		Amount:              reqInquiry.Amount,
		TransactionDatetime: time.Now(),
	}

	err = t.transactionRepository.Insert(ctx, &creditTransaction)
	if err != nil {
		return err
	}

	// balance transaction
	myAccount.Balance -= reqInquiry.Amount
	err = t.accountRespository.Update(ctx, &myAccount)
	if err != nil {
		return err
	}

	dofAcount.Balance += reqInquiry.Amount
	err = t.accountRespository.Update(ctx, &dofAcount)
	if err != nil {
		return err
	}

	go t.notificationAfterTransfer(myAccount, dofAcount, reqInquiry.Amount)
	return nil
}

func (t transactionService) notificationAfterTransfer(sofAccount domain.Account, dofAccount domain.Account, amount float64) {
	data := map[string]string{
		"amount": fmt.Sprintf("%.2f", amount),
	}

	_ = t.notificationService.Insert(context.Background(), sofAccount.UserId, "TRANSFER", data)
	_ = t.notificationService.Insert(context.Background(), dofAccount.UserId, "TRANSFER_DESTINATION", data)

	// notificationSender := domain.Notification{
	// 	UserID:    sofAccount.UserId,
	// 	Title:     "Transfer berhasil",
	// 	Body:      fmt.Sprintf("Transfer senilai %.2f", amount),
	// 	IsRead:    0,
	// 	Status:    1,
	// 	CreatedAt: time.Now(),
	// }

	// notificationReceiver := domain.Notification{
	// 	UserID:    dofAccount.UserId,
	// 	Title:     "Dana diterima",
	// 	Body:      fmt.Sprintf("Dana diterima senilai %.2f", amount),
	// 	IsRead:    0,
	// 	Status:    1,
	// 	CreatedAt: time.Now(),
	// }

	// _ = t.notificationRepository.Insert(context.Background(), &notificationSender)
	// if channel, ok := t.hub.NotificationChannel[sofAccount.ID]; ok {
	// 	channel <- dto.NotificationData{
	// 		ID:        notificationSender.ID,
	// 		Title:     notificationSender.Title,
	// 		Body:      notificationSender.Body,
	// 		Status:    notificationSender.Status,
	// 		IsRead:    notificationSender.IsRead,
	// 		CreatedAt: notificationSender.CreatedAt,
	// 	}
	// }

	// _ = t.notificationRepository.Insert(context.Background(), &notificationReceiver)
	// if channel, ok := t.hub.NotificationChannel[dofAccount.ID]; ok {
	// 	channel <- dto.NotificationData{
	// 		ID:        notificationReceiver.ID,
	// 		Title:     notificationReceiver.Title,
	// 		Body:      notificationReceiver.Body,
	// 		Status:    notificationReceiver.Status,
	// 		IsRead:    notificationReceiver.IsRead,
	// 		CreatedAt: notificationReceiver.CreatedAt,
	// 	}
	// }

}
