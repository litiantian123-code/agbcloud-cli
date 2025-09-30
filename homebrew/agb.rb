class Agb < Formula
  desc "Secure infrastructure for running AI-generated code"
  homepage "https://github.com/litiantian123-code/agbcloud-cli"
  url "https://github.com/litiantian123-code/agbcloud-cli/archive/refs/tags/v1.1.2.tar.gz"
  sha256 "57299a5ff70ecdb250a3b839a5da6cfd980bd2ece618a201644c08bf2df54b58"
  license "MIT"
  head "https://github.com/litiantian123-code/agbcloud-cli.git", branch: "main"

  depends_on "go" => :build

  def install
    # Set build variables matching the Makefile
    version = self.version
    # Handle git commit safely (archive tarball doesn't have .git directory)
    git_commit = begin
      Utils.safe_popen_read("git", "rev-parse", "--short", "HEAD").chomp
    rescue
      "unknown"
    end
    git_commit = Utils.safe_popen_read("git", "rev-parse", "--short", "HEAD").chomp
    build_date = Time.now.utc.strftime("%Y-%m-%dT%H:%M:%SZ")
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
