{ buildGoModule, fetchFromGitHub }:
buildGoModule rec {
  pname = "ory";
  version = "0.3.1";

  src = fetchFromGitHub {
    owner = "ory";
    repo = "cli";
    rev = "v${version}";
    hash = "sha256-dO595NzdkVug955dqji/ttAPb+sMGLxJftXHzHA37Lo=";
  };
  vendorHash = "sha256-H1dM/r7gJvjnexQwlA4uhJ7rUH15yg4AMRW/f0k1Ixw=";

  doCheck = false;
  installPhase = ''
    install -D $GOPATH/bin/cli $out/bin/ory
  '';
}
