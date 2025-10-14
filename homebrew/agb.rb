class Agb < Formula
  desc "Secure infrastructure for running AI-generated code"
  homepage "https://github.com/agbcloud/agbcloud-cli"
  url "https://github.com/agbcloud/agbcloud-cli/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "05bf8e62463c840f285bbb17b584ce3e55f2f8b07841ee280624444db5443229"
  license "MIT"
  head "https://github.com/agbcloud/agbcloud-cli.git", branch: "main"

  depends_on "go" => :build

  def install
    # Set build variables matching the Makefile
    version = self.version
    # Use embedded git commit from build time (since tarball has no .git directory)
    git_commit = "fcf862c"
    build_date = Time.now.utc.strftime("%Y-%m-%dT%H:%M:%SZ")

    # Set Go proxy for better network connectivity (especially in China)
    ENV["GOPROXY"] = "https://goproxy.cn,https://goproxy.io,https://proxy.golang.org,direct"
    ENV["GOSUMDB"] = "sum.golang.google.cn"
    ENV["GO111MODULE"] = "on"

    # Build flags matching your Makefile LDFLAGS (with optimization)
    ldflags = %W[
      -s
      -w
      -X github.com/agbcloud/agbcloud-cli/cmd.Version=#{version}
      -X github.com/agbcloud/agbcloud-cli/cmd.GitCommit=#{git_commit}
      -X github.com/agbcloud/agbcloud-cli/cmd.BuildDate=#{build_date}
    ]

    # Build from source using Go
    system "go", "build", *std_go_args(ldflags: ldflags), "."
  end

  test do
    # Test that binary is executable
    assert_predicate bin/"agb", :executable?

    # Test version command
    version_output = shell_output("#{bin}/agb version 2>&1")
    assert_match version.to_s, version_output

    # Test help command
    help_output = shell_output("#{bin}/agb --help")
    assert_match "agb", help_output
    assert_match "help", help_output
  end
end
