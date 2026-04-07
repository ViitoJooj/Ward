"""
Testes para os endpoints POST /login e POST /token.
"""

import requests
from config import URL_TOKEN, VALID_USER, TIMEOUT


class TestLogin:

    # ---------- Fluxos válidos ----------

    def test_login_com_email_e_senha(self, usuario_registrado):
        """Login com email + senha corretos deve retornar token."""
        resp = requests.post(URL_TOKEN, json={
            "email": VALID_USER["email"],
            "password": VALID_USER["password"],
        }, timeout=TIMEOUT)
        assert resp.status_code == 200
        data = resp.json()
        assert data["success"] is True
        assert data.get("token"), "Token ausente na resposta de login"
        assert len(data["token"]) > 20, "Token parece muito curto"

    def test_login_com_username_e_senha(self, usuario_registrado):
        """Login com username + senha corretos deve retornar token."""
        resp = requests.post(URL_TOKEN, json={
            "username": VALID_USER["username"],
            "password": VALID_USER["password"],
        }, timeout=TIMEOUT)
        assert resp.status_code == 200
        data = resp.json()
        assert data["success"] is True
        assert data.get("token")

    def test_login_retorna_dados_do_usuario(self, usuario_registrado):
        """O login deve retornar os dados do usuário (UserData)."""
        resp = requests.post(URL_TOKEN, json={
            "email": VALID_USER["email"],
            "password": VALID_USER["password"],
        }, timeout=TIMEOUT)
        data = resp.json()
        user = data.get("data", {})
        assert user.get("id"), "ID do usuário ausente"
        assert user.get("username") == VALID_USER["username"]
        assert user.get("email") == VALID_USER["email"]

    def test_login_nao_retorna_senha(self, usuario_registrado):
        """A resposta de login NUNCA deve conter a senha."""
        resp = requests.post(URL_TOKEN, json={
            "email": VALID_USER["email"],
            "password": VALID_USER["password"],
        }, timeout=TIMEOUT)
        assert VALID_USER["password"] not in resp.text

    # ---------- Fluxos inválidos ----------

    def test_login_senha_errada(self, usuario_registrado):
        """Login com senha incorreta deve ser rejeitado."""
        resp = requests.post(URL_TOKEN, json={
            "email": VALID_USER["email"],
            "password": "SenhaErrada!999",
        }, timeout=TIMEOUT)
        assert resp.status_code in (400, 401, 403)
        data = resp.json()
        assert data["success"] is False

    def test_login_usuario_inexistente(self):
        """Login com usuário inexistente deve ser rejeitado."""
        resp = requests.post(URL_TOKEN, json={
            "username": "usuario_fantasma_xyz",
            "password": "QualquerSenha123!",
        }, timeout=TIMEOUT)
        assert resp.status_code in (400, 401, 403, 404)

    def test_login_sem_senha(self, usuario_registrado):
        """Login sem campo de senha deve ser rejeitado."""
        resp = requests.post(URL_TOKEN, json={
            "email": VALID_USER["email"],
        }, timeout=TIMEOUT)
        assert resp.status_code in (400, 401, 422)

    def test_login_sem_identificador(self):
        """Login sem username nem email deve ser rejeitado."""
        resp = requests.post(URL_TOKEN, json={
            "password": "QualquerSenha123!",
        }, timeout=TIMEOUT)
        assert resp.status_code in (400, 401, 422)

    def test_login_body_vazio(self):
        """Login com body vazio deve ser rejeitado."""
        resp = requests.post(URL_TOKEN, json={}, timeout=TIMEOUT)
        assert resp.status_code in (400, 401, 422)

    def test_login_campos_string_vazia(self):
        """Login com strings vazias deve ser rejeitado."""
        resp = requests.post(URL_TOKEN, json={
            "username": "",
            "email": "",
            "password": "",
        }, timeout=TIMEOUT)
        assert resp.status_code in (400, 401, 422)
