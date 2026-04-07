"""
Configurações compartilhadas entre todos os módulos de teste.
"""

import uuid

BASE_URL = "http://localhost:7171/api/v1/auth"
URL_REGISTER = f"{BASE_URL}/register"
URL_LOGIN = f"{BASE_URL}/login"
URL_TOKEN = f"{BASE_URL}/token"
TIMEOUT = 10

# Sufixo único por execução para evitar conflitos
_UID = uuid.uuid4().hex[:8]

VALID_USER = {
    "username": f"testuser_{_UID}",
    "email": f"testuser_{_UID}@example.com",
    "password": "S3nh@F0rte!2026",
}
