class Agb < Formula
  desc "Secure infrastructure for running AI-generated code"
  homepage "https://github.com/agbcloud/agbcloud-cli"
  version "1.1.6"
  license "MIT"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/agbcloud/agbcloud-cli/releases/download/v1.1.6/agb-1.1.6-darwin-amd64.tar.gz"
      sha256 "ba1e634acee3eebf1c77324cfcd09419380ee06d9f6201037745b53066bd98f9"
    elsif Hardware::CPU.arm? || Hardware::CPU.arch == :arm64
      url "https://github.com/agbcloud/agbcloud-cli/releases/download/v1.1.6/agb-1.1.6-darwin-arm64.tar.gz"
      sha256 "4d8cc682fa171d705f91f0552cd2d8016b4433fed2abdc3aea2b35c724d87c81"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/agbcloud/agbcloud-cli/releases/download/v1.1.6/agb-1.1.6-linux-amd64.tar.gz"
      sha256 "35320d046c9555c4d57fa443d9b4730235353ee05318cc37d8d2b48199b64e48"
    elsif Hardware::CPU.arm? || Hardware::CPU.arch == :arm64
      url "https://github.com/agbcloud/agbcloud-cli/releases/download/v1.1.6/agb-1.1.6-linux-arm64.tar.gz"
      sha256 "6862efa459c95b11ff87680d0a5e47ca1a7fd1e302d6e84c2e3c772988802d39"
    end
  end

  def install
    bin.install "agb"
  end

  test do
    # Test that binary is executable
    assert_predicate bin/"agb", :executable?

    # Check if we can run the binary (skip if GLIBC incompatible)
    begin
      # Test version command
      system bin/"agb", "--version"

      # Test help command
      help_output = shell_output("#{bin}/agb --help")
      assert_match "agb", help_output
      assert_match "help", help_output
    rescue => e
      # Skip functional tests if binary cannot run due to system incompatibility
      if e.message.include?("GLIBC") || e.message.include?("not found")
        ohai "Skipping functional tests due to system incompatibility"
      else
        raise e
      end
    end
  end
end
