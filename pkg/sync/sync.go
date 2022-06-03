package sync

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"

	"errors"

	"github.com/gordallott/mcpsync/pkg/api"
)

var (
	errUnavailable = errors.New("unavailable")
	errUnchanged   = errors.New("unchanged")
)

var cardBuf = make([]byte, 0, 128*1024)

const (
	charTreeNode          = "├── "
	charTreeNodeEnd       = "└── "
	charTreeNodeIndent    = "│   "
	charTreeNodeIndentEnd = "    "
)

func syncCard(ctx context.Context, remoteReader io.ReadCloser, targetDest string) error {
	defer remoteReader.Close()

	buf := bytes.NewBuffer(cardBuf[:0])
	if _, err := io.Copy(buf, remoteReader); err != nil {
		return err
	}

	targetFile, err := os.ReadFile(targetDest)
	skipTargetCheck := false
	if errors.Is(err, os.ErrNotExist) {
		// file doesn't exist, no need to check consistancy
		skipTargetCheck = true
	}

	if !skipTargetCheck {
		shaRemote := sha256.Sum256(buf.Bytes())
		shaTarget := sha256.Sum256(targetFile)
		if bytes.Equal(shaRemote[:], shaTarget[:]) {
			return errUnchanged
		}
	}

	return ioutil.WriteFile(targetDest, buf.Bytes(), 0644)
}

func SyncOnce(ctx context.Context, targetIP, syncDest string, pollFrequency time.Duration) error {
	_, err := api.GetVersion(ctx, targetIP)
	if err != nil {
		return fmt.Errorf("%w: %s", errUnavailable, err)
	}

	fmt.Printf("Syncing %s => %s\n", targetIP, syncDest)

	cards, err := api.GetCards(ctx, targetIP)
	if err != nil {
		return err
	}

	fmt.Printf("Found %d cards\n", len(cards))

	for i, card := range cards {
		indent := charTreeNodeIndent
		node := charTreeNode
		if i == len(cards)-1 {
			indent = charTreeNodeIndentEnd
			node = charTreeNodeEnd
		}

		state, err := api.GetState(ctx, targetIP)
		if err != nil {
			return err
		}

		if state.GameID == card.GameID {
			fmt.Printf(node+"%s(%s): Skipping, Can not sync currently active saves\n", card.GameID, card.Name)
			continue
		}
		cardSyncDest := path.Join(syncDest, card.GameID)
		fmt.Printf(node+"%s: (%s) => %q\n", card.GameID, card.Name, cardSyncDest)

		err = os.MkdirAll(cardSyncDest, 0755)
		if err != nil {
			return fmt.Errorf("Can not make directory %q: %w", cardSyncDest, err)
		}

		for i := uint(1); i <= 8; i++ {
			filename := api.GetCardFilename(card.GameID, i)
			node := charTreeNode
			if i == 8 {
				node = charTreeNodeEnd
			}

			fmt.Printf(indent+node+"%q: ", filename)

			reader, err := api.GetCard(ctx, targetIP, card.GameID, 1)
			if err != nil {
				return err
			}

			err = syncCard(ctx, reader, path.Join(cardSyncDest, filename))
			switch {
			case errors.Is(err, errUnchanged):
				fmt.Printf("unchanged\n")
			case err != nil:
				fmt.Printf("error: %s\n", err)
			default:
				fmt.Printf("success!\n")
			}
		}

	}

	return nil
}

func Sync(ctx context.Context, targetIP, syncDest string, pollFrequency time.Duration) error {
	if targetIP == "" {
		return fmt.Errorf("targetIP must be set")
	}

	if syncDest == "" {
		return fmt.Errorf("syncDest must be set")
	}

	ticker := time.NewTicker(pollFrequency)
	defer ticker.Stop()

	err := SyncOnce(ctx, targetIP, syncDest, pollFrequency)
	switch {
	case errors.Is(err, errUnavailable):
		fmt.Printf("Sync unavailable: %s\n", err)
	case err != nil:
		return err
	}

main:
	for {
		select {
		case <-ctx.Done():
			break main
		case <-ticker.C:
			err := SyncOnce(ctx, targetIP, syncDest, pollFrequency)
			switch {
			case errors.Is(err, errUnavailable):
				fmt.Printf("Sync unavailable: %s\n", err)
			case err != nil:
				return err
			}

		}
	}

	return nil
}
