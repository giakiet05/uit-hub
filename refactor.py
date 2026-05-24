import os
import glob
import re

def replace_in_file(path, replacements):
    with open(path, 'r') as f:
        content = f.read()
    
    for old, new in replacements:
        content = content.replace(old, new)
        
    # Also replace gorm:"primaryKey" with nothing
    content = re.sub(r'\s*gorm:"primaryKey"\s*', ' ', content)
        
    with open(path, 'w') as f:
        f.write(content)

# 1. Update models.go
replace_in_file('apps/fake-uit-server/internal/model/models.go', [])

# 2. Update dto/types.go
replace_in_file('apps/fake-uit-server/internal/dto/types.go', [
    ('package usecase', 'package dto'),
    ('"server/models"', '"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/model"'),
    ('models.', 'model.'),
])

# 3. Update service/*.go
service_files = glob.glob('apps/fake-uit-server/internal/service/*.go')
for f in service_files:
    if f.endswith('auth_service.go') or f.endswith('student_service.go'):
        continue
        
    replacements = [
        ('package data', 'package service'),
        ('package usecase', 'package service'),
        ('"server/models"', '"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/model"'),
        ('"server/usecase"', ''),
        ('models.', 'model.'),
        ('usecase.', ''),
        ('SuccessResponse', 'dto.SuccessResponse'),
        ('AppError', 'apperror.AppError'),
        ('NewBadRequest', 'apperror.NewBadRequest'),
        ('NewUnauthorized', 'apperror.NewUnauthorized'),
        ('NewInternal', 'apperror.NewInternal'),
        ('NewNotFound', 'apperror.NewNotFound'),
    ]
    
    if f.endswith('usecase.go'):
        replacements.append(('"github.com/google/uuid"', '"github.com/google/uuid"\n\t"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"\n\t"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"'))
    
    replace_in_file(f, replacements)

# Fix specific things in usecase.go (dto prefix)
with open('apps/fake-uit-server/internal/service/usecase.go', 'r') as f:
    c = f.read()
for t in ['LoginRequest', 'StudentProfile', 'QueryYearSemester', 'TuitionFeeData', 'InsuranceData', 'TrainingPointsData', 'TranscriptRegisRequest', 'TuitionExtendRequest', 'MonthlyParkingRequest', 'GraduateRequest', 'GraduationThesisRequest', 'AssignmentSubmissionRequest', 'ConfirmLetterRequest', 'BankLoansRequest', 'TrainingPointConfirmRequest', 'LanguageCertificateUploadRequest', 'QueryRoomsAvailability', 'RoomItem', 'RoomAvailabilityData', 'ContactRequest', 'RequestStatusData']:
    c = c.replace(f' {t}', f' dto.{t}')
    c = c.replace(f'({t}', f'(dto.{t}')
    c = c.replace(f'[]{t}', f'[]dto.{t}')
with open('apps/fake-uit-server/internal/service/usecase.go', 'w') as f:
    f.write(c)

