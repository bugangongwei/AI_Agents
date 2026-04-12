package service

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/bugangongwei/job_loves_me/backend/internal/model"
	"github.com/bugangongwei/job_loves_me/backend/internal/repository"
	"github.com/ledongthuc/pdf"
)

type ResumeService struct {
	storagePath string
}

func NewResumeService() *ResumeService {
	path := os.Getenv("RESUME_STORAGE_PATH")
	if path == "" {
		path = "./storage/resumes"
	}
	// Create storage directory if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
	return &ResumeService{storagePath: path}
}

func (s *ResumeService) SaveAndParseResume(userID uint, fileName string, fileContent io.Reader) (*model.Resume, error) {
	// Create user storage directory
	userPath := filepath.Join(s.storagePath, string(rune(userID)))
	if _, err := os.Stat(userPath); os.IsNotExist(err) {
		os.MkdirAll(userPath, os.ModePerm)
	}

	// Always save as latest_resume.pdf for a simple versioning logic
	dstPath := filepath.Join(userPath, "latest_resume.pdf")
	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// Use a TeeReader to both save to file and keep content for parsing
	var buf bytes.Buffer
	tr := io.TeeReader(fileContent, &buf)

	if _, err := io.Copy(dst, tr); err != nil {
		return nil, err
	}

	// Parse PDF
	rawText, err := s.ExtractTextFromPDF(dstPath)
	if err != nil {
		return nil, err
	}

	// Deactivate previous latest resume
	repository.DB.Model(&model.Resume{}).Where("user_id = ? AND is_latest = ?", userID, true).Update("is_latest", false)

	// Create new resume entry
	resume := &model.Resume{
		UserID:   userID,
		FilePath: dstPath,
		RawText:  rawText,
		IsLatest: true,
	}

	if err := repository.DB.Create(resume).Error; err != nil {
		return nil, err
	}

	return resume, nil
}

func (s *ResumeService) ExtractTextFromPDF(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}

func (s *ResumeService) GetLatestResume(userID uint) (*model.Resume, error) {
	var resume model.Resume
	if err := repository.DB.Where("user_id = ? AND is_latest = ?", userID, true).First(&resume).Error; err != nil {
		return nil, err
	}
	return &resume, nil
}
