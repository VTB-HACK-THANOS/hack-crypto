package storage

import (
	"context"

	"github.com/VTB-HACK-THANOS/hack-crypto/models"
)

func (s *Storage) InsertNft(ctx context.Context, nft *models.Nft) (*models.Nft, error) {
	tx, err := s.contextTransaction(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := tx.
		NewSelect().
		Model(nft).
		Exec(ctx); err != nil {
		return nil, err
	}

	return nft, nil
}
