"""
Fixtures compartilhadas entre todos os módulos de teste (pytest).
"""

import pytest
import requests
from config import BASE_URL, URL_REGISTER, URL_TOKEN, VALID_USER, TIMEOUT


@pytest.fixture(scope="session", autouse=True)
def verificar_servidor():
    """Garante que o servidor está acessível antes de rodar os testes."""
    try:
        requests.get(BASE_URL, timeout=5)
    except requests.exceptions.ConnectionError:
        pytest.exit(
            "❌ Servidor não está acessível em "
            f"{BASE_URL}. Inicie a API antes de rodar os testes."
        )


@pytest.fixture(scope="session")
def usuario_registrado():
    """
    Registra um usuário de teste uma única vez para toda a sessão.
    Retorna os dados do registro (RegisterOutput).
    """
    resp = requests.post(URL_REGISTER, json=VALID_USER, timeout=TIMEOUT)
    assert resp.status_code in (200, 201), (
        f"Falha ao registrar usuário de teste: {resp.status_code} — {resp.text}"
    )
    data = resp.json()
    assert data["success"] is True
    return data


@pytest.fixture()
def token_valido(usuario_registrado):
    """
    Faz login e retorna um token JWT válido para testes que precisam
    de autenticação.
    """
    resp = requests.post(URL_TOKEN, json={
        "email": VALID_USER["email"],
        "password": VALID_USER["password"],
    }, timeout=TIMEOUT)
    assert resp.status_code == 200
    data = resp.json()
    assert data.get("token"), "Token não retornado no login"
    return data["token"]
