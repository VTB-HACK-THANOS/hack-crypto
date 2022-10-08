package crytowallet

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/VTB-HACK-THANOS/hack-crypto/models"
)

const (
	methodCreateNewWallet = "/v1/wallets/new"
	methodBalance         = "/v1/wallets/%s/balance" // public key
	methodHistory         = "/v1/wallets/%s/history" // public key
	methodTransferRuble   = "/v1/transfers/ruble"
	defaultTimeout        = 60 * time.Second
)

type Service struct {
	endpoint string
	client   *http.Client
}

func New(endpoint string) (*Service, error) {
	s := &Service{endpoint: endpoint, client: &http.Client{Timeout: defaultTimeout}}
	return s, nil
}

// CreateWallet creates new wallet.
func (s *Service) CreateWallet() (*models.WalletCredentials, error) {
	req, err := http.NewRequest(http.MethodPost, s.endpoint+methodCreateNewWallet, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	wallet := &models.WalletCredentials{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

// Balance shows user's balance.
func (s *Service) Balance(publicKey string) (*models.Balance, error) {
	if publicKey == "" {
		return nil, errors.New("empty public key")
	}

	requestURL := s.endpoint + fmt.Sprintf(methodBalance, publicKey)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	balance := &models.Balance{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, balance); err != nil {
		return nil, err
	}

	return balance, nil

}

func (s *Service) Transfer(fromPrivateKey, toPublicKey string, amount float64) error {
	type Request struct {
		FromPrivateKey string  `json:"fromPrivateKey "`
		ToPublicKey    string  `json:"toPublicKey"`
		Amount         float64 `json:"amount"`
	}
	requestURL := s.endpoint + methodTransferRuble

	r := &Request{
		FromPrivateKey: fromPrivateKey,
		ToPublicKey:    toPublicKey,
		Amount:         amount,
	}

	reqJson, err := json.Marshal(r)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(reqJson))
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("new transaction: %s\n", string(body))

	return nil
}

// History returns history of user's transactions.
func (s *Service) History(ctx context.Context, publicKey string, page, offset int, sort string) ([]*models.History, error) {
	type Request struct {
		Page   int
		Offset int
		Sort   string
	}

	type Response struct {
		History []*models.History `json:"history"`
	}

	r := &Request{
		Page:   page,
		Offset: offset,
		Sort:   sort,
	}

	reqJson, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	requestURL := s.endpoint + fmt.Sprintf(methodHistory, publicKey)

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	history := &Response{}
	if err := json.Unmarshal(body, history); err != nil {
		return nil, err
	}

	return history.History, nil
}
