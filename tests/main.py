import os
import subprocess
import sys

IGNORE_DIRS = {
    ".git",
    "__pycache__",
    ".pytest_cache",
    ".venv",
    "venv",
    "env",
    "node_modules",
}

def has_py_files(path: str) -> bool:
    for root, dirs, files in os.walk(path):
        dirs[:] = [d for d in dirs if d not in IGNORE_DIRS]

        for f in files:
            if f.endswith(".py"):
                return True
    return False


def find_folders_with_py(base_dir: str):
    folders = []

    for root, dirs, files in os.walk(base_dir):
        dirs[:] = [d for d in dirs if d not in IGNORE_DIRS]

        if any(f.endswith(".py") for f in files):
            folders.append(root)

    return folders


def run_pytest(folder: str) -> int:
    result = subprocess.run(
        ["pytest", folder, "-v", "-o", "python_files=*.py"],
        stdout=sys.stdout,
        stderr=sys.stderr
    )
    return result.returncode


def main():
    base_dir = os.getcwd()
    folders = find_folders_with_py(base_dir)

    if not folders:
        print("No python folders found.")
        return 0

    print(f"Found {len(folders)} folder(s) with python files:")
    for f in folders:
        print(f" - {f}")

    failed = False

    for folder in folders:
        code = run_pytest(folder)
        if code != 0:
            failed = True

    if failed:
        print("\nSome folders failed.")
        return 1

    print("\nAll folders passed.")
    return 0


if __name__ == "__main__":
    sys.exit(main())