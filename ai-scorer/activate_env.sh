#!/bin/bash
# Linux/Mac script to activate conda environment and run AI scorer

echo "Activating conda environment..."
source "$(conda info --base)/etc/profile.d/conda.sh"
conda activate ai-blockchain-scorer

if [ $? -ne 0 ]; then
    echo "ERROR: Failed to activate conda environment."
    echo "Make sure you have created it first: conda env create -f environment.yml"
    exit 1
fi

echo "Starting AI scorer..."
python app/api.py

