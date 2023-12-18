package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"time"

	"URLReminderBot/libraries/myErr"
	"URLReminderBot/storage"
)

const defaultPerm = 0774

var ErrNoPSavedPages = errors.New("no saved page")

type Storage struct {
	BasePath string
}

func New(basePath string) Storage {
	return Storage{BasePath: basePath}
}

func (s *Storage) Save(p *storage.Page) (err error) {
	defer func() { err = myErr.WrapIfErr("can't save: ", err) }()

	fPath := filepath.Join(s.BasePath, p.UserName)

	if err = os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(p)
	if err != nil {
		return err
	}

	fPath = path.Join(fPath, fName)

	file, err := os.Create(fPath)
	defer func() { _ = file.Close() }()

	if err = gob.NewEncoder(file).Encode(p); err != nil {
		return err
	}

	return nil
}

func (s *Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = myErr.WrapIfErr("can't pick a file", err) }()

	fPath := filepath.Join(s.BasePath, userName)

	files, err := os.ReadDir(fPath)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, ErrNoPSavedPages
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(fPath, file.Name()))
}

func (s *Storage) Remove(p *storage.Page) error {
	fName, err := fileName(p)
	if err != nil {
		return myErr.Wrap("can't remove file", err)
	}

	finalPath := filepath.Join(s.BasePath, p.UserName, fName)

	if err = os.Remove(finalPath); err != nil {
		msg := fmt.Sprintf("can't remove file %s", finalPath)
		return myErr.Wrap(msg, err)
	}

	return nil
}

func (s *Storage) IsExist(p *storage.Page) (bool, error) {
	fName, err := fileName(p)
	if err != nil {
		return false, myErr.Wrap("can't check if file exists", err)
	}

	finalPath := filepath.Join(s.BasePath, p.UserName, fName)

	switch _, err = os.Stat(finalPath); {
	case errors.Is(err, os.ErrExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", finalPath)
		return false, myErr.Wrap(msg, err)
	}

	return true, nil
}

func (s *Storage) decodePage(filePath string) (p *storage.Page, err error) {
	defer func() { err = myErr.WrapIfErr("can't decode page: ", err) }()
	f, err := os.Open(filePath)
	if err != nil {
		return nil, myErr.Wrap("can't decode page", err)
	}
	defer func() { _ = f.Close() }()

	if err = gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, err
	}

	return p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
