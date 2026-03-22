# Homebrew Formula for GitCode CLI
# This file is a template - Goreleaser will generate the actual formula

class Gc < Formula
  desc "GitCode CLI - Command line tool for GitCode"
  homepage "https://gitcode.com"
  version "0.1.0"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/gitcode-com/gitcode-cli/releases/download/v#{version}/gc_#{version}_darwin_amd64.tar.gz"
      sha256 ""
    end
    on_arm do
      url "https://github.com/gitcode-com/gitcode-cli/releases/download/v#{version}/gc_#{version}_darwin_arm64.tar.gz"
      sha256 ""
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/gitcode-com/gitcode-cli/releases/download/v#{version}/gc_#{version}_linux_amd64.tar.gz"
      sha256 ""
    end
    on_arm do
      url "https://github.com/gitcode-com/gitcode-cli/releases/download/v#{version}/gc_#{version}_linux_arm64.tar.gz"
      sha256 ""
    end
  end

  def install
    bin.install "gc"
    bash_completion.install "completions/gc.bash" => "gc"
    zsh_completion.install "completions/gc.zsh" => "_gc"
    fish_completion.install "completions/gc.fish"
  end

  test do
    assert_match "gc version", shell_output("#{bin}/gc version")
  end
end