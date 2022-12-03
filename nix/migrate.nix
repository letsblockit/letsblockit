{ buildGoModule, fetchFromGitHub, go_1_19 }:
buildGoModule.override { go = go_1_19; } rec {
  pname = "migrate";
  version = "4.15.2";

  src = fetchFromGitHub {
    owner = "golang-migrate";
    repo = "migrate";
    rev = "v${version}";
    sha256 = "sha256-nVR6zMG/a4VbGgR9a/6NqMNYwFTifAZW3F6rckvOEJM=";
  };
  vendorSha256 = "sha256-lPNPl6fqBT3XLQie9z93j91FLtrMjKbHnXUQ6b4lDb4=";

  tags = [ "pgx" ];
  subPackages = [ "cmd/migrate" ];
}
