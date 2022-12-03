{ buildGoModule, fetchFromGitHub, go_1_19 }:
buildGoModule.override { go = go_1_19; } rec {
  pname = "ory";
  version = "0.1.35";

  src = fetchFromGitHub {
    owner = "ory";
    repo = "cli";
    rev = "v${version}";
    sha256 = "sha256-I0STYR2KcRbQ/wv/Rb+vRm8gtWT6YdT8wV88yWCTjPc=";
  };
  vendorSha256 = "sha256-ds5SI5WmfW9n6yZ4fQAsFLEP88m1JsVrAJuxH53mcuE=";

  doCheck = false;
  installPhase = ''
    install -D $GOPATH/bin/cli $out/bin/ory
  '';
}
