"""
Wrapper module for GitCode CLI binary.

This module provides a Python entry point that calls the appropriate
pre-compiled binary based on the current platform.
"""

import os
import platform
import subprocess
import sys
from pathlib import Path


def get_binary_name() -> str:
    """Get the binary name for the current platform."""
    system = platform.system().lower()
    machine = platform.machine().lower()

    # Map machine architecture
    arch_map = {
        "x86_64": "amd64",
        "amd64": "amd64",
        "aarch64": "arm64",
        "arm64": "arm64",
    }
    arch = arch_map.get(machine, "amd64")

    # Map system to binary name
    if system == "linux":
        return f"gc-linux-{arch}"
    elif system == "darwin":
        return f"gc-darwin-{arch}"
    elif system == "windows":
        return "gc-windows-amd64.exe"
    else:
        raise RuntimeError(f"Unsupported platform: {system} {machine}")


def get_binary_path() -> Path:
    """Get the path to the binary for the current platform."""
    package_dir = Path(__file__).parent
    binary_name = get_binary_name()
    binary_path = package_dir / "bin" / binary_name

    if not binary_path.exists():
        raise FileNotFoundError(
            f"Binary not found for your platform: {binary_path}\n"
            f"Supported platforms: linux-amd64, linux-arm64, darwin-amd64, darwin-arm64, windows-amd64"
        )

    return binary_path


def main() -> int:
    """Main entry point for the gc command."""
    try:
        binary_path = get_binary_path()

        # Make sure the binary is executable
        if not os.access(binary_path, os.X_OK):
            binary_path.chmod(0o755)

        # Run the binary with all arguments
        result = subprocess.run(
            [str(binary_path)] + sys.argv[1:],
            cwd=os.getcwd(),
        )

        return result.returncode

    except FileNotFoundError as e:
        print(f"Error: {e}", file=sys.stderr)
        return 1
    except Exception as e:
        print(f"Error running gc: {e}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main())