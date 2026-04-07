"""
Testes de Timing Attack (side-channel).
"""

import time
import requests
from config import URL_TOKEN, VALID_USER, TIMEOUT


class TestTimingAttack:

    def test_tempo_consistente_entre_user_inexistente_e_senha_errada(
        self, usuario_registrado
    ):
        """
        A diferença de tempo entre 'user inexistente' e 'senha errada'
        deve ser menor que 200ms para dificultar timing attacks.
        """
        REPETICOES = 8

        tempos_inexistente = []
        for _ in range(REPETICOES):
            inicio = time.time()
            requests.post(URL_TOKEN, json={
                "username": "timing_fantasma_xyz",
                "password": "QualquerSenha!",
            }, timeout=TIMEOUT)
            tempos_inexistente.append(time.time() - inicio)

        tempos_senha_errada = []
        for _ in range(REPETICOES):
            inicio = time.time()
            requests.post(URL_TOKEN, json={
                "email": VALID_USER["email"],
                "password": "SenhaErrada!",
            }, timeout=TIMEOUT)
            tempos_senha_errada.append(time.time() - inicio)

        media_inexistente = sum(tempos_inexistente[1:]) / (REPETICOES - 1)
        media_senha_errada = sum(tempos_senha_errada[1:]) / (REPETICOES - 1)
        diferenca = abs(media_inexistente - media_senha_errada)

        THRESHOLD = 0.2  # 200ms
        assert diferenca < THRESHOLD, (
            f"Possível timing attack! Δ = {diferenca:.3f}s "
            f"(inexistente: {media_inexistente:.3f}s, "
            f"senha errada: {media_senha_errada:.3f}s)"
        )
