# Usage (the repo name is "gitcode-cli", not "homebrew-gitcode-cli",
# so you must pass the explicit URL when tapping):
#
#   brew tap zhongqixiao/gitcode-cli https://github.com/ZhongqiXiao/gitcode-cli
#   brew install zhongqixiao/gitcode-cli/gc
#
# After the release workflow runs, install simplifies to:
#   brew install gc
class Gc < Formula
  desc "GitCode CLI - Command-line tool for GitCode"
  homepage "https://gitcode.com"
  version "0.9.1"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/ZhongqiXiao/gitcode-cli/releases/download/v#{version}/gc_#{version}_darwin_amd64.tar.gz"
      sha256 "DARWIN_AMD64_SHA256"
    end
    on_arm do
      url "https://github.com/ZhongqiXiao/gitcode-cli/releases/download/v#{version}/gc_#{version}_darwin_arm64.tar.gz"
      sha256 "DARWIN_ARM64_SHA256"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/ZhongqiXiao/gitcode-cli/releases/download/v#{version}/gc_#{version}_linux_amd64.tar.gz"
      sha256 "LINUX_AMD64_SHA256"
    end
    on_arm do
      url "https://github.com/ZhongqiXiao/gitcode-cli/releases/download/v#{version}/gc_#{version}_linux_arm64.tar.gz"
      sha256 "LINUX_ARM64_SHA256"
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
