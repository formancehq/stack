package mangopay

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/formancehq/go-libs/pointer"
	"github.com/formancehq/payments/internal/connectors/plugins/public/mangopay/client"
	"github.com/formancehq/payments/internal/models"
)

func (p Plugin) createBankAccount(ctx context.Context, req models.CreateBankAccountRequest) (models.CreateBankAccountResponse, error) {
	userID := models.ExtractNamespacedMetadata(req.BankAccount.Metadata, client.MangopayUserIDMetadataKey)
	if userID == "" {
		return models.CreateBankAccountResponse{}, fmt.Errorf("missing userID in bank account metadata")
	}

	ownerAddress := client.OwnerAddress{
		AddressLine1: models.ExtractNamespacedMetadata(req.BankAccount.Metadata, models.BankAccountOwnerAddressLine1MetadataKey),
		AddressLine2: models.ExtractNamespacedMetadata(req.BankAccount.Metadata, models.BankAccountOwnerAddressLine2MetadataKey),
		City:         models.ExtractNamespacedMetadata(req.BankAccount.Metadata, models.BankAccountOwnerCityMetadataKey),
		Region:       models.ExtractNamespacedMetadata(req.BankAccount.Metadata, models.BankAccountOwnerRegionMetadataKey),
		PostalCode:   models.ExtractNamespacedMetadata(req.BankAccount.Metadata, models.BankAccountOwnerPostalCodeMetadataKey),
		Country: func() string {
			if req.BankAccount.Country == nil {
				return ""
			}
			return *req.BankAccount.Country
		}(),
	}

	var mangopayBankAccount *client.BankAccount
	if req.BankAccount.IBAN != nil {
		req := &client.CreateIBANBankAccountRequest{
			OwnerName:    req.BankAccount.Name,
			OwnerAddress: &ownerAddress,
			IBAN:         *req.BankAccount.IBAN,
			BIC: func() string {
				if req.BankAccount.SwiftBicCode == nil {
					return ""
				}
				return *req.BankAccount.SwiftBicCode
			}(),
			Tag: models.ExtractNamespacedMetadata(req.BankAccount.Metadata, client.MangopayTagMetadataKey),
		}

		var err error
		mangopayBankAccount, err = p.client.CreateIBANBankAccount(ctx, userID, req)
		if err != nil {
			return models.CreateBankAccountResponse{}, err
		}
	} else {
		if req.BankAccount.Country == nil {
			req.BankAccount.Country = pointer.For("")
		}
		switch *req.BankAccount.Country {
		case "US":
			if req.BankAccount.AccountNumber == nil {
				return models.CreateBankAccountResponse{}, fmt.Errorf("missing account number in bank account metadata")
			}

			req := &client.CreateUSBankAccountRequest{
				OwnerName:          req.BankAccount.Name,
				OwnerAddress:       &ownerAddress,
				AccountNumber:      *req.BankAccount.AccountNumber,
				ABA:                models.ExtractNamespacedMetadata(req.BankAccount.Metadata, client.MangopayABAMetadataKey),
				DepositAccountType: models.ExtractNamespacedMetadata(req.BankAccount.Metadata, client.MangopayDepositAccountTypeMetadataKey),
				Tag:                models.ExtractNamespacedMetadata(req.BankAccount.Metadata, client.MangopayTagMetadataKey),
			}

			var err error
			mangopayBankAccount, err = p.client.CreateUSBankAccount(ctx, userID, req)
			if err != nil {
				return models.CreateBankAccountResponse{}, err
			}

		case "CA":
			if req.BankAccount.AccountNumber == nil {
				return models.CreateBankAccountResponse{}, fmt.Errorf("missing account number in bank account metadata")
			}
			req := &client.CreateCABankAccountRequest{
				OwnerName:         req.BankAccount.Name,
				OwnerAddress:      &ownerAddress,
				AccountNumber:     *req.BankAccount.AccountNumber,
				InstitutionNumber: models.ExtractNamespacedMetadata(req.BankAccount.Metadata, client.MangopayInstitutionNumberMetadataKey),
				BranchCode:        models.ExtractNamespacedMetadata(req.BankAccount.Metadata, client.MangopayBranchCodeMetadataKey),
				BankName:          models.ExtractNamespacedMetadata(req.BankAccount.Metadata, client.MangopayBankNameMetadataKey),
				Tag:               models.ExtractNamespacedMetadata(req.BankAccount.Metadata, client.MangopayTagMetadataKey),
			}

			var err error
			mangopayBankAccount, err = p.client.CreateCABankAccount(ctx, userID, req)
			if err != nil {
				return models.CreateBankAccountResponse{}, err
			}

		case "GB":
			if req.BankAccount.AccountNumber == nil {
				return models.CreateBankAccountResponse{}, fmt.Errorf("missing account number in bank account metadata")
			}

			req := &client.CreateGBBankAccountRequest{
				OwnerName:     req.BankAccount.Name,
				OwnerAddress:  &ownerAddress,
				AccountNumber: *req.BankAccount.AccountNumber,
				SortCode:      models.ExtractNamespacedMetadata(req.BankAccount.Metadata, client.MangopaySortCodeMetadataKey),
				Tag:           models.ExtractNamespacedMetadata(req.BankAccount.Metadata, client.MangopayTagMetadataKey),
			}

			var err error
			mangopayBankAccount, err = p.client.CreateGBBankAccount(ctx, userID, req)
			if err != nil {
				return models.CreateBankAccountResponse{}, err
			}

		default:
			if req.BankAccount.AccountNumber == nil {
				return models.CreateBankAccountResponse{}, fmt.Errorf("missing account number in bank account metadata")
			}

			req := &client.CreateOtherBankAccountRequest{
				OwnerName:     req.BankAccount.Name,
				OwnerAddress:  &ownerAddress,
				AccountNumber: *req.BankAccount.AccountNumber,
				BIC: func() string {
					if req.BankAccount.SwiftBicCode == nil {
						return ""
					}
					return *req.BankAccount.SwiftBicCode
				}(),
				Country: *req.BankAccount.Country,
				Tag:     models.ExtractNamespacedMetadata(req.BankAccount.Metadata, client.MangopayTagMetadataKey),
			}

			var err error
			mangopayBankAccount, err = p.client.CreateOtherBankAccount(ctx, userID, req)
			if err != nil {
				return models.CreateBankAccountResponse{}, err
			}
		}
	}

	var account models.PSPAccount
	if mangopayBankAccount != nil {
		raw, err := json.Marshal(mangopayBankAccount)
		if err != nil {
			return models.CreateBankAccountResponse{}, err
		}

		account = models.PSPAccount{
			Reference: mangopayBankAccount.ID,
			CreatedAt: time.Unix(mangopayBankAccount.CreationDate, 0),
			Name:      &mangopayBankAccount.OwnerName,
			Metadata: map[string]string{
				"user_id": userID,
			},
			Raw: raw,
		}

	}

	return models.CreateBankAccountResponse{
		RelatedAccount: account,
	}, nil
}
