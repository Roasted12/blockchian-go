"""
AI SCORING SERVICE – ADVISORY ML MODELS

This service provides advisory scoring for:
- Transaction anomaly detection (IsolationForest)
- Fee adequacy estimation (simple regression)
- Peer reliability scoring (future)

Important:
- This is ADVISORY ONLY
- Does NOT affect blockchain consensus
- Helps prioritize transactions and detect suspicious activity
- If service is down, blockchain continues operating normally
"""

from flask import Flask, request, jsonify
from flask_cors import CORS
import numpy as np
from sklearn.ensemble import IsolationForest
import joblib
import os
import logging

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = Flask(__name__)
CORS(app)  # Allow cross-origin requests

# Global model storage
tx_anomaly_model = None
model_path = "models/tx_anomaly_model.pkl"


def load_or_create_model():
    """
    Load existing model or create a new one.
    
    IsolationForest is an unsupervised model, so we can:
    - Train on normal transactions
    - Detect anomalies (outliers) automatically
    - No labels needed!
    """
    global tx_anomaly_model
    
    if os.path.exists(model_path):
        logger.info(f"Loading model from {model_path}")
        tx_anomaly_model = joblib.load(model_path)
    else:
        logger.info("Creating new IsolationForest model")
        # Create model with default parameters
        # contamination: expected proportion of anomalies (0.1 = 10%)
        # random_state: for reproducibility
        tx_anomaly_model = IsolationForest(
            contamination=0.1,
            random_state=42,
            n_estimators=100
        )
        # Train on dummy data (in production, train on historical data)
        # For now, we'll use a simple initialization
        dummy_data = np.array([
            [2, 2, 100.0, 100.0, 1.0, 0.01, 0.95, 1],  # Normal transaction
            [1, 1, 50.0, 50.0, 0.5, 0.01, 1.0, 1],     # Normal transaction
            [3, 3, 200.0, 200.0, 2.0, 0.01, 0.9, 1],   # Normal transaction
        ])
        tx_anomaly_model.fit(dummy_data)
        
        # Save model
        os.makedirs("models", exist_ok=True)
        joblib.dump(tx_anomaly_model, model_path)
        logger.info(f"Model saved to {model_path}")


@app.route('/health', methods=['GET'])
def health():
    """
    Health check endpoint.
    
    Returns:
        JSON with service status
    """
    return jsonify({
        "status": "healthy",
        "service": "ai-scorer",
        "model_loaded": tx_anomaly_model is not None
    })


@app.route('/score/tx', methods=['POST'])
def score_transaction():
    """
    Score a transaction for anomaly detection and fee adequacy.
    
    Request body:
        {
            "num_inputs": 2,
            "num_outputs": 2,
            "total_input": 100.0,
            "total_output": 99.0,
            "fee": 1.0,
            "fee_rate": 0.01,
            "change_ratio": 0.99,
            "input_diversity": 1
        }
    
    Response:
        {
            "anomaly_score": 0.2,  # 0.0 = normal, 1.0 = highly anomalous
            "fee_adequacy": 0.8,    # 0.0 = low fee, 1.0 = high fee
            "message": "Transaction scored successfully"
        }
    """
    try:
        # Get features from request
        data = request.get_json()
        if not data:
            return jsonify({"error": "No JSON data provided"}), 400
        
        # Extract features (in same order as model expects)
        features = np.array([[
            data.get("num_inputs", 0),
            data.get("num_outputs", 0),
            data.get("total_input", 0.0),
            data.get("total_output", 0.0),
            data.get("fee", 0.0),
            data.get("fee_rate", 0.0),
            data.get("change_ratio", 0.0),
            data.get("input_diversity", 0)
        ]])
        
        # Predict anomaly score
        # IsolationForest returns: -1 (anomaly) or 1 (normal)
        # We convert to 0.0-1.0 scale where 1.0 = most anomalous
        prediction = tx_anomaly_model.predict(features)[0]
        anomaly_score = 1.0 if prediction == -1 else 0.0
        
        # Get anomaly score (decision function gives confidence)
        # Lower values = more anomalous
        decision_score = tx_anomaly_model.decision_function(features)[0]
        # Normalize to 0.0-1.0 (inverse: higher = more anomalous)
        # decision_score is typically in range [-0.5, 0.5]
        normalized_score = max(0.0, min(1.0, 0.5 - decision_score))
        anomaly_score = normalized_score
        
        # Calculate fee adequacy (simple heuristic)
        # Higher fee rate = better adequacy
        fee_rate = data.get("fee_rate", 0.0)
        fee_adequacy = min(1.0, max(0.0, fee_rate * 100))  # Scale fee_rate
        
        # If fee is very low, reduce adequacy
        fee = data.get("fee", 0.0)
        if fee < 0.1:
            fee_adequacy *= 0.5
        
        response = {
            "anomaly_score": float(anomaly_score),
            "fee_adequacy": float(fee_adequacy),
            "message": "Transaction scored successfully"
        }
        
        logger.info(f"Scored transaction: anomaly={anomaly_score:.2f}, fee={fee_adequacy:.2f}")
        return jsonify(response)
        
    except Exception as e:
        logger.error(f"Error scoring transaction: {e}")
        return jsonify({"error": str(e)}), 500


@app.route('/score/peer', methods=['POST'])
def score_peer():
    """
    Score a peer for reliability (future implementation).
    
    For now, returns a placeholder response.
    """
    return jsonify({
        "reliability_score": 0.5,
        "message": "Peer scoring not yet implemented"
    })


if __name__ == '__main__':
    # Load or create model on startup
    load_or_create_model()
    
    # Start server
    port = int(os.getenv('PORT', 5000))
    logger.info(f"Starting AI scorer service on port {port}")
    app.run(host='0.0.0.0', port=port, debug=True)

