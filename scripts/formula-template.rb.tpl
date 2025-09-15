class AgbcloudAT<%= sanitized_timestamp %> < Formula
  desc "AgbCloud CLI - Test Build <%= timestamp %>"
  homepage "https://your-company.com/agbcloud"
  version "<%= version %>"

  if OS.mac?
    if Hardware::CPU.arm?
      url "https://agbcloud-internal.oss-cn-hangzhou.aliyuncs.com/agbcloud-<%= version %>-darwin-arm64.tar.gz"
      sha256 "<%= darwin_arm64_sha256 %>"
    else
      url "https://agbcloud-internal.oss-cn-hangzhou.aliyuncs.com/agbcloud-<%= version %>-darwin-amd64.tar.gz"
      sha256 "<%= darwin_amd64_sha256 %>"
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://agbcloud-internal.oss-cn-hangzhou.aliyuncs.com/agbcloud-<%= version %>-linux-arm64.tar.gz"
      sha256 "<%= linux_arm64_sha256 %>"
    else
      url "https://agbcloud-internal.oss-cn-hangzhou.aliyuncs.com/agbcloud-<%= version %>-linux-amd64.tar.gz"
      sha256 "<%= linux_amd64_sha256 %>"
    end
  end

  def install
    bin.install "agbcloud"
  end

  test do
    system "#{bin}/agbcloud", "version"
  end

  def caveats
    <<~EOS
      This is a test build of AgbCloud CLI.
      Build timestamp: <%= timestamp %>
      Git commit: <%= git_commit %>
      
      To switch between versions:
        brew unlink agbcloud@<%= sanitized_timestamp %>
        brew link agbcloud@other-version
    EOS
  end
end 