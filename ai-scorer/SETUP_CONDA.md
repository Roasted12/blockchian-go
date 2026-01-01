# Setting Up Conda Environment for AI Scorer

## Option 1: Create New Conda Environment (Recommended)

### Step 1: Create the environment from environment.yml
```bash
cd ai-scorer
conda env create -f environment.yml
```

### Step 2: Activate the environment
```bash
conda activate ai-blockchain-scorer
```

### Step 3: Verify installation
```bash
python --version  # Should show Python 3.9
pip list          # Should show flask, numpy, scikit-learn, etc.
```

### Step 4: Run the AI scorer
```bash
python app/api.py
```

## Option 2: Use Existing Conda Environment

If you already have a conda environment you want to use:

### Step 1: Activate your environment
```bash
conda activate your-env-name
```

### Step 2: Install dependencies
```bash
cd ai-scorer
pip install -r requirements.txt
```

## Setting Python Interpreter in Your IDE (Cursor/VS Code)

### Method 1: Command Palette
1. Press `Ctrl+Shift+P` (Windows) or `Cmd+Shift+P` (Mac)
2. Type "Python: Select Interpreter"
3. Choose the conda environment:
   - Look for `ai-blockchain-scorer` or your conda env name
   - Path will look like: `C:\Users\YourName\anaconda3\envs\ai-blockchain-scorer\python.exe`

### Method 2: Bottom Status Bar
1. Click on the Python version in the bottom-right status bar
2. Select your conda environment from the list

### Method 3: Create .vscode/settings.json
Create a `.vscode` folder in the `ai-scorer` directory with:

```json
{
    "python.defaultInterpreterPath": "${env:CONDA_PREFIX}/python.exe",
    "python.terminal.activateEnvironment": true
}
```

Or specify the full path:
```json
{
    "python.defaultInterpreterPath": "C:\\Users\\YourName\\anaconda3\\envs\\ai-blockchain-scorer\\python.exe"
}
```

## Verify Environment is Selected

1. Open a Python file (e.g., `app/api.py`)
2. Check the bottom-right corner - it should show your conda environment
3. Open terminal in IDE - it should auto-activate the conda environment

## Troubleshooting

### If conda command not found:
- Add Anaconda to PATH:
  - Windows: Add `C:\Users\YourName\anaconda3\Scripts` to PATH
  - Or use Anaconda Prompt instead of regular terminal

### If IDE doesn't see conda environments:
- Make sure Anaconda is installed correctly
- Restart the IDE after installing conda
- Manually specify the Python path in IDE settings

### To list all conda environments:
```bash
conda env list
```

### To remove and recreate environment:
```bash
conda env remove -n ai-blockchain-scorer
conda env create -f environment.yml
```

