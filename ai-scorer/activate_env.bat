@echo off
REM Windows batch script to activate conda environment and run AI scorer

echo Activating conda environment...
call conda activate ai-blockchain-scorer

if errorlevel 1 (
    echo ERROR: Failed to activate conda environment.
    echo Make sure you have created it first: conda env create -f environment.yml
    pause
    exit /b 1
)

echo Starting AI scorer...
python app/api.py

pause

