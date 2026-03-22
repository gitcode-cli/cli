#!/bin/bash
# Post-install script for gc package

set -e

# Update shell completions if available
if command -v update-shell-completions &> /dev/null; then
    update-shell-completions
fi

# Inform user about setup
echo "GitCode CLI (gc) has been installed successfully!"
echo ""
echo "To get started:"
echo "  1. Run 'gc auth login' to authenticate"
echo "  2. Set your token: export GC_TOKEN=your_token"
echo ""
echo "Documentation: https://gitcode.com/help/cli"

exit 0