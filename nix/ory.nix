{ buildGoModule, fetchFromGitHub, go_1_17 }:
buildGoModule.override { go = go_1_17; } rec {
  pname = "ory";
  version = "0.1.22";

  src = fetchFromGitHub {
    owner = "ory";
    repo = "cli";
    rev = "v${version}";
    sha256 = "sha256-7wl2fegTGfxhN9DZCmfho5qZCcGT3ssUkm6SlwRSU+M=";
  };
  vendorSha256 = "sha256-Ye5lNgWXhInVrJbXjwFIBFdrZsZfIdLIVYXFz4SBJNE=";

  doCheck = false;
  installPhase = ''
    install -D $GOPATH/bin/cli $out/bin/ory
  '';
}
