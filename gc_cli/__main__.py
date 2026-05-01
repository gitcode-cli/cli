"""Allow running GitCode CLI with ``python -m gc_cli``."""

import sys

from gc_cli.wrapper import main


if __name__ == "__main__":
    sys.exit(main())
